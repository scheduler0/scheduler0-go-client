package scheduler0_go_client

import "fmt"

// AddFeatureToAccount adds a feature to an account
func (c *Client) AddFeatureToAccount(accountID string, body *FeatureRequest) (*FeatureRequestResponse, error) {
	req, err := c.newRequest("PUT", fmt.Sprintf("/accounts/%s/feature", accountID), body, accountID)
	if err != nil {
		return nil, err
	}

	var result FeatureRequestResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RemoveFeatureFromAccount removes a feature from an account
func (c *Client) RemoveFeatureFromAccount(accountID string, body *FeatureRequest) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/accounts/%s/feature", accountID), body, accountID)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

// AddAllFeaturesToAccount adds all features to an account
func (c *Client) AddAllFeaturesToAccount(accountID string) error {
	req, err := c.newRequest("PUT", fmt.Sprintf("/accounts/%s/features/all", accountID), nil, accountID)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

// RemoveAllFeaturesFromAccount removes all features from an account
func (c *Client) RemoveAllFeaturesFromAccount(accountID string) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/accounts/%s/features/all", accountID), nil, accountID)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}
