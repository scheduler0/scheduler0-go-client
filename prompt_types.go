package scheduler0_go_client

// PromptJobRequest represents the request body for creating jobs from AI prompt
type PromptJobRequest struct {
	Prompt     string   `json:"prompt"`
	Purposes   []string `json:"purposes,omitempty"`
	Events     []string `json:"events,omitempty"`
	Recipients []string `json:"recipients,omitempty"`
	Channels   []string `json:"channels,omitempty"`
	Timezone   string   `json:"timezone,omitempty"`
}

// PromptJobResponse represents a job configuration generated from AI prompt
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

