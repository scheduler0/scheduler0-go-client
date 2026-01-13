package scheduler0_go_client

// BackupDatabase initiates an automatic timestamped backup
func (c *Client) BackupDatabase() (*BackupRestoreResponse, error) {
	req, err := c.newRequest("POST", "/cluster/backup", nil, "")
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
