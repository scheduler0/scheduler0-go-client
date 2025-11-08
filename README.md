# Scheduler0 Go Client

A Go client library for interacting with the [Scheduler0 API](https://scheduler0.com/api). This client provides a convenient way to manage accounts, credentials, executions, executors, projects, jobs, features, and monitor the health of your Scheduler0 cluster.

## Features

- **Account Management**
  - Create accounts
  - Get account details
  - Add/remove features from accounts

- **Feature Management**
  - List available features

- **Credentials Management**
  - List credentials with pagination and ordering
  - Create new credentials
  - Get credential details
  - Update credentials
  - Delete credentials

- **Executions Management**
  - List job executions with date filtering
  - Filter by project ID and job ID
  - View execution details and logs

- **Executors Management**
  - List executors with pagination and ordering
  - Create new executors (webhook, cloud function, container)
  - Get executor details
  - Update executors
  - Delete executors

- **Projects Management**
  - List projects with pagination
  - Create new projects
  - Get project details
  - Update projects
  - Delete projects

- **Jobs Management**
  - List jobs with pagination and ordering
  - Create new jobs with comprehensive scheduling options
  - Batch create multiple jobs in a single request
  - Get job details
  - Update jobs
  - Delete jobs

- **Async Tasks Management**
  - Get async task status by request ID

- **Health Monitoring**
  - Check cluster health
  - View raft statistics
  - Monitor leader status

## Installation

```bash
go get github.com/scheduler0/scheduler0-go-client
```

## API Documentation

- **OpenAPI Specification**: [openapi.json](https://scheduler0.com/api-reference) - Complete API specification

## Authentication

The Scheduler0 Go client supports multiple authentication methods:

### 1. API Key + Secret Authentication (Default)
Most endpoints require API Key and Secret authentication with an Account ID:

```go
client, err := scheduler0_go_client.NewAPIClientWithAccount(
    "http://localhost:7070",  // Base URL
    "v1",                     // API Version
    "your-api-key",           // API Key
    "your-api-secret",        // API Secret
    "123",                    // Account ID
)
```

### 2. Basic Authentication (Peer Communication)
For peer-to-peer communication:

```go
client, err := scheduler0_go_client.NewBasicAuthClient(
    "http://localhost:7070",  // Base URL
    "v1",                     // API Version
    "username",               // Username
    "password",               // Password
)
```

### 3. Options Pattern
For more flexibility, use the options pattern:

```go
client, err := scheduler0_go_client.NewClient(
    "http://localhost:7070",  // Base URL
    "v1",                     // API Version
    scheduler0_go_client.WithAPIKey("api-key", "api-secret"),
    scheduler0_go_client.WithAccountID("123"),
)
```

## Usage

### Managing Accounts

```go
// Create a new account
account := &scheduler0_go_client.AccountCreateRequestBody{
    Name: "My Account",
}
result, err := client.CreateAccount(account)

// Get account details
account, err := client.GetAccount("account-id")

// Add feature to account
feature := &scheduler0_go_client.FeatureRequest{
    FeatureID: 1,
}
result, err := client.AddFeatureToAccount("account-id", feature)

// Remove feature from account
err := client.RemoveFeatureFromAccount("account-id", feature)
```

### Managing Features

```go
// List all available features
features, err := client.ListFeatures()
```

### Managing Credentials

```go
// List credentials with pagination and ordering
credentials, err := client.ListCredentials(10, 0, "date_created", "desc")

// Create a new credential
credential, err := client.CreateCredential()

// Get a specific credential
credential, err := client.GetCredential("credential-id")

// Update a credential
credential, err := client.UpdateCredential("credential-id")

// Delete a credential
err := client.DeleteCredential("credential-id")
```

### Managing Executions

```go
// List executions with date filtering
executions, err := client.ListExecutions(
    "2024-01-01T00:00:00Z",  // Start date
    "2024-12-31T23:59:59Z",  // End date
    0,                       // Project ID (0 for all)
    0,                       // Job ID (0 for all)
    10,                      // Limit
    0,                       // Offset
)
```

### Managing Executors

```go
// List executors with pagination and ordering
executors, err := client.ListExecutors(10, 0, "date_created", "desc")

// Create a webhook executor
executor := &scheduler0_go_client.ExecutorRequestBody{
    Name:           "webhook-executor",
    Type:           "webhook_url",
    WebhookURL:     "https://example.com/webhook",
    WebhookMethod:  "POST",
    WebhookSecret:  "secret-key",
}

// Create a cloud function executor
executor := &scheduler0_go_client.ExecutorRequestBody{
    Name:             "cloud-function-executor",
    Type:             "cloud_function",
    Region:           "us-west-1",
    CloudProvider:    "aws",
    CloudResourceURL: "https://example.com/function",
    CloudAPIKey:      "api-key",
    CloudAPISecret:   "api-secret",
}

result, err := client.CreateExecutor(executor)

// Get a specific executor
executor, err := client.GetExecutor("executor-id")

// Update an executor
update := &scheduler0_go_client.ExecutorRequestBody{
    Name: "updated-executor",
    // ... other fields
}
result, err := client.UpdateExecutor("executor-id", update)

// Delete an executor
err := client.DeleteExecutor("executor-id")
```

### Managing Projects

```go
// List projects with pagination
projects, err := client.ListProjects(10, 0)

// Create a new project
project := &scheduler0_go_client.ProjectRequestBody{
    Name:        "My Project",
    Description: "Project description",
}
result, err := client.CreateProject(project)

// Get a specific project
project, err := client.GetProject("project-id")

// Update a project
update := &scheduler0_go_client.ProjectUpdateRequestBody{
    Description: "Updated description",
}
result, err := client.UpdateProject("project-id", update)

// Delete a project
err := client.DeleteProject("project-id")
```

### Managing Jobs

```go
// List jobs with pagination and ordering
jobs, err := client.ListJobs("project-id", 10, 0, "date_created", "desc")

// Create a single job
job := &scheduler0_go_client.JobRequestBody{
    ProjectID:     123,                    // Required
    Timezone:      "UTC",                  // Required
    ExecutorID:    &executorID,            // Optional
    Data:          "job payload data",     // Optional
    Spec:          "0 30 * * * *",         // Optional
    StartDate:     "2024-01-01T00:00:00Z", // Optional
    EndDate:       "2024-12-31T23:59:59Z", // Optional
    TimezoneOffset: 0,                     // Optional
    RetryMax:      3,                      // Optional
    Status:        "active",               // Optional
}
result, err := client.CreateJob(job)

// Create multiple jobs in a single batch request
jobs := []scheduler0_go_client.JobRequestBody{
    {
        ProjectID:     123,
        Timezone:      "UTC",
        Data:          "job 1 payload",
        Spec:          "0 30 * * * *",
        StartDate:     "2024-01-01T00:00:00Z",
        RetryMax:      3,
    },
    {
        ProjectID:     123,
        Timezone:      "UTC",
        Data:          "job 2 payload",
        Spec:          "0 0 * * * *",
        StartDate:     "2024-01-01T00:00:00Z",
        RetryMax:      5,
    },
}
batchResult, err := client.BatchCreateJobs(jobs)

// Get a specific job
job, err := client.GetJob("job-id")

// Update a job
update := &scheduler0_go_client.JobUpdateRequestBody{
    Data:   "updated payload",
    Spec:   "0 0 * * * *",
    Status: "inactive",
}
result, err := client.UpdateJob("job-id", update)

// Delete a job
err := client.DeleteJob("job-id")
```

### Managing Async Tasks

```go
// Get async task status
task, err := client.GetAsyncTask("request-id")
```

### Health Monitoring

```go
// Check cluster health (no authentication required)
health, err := client.Healthcheck()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Leader: %s\n", health.Data.LeaderAddress)
fmt.Printf("Raft State: %s\n", health.Data.RaftStats.State)
```

## Data Types

### Job Status
- `"active"` - Job is active and will be executed
- `"inactive"` - Job is inactive and will not be executed

### Executor Types
- `"webhook_url"` - HTTP webhook executor
- `"cloud_function"` - Cloud function executor
- `"container"` - Container executor

### Webhook Methods
- `"GET"`, `"POST"`, `"PUT"`, `"DELETE"`

### Job Creation Behavior
- **Single Job Creation**: `CreateJob()` internally uses batch creation with a single job
- **Batch Job Creation**: `BatchCreateJobs()` allows creating multiple jobs in one API call
- **Backend API**: The `/api/v1/jobs` POST endpoint expects an array of jobs for batch processing
- **Response Format**: Batch creation returns `BatchJobResponse` with `Data []Job` array

## Error Handling

The client returns Go errors for API errors. Check the error message for details:

```go
result, err := client.CreateJob(job)
if err != nil {
    if strings.Contains(err.Error(), "API error: 400") {
        // Handle bad request
    } else if strings.Contains(err.Error(), "API error: 401") {
        // Handle unauthorized
    } else if strings.Contains(err.Error(), "API error: 403") {
        // Handle forbidden
    } else if strings.Contains(err.Error(), "API error: 404") {
        // Handle not found
    }
    log.Fatal(err)
}
```

## Account ID Requirements

Most endpoints require the `X-Account-ID` header. The following endpoints require account ID:
- `/api/v1/jobs/*`
- `/api/v1/projects/*`
- `/api/v1/credentials/*`
- `/api/v1/executors/*`
- `/api/v1/async-tasks/*`
- `/api/v1/executions`

Account endpoints (`/api/v1/accounts/*`) and features (`/api/v1/features`) do not require account ID.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.