package scheduler0_go_client

// GetBackupRestoreProgress retrieves the current backup/restore progress
func (c *Client) GetBackupRestoreProgress() (*BackupRestoreProgressResponse, error) {
	req, err := c.newRequest("GET", "/cluster/backup-restore-progress", nil, "")
	if err != nil {
		return nil, err
	}

	var result BackupRestoreProgressResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
