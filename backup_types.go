package scheduler0_go_client

import "time"

// BackupRestoreProgress represents the progress of a backup or restore operation
type BackupRestoreProgress struct {
	OperationType string    `json:"operationType"` // "backup" | "restore" | "backup-to-file"
	Status        string    `json:"status"`        // "idle" | "in-progress" | "completed" | "failed"
	Progress      int       `json:"progress"`      // 0-100
	Message       string    `json:"message"`       // Progress message or error
	StartTime     time.Time `json:"startTime"`
	EndTime       time.Time `json:"endTime"`
	BackupPath    string    `json:"backupPath"` // Result path for backup operations
}

// BackupRestoreResponse is the response from backup/restore initiation
type BackupRestoreResponse struct {
	Data    map[string]string `json:"data"`
	Success bool              `json:"success"`
}

// BackupRestoreProgressResponse wraps the progress data
type BackupRestoreProgressResponse struct {
	Data    BackupRestoreProgress `json:"data"`
	Success bool                  `json:"success"`
}

// BackupToFileRequest is the request body for backup-to-file endpoint
type BackupToFileRequest struct {
	DestPath string `json:"destPath"`
}

// RestoreRequest is the request body for restore endpoint
type RestoreRequest struct {
	BackupPath string `json:"backupPath"`
}
