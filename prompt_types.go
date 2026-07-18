package scheduler0_go_client

// PromptJobRequest represents the request body for creating jobs from AI prompt
type PromptJobRequest struct {
	Prompt     string   `json:"prompt"`
	Purposes   []string `json:"purposes,omitempty"`
	Events     []string `json:"events,omitempty"`
	Recipients []string `json:"recipients,omitempty"`
	Channels   []string `json:"channels,omitempty"`
	Timezone   string   `json:"timezone,omitempty"`

	// Locale is an optional BCP 47 language tag (e.g. "en", "en-US"). The AI
	// prompt flow is English-only; a non-en* locale bypasses the intent guardrail.
	Locale string `json:"locale,omitempty"`
}

// ClassifyPromptRequest is the request body for POST /prompt/classify.
type ClassifyPromptRequest struct {
	Prompt string `json:"prompt"`
	// Locale is an optional BCP-47 locale. Only English (en*) is currently
	// supported; other locales are rejected by the server with a 400.
	Locale string `json:"locale,omitempty"`
}

// PromptJobResponse represents a single job configuration generated from AI prompt.
// It maps to the objects inside each provider result's "jobs" array.
type PromptJobResponse struct {
	Kind           string                 `json:"kind,omitempty"`
	Purpose        string                 `json:"purpose,omitempty"`
	Subject        string                 `json:"subject,omitempty"`
	NextRunAt      *string                `json:"nextRunAt,omitempty"`
	Recurrence     string                 `json:"recurrence,omitempty"`
	Event          string                 `json:"event,omitempty"`
	Delivery       string                 `json:"delivery,omitempty"`
	CronExpression string                 `json:"cronExpression,omitempty"`
	Channel        string                 `json:"channel,omitempty"`
	Recipients     []string               `json:"recipients,omitempty"`
	StartDate      *string                `json:"startDate,omitempty"`
	EndDate        *string                `json:"endDate,omitempty"`
	Timezone       string                 `json:"timezone,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// PromptProviderResult wraps the per-provider metadata together with the jobs
// it produced. The API returns one entry per AI provider that was consulted.
type PromptProviderResult struct {
	Provider     string              `json:"provider"`
	Model        string              `json:"model"`
	Jobs         []PromptJobResponse `json:"jobs"`
	InputTokens  int                 `json:"inputTokens"`
	OutputTokens int                 `json:"outputTokens"`
	TotalTokens  int                 `json:"totalTokens"`
	DurationMs   int64               `json:"durationMs"`
}

// IntentClassification is the response from the edge intent classifier.
type IntentClassification struct {
	Text     string `json:"text"`
	Decision string `json:"decision"`
	Reason   string `json:"reason"`
}

// PromptResult is the top-level response body for POST /prompt.
// Classification is nil when the intent guardrail is disabled or errored (fail-open).
type PromptResult struct {
	Providers      []PromptProviderResult `json:"providers"`
	Classification *IntentClassification  `json:"classification,omitempty"`
}

// PromptSkippedError is returned when the intent guardrail rejected or clarified the
// prompt (HTTP 422). Classification carries the full classifier output.
type PromptSkippedError struct {
	Message        string
	Classification *IntentClassification
}

func (e *PromptSkippedError) Error() string {
	return e.Message
}
