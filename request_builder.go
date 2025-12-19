package scheduler0_go_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"reflect"
)

func (c *Client) newRequest(method, endpoint string, body interface{}, accountIDOverride ...string) (*http.Request, error) {
	versionPrefix := fmt.Sprintf("/api/%s/", c.Version)

	rel := &url.URL{Path: path.Join(fmt.Sprintf("%s%s", c.BaseURL.Path, versionPrefix), endpoint)}
	u := c.BaseURL.ResolveReference(rel)

	var buf bytes.Buffer
	if body != nil {
		err := json.NewEncoder(&buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Set authentication based on client type
	if c.Username != "" && c.Password != "" {
		// Basic Auth for peer communication
		req.SetBasicAuth(c.Username, c.Password)
		req.Header.Set("X-Peer", "cmd")
	} else if c.APIKey != "" && c.APISecret != "" {
		// API Key + Secret authentication
		req.Header.Set("X-API-Key", c.APIKey)
		req.Header.Set("X-API-Secret", c.APISecret)
	}

	// Add account ID based on override/body/client default preferences
	accountID := c.resolveAccountID(body, accountIDOverride)
	if accountID != "" {
		req.Header.Set("X-Account-ID", accountID)
	}

	return req, nil
}

func (c *Client) newRequestWithQuery(method, endpoint string, body interface{}, queryParams map[string]string, accountIDOverride ...string) (*http.Request, error) {
	versionPrefix := fmt.Sprintf("/api/%s/", c.Version)

	rel := &url.URL{Path: path.Join(fmt.Sprintf("%s%s", c.BaseURL.Path, versionPrefix), endpoint)}
	u := c.BaseURL.ResolveReference(rel)

	// Add query parameters
	if len(queryParams) > 0 {
		q := u.Query()
		for key, value := range queryParams {
			q.Set(key, value)
		}
		u.RawQuery = q.Encode()
	}

	var buf bytes.Buffer
	if body != nil {
		err := json.NewEncoder(&buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Set authentication based on client type
	if c.Username != "" && c.Password != "" {
		// Basic Auth for peer communication
		req.SetBasicAuth(c.Username, c.Password)
		req.Header.Set("X-Peer", "cmd")
	} else if c.APIKey != "" && c.APISecret != "" {
		// API Key + Secret authentication
		req.Header.Set("X-API-Key", c.APIKey)
		req.Header.Set("X-API-Secret", c.APISecret)
	}

	// Add account ID based on override/body/client default preferences
	accountID := c.resolveAccountID(body, accountIDOverride)
	if accountID != "" {
		req.Header.Set("X-Account-ID", accountID)
	}

	return req, nil
}

func (c *Client) resolveAccountID(body interface{}, accountIDOverride []string) string {
	if len(accountIDOverride) > 0 && accountIDOverride[0] != "" {
		return accountIDOverride[0]
	}

	if accountID := extractAccountIDFromBody(body); accountID != "" {
		return accountID
	}

	return c.AccountID
}

func extractAccountIDFromBody(body interface{}) string {
	if body == nil {
		return ""
	}
	return extractAccountIDFromValue(reflect.ValueOf(body))
}

func extractAccountIDFromValue(val reflect.Value) string {
	if !val.IsValid() {
		return ""
	}

	switch val.Kind() {
	case reflect.Interface, reflect.Ptr:
		if val.IsNil() {
			return ""
		}
		return extractAccountIDFromValue(val.Elem())
	case reflect.Struct:
		field := val.FieldByName("accountId")
		if field.IsValid() && field.CanInterface() {
			if accountID := accountIDValue(field); accountID != "" {
				return accountID
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if accountID := extractAccountIDFromValue(val.Index(i)); accountID != "" {
				return accountID
			}
		}
	}

	return ""
}

func accountIDValue(field reflect.Value) string {
	switch field.Kind() {
	case reflect.String:
		if field.String() != "" {
			return field.String()
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Int() != 0 {
			return fmt.Sprintf("%d", field.Int())
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if field.Uint() != 0 {
			return fmt.Sprintf("%d", field.Uint())
		}
	}
	return ""
}
