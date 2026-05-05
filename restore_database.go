package scheduler0_go_client

// RestoreDatabase initiates a restore from a backup file.
// fileName is the backup file name (S3 object key when S3 is configured, or local file path otherwise).
func (c *Client) RestoreDatabase(fileName string) (*BackupRestoreResponse, error) {
	reqBody := RestoreRequest{
		FilePath: fileName,
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
