package scheduler0_go_client

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

type ExecutionResponse struct {
	Success bool      `json:"success"`
	Data    Execution `json:"data"`
}

type PaginatedExecutionsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Total      int         `json:"total"`
		Offset     int         `json:"offset"`
		Limit      int         `json:"limit"`
		Executions []Execution `json:"executions"`
	} `json:"data"`
}

type ListExecutionsParams struct {
	StartDate      string
	EndDate        string
	ProjectID      int64
	JobID          int64
	AccountID      int64
	Limit          int
	Offset         int
	State          string
	OrderBy        string
	OrderDirection string
}

type DateRangeAnalyticsPoint struct {
	Date      string `json:"date"`
	Time      string `json:"time"`
	Scheduled uint64 `json:"scheduled"`
	Success   uint64 `json:"success"`
	Failed    uint64 `json:"failed"`
}

type DateRangeAnalyticsResponse struct {
	AccountID uint64                    `json:"accountId"`
	Timezone  string                    `json:"timezone"`
	StartDate string                    `json:"startDate"`
	StartTime string                    `json:"startTime"`
	EndDate   string                    `json:"endDate"`
	EndTime   string                    `json:"endTime"`
	Points    []DateRangeAnalyticsPoint `json:"points"`
}

type GetDateRangeAnalyticsParams struct {
	StartDate string `json:"startDate"`
	StartTime string `json:"startTime"`
	AccountID int64  `json:"accountId"`
}

type DateRangeAnalyticsAPIResponse struct {
	Success bool                       `json:"success"`
	Data    DateRangeAnalyticsResponse `json:"data"`
}

type ExecutionTotalsResponse struct {
	AccountID uint64 `json:"accountId"`
	Scheduled uint64 `json:"scheduled"`
	Success   uint64 `json:"success"`
	Failed    uint64 `json:"failed"`
}

type ExecutionTotalsAPIResponse struct {
	Success bool                    `json:"success"`
	Data    ExecutionTotalsResponse `json:"data"`
}

type CleanupOldLogsRequestBody struct {
	AccountID       string `json:"accountId"`
	RetentionMonths int    `json:"retentionMonths"`
}

type CleanupOldLogsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Message string `json:"message"`
	} `json:"data"`
}
