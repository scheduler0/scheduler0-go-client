package scheduler0_go_client

// AccountAISettings holds per-account AI provider credentials.
type AccountAISettings struct {
	AccountID          uint64 `json:"account_id,omitempty"`
	Provider           string `json:"provider,omitempty"`
	Model              string `json:"model,omitempty"`
	OpenAIAPIKey       string `json:"openai_api_key,omitempty"`
	AnthropicAPIKey    string `json:"anthropic_api_key,omitempty"`
	BedrockAccessKeyID string `json:"bedrock_access_key_id,omitempty"`
	BedrockSecretKey   string `json:"bedrock_secret_key,omitempty"`
	BedrockRegion      string `json:"bedrock_region,omitempty"`
}

// AccountAISettingsResponse wraps the standard API envelope for a single settings object.
type AccountAISettingsResponse struct {
	Success bool              `json:"success"`
	Data    AccountAISettings `json:"data"`
}
