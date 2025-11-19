package scheduler0_go_client

// Execution represents a job execution log
type Execution struct {
	ID                    int64   `json:"id"`
	AccountID             int64   `json:"accountId"`
	UniqueID              string  `json:"uniqueId"`
	State                 int64   `json:"state"`
	NodeID                int64   `json:"nodeId"`
	JobID                 int64   `json:"jobId"`
	LastExecutionDatetime string  `json:"lastExecutionDatetime"`
	NextExecutionDatetime string  `json:"nextExecutionDatetime"`
	JobQueueVersion       int64   `json:"jobQueueVersion"`
	ExecutionVersion      int64   `json:"executionVersion"`
	DateCreated           string  `json:"dateCreated"`
	DateModified          *string `json:"dateModified"`
}

// ExecutionResponse represents the response for a single execution
type ExecutionResponse struct {
	Success bool      `json:"success"`
	Data    Execution `json:"data"`
}

// PaginatedExecutionsResponse represents a paginated list of executions
type PaginatedExecutionsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Total      int         `json:"total"`
		Offset     int         `json:"offset"`
		Limit      int         `json:"limit"`
		Executions []Execution `json:"executions"`
	} `json:"data"`
}

