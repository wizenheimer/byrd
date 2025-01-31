package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"math/rand"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

type result struct {
	summary models.ChangeSummary
	err     error
}

var fallbackTemplates = []string{
	// Bold statements
	"👀 Look who's making moves in %s",
	"Well well well... what do we have in %s",
	"Your competition didn't think we'd notice their %s changes",
	"Spotted: Your rival trying to be sneaky with %s",

	// Direct callouts
	"Drop what you're doing - %s just got interesting",
	"This %s update? Yeah, you'll want to see this",
	"Your competitor thinks they're slick with these %s changes",
	"Hot take: Someone's shaking up %s",

	// Business sass
	"The tea is hot: Fresh drama in %s",
	"Plot twist in %s (and it's a good one)",
	"Your competition's %s move is... certainly a choice",
	"Someone's feeling brave in the %s space",

	// Market commentary
	"The competition chose chaos in %s today",
	"Things are getting spicy in %s",
	"You'll never guess what changed in %s (or maybe you will)",

	// Competitive edge
	"Heads up: The game just changed in %s",
	"Your rival's latest %s move is actually pretty bold",
	"This %s update? Worth your coffee break",
	"Someone's trying to steal your %s thunder",

	// Industry chatter
	"The industry's buzzing about these %s changes",
	"We need to talk about what's happening in %s",
	"PSA: Your competition is up to something in %s",
	"Breaking: Things just got interesting in %s",

	// Peak snark
	"Oh, you'll love this %s update (or not)",
	"Looks like someone's been busy with %s",
	"The plot thickens in %s land",
	"Your competition picked quite a day for %s changes",
}

var noChangeTemplates = []string{
	// All original templates properly articled
	"All quiet on the %s front",
	"No noteworthy changes on the %s front",
	"No changes detected in %s",
	"No new activity in %s",
	"The %s remains unchanged",

	"No activity on the %s front",
	"The situation remains unchanged in %s",
	"No developments in %s today",
	"The %s status: no changes",
	"No notable updates in %s",
	"The %s remains static",
	"No movement detected in %s",
	"The %s remains constant",
	"No shifts observed in %s",
	"All steady on the %s side",

	"Nothing to report on %s",
	"The %s shows no changes",
	"No changes found in %s",
	"All clear on the %s front",
	"No movement in the %s space",
	"The status quo maintains in %s",
	"The %s remains static",
	"No updates detected in %s",
	"No new developments for %s",
	"No developments on the %s end",
	"The %s continues unchanged",
	"No activity detected in %s",
	"All quiet in the %s sphere",
	"No shifts spotted in %s",
	"Nothing notable in %s",
	"The %s stays consistent",
	"No changes observed in %s",
	"The situation remains stable in %s",
	"No modifications to %s",
	"The %s holds steady",
	"No variation found in %s",
	"The status remains unchanged for %s",
	"No alterations in %s",
	"The %s remains the same",
	"Nothing different in %s",
}

func GenerateSchema[T any]() interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

var ChangeSummaryResponseSchema = GenerateSchema[models.ChangeSummary]()

func getNoChangeSummary(category string) string {
	// Create a new random seed for each function call
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Get a truly random index for no changes template
	templateIndex := r.Intn(len(noChangeTemplates))
	template := noChangeTemplates[templateIndex]
	return fmt.Sprintf(template, category)
}

func getFallbackSummary(category string, changes []string) string {
	// Create a new random seed for each function call
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	if len(changes) == 0 {
		// Get a truly random index for no changes template
		templateIndex := r.Intn(len(noChangeTemplates))
		template := noChangeTemplates[templateIndex]
		return fmt.Sprintf(template, category)
	}

	// Get a truly random index for fallback template
	templateIndex := r.Intn(len(fallbackTemplates))
	template := fallbackTemplates[templateIndex]
	return fmt.Sprintf(template, category)
}

func generateCategorySummary(ctx context.Context, client *openai.Client, category string, changes []string) (models.ChangeSummary, error) {
	if len(changes) == 0 {
		return models.ChangeSummary{
			Category: category,
			Summary:  getNoChangeSummary(category),
		}, nil
	}

	prompt := fmt.Sprintf("Give a brief 1-2 line summary of these changes for %s category:\n\n%s",
		category,
		strings.Join(changes, "\n"),
	)

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("changes_summary"),
		Description: openai.F("Brief summary of changes"),
		Schema:      openai.F(ChangeSummaryResponseSchema),
		Strict:      openai.Bool(true),
	}

	chat, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		}),
		ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
			openai.ResponseFormatJSONSchemaParam{
				Type:       openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
				JSONSchema: openai.F(schemaParam),
			},
		),
		Model: openai.F(openai.ChatModelGPT4o2024_08_06),
	})

	if err != nil {
		return models.ChangeSummary{}, fmt.Errorf("OpenAI API error: %v", err)
	}

	if chat == nil {
		return models.ChangeSummary{}, errors.New("received nil response from OpenAI")
	}

	var summary models.ChangeSummary
	choices := chat.Choices
	if len(choices) == 0 || choices == nil {
		return models.ChangeSummary{}, errors.New("no choices returned from ai service")
	}

	err = json.Unmarshal([]byte(choices[0].Message.Content), &summary)
	if err != nil {
		return models.ChangeSummary{}, fmt.Errorf("JSON parsing error: %v", err)
	}

	return summary, nil
}

func processCategoryAsync(ctx context.Context, client *openai.Client, category string, changes []string, resultChan chan<- result) {
	summary, err := generateCategorySummary(ctx, client, category, changes)
	select {
	case <-ctx.Done():
		return
	case resultChan <- result{summary, err}:
	}
}
