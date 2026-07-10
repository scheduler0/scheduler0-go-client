package scheduler0_go_client

import "fmt"

// PromptJobRequest represents the request body for creating jobs from AI prompt
type PromptJobRequest struct {
	Prompt     string   `json:"prompt"`
	Purposes   []string `json:"purposes,omitempty"`
	Events     []string `json:"events,omitempty"`
	Recipients []string `json:"recipients,omitempty"`
	Channels   []string `json:"channels,omitempty"`
	Timezone   string   `json:"timezone,omitempty"`
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

// PromptClassification contains the intent-guardrail decision for a prompt.
type PromptClassification struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason"`
}

// PromptResponse is the structured result of CreateJobFromPrompt. It exposes
// the per-provider results and an optional intent-guardrail classification.
type PromptResponse struct {
	Providers      []PromptProviderResult `json:"providers"`
	Classification *PromptClassification  `json:"classification,omitempty"`
}

// PromptSkippedError is returned by CreateJobFromPrompt when the intent
// guardrail rejects the prompt (HTTP 422). It carries the guardrail's decision
// and reason so the caller can surface them without parsing raw HTTP bodies.
type PromptSkippedError struct {
	Message        string
	Classification *PromptClassification
}

func (e *PromptSkippedError) Error() string {
	if e.Classification != nil {
		return fmt.Sprintf("prompt skipped: %s — %s", e.Classification.Decision, e.Classification.Reason)
	}
	return fmt.Sprintf("prompt skipped: %s", e.Message)
}

