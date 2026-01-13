package scheduler0_go_client

// RestoreDatabase initiates a restore from a backup file
func (c *Client) RestoreDatabase(backupPath string) (*BackupRestoreResponse, error) {
	reqBody := RestoreRequest{
		BackupPath: backupPath,
	}

	req, err := c.newRequest("POST", "/cluster/restore", reqBody, "")
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
