package scheduler0_go_client

// Job represents a scheduled job
type Job struct {
	ID                int64   `json:"id,omitempty"`
	AccountID         int64   `json:"accountId,omitempty"`
	ProjectID         int64   `json:"projectId,omitempty"`
	ExecutorID        *int64  `json:"executorId,omitempty"`
	Data              string  `json:"data,omitempty"`
	Spec              string  `json:"spec,omitempty"`
	StartDate         string  `json:"startDate,omitempty"`
	EndDate           string  `json:"endDate,omitempty"`
	LastExecutionDate string  `json:"lastExecutionDate,omitempty"`
	Timezone          string  `json:"timezone,omitempty"`
	TimezoneOffset    int64   `json:"timezoneOffset,omitempty"`
	RetryMax          int     `json:"retryMax,omitempty"`
	ExecutionID       string  `json:"executionId,omitempty"`
	Status            string  `json:"status,omitempty"`
	DateCreated       string  `json:"dateCreated,omitempty"`
	DateModified      *string `json:"dateModified,omitempty"`
	CreatedBy         string  `json:"createdBy,omitempty"`
	ModifiedBy        *string `json:"modifiedBy,omitempty"`
	DeletedBy         *string `json:"deletedBy,omitempty"`
}

// JobResponse represents the response for a single job
type JobResponse struct {
	Success bool `json:"success"`
	Data    Job  `json:"data"`
}

// BatchJobResponse represents the response for batch job creation
type BatchJobResponse struct {
	Success bool   `json:"success"`
	Data    string `json:"data"`
}

// JobRequestBody represents the request body for creating a job
type JobRequestBody struct {
	AccountID      int64  `json:"-"`
	ProjectID      int64  `json:"projectId"`
	Timezone       string `json:"timezone"`
	ExecutorID     *int64 `json:"executorId,omitempty"`
	Data           string `json:"data,omitempty"`
	Spec           string `json:"spec,omitempty"`
	StartDate      string `json:"startDate,omitempty"`
	EndDate        string `json:"endDate,omitempty"`
	TimezoneOffset int64  `json:"timezoneOffset,omitempty"`
	RetryMax       int    `json:"retryMax,omitempty"`
	Status         string `json:"status,omitempty"`
	CreatedBy      string `json:"createdBy"`
}

// JobUpdateRequestBody represents the request body for updating a job
type JobUpdateRequestBody struct {
	AccountID      int64  `json:"-"`
	ProjectID      int64  `json:"projectId,omitempty"`
	ExecutorID     *int64 `json:"executorId,omitempty"`
	Data           string `json:"data,omitempty"`
	Spec           string `json:"spec,omitempty"`
	StartDate      string `json:"startDate,omitempty"`
	EndDate        string `json:"endDate,omitempty"`
	Timezone       string `json:"timezone,omitempty"`
	TimezoneOffset int64  `json:"timezoneOffset,omitempty"`
	RetryMax       int    `json:"retryMax,omitempty"`
	Status         string `json:"status,omitempty"`
	ModifiedBy     string `json:"modifiedBy"`
}

// JobDeleteRequestBody represents the request body for deleting a job
type JobDeleteRequestBody struct {
	AccountID int64  `json:"-"`
	DeletedBy string `json:"deletedBy"`
}

// PaginatedJobsResponse represents a paginated list of jobs
type PaginatedJobsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Total  int   `json:"total"`
		Offset int   `json:"offset"`
		Limit  int   `json:"limit"`
		Jobs   []Job `json:"jobs"`
	} `json:"data"`
}

// ListJobsParams represents parameters for listing jobs
type ListJobsParams struct {
	ProjectID        string // Project ID to filter by (empty string for all projects)
	AccountID        string // Account ID override (empty string to use client default)
	Limit            int    // Maximum number of items to return
	Offset           int    // Number of items to skip
	OrderBy          string // Field to order by (e.g., "date_created", "date_modified")
	OrderByDirection string // Direction to order ("asc" or "desc")
}
