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

func TestCreateCredential(t *testing.T) {
	mockResponse := CredentialResponse{
		Success: true,
		Data: Credential{
			ID:        1,
			Archived:  false,
			APIKey:    "new-key",
			APISecret: "new-secret",
			CreatedAt: "2025-01-01T00:00:00Z",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/credentials", r.URL.Path)
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

	result, err := client.CreateCredential()
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "new-key", result.Data.APIKey)
}

func TestGetCredential(t *testing.T) {
	mockResponse := CredentialResponse{
		Success: true,
		Data: Credential{
			ID:        1,
			Archived:  false,
			APIKey:    "get-key",
			APISecret: "get-secret",
			CreatedAt: "2025-01-01T00:00:00Z",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/credentials/1", r.URL.Path)
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

	result, err := client.GetCredential("1")
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "get-key", result.Data.APIKey)
}

func TestUpdateCredential(t *testing.T) {
	mockResponse := CredentialResponse{
		Success: true,
		Data: Credential{
			ID:        1,
			Archived:  false,
			APIKey:    "updated-key",
			APISecret: "updated-secret",
			CreatedAt: "2025-01-01T00:00:00Z",
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

	result, err := client.UpdateCredential("1")
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

	err := client.DeleteCredential("1")
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
					ID:                1,
					UniqueID:          "exec-1",
					State:             "success",
					NodeID:            "node-1",
					JobID:             "job-1",
					LastExecutionTime: "2025-01-01T00:00:00Z",
					NextExecutionTime: "2025-01-02T00:00:00Z",
					JobQueueVersion:   1,
					ExecutionVersion:  1,
					Logs:              "Job executed successfully",
					CreatedAt:         "2025-01-01T00:00:00Z",
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

	result, err := client.ListExecutions()
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 1, result.Data.Total)
	assert.Equal(t, "success", result.Data.Executions[0].State)
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
					FilePath:         "/path/to/function",
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

	result, err := client.ListExecutors()
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
			FilePath:         "/path/to/function",
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
		FilePath:         "/path/to/function",
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
			FilePath:         "/path/to/function",
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
			FilePath:         "/path/to/function",
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

	body := &ExecutorRequestBody{
		Name:             "updated-executor",
		Type:             "cloud_function",
		Region:           "us-west-1",
		CloudProvider:    "aws",
		FilePath:         "/path/to/function",
		CloudResourceURL: "https://example.com/function",
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

	err := client.DeleteExecutor("1")
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
					ID:          "proj-1",
					Name:        "Test Project",
					Description: "Test Description",
					CreatedAt:   "2025-01-01T00:00:00Z",
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

	result, err := client.ListProjects()
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 1, result.Data.Total)
	assert.Equal(t, "Test Project", result.Data.Projects[0].Name)
}

func TestCreateProject(t *testing.T) {
	mockResponse := ProjectResponse{
		Success: true,
		Data: Project{
			ID:          "proj-1",
			Name:        "New Project",
			Description: "New Description",
			CreatedAt:   "2025-01-01T00:00:00Z",
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
			ID:          "proj-1",
			Name:        "Get Project",
			Description: "Get Description",
			CreatedAt:   "2025-01-01T00:00:00Z",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/projects/proj-1", r.URL.Path)
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

	result, err := client.GetProject("proj-1")
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "Get Project", result.Data.Name)
}

func TestUpdateProject(t *testing.T) {
	mockResponse := ProjectResponse{
		Success: true,
		Data: Project{
			ID:          "proj-1",
			Name:        "Updated Project",
			Description: "Updated Description",
			CreatedAt:   "2025-01-01T00:00:00Z",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/projects/proj-1", r.URL.Path)
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
	}

	result, err := client.UpdateProject("proj-1", body)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "Updated Description", result.Data.Description)
}

func TestDeleteProject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/projects/proj-1", r.URL.Path)
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

	err := client.DeleteProject("proj-1")
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
					ID:          "job-1",
					ProjectID:   "proj-1",
					Description: "Test Job",
					ExecutorID:  "exec-1",
					Data:        "job data",
					Spec:        "0 30 * * * *",
					StartDate:   "2025-01-01T00:00:00Z",
					EndDate:     "2025-12-31T00:00:00Z",
					Timezone:    "UTC",
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

	result, err := client.ListJobs()
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 1, result.Data.Total)
	assert.Equal(t, "Test Job", result.Data.Jobs[0].Description)
}

func TestCreateJob(t *testing.T) {
	mockResponse := JobResponse{
		Success: true,
		Data: Job{
			ID:          "job-1",
			ProjectID:   "proj-1",
			Description: "New Job",
			ExecutorID:  "exec-1",
			Data:        "job data",
			Spec:        "0 30 * * * *",
			StartDate:   "2025-01-01T00:00:00Z",
			EndDate:     "2025-12-31T00:00:00Z",
			Timezone:    "UTC",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/jobs", r.URL.Path)
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

	body := &JobRequestBody{
		Description: "New Job",
		Timezone:    "UTC",
		CallbackURL: "https://example.com/callback",
		StartDate:   "2025-01-01T00:00:00Z",
		EndDate:     "2025-12-31T00:00:00Z",
	}

	result, err := client.CreateJob(body)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "New Job", result.Data.Description)
}

func TestGetJob(t *testing.T) {
	mockResponse := JobResponse{
		Success: true,
		Data: Job{
			ID:          "job-1",
			ProjectID:   "proj-1",
			Description: "Get Job",
			ExecutorID:  "exec-1",
			Data:        "job data",
			Spec:        "0 30 * * * *",
			StartDate:   "2025-01-01T00:00:00Z",
			EndDate:     "2025-12-31T00:00:00Z",
			Timezone:    "UTC",
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
	assert.Equal(t, "Get Job", result.Data.Description)
}

func TestUpdateJob(t *testing.T) {
	mockResponse := JobResponse{
		Success: true,
		Data: Job{
			ID:          "job-1",
			ProjectID:   "proj-1",
			Description: "Updated Job",
			ExecutorID:  "exec-1",
			Data:        "job data",
			Spec:        "0 45 * * * *",
			StartDate:   "2025-01-01T00:00:00Z",
			EndDate:     "2025-12-31T00:00:00Z",
			Timezone:    "UTC",
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
		ProjectID:   "proj-1",
		Spec:        "0 45 * * * *",
		CallbackURL: "https://example.com/callback",
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

	err := client.DeleteJob("job-1")
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
