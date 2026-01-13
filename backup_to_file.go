package scheduler0_go_client

// BackupDatabaseToFile initiates a backup to a specific file path
func (c *Client) BackupDatabaseToFile(destPath string) (*BackupRestoreResponse, error) {
	reqBody := BackupToFileRequest{
		DestPath: destPath,
	}

	req, err := c.newRequest("POST", "/cluster/backup-to-file", reqBody, "")
	if err != nil {
		return nil, err
	}

	var result BackupRestoreResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
