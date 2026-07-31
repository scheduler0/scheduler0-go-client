package scheduler0_go_client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper function to create a test client with API key authentication
func createTestAPIClient(server *httptest.Server) *Client {
	u, _ := url.Parse(server.URL)
	client, _ := NewClient(server.URL, "v1",
		WithAPIKey("mock-api-key", "mock-api-secret"),
		WithAccountID("123"))
	client.BaseURL = u
	client.HTTPClient = server.Client()
	return client
}

// Helper function to create a test client with basic authentication
func createTestBasicAuthClient(server *httptest.Server) *Client {
	u, _ := url.Parse(server.URL)
	client, _ := NewClient(server.URL, "v1",
		WithBasicAuth("testuser", "testpass"))
	client.BaseURL = u
	client.HTTPClient = server.Client()
	return client
}

// Helper function to create a test client without authentication
func createTestNoAuthClient(server *httptest.Server) *Client {
	u, _ := url.Parse(server.URL)
	client, _ := NewClient(server.URL, "v1")
	client.BaseURL = u
	client.HTTPClient = server.Client()
	return client
}

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
					ID:          1,
					AccountID:   123,
					Archived:    false,
					APIKey:      "mock-key",
					DateCreated: "2025-01-01T00:00:00Z",
					CreatedBy:   "user-1",
				},
			},
		},
	}

	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/credentials", r.URL.Path)
		assert.Equal(t, "mock-api-key", r.Header.Get("X-API-Key"))
		assert.Equal(t, "mock-api-secret", r.Header.Get("X-Secret-Key"))
		assert.Equal(t, "123", r.Header.Get("X-Account-ID"))
		// Check query parameters
		assert.Equal(t, "10", r.URL.Query().Get("limit"))
		assert.Equal(t, "0", r.URL.Query().Get("offset"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Create client using helper
	client := createTestAPIClient(server)

	// Make call with parameters
	result, err := client.ListCredentials(ListCredentialsParams{
		Limit:            10,
		Offset:           0,
		OrderBy:          "date_created",
		OrderByDirection: "desc",
	})
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 1, result.Data.Total)
	assert.Equal(t, "mock-key", result.Data.Credentials[0].APIKey)
}

func TestCreateCredential(t *testing.T) {
	mockResponse := CredentialResponse{
		Success: true,
		Data: Credential{
			ID:              1,
			AccountID:       123,
			Archived:        false,
			APIKey:          "new-key",
			PlaintextSecret: "new-plaintext-secret",
			DateCreated:     "2025-01-01T00:00:00Z",
			CreatedBy:       "user-1",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/credentials", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "mock-api-key", r.Header.Get("X-API-Key"))
		assert.Equal(t, "mock-api-secret", r.Header.Get("X-Secret-Key"))
		assert.Equal(t, "123", r.Header.Get("X-Account-ID"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := createTestAPIClient(server)

	body := &CredentialCreateRequestBody{
		CreatedBy: "user-1",
	}
	result, err := client.CreateCredential(body)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "new-key", result.Data.APIKey)
}

// TestCreateCredential_AccountIDFromBody tests that AccountID is extracted from the request body
// when the client doesn't have AccountID set. This reproduces the bug where reflection
// looks for "accountId" (camelCase) but the field is "AccountID" (PascalCase).
func TestCreateCredential_AccountIDFromBody(t *testing.T) {
	mockResponse := CredentialResponse{
		Success: true,
		Data: Credential{
			ID:          1,
			AccountID:   456,
			Archived:    false,
			APIKey:      "new-key",
			DateCreated: "2025-01-01T00:00:00Z",
			CreatedBy:   "user-1",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/credentials", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		// Verify basic auth headers are present
		username, password, ok := r.BasicAuth()
		assert.True(t, ok, "Basic auth should be present")
		assert.Equal(t, "testuser", username)
		assert.Equal(t, "testpass", password)
		assert.Equal(t, "cmd", r.Header.Get("X-Peer"))
		// This assertion should fail initially, demonstrating the bug
		// AccountID should be extracted from the request body
		accountIDHeader := r.Header.Get("X-Account-ID")
		assert.Equal(t, "456", accountIDHeader, "X-Account-ID header should be extracted from request body")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Create client WITHOUT AccountID set (like in auth.go)
	client := createTestBasicAuthClient(server)

	// Create credential with AccountID in the body
	body := &CredentialCreateRequestBody{
		AccountID: 456,
		CreatedBy: "user-1",
	}
	result, err := client.CreateCredential(body)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "new-key", result.Data.APIKey)
}

func TestGetCredential(t *testing.T) {
	mockResponse := CredentialResponse{
		Success: true,
		Data: Credential{
			ID:          1,
			AccountID:   123,
			Archived:    false,
			APIKey:      "get-key",
			DateCreated: "2025-01-01T00:00:00Z",
			CreatedBy:   "user-1",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/credentials/1", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "mock-api-key", r.Header.Get("X-API-Key"))
		assert.Equal(t, "mock-api-secret", r.Header.Get("X-Secret-Key"))
		assert.Equal(t, "123", r.Header.Get("X-Account-ID"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := createTestAPIClient(server)

	result, err := client.GetCredential("1")
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "get-key", result.Data.APIKey)
}

func TestUpdateCredential(t *testing.T) {
	mockResponse := CredentialResponse{
		Success: true,
		Data: Credential{
			ID:          1,
			AccountID:   123,
			Archived:    false,
			APIKey:      "updated-key",
			DateCreated: "2025-01-01T00:00:00Z",
			CreatedBy:   "user-1",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/credentials/1", r.URL.Path)
		assert.Equal(t, "PUT", r.Method)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := &Client{
		BaseURL:    u,
		HTTPClient: server.Client(),
		APIKey:     "mock-api-key",
		APISecret:  "mock-api-secret",
		Version:    "v1",
	}

	body := &CredentialUpdateRequestBody{
		ModifiedBy: "user-1",
	}
	result, err := client.UpdateCredential("1", body)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "updated-key", result.Data.APIKey)
}

func TestDeleteCredential(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/credentials/1", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := &Client{
		BaseURL:    u,
		HTTPClient: server.Client(),
		APIKey:     "mock-api-key",
		APISecret:  "mock-api-secret",
		Version:    "v1",
	}

	body := &CredentialDeleteRequestBody{
		DeletedBy: "user-1",
	}
	err := client.DeleteCredential("1", body)
	assert.NoError(t, err)
}

func TestListExecutions(t *testing.T) {
	mockResponse := PaginatedExecutionsResponse{
		Success: true,
		Data: struct {
			Total      int         `json:"total"`
			Offset     int         `json:"offset"`
			Limit      int         `json:"limit"`
			Executions []Execution `json:"executions"`
		}{
			Total:  1,
			Offset: 0,
			Limit:  10,
			Executions: []Execution{
				{
					ID:                    1,
					AccountID:             123,
					UniqueID:              "exec-1",
					State:                 1, // 1 = success
					NodeID:                1,
					JobID:                 1,
					LastExecutionDatetime: "2025-01-01T00:00:00Z",
					NextExecutionDatetime: "2025-01-02T00:00:00Z",
					JobQueueVersion:       1,
					ExecutionVersion:      1,
					DateCreated:           "2025-01-01T00:00:00Z",
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/executions", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := &Client{
		BaseURL:    u,
		HTTPClient: server.Client(),
		APIKey:     "mock-api-key",
		APISecret:  "mock-api-secret",
		Version:    "v1",
	}

	result, err := client.ListExecutions(ListExecutionsParams{
		StartDate: "2025-01-01T00:00:00Z",
		EndDate:   "2025-01-01T23:59:59Z",
		ProjectID: 0,
		JobID:     0,
		AccountID: 0,
		Limit:     10,
		Offset:    0,
	})
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 1, result.Data.Total)
	assert.Equal(t, int64(1), result.Data.Executions[0].State) // 1 = success
}

func TestListExecutors(t *testing.T) {
	mockResponse := PaginatedExecutorsResponse{
		Success: true,
		Data: struct {
			Total     int        `json:"total"`
			Offset    int        `json:"offset"`
			Limit     int        `json:"limit"`
			Executors []Executor `json:"executors"`
		}{
			Total:  1,
			Offset: 0,
			Limit:  10,
			Executors: []Executor{
				{
					ID:               1,
					Name:             "test-executor",
					Type:             "cloud_function",
					Region:           "us-west-1",
					CloudProvider:    "aws",
					CloudResourceURL: "https://example.com/function",
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/executors", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := &Client{
		BaseURL:    u,
		HTTPClient: server.Client(),
		APIKey:     "mock-api-key",
		APISecret:  "mock-api-secret",
		Version:    "v1",
	}

	result, err := client.ListExecutors(ListExecutorsParams{
		AccountID:        0,
		Limit:            10,
		Offset:           0,
		OrderBy:          "date_created",
		OrderByDirection: "desc",
	})
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 1, result.Data.Total)
	assert.Equal(t, "test-executor", result.Data.Executors[0].Name)
}

func TestCreateExecutor(t *testing.T) {
	mockResponse := ExecutorResponse{
		Success: true,
		Data: Executor{
			ID:               1,
			Name:             "new-executor",
			Type:             "cloud_function",
			Region:           "us-west-1",
			CloudProvider:    "aws",
			CloudResourceURL: "https://example.com/function",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/executors", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := &Client{
		BaseURL:    u,
		HTTPClient: server.Client(),
		APIKey:     "mock-api-key",
		APISecret:  "mock-api-secret",
		Version:    "v1",
	}

	body := &ExecutorRequestBody{
		Name:             "new-executor",
		Type:             "cloud_function",
		Region:           "us-west-1",
		CloudProvider:    "aws",
		CloudResourceURL: "https://example.com/function",
	}

	result, err := client.CreateExecutor(body)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "new-executor", result.Data.Name)
}

func TestGetExecutor(t *testing.T) {
	mockResponse := ExecutorResponse{
		Success: true,
		Data: Executor{
			ID:               1,
			Name:             "get-executor",
			Type:             "cloud_function",
			Region:           "us-west-1",
			CloudProvider:    "aws",
			CloudResourceURL: "https://example.com/function",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/executors/1", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := &Client{
		BaseURL:    u,
		HTTPClient: server.Client(),
		APIKey:     "mock-api-key",
		APISecret:  "mock-api-secret",
		Version:    "v1",
	}

	result, err := client.GetExecutor("1")
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "get-executor", result.Data.Name)
}

func TestUpdateExecutor(t *testing.T) {
	mockResponse := ExecutorResponse{
		Success: true,
		Data: Executor{
			ID:               1,
			Name:             "updated-executor",
			Type:             "cloud_function",
			Region:           "us-west-1",
			CloudProvider:    "aws",
			CloudResourceURL: "https://example.com/function",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/executors/1", r.URL.Path)
		assert.Equal(t, "PUT", r.Method)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := &Client{
		BaseURL:    u,
		HTTPClient: server.Client(),
		APIKey:     "mock-api-key",
		APISecret:  "mock-api-secret",
		Version:    "v1",
	}

	body := &ExecutorUpdateRequestBody{
		Name:             "updated-executor",
		Type:             "cloud_function",
		Region:           "us-west-1",
		CloudProvider:    "aws",
		CloudResourceURL: "https://example.com/function",
		ModifiedBy:       "user-1",
	}

	result, err := client.UpdateExecutor("1", body)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "updated-executor", result.Data.Name)
}

func TestDeleteExecutor(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/executors/1", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := &Client{
		BaseURL:    u,
		HTTPClient: server.Client(),
		APIKey:     "mock-api-key",
		APISecret:  "mock-api-secret",
		Version:    "v1",
	}

	body := &ExecutorDeleteRequestBody{
		DeletedBy: "user-1",
	}
	err := client.DeleteExecutor("1", body)
	assert.NoError(t, err)
}

func TestTestInvokeExecutor(t *testing.T) {
	mockResponse := TestInvocationResponse{
		Success: true,
		Data: TestInvocationResult{
			Test:         true,
			ExecutorID:   1,
			ExecutorType: "webhook_url",
			Success:      true,
			DurationMs:   142,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/executors/1/test-invoke", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := &Client{
		BaseURL:    u,
		HTTPClient: server.Client(),
		APIKey:     "mock-api-key",
		APISecret:  "mock-api-secret",
		Version:    "v1",
	}

	body := &TestInvocationRequestBody{
		Job: &Job{
			Spec:     "0 2 * * *",
			Data:     "{\"action\":\"process_data\"}",
			Timezone: "UTC",
			RetryMax: 2,
		},
		Age: "24h",
	}

	result, err := client.TestInvokeExecutor("1", body)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.True(t, result.Data.Success)
	assert.Equal(t, "webhook_url", result.Data.ExecutorType)
}

func TestListProjects(t *testing.T) {
	mockResponse := PaginatedProjectsResponse{
		Success: true,
		Data: struct {
			Total    int       `json:"total"`
			Offset   int       `json:"offset"`
			Limit    int       `json:"limit"`
			Projects []Project `json:"projects"`
		}{
			Total:  1,
			Offset: 0,
			Limit:  10,
			Projects: []Project{
				{
					ID:          1,
					AccountID:   123,
					Name:        "Test Project",
					Description: "Test Description",
					DateCreated: "2025-01-01T00:00:00Z",
					CreatedBy:   "user-1",
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/projects", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := &Client{
		BaseURL:    u,
		HTTPClient: server.Client(),
		APIKey:     "mock-api-key",
		APISecret:  "mock-api-secret",
		Version:    "v1",
	}

	result, err := client.ListProjects(ListProjectsParams{
		AccountID:        0,
		Limit:            10,
		Offset:           0,
		OrderBy:          "date_created",
		OrderByDirection: "desc",
	})
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 1, result.Data.Total)
	assert.Equal(t, "Test Project", result.Data.Projects[0].Name)
}

func TestCreateProject(t *testing.T) {
	mockResponse := ProjectResponse{
		Success: true,
		Data: Project{
			ID:          1,
			AccountID:   123,
			Name:        "New Project",
			Description: "New Description",
			DateCreated: "2025-01-01T00:00:00Z",
			CreatedBy:   "user-1",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/projects", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := &Client{
		BaseURL:    u,
		HTTPClient: server.Client(),
		APIKey:     "mock-api-key",
		APISecret:  "mock-api-secret",
		Version:    "v1",
	}

	body := &ProjectRequestBody{
		Name:        "New Project",
		Description: "New Description",
		CreatedBy:   "user-1",
	}

	result, err := client.CreateProject(body)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "New Project", result.Data.Name)
}

func TestGetProject(t *testing.T) {
	mockResponse := ProjectResponse{
		Success: true,
		Data: Project{
			ID:          1,
			AccountID:   123,
			Name:        "Get Project",
			Description: "Get Description",
			DateCreated: "2025-01-01T00:00:00Z",
			CreatedBy:   "user-1",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/projects/1", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := &Client{
		BaseURL:    u,
		HTTPClient: server.Client(),
		APIKey:     "mock-api-key",
		APISecret:  "mock-api-secret",
		Version:    "v1",
	}

	result, err := client.GetProject(1)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "Get Project", result.Data.Name)
}

func TestUpdateProject(t *testing.T) {
	mockResponse := ProjectResponse{
		Success: true,
		Data: Project{
			ID:          1,
			AccountID:   123,
			Name:        "Updated Project",
			Description: "Updated Description",
			DateCreated: "2025-01-01T00:00:00Z",
			CreatedBy:   "user-1",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/projects/1", r.URL.Path)
		assert.Equal(t, "PUT", r.Method)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := &Client{
		BaseURL:    u,
		HTTPClient: server.Client(),
		APIKey:     "mock-api-key",
		APISecret:  "mock-api-secret",
		Version:    "v1",
	}

	body := &ProjectUpdateRequestBody{
		Description: "Updated Description",
		ModifiedBy:  "user-1",
	}

	result, err := client.UpdateProject(1, body)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "Updated Description", result.Data.Description)
}

func TestDeleteProject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/projects/1", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := &Client{
		BaseURL:    u,
		HTTPClient: server.Client(),
		APIKey:     "mock-api-key",
		APISecret:  "mock-api-secret",
		Version:    "v1",
	}

	body := &ProjectDeleteRequestBody{
		DeletedBy: "user-1",
	}
	err := client.DeleteProject(1, body)
	assert.NoError(t, err)
}

func TestListJobs(t *testing.T) {
	mockResponse := PaginatedJobsResponse{
		Success: true,
		Data: struct {
			Total  int   `json:"total"`
			Offset int   `json:"offset"`
			Limit  int   `json:"limit"`
			Jobs   []Job `json:"jobs"`
		}{
			Total:  1,
			Offset: 0,
			Limit:  10,
			Jobs: []Job{
				{
					ID:         1,
					AccountID:  123,
					ProjectID:  1,
					Data:       "job data",
					ExecutorID: &[]int64{1}[0],
					Spec:       "0 30 * * * *",
					StartDate:  "2025-01-01T00:00:00Z",
					EndDate:    "2025-12-31T00:00:00Z",
					Timezone:   "UTC",
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/jobs", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := &Client{
		BaseURL:    u,
		HTTPClient: server.Client(),
		APIKey:     "mock-api-key",
		APISecret:  "mock-api-secret",
		Version:    "v1",
	}

	result, err := client.ListJobs(ListJobsParams{
		ProjectID:        "proj-1",
		Limit:            10,
		Offset:           0,
		OrderBy:          "date_created",
		OrderByDirection: "desc",
	})
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 1, result.Data.Total)
	assert.Equal(t, "job data", result.Data.Jobs[0].Data)
}

func TestCreateJob(t *testing.T) {
	mockResponse := BatchJobResponse{
		Success: true,
		Data:    "request-id-123",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/jobs", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", "/async-tasks/request-id-123")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := &Client{
		BaseURL:    u,
		HTTPClient: server.Client(),
		APIKey:     "mock-api-key",
		APISecret:  "mock-api-secret",
		Version:    "v1",
	}

	body := &JobRequestBody{
		ProjectID: 1,
		Timezone:  "UTC",
		Data:      "New Job",
		StartDate: "2025-01-01T00:00:00Z",
		EndDate:   "2025-12-31T00:00:00Z",
		CreatedBy: "user-1",
	}

	result, err := client.CreateJob(body)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "request-id-123", result.Data)
}

func TestGetJob(t *testing.T) {
	mockResponse := JobResponse{
		Success: true,
		Data: Job{
			ID:         1,
			AccountID:  123,
			ProjectID:  1,
			Data:       "job data",
			ExecutorID: &[]int64{1}[0],
			Spec:       "0 30 * * * *",
			StartDate:  "2025-01-01T00:00:00Z",
			EndDate:    "2025-12-31T00:00:00Z",
			Timezone:   "UTC",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/jobs/job-1", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := &Client{
		BaseURL:    u,
		HTTPClient: server.Client(),
		APIKey:     "mock-api-key",
		APISecret:  "mock-api-secret",
		Version:    "v1",
	}

	result, err := client.GetJob("job-1")
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "job data", result.Data.Data)
}

func TestUpdateJob(t *testing.T) {
	mockResponse := JobResponse{
		Success: true,
		Data: Job{
			ID:         1,
			AccountID:  123,
			ProjectID:  1,
			Data:       "job data",
			ExecutorID: &[]int64{1}[0],
			Spec:       "0 45 * * * *",
			StartDate:  "2025-01-01T00:00:00Z",
			EndDate:    "2025-12-31T00:00:00Z",
			Timezone:   "UTC",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/jobs/job-1", r.URL.Path)
		assert.Equal(t, "PUT", r.Method)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := &Client{
		BaseURL:    u,
		HTTPClient: server.Client(),
		APIKey:     "mock-api-key",
		APISecret:  "mock-api-secret",
		Version:    "v1",
	}

	body := &JobUpdateRequestBody{
		ProjectID: 1,
		Spec:      "0 45 * * * *",
		Data:      "Updated Job",
	}

	result, err := client.UpdateJob("job-1", body)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "0 45 * * * *", result.Data.Spec)
}

func TestDeleteJob(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/jobs/job-1", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := &Client{
		BaseURL:    u,
		HTTPClient: server.Client(),
		APIKey:     "mock-api-key",
		APISecret:  "mock-api-secret",
		Version:    "v1",
	}

	body := &JobDeleteRequestBody{
		DeletedBy: "user-1",
	}
	err := client.DeleteJob("job-1", body)
	assert.NoError(t, err)
}

func TestHealthcheck(t *testing.T) {
	mockResponse := HealthcheckResponse{
		Success: true,
		Data: HealthcheckData{
			LeaderAddress: "127.0.0.1:7070",
			LeaderID:      "1",
			RaftStats: RaftStats{
				AppliedIndex:        "162",
				CommitIndex:         "162",
				FSMPending:          "0",
				LastContact:         "0",
				LastLogIndex:        "162",
				LastLogTerm:         "7",
				LastSnapshotIndex:   "55",
				LastSnapshotTerm:    "5",
				LatestConfiguration: "[{Suffrage:Voter ID:1 Address:127.0.0.1:7070}]",
				LatestConfigIndex:   "0",
				NumPeers:            "0",
				ProtocolVersion:     "3",
				ProtocolVersionMax:  "3",
				ProtocolVersionMin:  "0",
				SnapshotVersionMax:  "1",
				SnapshotVersionMin:  "0",
				State:               "Leader",
				Term:                "7",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/healthcheck", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := &Client{
		BaseURL:    u,
		HTTPClient: server.Client(),
		APIKey:     "mock-api-key",
		APISecret:  "mock-api-secret",
		Version:    "v1",
	}

	result, err := client.Healthcheck()
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "127.0.0.1:7070", result.Data.LeaderAddress)
	assert.Equal(t, "Leader", result.Data.RaftStats.State)
}

// Test new authentication methods

func TestNewClientWithOptions(t *testing.T) {
	// Test API key authentication
	client, err := NewClient("https://api.test.com", "v1",
		WithAPIKey("test-key", "test-secret"),
		WithAccountID("456"))
	assert.NoError(t, err)
	assert.Equal(t, "test-key", client.APIKey)
	assert.Equal(t, "test-secret", client.APISecret)
	assert.Equal(t, "456", client.AccountID)

	// Test basic authentication
	client, err = NewClient("https://api.test.com", "v1",
		WithBasicAuth("user", "pass"))
	assert.NoError(t, err)
	assert.Equal(t, "user", client.Username)
	assert.Equal(t, "pass", client.Password)

	// Test no authentication
	client, err = NewClient("https://api.test.com", "v1")
	assert.NoError(t, err)
	assert.Equal(t, "", client.APIKey)
	assert.Equal(t, "", client.Username)
}

func TestConvenienceFunctions(t *testing.T) {
	// Test NewAPIClient
	client, err := NewAPIClient("https://api.test.com", "v1", "key", "secret")
	assert.NoError(t, err)
	assert.Equal(t, "key", client.APIKey)
	assert.Equal(t, "secret", client.APISecret)

	// Test NewAPIClientWithAccount
	client, err = NewAPIClientWithAccount("https://api.test.com", "v1", "key", "secret", "123")
	assert.NoError(t, err)
	assert.Equal(t, "key", client.APIKey)
	assert.Equal(t, "secret", client.APISecret)
	assert.Equal(t, "123", client.AccountID)

	// Test NewBasicAuthClient
	client, err = NewBasicAuthClient("https://api.test.com", "v1", "user", "pass")
	assert.NoError(t, err)
	assert.Equal(t, "user", client.Username)
	assert.Equal(t, "pass", client.Password)
}

func TestBasicAuthClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check basic auth
		username, password, ok := r.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "testuser", username)
		assert.Equal(t, "testpass", password)
		assert.Equal(t, "cmd", r.Header.Get("X-Peer"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	client := createTestBasicAuthClient(server)

	// Test that basic auth is set correctly
	assert.Equal(t, "testuser", client.Username)
	assert.Equal(t, "testpass", client.Password)
}

func TestNoAuthClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Should not have any auth headers
		assert.Equal(t, "", r.Header.Get("X-API-Key"))
		assert.Equal(t, "", r.Header.Get("X-Secret-Key"))
		assert.Equal(t, "", r.Header.Get("X-Account-ID"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"status": "ok", "success": true})
	}))
	defer server.Close()

	client := createTestNoAuthClient(server)

	// Test healthcheck without auth
	result, err := client.Healthcheck()
	assert.NoError(t, err)
	assert.True(t, result.Success)
}

// Test new methods added

func TestArchiveCredential(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/credentials/1/archive", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := createTestAPIClient(server)

	err := client.ArchiveCredential("1", "user-1")
	assert.NoError(t, err)
}

func TestAddAllFeaturesToAccount(t *testing.T) {
	mockResponse := map[string]interface{}{
		"success": true,
		"data": map[string]string{
			"message": "All features added successfully",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/accounts/123/features/all", r.URL.Path)
		assert.Equal(t, "PUT", r.Method)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := createTestAPIClient(server)

	err := client.AddAllFeaturesToAccount("123")
	assert.NoError(t, err)
}

func TestRemoveAllFeaturesFromAccount(t *testing.T) {
	mockResponse := map[string]interface{}{
		"success": true,
		"data": map[string]string{
			"message": "All features removed successfully",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/accounts/123/features/all", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := createTestAPIClient(server)

	err := client.RemoveAllFeaturesFromAccount("123")
	assert.NoError(t, err)
}

func TestCreateJobFromPrompt(t *testing.T) {
	mockResponse := promptJobsEnvelope{
		Success: true,
		Data: PromptResult{
			Providers: []PromptProviderResult{
				{
					Provider: "openai",
					Model:    "gpt-4",
					Jobs: []PromptJobResponse{
						{
							Kind:           "FOLLOW_UP",
							Purpose:        "Send follow-up email",
							Subject:        "Follow up on your request",
							CronExpression: "0 9 * * *",
							Recurrence:     "every day",
							Timezone:       "UTC",
						},
					},
					InputTokens:  100,
					OutputTokens: 50,
					TotalTokens:  150,
					DurationMs:   300,
				},
			},
			Classification: &IntentClassification{
				Text:     "Create a job to send follow-up emails daily at 9 AM",
				Decision: "allow",
				Reason:   "request_with_temporal_signal",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/ai/prompt", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "123", r.Header.Get("X-Account-ID"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := createTestAPIClient(server)

	body := &PromptJobRequest{
		Prompt:   "Create a job to send follow-up emails daily at 9 AM",
		Purposes: []string{"follow-up"},
		Timezone: "UTC",
	}

	result, err := client.CreateJobFromPrompt(body)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Providers, 1)
	assert.Equal(t, "FOLLOW_UP", result.Providers[0].Jobs[0].Kind)
	assert.Equal(t, "Send follow-up email", result.Providers[0].Jobs[0].Purpose)
	assert.NotNil(t, result.Classification)
	assert.Equal(t, "allow", result.Classification.Decision)
}

func TestCreateJobFromPrompt_SkippedError(t *testing.T) {
	mockBody := map[string]any{
		"success": false,
		"data": map[string]any{
			"message": "prompt skipped (reject): informational_question_not_schedule_request",
			"classification": map[string]any{
				"text":     "What is Kubernetes?",
				"decision": "reject",
				"reason":   "informational_question_not_schedule_request",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(mockBody)
	}))
	defer server.Close()

	client := createTestAPIClient(server)

	_, err := client.CreateJobFromPrompt(&PromptJobRequest{Prompt: "What is Kubernetes?"})
	assert.Error(t, err)
	assert.True(t, IsPromptSkippedError(err), "expected PromptSkippedError")

	var skipped *PromptSkippedError
	assert.ErrorAs(t, err, &skipped)
	assert.NotNil(t, skipped.Classification)
	assert.Equal(t, "reject", skipped.Classification.Decision)
}

func TestClassifyPrompt(t *testing.T) {
	mockResponse := classifyEnvelope{
		Success: true,
		Data: struct {
			Classification IntentClassification `json:"classification"`
		}{
			Classification: IntentClassification{
				Text:     "Remind me every Monday at 9am",
				Decision: "allow",
				Reason:   "request_with_temporal_signal",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/ai/prompt/classify", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := createTestAPIClient(server)

	result, err := client.ClassifyPrompt(&ClassifyPromptRequest{Prompt: "Remind me every Monday at 9am"})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "allow", result.Decision)
	assert.Equal(t, "request_with_temporal_signal", result.Reason)
}

func TestAnalyzeSuggestions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/ai/suggestions/analyze", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"success": true,
			"data": {
				"request_id": "req_1",
				"conversation_id": "conv_123",
				"suggestions": [{"id":"sug_001","type":"COMMITMENT","status":"OPEN","confidence":0.95}],
				"obligations": [{"id":"obl_001","status":"OPEN","suggestion_id":"sug_001"}],
				"warnings": [],
				"engine": {"engine_version":"1.0.0"}
			}
		}`))
	}))
	defer server.Close()

	client := createTestAPIClient(server)

	result, err := client.AnalyzeSuggestions(&AnalyzeSuggestionsRequest{
		ConversationID: "conv_123",
		Messages: []SuggestionMessage{
			{Speaker: "Victor", Timestamp: "2026-07-17T10:00:00-04:00", Message: "I'll send the proposal tomorrow."},
		},
		Options: &SuggestionOptions{Locale: "en", DefaultTimezone: "America/Toronto"},
	})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "conv_123", result.ConversationID)
	assert.Len(t, result.Suggestions, 1)
	assert.Equal(t, "COMMITMENT", result.Suggestions[0]["type"])
	assert.Len(t, result.Obligations, 1)
}

func TestSendTimeSuggestions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/ai/suggestions/time", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"success": true,
			"data": {
				"request_id": "req_1",
				"reference_time": "2026-07-17T17:45:00-04:00",
				"policy": {"id":"default_send_time","version":"1.0.0"},
				"engine": {"version":"1.0.0"},
				"suggestions": [{"id":"sts_001","send_at":"2026-07-20T12:00:00-04:00","label":"Monday morning","score":0.94,"rank":1}],
				"search": {"candidates_generated": 143, "candidates_scored": 16},
				"warnings": []
			}
		}`))
	}))
	defer server.Close()

	client := createTestAPIClient(server)

	result, err := client.SendTimeSuggestions(&SendTimeSuggestionsRequest{
		Sender: &SendTimeParticipant{ID: "user_123", Timezone: "America/Toronto"},
		Recipients: []SendTimeParticipant{
			{ID: "user_456", Timezone: "America/Los_Angeles", Role: "primary"},
		},
		Message: &SendTimeMessage{Priority: "normal"},
	})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "req_1", result.RequestID)
	assert.Equal(t, "2026-07-17T17:45:00-04:00", result.ReferenceTime)
	assert.Len(t, result.Suggestions, 1)
	assert.Equal(t, "sts_001", result.Suggestions[0]["id"])
}

func TestScheduleFromPrompt(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/ai/schedule", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "123", r.Header.Get("X-Account-ID"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"success": true,
			"data": {
				"classification": {"text":"remind the team every monday","decision":"allow","reason":"request_with_temporal_signal"},
				"project": {"id": 7, "name": "Team reminders", "description": "auto"},
				"projectCreated": true,
				"executor": {"id": 3, "name": "Email sender", "description": "sends email", "tags": ["email"]},
				"executorMatchedBy": "llm",
				"executorMatchReason": "matches email channel",
				"jobs": [{"id": 11, "projectId": 7, "executorId": 3, "spec": "0 9 * * 1", "timezone": "UTC", "status": "active"}],
				"provider": "openai",
				"model": "gpt-4"
			}
		}`))
	}))
	defer server.Close()

	client := createTestAPIClient(server)

	result, err := client.ScheduleFromPrompt(&SchedulePromptRequest{
		Prompt:    "Remind the team every Monday at 9am",
		Channels:  []string{"email"},
		CreatedBy: "victor",
	})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(7), result.Project.ID)
	assert.True(t, result.ProjectCreated)
	assert.Equal(t, int64(3), result.Executor.ID)
	assert.Equal(t, "llm", result.ExecutorMatchedBy)
	assert.Len(t, result.Jobs, 1)
	assert.Equal(t, int64(11), result.Jobs[0].ID)
	assert.NotNil(t, result.Classification)
	assert.Equal(t, "allow", result.Classification.Decision)
}

func TestScheduleFromPrompt_SkippedError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/ai/schedule", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{
			"success": false,
			"data": {
				"message": "prompt skipped (reject): informational_question_not_schedule_request",
				"classification": {"text":"what is kubernetes?","decision":"reject","reason":"informational_question_not_schedule_request"}
			}
		}`))
	}))
	defer server.Close()

	client := createTestAPIClient(server)

	_, err := client.ScheduleFromPrompt(&SchedulePromptRequest{Prompt: "What is Kubernetes?", CreatedBy: "victor"})
	assert.Error(t, err)
	assert.True(t, IsPromptSkippedError(err), "expected PromptSkippedError")
}

func TestBatchCreateJobs(t *testing.T) {
	mockResponse := BatchJobResponse{
		Success: true,
		Data:    "request-id-456",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/jobs", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", "/async-tasks/request-id-456")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := createTestAPIClient(server)

	jobs := []JobRequestBody{
		{
			ProjectID: 1,
			Timezone:  "UTC",
			Data:      "Job 1",
			CreatedBy: "user-1",
		},
		{
			ProjectID: 1,
			Timezone:  "UTC",
			Data:      "Job 2",
			CreatedBy: "user-1",
		},
	}

	result, err := client.BatchCreateJobs(jobs)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "request-id-456", result.Data)
}

func TestListFeatures(t *testing.T) {
	mockResponse := FeaturesResponse{
		Success: true,
		Data: []Feature{
			{
				ID:          1,
				Name:        "feature-1",
				DateCreated: "2025-01-01T00:00:00Z",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/features", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := createTestAPIClient(server)

	result, err := client.ListFeatures()
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Len(t, result.Data, 1)
	assert.Equal(t, "feature-1", result.Data[0].Name)
}

func TestGetAccount(t *testing.T) {
	mockResponse := AccountResponse{
		Success: true,
		Data: Account{
			ID:          123,
			Name:        "Test Account",
			DateCreated: "2025-01-01T00:00:00Z",
			Features: []AccountFeature{
				{
					AccountID: 123,
					FeatureID: 1,
					Feature:   "feature-1",
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/accounts/123", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := createTestAPIClient(server)

	result, err := client.GetAccount("123")
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "Test Account", result.Data.Name)
	assert.Len(t, result.Data.Features, 1)
}

func TestCreateAccount(t *testing.T) {
	mockResponse := AccountResponse{
		Success: true,
		Data: Account{
			ID:          123,
			Name:        "New Account",
			DateCreated: "2025-01-01T00:00:00Z",
			Features:    []AccountFeature{},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/accounts", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := createTestAPIClient(server)

	body := &AccountCreateRequestBody{
		Name: "New Account",
	}

	result, err := client.CreateAccount(body)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "New Account", result.Data.Name)
}

func TestAddFeatureToAccount(t *testing.T) {
	mockResponse := FeatureRequestResponse{
		Success: true,
		Data: FeatureRequest{
			FeatureID: 1,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/accounts/123/feature", r.URL.Path)
		assert.Equal(t, "PUT", r.Method)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := createTestAPIClient(server)

	body := &FeatureRequest{
		FeatureID: 1,
	}

	result, err := client.AddFeatureToAccount("123", body)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, int64(1), result.Data.FeatureID)
}

func TestRemoveFeatureFromAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/accounts/123/feature", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := createTestAPIClient(server)

	body := &FeatureRequest{
		FeatureID: 1,
	}

	err := client.RemoveFeatureFromAccount("123", body)
	assert.NoError(t, err)
}

func TestGetAsyncTask(t *testing.T) {
	mockResponse := AsyncTaskResponse{
		Success: true,
		Data: AsyncTask{
			ID:          1,
			RequestID:   "request-123",
			Input:       "input data",
			Output:      "output data",
			Service:     "job-service",
			State:       2, // Success
			DateCreated: "2025-01-01T00:00:00Z",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/async-tasks/request-123", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := createTestAPIClient(server)

	result, err := client.GetAsyncTask("request-123")
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "request-123", result.Data.RequestID)
	assert.Equal(t, 2, result.Data.State)
}

func TestGetDateRangeAnalytics(t *testing.T) {
	mockResponse := DateRangeAnalyticsAPIResponse{
		Success: true,
		Data: DateRangeAnalyticsResponse{
			AccountID: uint64(123),
			Timezone:  "UTC",
			StartDate: "2025-01-01",
			StartTime: "00:00:00",
			EndDate:   "2025-01-01",
			EndTime:   "23:59:59",
			Points: []DateRangeAnalyticsPoint{
				{
					Date:      "2025-01-01",
					Time:      "00:00:00",
					Scheduled: uint64(10),
					Success:   uint64(8),
					Failed:    uint64(2),
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/executions/analytics", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "2025-01-01", r.URL.Query().Get("startDate"))
		assert.Equal(t, "00:00:00", r.URL.Query().Get("startTime"))
		assert.Equal(t, "123", r.Header.Get("X-Account-ID"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := &Client{
		BaseURL:    u,
		HTTPClient: server.Client(),
		APIKey:     "mock-api-key",
		APISecret:  "mock-api-secret",
		Version:    "v1",
	}

	result, err := client.GetDateRangeAnalytics(GetDateRangeAnalyticsParams{
		StartDate: "2025-01-01",
		StartTime: "00:00:00",
		AccountID: 123,
	})
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, uint64(123), result.Data.AccountID)
	assert.Equal(t, 1, len(result.Data.Points))
	assert.Equal(t, uint64(10), result.Data.Points[0].Scheduled)
}

func TestGetExecutionTotals(t *testing.T) {
	mockResponse := ExecutionTotalsAPIResponse{
		Success: true,
		Data: ExecutionTotalsResponse{
			AccountID: 123,
			Scheduled: 100,
			Success:   80,
			Failed:    20,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/executions/totals", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "123", r.Header.Get("X-Account-ID"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := &Client{
		BaseURL:    u,
		HTTPClient: server.Client(),
		APIKey:     "mock-api-key",
		APISecret:  "mock-api-secret",
		Version:    "v1",
	}

	result, err := client.GetExecutionTotals(123)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, uint64(123), result.Data.AccountID)
	assert.Equal(t, uint64(100), result.Data.Scheduled)
	assert.Equal(t, uint64(80), result.Data.Success)
	assert.Equal(t, uint64(20), result.Data.Failed)
}

func TestCleanupOldExecutionLogs(t *testing.T) {
	mockResponse := CleanupOldLogsResponse{
		Success: true,
		Data: struct {
			Message string `json:"message"`
		}{
			Message: "Old execution logs cleaned up successfully for account 123",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/executions/cleanup-old-logs", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "123", r.Header.Get("X-Account-ID"))

		// Verify request body
		var requestBody CleanupOldLogsRequestBody
		json.NewDecoder(r.Body).Decode(&requestBody)
		assert.Equal(t, "123", requestBody.AccountID)
		assert.Equal(t, 6, requestBody.RetentionMonths)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := &Client{
		BaseURL:    u,
		HTTPClient: server.Client(),
		APIKey:     "mock-api-key",
		APISecret:  "mock-api-secret",
		Version:    "v1",
	}

	result, err := client.CleanupOldExecutionLogs("123", 6)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Contains(t, result.Data.Message, "cleaned up successfully")
}

func TestRegisterLocalExecutor(t *testing.T) {
	mockResponse := LocalExecutorRegisterResponse{Success: true}
	mockResponse.Data.ID = 42

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/local-executors", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := createTestAPIClient(server)

	result, err := client.RegisterLocalExecutor(&LocalExecutorRegisterRequest{
		Name:       "My Local Executor",
		Command:    "/usr/local/bin/process-job.sh",
		WorkingDir: "/home/deploy/app",
		CreatedBy:  "user-1",
	})
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, int64(42), result.Data.ID)
}

func TestPullLocalExecutorJobs(t *testing.T) {
	mockResponse := LocalExecutorJobsResponse{
		Success: true,
		Data: []Job{
			{ID: 1, AccountID: 123, ProjectID: 456, Spec: "* * * * *", Timezone: "UTC"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/local-executors/42/jobs", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := createTestAPIClient(server)

	result, err := client.PullLocalExecutorJobs(42)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Len(t, result.Data, 1)
	assert.Equal(t, int64(1), result.Data[0].ID)
}

func TestReportLocalExecutions(t *testing.T) {
	mockResponse := ReportLocalExecutionsResponse{Success: true}
	mockResponse.Data.Committed = 2

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/local-executors/42/executions", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var reports []LocalExecutionReport
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&reports))
		assert.Len(t, reports, 2)
		assert.Equal(t, "exec-1", reports[0].UniqueID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := createTestAPIClient(server)

	reports := []LocalExecutionReport{
		{
			JobID:             1,
			UniqueID:          "exec-1",
			State:             1,
			LastExecutionTime: "2025-01-01T00:00:00Z",
			NextExecutionTime: "2025-01-02T00:00:00Z",
			ExecutionVersion:  5,
			JobQueueVersion:   2,
		},
		{JobID: 2, UniqueID: "exec-2", State: 2},
	}

	result, err := client.ReportLocalExecutions(42, reports)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 2, result.Data.Committed)
}

// TestAIUsageResponse_DecodesUsage verifies the client decodes the log-derived AI usage
// payload: per-dimension limit/used/remaining plus the period boundary.
func TestAIUsageResponse_DecodesUsage(t *testing.T) {
	payload := `{"success":true,"data":{"accountId":123,"periodStart":"2025-01-01T00:00:00Z","nextResetDate":"2025-02-01T00:00:00Z","prompt":{"limit":100000,"used":42,"remaining":99958},"classify":{"limit":1000,"used":250,"remaining":750}}}`

	var resp AIUsageResponse
	assert.NoError(t, json.Unmarshal([]byte(payload), &resp))
	assert.True(t, resp.Success)
	assert.Equal(t, int64(123), resp.Data.AccountID)
	assert.Equal(t, "2025-02-01T00:00:00Z", resp.Data.NextResetDate)
	assert.Equal(t, int64(100000), resp.Data.Prompt.Limit)
	assert.Equal(t, int64(42), resp.Data.Prompt.Used)
	assert.Equal(t, int64(99958), resp.Data.Prompt.Remaining)
	assert.Equal(t, int64(1000), resp.Data.Classify.Limit)
	assert.Equal(t, int64(250), resp.Data.Classify.Used)
	assert.Equal(t, int64(750), resp.Data.Classify.Remaining)
}

// TestListPromptRequests_SendsLimit verifies that a non-zero limit is forwarded verbatim.
func TestListPromptRequests_SendsLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/ai/prompt-requests", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "50", r.URL.Query().Get("limit"))
		assert.Equal(t, "10", r.URL.Query().Get("offset"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PromptRequestsResponse{Success: true, Data: PromptRequestsData{Limit: 50, Offset: 10}})
	}))
	defer server.Close()

	client := createTestAPIClient(server)
	result, err := client.ListPromptRequests(ListPromptRequestsParams{Limit: 50, Offset: 10})
	assert.NoError(t, err)
	assert.True(t, result.Success)
}

// TestListPromptRequests_OmitsZeroLimit verifies that a zero limit is omitted from the query
// so the server applies its authoritative default page size.
func TestListPromptRequests_OmitsZeroLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/ai/prompt-requests", r.URL.Path)
		_, hasLimit := r.URL.Query()["limit"]
		assert.False(t, hasLimit, "limit query param should be omitted when Limit is zero")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PromptRequestsResponse{Success: true, Data: PromptRequestsData{Limit: 25}})
	}))
	defer server.Close()

	client := createTestAPIClient(server)
	result, err := client.ListPromptRequests(ListPromptRequestsParams{Limit: 0, Offset: 0})
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, uint64(25), result.Data.Limit)
}
