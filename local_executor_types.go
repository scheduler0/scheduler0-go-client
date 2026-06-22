package scheduler0_go_client

// LocalExecutorRegisterRequest is the body for POST /local-executors.
type LocalExecutorRegisterRequest struct {
	Name       string `json:"name"`
	Command    string `json:"command"`
	WorkingDir string `json:"workingDir,omitempty"`
	CreatedBy  string `json:"createdBy,omitempty"`
}

// LocalExecutorRegisterResponse is the server response for executor registration.
type LocalExecutorRegisterResponse struct {
	Success bool `json:"success"`
	Data    struct {
		ID int64 `json:"id"`
	} `json:"data"`
}

// LocalExecutorJobsResponse is the server response for pulling assigned jobs.
type LocalExecutorJobsResponse struct {
	Success bool  `json:"success"`
	Data    []Job `json:"data"`
}

// LocalExecutionReport is a single execution event sent to the server.
type LocalExecutionReport struct {
	JobID             int64  `json:"jobId"`
	UniqueID          string `json:"uniqueId"`
	State             int    `json:"state"` // 0=scheduled, 1=success, 2=failed
	LastExecutionTime string `json:"lastExecutionTime"` // RFC3339
	NextExecutionTime string `json:"nextExecutionTime"` // RFC3339
	ExecutionVersion  uint64 `json:"executionVersion"`
	JobQueueVersion   uint64 `json:"jobQueueVersion"`
}

// ReportLocalExecutionsResponse is the server response for reporting executions.
type ReportLocalExecutionsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Committed int `json:"committed"`
	} `json:"data"`
}
