package scheduler0_go_client

import "fmt"

// ListCredentials retrieves all credentials with optional query parameters
func (c *Client) ListCredentials(params ListCredentialsParams) (*PaginatedCredentialsResponse, error) {
	queryParams := map[string]string{
		"limit":  fmt.Sprintf("%d", params.Limit),
		"offset": fmt.Sprintf("%d", params.Offset),
	}

	if params.OrderBy != "" {
		queryParams["orderBy"] = params.OrderBy
	}
	if params.OrderByDirection != "" {
		queryParams["orderByDirection"] = params.OrderByDirection
	}

	var accountIDOverride string
	if params.AccountID > 0 {
		accountIDOverride = fmt.Sprintf("%d", params.AccountID)
	}

	req, err := c.newRequestWithQuery("GET", "/credentials", nil, queryParams, accountIDOverride)
	if err != nil {
		return nil, err
	}

	var result PaginatedCredentialsResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
