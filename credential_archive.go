package scheduler0_go_client

import "fmt"

// ArchiveCredential archives a credential by ID
// accountIDOverride is optional - if provided, overrides the client's default account ID
func (c *Client) ArchiveCredential(id string, archivedBy string, accountIDOverride ...string) error {
	requestBody := map[string]string{
		"archivedBy": archivedBy,
	}

	var accountID string
	if len(accountIDOverride) > 0 {
		accountID = accountIDOverride[0]
	}

	req, err := c.newRequest("POST", fmt.Sprintf("/credentials/%s/archive", id), requestBody, accountID)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}
