package scheduler0_go_client

// Executor represents a job executor
type Executor struct {
	ID               int64   `json:"id"`
	AccountID        int64   `json:"accountId"`
	Name             string  `json:"name"`
	Type             string  `json:"type"`
	Region           string  `json:"region"`
	CloudProvider    string  `json:"cloudProvider"`
	CloudResourceURL string  `json:"cloudResourceUrl"`
	CloudAPIKey      string  `json:"cloudApiKey"`
	CloudAPISecret   string  `json:"cloudApiSecret"`
	WebhookURL       string  `json:"webhookUrl"`
	WebhookSecret    string  `json:"webhookSecret"`
	WebhookMethod    string  `json:"webhookMethod"`
	DateCreated      string  `json:"dateCreated"`
	DateModified     *string `json:"dateModified"`
	DateDeleted      *string `json:"dateDeleted"`
	CreatedBy        string  `json:"createdBy"`
	ModifiedBy       *string `json:"modifiedBy"`
	DeletedBy        *string `json:"deletedBy"`
}

// ExecutorResponse represents the response for a single executor
type ExecutorResponse struct {
	Success bool     `json:"success"`
	Data    Executor `json:"data"`
}

// ExecutorRequestBody represents the request body for creating an executor
type ExecutorRequestBody struct {
	AccountID        int64  `json:"-"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	Region           string `json:"region"`
	CloudProvider    string `json:"cloudProvider"`
	CloudResourceURL string `json:"cloudResourceUrl"`
	CloudAPIKey      string `json:"cloudApiKey,omitempty"`
	CloudAPISecret   string `json:"cloudApiSecret,omitempty"`
	WebhookURL       string `json:"webhookUrl,omitempty"`
	WebhookSecret    string `json:"webhookSecret,omitempty"`
	WebhookMethod    string `json:"webhookMethod,omitempty"`
	CreatedBy        string `json:"createdBy"`
}

// ExecutorUpdateRequestBody represents the request body for updating an executor
type ExecutorUpdateRequestBody struct {
	AccountID        int64  `json:"-"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	Region           string `json:"region"`
	CloudProvider    string `json:"cloudProvider"`
	CloudResourceURL string `json:"cloudResourceUrl"`
	CloudAPIKey      string `json:"cloudApiKey,omitempty"`
	CloudAPISecret   string `json:"cloudApiSecret,omitempty"`
	WebhookURL       string `json:"webhookUrl,omitempty"`
	WebhookSecret    string `json:"webhookSecret,omitempty"`
	WebhookMethod    string `json:"webhookMethod,omitempty"`
	ModifiedBy       string `json:"modifiedBy"`
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
