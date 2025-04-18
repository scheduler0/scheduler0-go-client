package scheduler0_go_client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListCredentials(t *testing.T) {
	// Prepare mock response
	mockResponse := PaginatedCredentialsResponse{
		Success: true,
		Data: struct {
			Total       int          `json:"total"`
			Offset      int          `json:"offset"`
			Limit       int          `json:"limit"`
			Credentials []Credential `json:"credentials"`
		}{
			Total:  1,
			Offset: 0,
			Limit:  10,
			Credentials: []Credential{
				{
					ID:        1,
					Archived:  false,
					APIKey:    "mock-key",
					APISecret: "mock-secret",
					CreatedAt: "2025-01-01T00:00:00Z",
				},
			},
		},
	}

	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/credentials", r.URL.Path)
		assert.Equal(t, "mock-api-key", r.Header.Get("X-API-Key"))
		assert.Equal(t, "mock-api-secret", r.Header.Get("X-API-Secret"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Parse mock URL
	u, _ := url.Parse(server.URL)

	// Create client
	client := &Client{
		BaseURL:    u,
		HTTPClient: server.Client(),
		APIKey:     "mock-api-key",
		APISecret:  "mock-api-secret",
		Version:    "v1",
	}

	// Make call
	result, err := client.ListCredentials()
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 1, result.Data.Total)
	assert.Equal(t, "mock-key", result.Data.Credentials[0].APIKey)
}
