package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"math/rand"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
)

type CategoryChange struct {
	Category string   `json:"category" jsonschema_description:"The category of changes"`
	Summary  string   `json:"summary" jsonschema_description:"Brief summary of the changes"`
	Changes  []string `json:"changes" jsonschema_description:"List of detailed changes"`
}

type ChangeResponse struct {
	Changes []CategoryChange `json:"changes"`
}

type result struct {
	summary ChangeSummary
	err     error
}

type ChangeSummary struct {
	Category string `json:"category"`
	Summary  string `json:"summary"`
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
	// Sassy nothing
	"Your competition's %s game? Still snoozing",
	"Breaking: Absolutely nothing happening in %s",
	"Crickets chirping in the %s space",

	// Dramatic silence
	"Plot twist: There is no plot twist in %s",
	"The %s update we've all been waiting for... isn't here",
	"Your competition's %s moves? Still loading...",
	"Some suspense in %s (just kidding, nothing changed)",

	// Business zen
	"The %s space is channeling peace and quiet today",
	"Your competition chose meditation for their %s strategy",
	"Someone's taking a power nap in %s wing, apparently",

	// Market commentary
	"Competitive intelligence report for %s: *tumbleweeds*",
	"In today's riveting %s news: nothing",
	"Your %s radar is clear (almost too clear 🤔)",
	"The %s space is giving strong 'do not disturb' energy",
	"Groundbreaking update: Nothing changed in %s",
	"Your rival is really committed to keeping %s exactly the same",

	// Playful shade
	"The %s space is taking a power nap",
	"Someone hit pause on %s changes",
	"Plot twist: There is no plot twist in %s",
	"Your competition chose peace in %s today",

	// Dramatic nothing
	"Breaking news: Still nothing new in %s",
	"The suspense continues in %s (jk, nothing changed)",
	"Things are getting wild in %s (they're not)",
	"Alert: Your competition is living their best %s life (unchanged)",

	// Chef's kiss
	"Your %s radar is picking up pure silence",
	"The most consistent thing about %s? No updates",
	"Some things never change (like %s)",
	"Your competition's really nailing the whole 'stable %s' vibe",
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

var ChangeSummaryResponseSchema = GenerateSchema[ChangeSummary]()

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

func generateCategorySummary(ctx context.Context, client *openai.Client, category string, changes []string) (ChangeSummary, error) {
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
		return ChangeSummary{}, fmt.Errorf("OpenAI API error: %v", err)
	}

	var summary ChangeSummary
	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &summary)
	if err != nil {
		return ChangeSummary{}, fmt.Errorf("JSON parsing error: %v", err)
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
