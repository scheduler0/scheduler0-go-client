package scheduler0_go_client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s", string(body))
	}

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}
	return nil
}

