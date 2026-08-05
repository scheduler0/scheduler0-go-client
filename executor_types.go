package scheduler0_go_client

// Executor represents a job executor
type Executor struct {
	ID        int64  `json:"id"`
	AccountID int64  `json:"accountId"`
	Name      string `json:"name"`
	// Description and Tags describe what the executor does. They are used by the
	// /ai/schedule endpoint to match an executor to a prompt's purpose and channels.
	Description      string   `json:"description,omitempty"`
	Tags             []string `json:"tags,omitempty"`
	Type             string   `json:"type"`
	Region           string   `json:"region"`
	CloudProvider    string   `json:"cloudProvider"`
	CloudResourceURL string   `json:"cloudResourceUrl"`
	// CloudAPIKey, CloudAPISecret and WebhookSecret are secrets. The server stores them
	// encrypted and only returns them once, in the CreateExecutor response. They are
	// always empty on GetExecutor, ListExecutors and UpdateExecutor responses.
	CloudAPIKey        string  `json:"cloudApiKey"`
	CloudAPISecret     string  `json:"cloudApiSecret"`
	WebhookURL         string  `json:"webhookUrl"`
	WebhookSecret      string  `json:"webhookSecret"`
	WebhookMethod      string  `json:"webhookMethod"`
	Command            string  `json:"command,omitempty"`
	WorkingDir         string  `json:"workingDir,omitempty"`
	PayloadAggregation bool    `json:"payloadAggregation,omitempty"`
	DateCreated        string  `json:"dateCreated"`
	DateModified       *string `json:"dateModified"`
	DateDeleted        *string `json:"dateDeleted"`
	CreatedBy          string  `json:"createdBy"`
	ModifiedBy         *string `json:"modifiedBy"`
	DeletedBy          *string `json:"deletedBy"`
}

// ExecutorResponse represents the response for a single executor
type ExecutorResponse struct {
	Success bool     `json:"success"`
	Data    Executor `json:"data"`
}

// ExecutorRequestBody represents the request body for creating an executor
type ExecutorRequestBody struct {
	AccountID          int64    `json:"-"`
	Name               string   `json:"name"`
	Description        string   `json:"description,omitempty"`
	Tags               []string `json:"tags,omitempty"`
	Type               string   `json:"type"`
	Region             string   `json:"region"`
	CloudProvider      string   `json:"cloudProvider"`
	CloudResourceURL   string   `json:"cloudResourceUrl"`
	CloudAPIKey        string   `json:"cloudApiKey,omitempty"`
	CloudAPISecret     string   `json:"cloudApiSecret,omitempty"`
	WebhookURL         string   `json:"webhookUrl,omitempty"`
	WebhookSecret      string   `json:"webhookSecret,omitempty"`
	WebhookMethod      string   `json:"webhookMethod,omitempty"`
	Command            string   `json:"command,omitempty"`
	WorkingDir         string   `json:"workingDir,omitempty"`
	PayloadAggregation bool     `json:"payloadAggregation,omitempty"`
	CreatedBy          string   `json:"createdBy"`
}

// ExecutorUpdateRequestBody represents the request body for updating an executor
type ExecutorUpdateRequestBody struct {
	AccountID          int64    `json:"-"`
	Name               string   `json:"name"`
	Description        string   `json:"description,omitempty"`
	Tags               []string `json:"tags,omitempty"`
	Type               string   `json:"type"`
	Region             string   `json:"region"`
	CloudProvider      string   `json:"cloudProvider"`
	CloudResourceURL   string   `json:"cloudResourceUrl"`
	CloudAPIKey        string   `json:"cloudApiKey,omitempty"`
	CloudAPISecret     string   `json:"cloudApiSecret,omitempty"`
	WebhookURL         string   `json:"webhookUrl,omitempty"`
	WebhookSecret      string   `json:"webhookSecret,omitempty"`
	WebhookMethod      string   `json:"webhookMethod,omitempty"`
	Command            string   `json:"command,omitempty"`
	WorkingDir         string   `json:"workingDir,omitempty"`
	PayloadAggregation bool     `json:"payloadAggregation,omitempty"`
	ModifiedBy         string   `json:"modifiedBy"`
}

// ExecutorDeleteRequestBody represents the request body for deleting an executor
type ExecutorDeleteRequestBody struct {
	AccountID int64  `json:"-"`
	DeletedBy string `json:"deletedBy"`
}

// PaginatedExecutorsResponse represents a paginated list of executors
type PaginatedExecutorsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Total     int        `json:"total"`
		Offset    int        `json:"offset"`
		Limit     int        `json:"limit"`
		Executors []Executor `json:"executors"`
	} `json:"data"`
}

// ListExecutorsParams represents parameters for listing executors
type ListExecutorsParams struct {
	AccountID        int64  // Account ID override (0 to use client default)
	Limit            int    // Maximum number of items to return
	Offset           int    // Number of items to skip
	OrderBy          string // Field to order by (e.g., "date_created", "date_modified")
	OrderByDirection string // Direction to order ("asc" or "desc")
}

// TestInvocationRequestBody is the (optional) body for TestInvokeExecutor. All
// fields are optional; an empty body invokes the executor with a default
// synthetic job.
type TestInvocationRequestBody struct {
	AccountID int64 `json:"-"`
	// Job carries the standard job attributes to include in the invocation
	// payload. Server-managed fields (id, accountId, executorId, dateCreated,
	// lastExecutionDate) are ignored and overridden by the server.
	Job *Job `json:"job,omitempty"`
	// Age is how old the synthetic job entry should appear, as a Go duration
	// string (e.g. "24h", "1h30m", "15m"). Sets the job's dateCreated and
	// lastExecutionDate to now-age.
	Age string `json:"age,omitempty"`
	// ExecutionTime is the moment the job is treated as executing at (RFC3339).
	// Becomes the payload's lastExecutionDateTime. Defaults to now.
	ExecutionTime string `json:"executionTime,omitempty"`
}

// JobInvocationPayload is the payload delivered to an executor when a job fires.
type JobInvocationPayload struct {
	Job                   Job    `json:"job"`
	LastExecutionDateTime string `json:"lastExecutionDateTime,omitempty"`
	LastExecutionStatus   string `json:"lastExecutionStatus,omitempty"`
}

// AggregatedJobInvocationPayload is delivered when payload aggregation is enabled
// and multiple jobs sharing an executor fire at the same scheduled time.
type AggregatedJobInvocationPayload struct {
	Aggregated bool                   `json:"aggregated"`
	Jobs       []JobInvocationPayload `json:"jobs"`
}

// TestInvocationResult reports the outcome of a test invocation.
type TestInvocationResult struct {
	Test         bool                 `json:"test"`
	ExecutorID   int64                `json:"executorId"`
	ExecutorType string               `json:"executorType"`
	Success      bool                 `json:"success"`
	Error        string               `json:"error,omitempty"`
	StartedAt    string               `json:"startedAt"`
	FinishedAt   string               `json:"finishedAt"`
	DurationMs   int64                `json:"durationMs"`
	Payload      JobInvocationPayload `json:"payload"`
}

// TestInvocationResponse represents the response for a test invocation.
type TestInvocationResponse struct {
	Success bool                 `json:"success"`
	Data    TestInvocationResult `json:"data"`
}
