package ai

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	svc "github.com/wizenheimer/byrd/src/internal/interfaces/service"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/errs"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type openAIService struct {
	client  *openai.Client
	logger  *logger.Logger
	builder *ProfileBuilder
}

func NewOpenAIService(apiKey string, logger *logger.Logger) (svc.AIService, error) {
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

// AnalyzeContentDifferences analyzes the content differences between two versions of a URL
func (s *openAIService) AnalyzeContentDifferences(ctx context.Context, version1, version2 string, fields []string) (*models.DynamicChanges, errs.Error) {
	pErr := errs.New()
	profileRequest := ProfileRequest{
		Name:        "competitor_updates",
		Description: "Carefully compare these two versions of content, identify and surface changes",
		FieldNames:  fields,
	}

	profile, err := s.builder.BuildProfile(profileRequest, true)
	if err != nil {
		pErr.Add(svc.ErrBuildingProfile, map[string]interface{}{
			"profileRequest": profileRequest,
		})
		return nil, pErr
	}

	s.logger.Debug("analyzing content differences", zap.String("version1", version1), zap.String("version2", version2), zap.Any("profile", profile.Fields), zap.Strings("requested_fields", fields))

	chat, err := s.prepareTextCompletion(ctx, version1, version2, profile)
	if err != nil && err.HasErrors() {
		pErr.Merge(err)
		return nil, pErr
	}

	return s.parseCompletion(chat)
}

// AnalyzeVisualDifferences analyzes the visual differences between two screenshots
func (s *openAIService) AnalyzeVisualDifferences(ctx context.Context, screenshot1, screenshot2 image.Image, fields []string) (*models.DynamicChanges, errs.Error) {
	dErr := errs.New()
	profileRequest := ProfileRequest{
		Name:        "competitor_updates",
		Description: "Carefully compare and contrast visual changes in the webpage",
		FieldNames:  fields,
	}

	profile, err := s.builder.BuildProfile(profileRequest, true)
	if err != nil {
		dErr.Add(svc.ErrBuildingProfile, map[string]interface{}{
			"profileRequest": profileRequest,
		})
		return nil, dErr
	}

	chat, err := s.prepareImageCompletion(ctx, screenshot1, screenshot2, profile)
	if err != nil && err.HasErrors() {
		dErr.Merge(err)
		return nil, dErr
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

func (s *openAIService) prepareTextCompletion(ctx context.Context, version1, version2 string, profile models.Profile) (*openai.ChatCompletion, errs.Error) {
	dErr := errs.New()
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

	dErr.Add(err, map[string]interface{}{
		"profile": profile,
	})

	return chat, dErr
}

func (s *openAIService) prepareImageCompletion(ctx context.Context, version1, version2 image.Image, profile models.Profile) (*openai.ChatCompletion, errs.Error) {
	dErr := errs.New()
	// convert images to base64
	version1Base64, err := imageToBase64URL(version1)
	if err != nil {
		dErr.Add(svc.ErrConvertingImageToBase64, map[string]interface{}{
			"profile": profile,
		})
		return nil, dErr
	}

	version2Base64, err := imageToBase64URL(version2)
	if err != nil {
		dErr.Add(svc.ErrConvertingImageToBase64, map[string]interface{}{
			"profile": profile,
		})
		return nil, dErr
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

	dErr.Add(err, map[string]interface{}{
		"profile": profile,
	})
	return chat, dErr
}

func (s *openAIService) parseCompletion(chat *openai.ChatCompletion) (*models.DynamicChanges, errs.Error) {
	dErr := errs.New()
	if chat.Choices[0].Message.Refusal != "" {
		dErr.Add(svc.ErrEncounteredRefusal, map[string]interface{}{
			"refusal": chat.Choices[0].Message.Refusal,
		})
		return nil, dErr
	}

	changes := &models.DynamicChanges{
		Fields: make(map[string]interface{}),
	}

	err := json.Unmarshal([]byte(chat.Choices[0].Message.Content), changes)
	if err != nil {
		dErr.Add(svc.ErrParsingChanges, map[string]interface{}{
			"content": chat.Choices[0].Message.Content,
		})
		return nil, dErr
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
		return "", svc.ErrEncodingImage
	}

	// Convert the buffer bytes to base64 string
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Add prefix for base64 URL
	base64Str = "data:image/png;base64," + base64Str

	return base64Str, nil
}
