package scheduler0_go_client

import "fmt"

// ListExecutors retrieves all executors with optional query parameters
func (c *Client) ListExecutors(accountID int64, limit, offset int, orderBy, orderByDirection string) (*PaginatedExecutorsResponse, error) {
	queryParams := map[string]string{
		"limit":  fmt.Sprintf("%d", limit),
		"offset": fmt.Sprintf("%d", offset),
	}

	if orderBy != "" {
		queryParams["orderBy"] = orderBy
	}
	if orderByDirection != "" {
		queryParams["orderByDirection"] = orderByDirection
	}

	var accountIDOverride string
	if accountID > 0 {
		accountIDOverride = fmt.Sprintf("%d", accountID)
	}

	req, err := c.newRequestWithQuery("GET", "/executors", nil, queryParams, accountIDOverride)
	if err != nil {
		return nil, err
	}

	var result PaginatedExecutorsResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
