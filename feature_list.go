package scheduler0_go_client

// ListFeatures retrieves all available features
func (c *Client) ListFeatures() (*FeaturesResponse, error) {
	req, err := c.newRequest("GET", "/features", nil)
	if err != nil {
		return nil, err
	}

	var result FeaturesResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

