package scheduler0_go_client

// Async task states returned in AsyncTask.State.
const (
	AsyncTaskNotStarted = 0
	AsyncTaskInProgress = 1
	AsyncTaskSuccess    = 2
	AsyncTaskFailed     = 3
)

// AsyncTask represents an async task
type AsyncTask struct {
	ID           int64  `json:"id"`
	RequestID    string `json:"requestId"`
	Input        string `json:"input"`
	Output       string `json:"output"`
	Service      string `json:"service"`
	State        int    `json:"state"` // see AsyncTask* constants
	DateCreated  string `json:"dateCreated"`
	AccountID    int64  `json:"accountId"`
	DateModified string `json:"dateModified"`
}

// AsyncTaskResponse represents the response for a single async task
type AsyncTaskResponse struct {
	Success bool      `json:"success"`
	Data    AsyncTask `json:"data"`
}
