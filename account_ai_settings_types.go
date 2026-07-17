package scheduler0_go_client

// ActiveModel is a single entry in the ordered active-model list for an account.
// Provider is one of "openai", "anthropic", "bedrock", "openrouter".
type ActiveModel struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
}

// ModelInfo describes a single approved model returned by GET /ai/models.
type ModelInfo struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Default     bool   `json:"default,omitempty"`
}

// AccountAISettings holds per-account AI provider credentials.
type AccountAISettings struct {
	AccountID          uint64        `json:"account_id,omitempty"`
	Provider           string        `json:"provider,omitempty"`
	Model              string        `json:"model,omitempty"`
	ActiveModels       []ActiveModel `json:"active_models,omitempty"`
	OpenAIAPIKey       string        `json:"openai_api_key,omitempty"`
	AnthropicAPIKey    string        `json:"anthropic_api_key,omitempty"`
	BedrockAccessKeyID string        `json:"bedrock_access_key_id,omitempty"`
	BedrockSecretKey   string        `json:"bedrock_secret_key,omitempty"`
	BedrockRegion      string        `json:"bedrock_region,omitempty"`
	OpenRouterAPIKey   string        `json:"openrouter_api_key,omitempty"`
}

// AccountAISettingsResponse wraps the standard API envelope for a single settings object.
type AccountAISettingsResponse struct {
	Success bool              `json:"success"`
	Data    AccountAISettings `json:"data"`
}

// AIModelsResponse wraps the standard API envelope for the GET /ai/models catalog.
// Data maps provider names to their approved ModelInfo lists.
type AIModelsResponse struct {
	Success bool                     `json:"success"`
	Data    map[string][]ModelInfo   `json:"data"`
}
