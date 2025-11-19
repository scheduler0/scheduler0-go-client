package scheduler0_go_client

import "fmt"

// ArchiveCredential archives a credential by ID
func (c *Client) ArchiveCredential(id string, archivedBy string) error {
	requestBody := map[string]string{
		"archivedBy": archivedBy,
	}

	req, err := c.newRequest("POST", fmt.Sprintf("/credentials/%s/archive", id), requestBody)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

