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
					APISecret:   "mock-secret",
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
		assert.Equal(t, "mock-api-secret", r.Header.Get("X-API-Secret"))
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
	result, err := client.ListCredentials(10, 0, "date_created", "desc")
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 1, result.Data.Total)
	assert.Equal(t, "mock-key", result.Data.Credentials[0].APIKey)
}

func TestCreateCredential(t *testing.T) {
	mockResponse := CredentialResponse{
		Success: true,
		Data: Credential{
			ID:          1,
			AccountID:   123,
			Archived:    false,
			APIKey:      "new-key",
			APISecret:   "new-secret",
			DateCreated: "2025-01-01T00:00:00Z",
			CreatedBy:   "user-1",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/credentials", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "mock-api-key", r.Header.Get("X-API-Key"))
		assert.Equal(t, "mock-api-secret", r.Header.Get("X-API-Secret"))
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

func TestGetCredential(t *testing.T) {
	mockResponse := CredentialResponse{
		Success: true,
		Data: Credential{
			ID:          1,
			AccountID:   123,
			Archived:    false,
			APIKey:      "get-key",
			APISecret:   "get-secret",
			DateCreated: "2025-01-01T00:00:00Z",
			CreatedBy:   "user-1",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/credentials/1", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "mock-api-key", r.Header.Get("X-API-Key"))
		assert.Equal(t, "mock-api-secret", r.Header.Get("X-API-Secret"))
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
			APISecret:   "updated-secret",
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

	result, err := client.ListExecutions("2025-01-01T00:00:00Z", "2025-01-01T23:59:59Z", 0, 0, 10, 0)
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

	result, err := client.ListExecutors(10, 0, "date_created", "desc")
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

	result, err := client.ListProjects(10, 0, "date_created", "desc")
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

	result, err := client.ListJobs("proj-1", 10, 0, "date_created", "desc")
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
		assert.Equal(t, "", r.Header.Get("X-API-Secret"))
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
	mockResponse := []PromptJobResponse{
		{
			Kind:           "FOLLOW_UP",
			Purpose:        "Send follow-up email",
			Subject:        "Follow up on your request",
			CronExpression: "0 9 * * *",
			Recurrence:     "every day",
			Timezone:       "UTC",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/prompt", r.URL.Path)
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
	assert.Len(t, result, 1)
	assert.Equal(t, "FOLLOW_UP", result[0].Kind)
	assert.Equal(t, "Send follow-up email", result[0].Purpose)
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
