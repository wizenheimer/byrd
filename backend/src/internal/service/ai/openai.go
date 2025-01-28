// ./src/internal/service/ai/openai.go
package ai

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"sync"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type openAIService struct {
	client  *openai.Client
	logger  *logger.Logger
	builder *ProfileBuilder
}

func NewOpenAIService(apiKey string, logger *logger.Logger) (AIService, error) {
	client := openai.NewClient(
		option.WithAPIKey(
			apiKey,
		),
		option.WithMaxRetries(
			3,
		),
	)

	fieldRegistry := NewFieldRegistry()
	builder := NewProfileBuilder(fieldRegistry)

	logger.Debug("creating new openAI service", zap.Any("available_fields", fieldRegistry.ListAvailableFields()))

	s := openAIService{
		client:  client,
		logger:  logger.WithFields(map[string]interface{}{"module": "ai_service"}),
		builder: builder,
	}

	return &s, s.validate()
}

func (s *openAIService) SummarizeChanges(ctx context.Context, changeList []*models.DynamicChanges) ([]models.CategoryChange, error) {
	if len(changeList) == 0 {
		return []models.CategoryChange{}, nil
	}

	changes, err := models.MergeDynamicChanges(changeList)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	numCategories := len(changes.Fields)
	resultChan := make(chan result, numCategories)
	var wg sync.WaitGroup

	processedCategories := make(map[string]bool)
	var mu sync.Mutex

	// Fixed type conversion in goroutines
	for category, changesList := range changes.Fields {
		wg.Add(1)
		go func(cat string, list interface{}) {
			defer wg.Done()

			interfaceList := list.([]interface{})
			stringList := make([]string, len(interfaceList))
			for i, v := range interfaceList {
				stringList[i] = v.(string)
			}

			processCategoryAsync(ctx, s.client, cat, stringList, resultChan)
		}(category, changesList)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	response := models.ChangeResponse{
		Changes: make([]models.CategoryChange, 0, numCategories),
	}

	for res := range resultChan {
		if res.err != nil {
			s.logger.Error("error processing category", zap.Error(res.err))
			continue
		}

		mu.Lock()
		processedCategories[res.summary.Category] = true
		mu.Unlock()

		// Fixed type conversion here too
		interfaceList := changes.Fields[res.summary.Category].([]interface{})
		stringList := make([]string, len(interfaceList))
		for i, v := range interfaceList {
			stringList[i] = v.(string)
		}

		categoryChange := models.CategoryChange{
			Category: res.summary.Category,
			Summary:  res.summary.Summary,
			Changes:  stringList,
		}

		response.Changes = append(response.Changes, categoryChange)
	}

	for category, changesList := range changes.Fields {
		mu.Lock()
		if !processedCategories[category] {
			// And here
			interfaceList := changesList.([]interface{})
			stringList := make([]string, len(interfaceList))
			for i, v := range interfaceList {
				stringList[i] = v.(string)
			}

			categoryChange := models.CategoryChange{
				Category: category,
				Summary:  getFallbackSummary(category, stringList),
				Changes:  stringList,
			}

			response.Changes = append(response.Changes, categoryChange)
			s.logger.Debug("fallback summary", zap.String("category", category), zap.Strings("changes", stringList))
		}
		mu.Unlock()
	}

	return response.Changes, nil
}

// AnalyzeContentDifferences analyzes the content differences between two versions of a URL
func (s *openAIService) AnalyzeContentDifferences(ctx context.Context, version1, version2 string, fields []string) (*models.DynamicChanges, error) {
	s.logger.Debug("analyzing content differences", zap.Int("version1", len(version1)), zap.Int("version2", len(version2)), zap.Strings("requested_fields", fields))

	profileRequest := ProfileRequest{
		Name:        "competitor_updates",
		Description: "Carefully compare these two versions of content, identify and surface changes",
		FieldNames:  fields,
	}

	profile, err := s.builder.BuildProfile(profileRequest, true)
	if err != nil {
		return nil, err
	}

	s.logger.Debug("analyzing content differences", zap.String("version1", version1), zap.String("version2", version2), zap.Any("profile", profile.Fields), zap.Strings("requested_fields", fields))

	chat, err := s.prepareTextCompletion(ctx, version1, version2, profile)
	if err != nil {
		return nil, err
	}

	return s.parseCompletion(chat)
}

// AnalyzeVisualDifferences analyzes the visual differences between two screenshots
func (s *openAIService) AnalyzeVisualDifferences(ctx context.Context, screenshot1, screenshot2 image.Image, fields []string) (*models.DynamicChanges, error) {
	s.logger.Debug("analyzing visual differences", zap.Any("requested_fields", fields))

	profileRequest := ProfileRequest{
		Name:        "competitor_updates",
		Description: "Carefully compare and contrast visual changes in the webpage",
		FieldNames:  fields,
	}

	profile, err := s.builder.BuildProfile(profileRequest, true)
	if err != nil {
		return nil, err
	}

	chat, err := s.prepareImageCompletion(ctx, screenshot1, screenshot2, profile)
	if err != nil {
		return nil, err
	}

	return s.parseCompletion(chat)
}

func (s *openAIService) Close() {
	s.logger.Debug("closing openAI service")
}

// Check if the client is valid by sending a test request
func (s *openAIService) validate() error {
	_, err := s.client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("Yo"),
		}),
		Model:     openai.F(openai.ChatModelGPT4o),
		MaxTokens: openai.Int(3),
	})

	return err
}

func (s *openAIService) prepareTextCompletion(ctx context.Context, version1, version2 string, profile models.Profile) (*openai.ChatCompletion, error) {
	opts := s.prepareCompareOptions(&profile)

	userPrompt := fmt.Sprintf("Compare these two versions of content and identify changes:\n\nVersion 1:\n%s\n\nVersion 2:\n%s", version1, version2)

	schema := models.GenerateDynamicSchema(profile.Fields)
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("dynamic_intelligence_tracking"),
		Description: openai.F("Track changes in specified intelligence categories"),
		Schema:      openai.F(schema),
		Strict:      openai.Bool(true),
	}

	chat, err := s.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(opts.SystemPrompt),
			openai.UserMessage(userPrompt),
		}),
		ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
			openai.ResponseFormatJSONSchemaParam{
				Type:       openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
				JSONSchema: openai.F(schemaParam),
			},
		),
		Model:       openai.F(opts.Model),
		Temperature: openai.F(opts.Temperature),
		MaxTokens:   openai.F(opts.MaxTokens),
	})

	return chat, err
}

func (s *openAIService) prepareImageCompletion(ctx context.Context, version1, version2 image.Image, profile models.Profile) (*openai.ChatCompletion, error) {
	// convert images to base64
	version1Base64, err := imageToBase64URL(version1)
	if err != nil {
		return nil, ErrConvertingImageToBase64
	}

	version2Base64, err := imageToBase64URL(version2)
	if err != nil {
		return nil, ErrConvertingImageToBase64
	}

	opts := s.prepareCompareOptions(&profile)

	userPrompt := "Carefully compare these two versions of images and identify changes"

	schema := models.GenerateDynamicSchema(profile.Fields)
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("dynamic_intelligence_tracking"),
		Description: openai.F("Track changes in specified intelligence categories"),
		Schema:      openai.F(schema),
		Strict:      openai.Bool(true),
	}

	chat, err := s.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(opts.SystemPrompt),
			openai.UserMessageParts(
				openai.TextPart(userPrompt),
				openai.ImagePart(version1Base64),
				openai.ImagePart(version2Base64),
			),
		}),
		ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
			openai.ResponseFormatJSONSchemaParam{
				Type:       openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
				JSONSchema: openai.F(schemaParam),
			},
		),
		Model:       openai.F(opts.Model),
		Temperature: openai.F(opts.Temperature),
		MaxTokens:   openai.F(opts.MaxTokens),
	})

	return chat, err
}

func (s *openAIService) parseCompletion(chat *openai.ChatCompletion) (*models.DynamicChanges, error) {
	if chat.Choices[0].Message.Refusal != "" {
		return nil, ErrEncounteredRefusal
	}

	changes := &models.DynamicChanges{
		Fields: make(map[string]interface{}),
	}

	err := json.Unmarshal([]byte(chat.Choices[0].Message.Content), changes)
	if err != nil {
		return nil, ErrParsingChanges
	}

	return changes, nil
}

func (s *openAIService) prepareCompareOptions(profile *models.Profile) models.CompareOptions {
	return models.CompareOptions{
		SystemPrompt: models.BuildCompetitorSystemPrompt(profile.Fields),
		Model:        openai.ChatModelGPT4oMini,
		Temperature:  0.7,
		MaxTokens:    2048,
	}
}

func imageToBase64URL(img image.Image) (string, error) {
	// Create a buffer to store the image
	var buf bytes.Buffer

	// Encode the image in PNG format to the buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return "", ErrEncodingImage
	}

	// Convert the buffer bytes to base64 string
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Add prefix for base64 URL
	base64Str = "data:image/png;base64," + base64Str

	return base64Str, nil
}
