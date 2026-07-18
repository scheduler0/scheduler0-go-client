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
	CloudAPIKey    string  `json:"cloudApiKey"`
	CloudAPISecret string  `json:"cloudApiSecret"`
	WebhookURL     string  `json:"webhookUrl"`
	WebhookSecret  string  `json:"webhookSecret"`
	WebhookMethod  string  `json:"webhookMethod"`
	Command        string  `json:"command,omitempty"`
	WorkingDir     string  `json:"workingDir,omitempty"`
	DateCreated    string  `json:"dateCreated"`
	DateModified   *string `json:"dateModified"`
	DateDeleted    *string `json:"dateDeleted"`
	CreatedBy      string  `json:"createdBy"`
	ModifiedBy     *string `json:"modifiedBy"`
	DeletedBy      *string `json:"deletedBy"`
}

// ExecutorResponse represents the response for a single executor
type ExecutorResponse struct {
	Success bool     `json:"success"`
	Data    Executor `json:"data"`
}

// ExecutorRequestBody represents the request body for creating an executor
type ExecutorRequestBody struct {
	AccountID        int64    `json:"-"`
	Name             string   `json:"name"`
	Description      string   `json:"description,omitempty"`
	Tags             []string `json:"tags,omitempty"`
	Type             string   `json:"type"`
	Region           string   `json:"region"`
	CloudProvider    string   `json:"cloudProvider"`
	CloudResourceURL string   `json:"cloudResourceUrl"`
	CloudAPIKey      string   `json:"cloudApiKey,omitempty"`
	CloudAPISecret   string   `json:"cloudApiSecret,omitempty"`
	WebhookURL       string   `json:"webhookUrl,omitempty"`
	WebhookSecret    string   `json:"webhookSecret,omitempty"`
	WebhookMethod    string   `json:"webhookMethod,omitempty"`
	Command          string   `json:"command,omitempty"`
	WorkingDir       string   `json:"workingDir,omitempty"`
	CreatedBy        string   `json:"createdBy"`
}

// ExecutorUpdateRequestBody represents the request body for updating an executor
type ExecutorUpdateRequestBody struct {
	AccountID        int64    `json:"-"`
	Name             string   `json:"name"`
	Description      string   `json:"description,omitempty"`
	Tags             []string `json:"tags,omitempty"`
	Type             string   `json:"type"`
	Region           string   `json:"region"`
	CloudProvider    string   `json:"cloudProvider"`
	CloudResourceURL string   `json:"cloudResourceUrl"`
	CloudAPIKey      string   `json:"cloudApiKey,omitempty"`
	CloudAPISecret   string   `json:"cloudApiSecret,omitempty"`
	WebhookURL       string   `json:"webhookUrl,omitempty"`
	WebhookSecret    string   `json:"webhookSecret,omitempty"`
	WebhookMethod    string   `json:"webhookMethod,omitempty"`
	Command          string   `json:"command,omitempty"`
	WorkingDir       string   `json:"workingDir,omitempty"`
	ModifiedBy       string   `json:"modifiedBy"`
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
