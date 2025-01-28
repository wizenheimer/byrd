package models

type CategoryChange struct {
	Category string   `json:"category" jsonschema_description:"The category of changes"`
	Summary  string   `json:"summary" jsonschema_description:"Brief summary of the changes"`
	Changes  []string `json:"changes" jsonschema_description:"List of detailed changes"`
}

type ChangeResponse struct {
	Changes []CategoryChange `json:"changes"`
}

type ChangeSummary struct {
	Category string `json:"category"`
	Summary  string `json:"summary"`
}
