package scheduler0_go_client

// BackupRestoreResponse is the response from backup/restore initiation
type BackupRestoreResponse struct {
	Data    map[string]string `json:"data"`
	Success bool              `json:"success"`
}

// RestoreRequest is the request body for restore endpoint
type RestoreRequest struct {
	FilePath string `json:"filePath"` // Backup file name (S3 object key when S3 configured, or local path)
}
