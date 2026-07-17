package scheduler0_go_client

import "fmt"

// GetAccountClassifyCount retrieves the remaining monthly AI classify-request quota for an account.
// accountID is used both in the URL path and as the X-Account-ID header for authentication.
func (c *Client) GetAccountClassifyCount(accountID string) (*AccountClassifyCountResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/accounts/%s/classify-count", accountID), nil, accountID)
	if err != nil {
		return nil, err
	}

	var result AccountClassifyCountResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// IncreaseAccountClassifyCount increases the remaining classify-request quota for an account.
// accountID is used both in the URL path and as the X-Account-ID header for authentication.
func (c *Client) IncreaseAccountClassifyCount(accountID string, count uint64) (*AccountClassifyCountIncreaseResponse, error) {
	body := map[string]uint64{
		"count": count,
	}
	req, err := c.newRequest("PUT", fmt.Sprintf("/accounts/%s/classify-count", accountID), body, accountID)
	if err != nil {
		return nil, err
	}

	var result AccountClassifyCountIncreaseResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
