package scheduler0_go_client

// ScheduleProjectInput creates-or-reuses a project by name for POST /ai/schedule.
// Ignored when ProjectID is set on the request.
type ScheduleProjectInput struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// SchedulePromptRequest is the request body for POST /ai/schedule. It carries the same
// natural-language hints as the prompt endpoint plus the target project and executor.
type SchedulePromptRequest struct {
	// AccountID overrides the X-Account-ID header for this request (0 uses the client default).
	AccountID  int64    `json:"-"`
	Prompt     string   `json:"prompt"`
	Purposes   []string `json:"purposes,omitempty"`
	Events     []string `json:"events,omitempty"`
	Recipients []string `json:"recipients,omitempty"`
	Channels   []string `json:"channels,omitempty"`
	Timezone   string   `json:"timezone,omitempty"`
	Locale     string   `json:"locale,omitempty"`
	// ProjectID reuses an existing project. Takes precedence over Project.
	ProjectID *int64 `json:"projectId,omitempty"`
	// Project creates-or-reuses a project by name when ProjectID is absent.
	Project *ScheduleProjectInput `json:"project,omitempty"`
	// ExecutorID pins a specific executor and skips LLM matching.
	ExecutorID *int64 `json:"executorId,omitempty"`
	// CreatedBy is required; it is stamped on the created project and jobs.
	CreatedBy string `json:"createdBy"`
}

// ScheduleResult is the response body for POST /ai/schedule.
type ScheduleResult struct {
	Classification      *IntentClassification `json:"classification,omitempty"`
	Project             Project               `json:"project"`
	ProjectCreated      bool                  `json:"projectCreated"`
	Executor            Executor              `json:"executor"`
	ExecutorMatchedBy   string                `json:"executorMatchedBy"`
	ExecutorMatchReason string                `json:"executorMatchReason,omitempty"`
	Jobs                []Job                 `json:"jobs"`
	Provider            string                `json:"provider,omitempty"`
	Model               string                `json:"model,omitempty"`
}
