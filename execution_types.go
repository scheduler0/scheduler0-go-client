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

// ListExecutionsParams represents parameters for listing executions
type ListExecutionsParams struct {
	StartDate      string // Start date for filtering (RFC3339 format, required)
	EndDate        string // End date for filtering (RFC3339 format, required)
	ProjectID      int64  // Project ID to filter by (0 for all)
	JobID          int64  // Job ID to filter by (0 for all)
	AccountID      int64  // Account ID override (0 to use client default)
	Limit          int    // Maximum number of items to return
	Offset         int    // Number of items to skip
	State          string // State filter: "scheduled", "completed", "failed", or "" for all
	OrderBy        string // Sort field: "dateCreated", "lastExecutionDateTime", "nextExecutionDateTime"
	OrderDirection string // Sort direction: "ASC" or "DESC"
}

// ExecutionMinuteBucket represents execution counts grouped by minute
type ExecutionMinuteBucket struct {
	Minute    string `json:"minute"`    // RFC3339 formatted time
	Total     uint64 `json:"total"`     // Total executions in this minute
	Scheduled uint64 `json:"scheduled"` // Scheduled executions in this minute
	Success   uint64 `json:"success"`   // Successful executions in this minute
	Failed    uint64 `json:"failed"`   // Failed executions in this minute
}

// ExecutionMinuteBucketsResponse represents the response for minute buckets
type ExecutionMinuteBucketsResponse struct {
	Success bool                  `json:"success"`
	Data    []ExecutionMinuteBucket `json:"data"`
}

// GetExecutionMinuteBucketsParams represents parameters for getting execution minute buckets
type GetExecutionMinuteBucketsParams struct {
	StartDate string // Start date for filtering (RFC3339 format, required)
	EndDate   string // End date for filtering (RFC3339 format, required)
	JobID     int64  // Job ID to filter by (0 for all jobs in account)
	AccountID int64  // Account ID override (0 to use client default)
}

