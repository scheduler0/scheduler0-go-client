package scheduler0_go_client

// AsyncTask represents an async task
type AsyncTask struct {
	ID          int64  `json:"id"`
	RequestID   string `json:"requestId"`
	Input       string `json:"input"`
	Output      string `json:"output"`
	Service     string `json:"service"`
	State       int    `json:"state"`
	DateCreated string `json:"dateCreated"`
}

// AsyncTaskResponse represents the response for a single async task
type AsyncTaskResponse struct {
	Success bool      `json:"success"`
	Data    AsyncTask `json:"data"`
}

