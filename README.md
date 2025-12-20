# Scheduler0 Go Client

A Go client library for interacting with the [Scheduler0 API](https://scheduler0.com/api). This client provides a convenient way to manage accounts, credentials, executions, executors, projects, jobs, features, create jobs from AI prompts, and monitor the health of your Scheduler0 cluster.

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

- **AI-Powered Job Creation**
  - Create job configurations from natural language prompts
  - AI generates cron expressions, scheduling, and job metadata
  - Supports purposes, events, recipients, and channels

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

- **OpenAPI Specification**: [openapi.json](https://api-reference.scheduler0.com) - Complete API specification

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
credentials, err := client.ListCredentials(scheduler0_go_client.ListCredentialsParams{
    Limit:            10,
    Offset:           0,
    OrderBy:          "date_created",
    OrderByDirection: "desc",
})

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
executions, err := client.ListExecutions(scheduler0_go_client.ListExecutionsParams{
    StartDate: "2024-01-01T00:00:00Z",  // Required: Start date (RFC3339 format)
    EndDate:   "2024-12-31T23:59:59Z",  // Required: End date (RFC3339 format)
    ProjectID: 0,                       // Optional: Project ID (0 for all)
    JobID:     0,                       // Optional: Job ID (0 for all)
    AccountID: 0,                       // Optional: Account ID override (0 uses client default)
    Limit:     10,                      // Required: Maximum number of items
    Offset:    0,                       // Required: Number of items to skip
})
```

### Managing Executors

```go
// List executors with pagination and ordering
executors, err := client.ListExecutors(scheduler0_go_client.ListExecutorsParams{
    AccountID:        0,                // Optional: Account ID override (0 uses client default)
    Limit:            10,
    Offset:           0,
    OrderBy:          "date_created",
    OrderByDirection: "desc",
})

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
// List projects with pagination and ordering
projects, err := client.ListProjects(scheduler0_go_client.ListProjectsParams{
    AccountID:        0,                // Optional: Account ID override (0 uses client default)
    Limit:            10,
    Offset:           0,
    OrderBy:          "date_created",    // Optional: Field to order by
    OrderByDirection: "desc",            // Optional: Order direction ("asc" or "desc")
})

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
jobs, err := client.ListJobs(scheduler0_go_client.ListJobsParams{
    ProjectID:        "",               // Optional: Project ID to filter by (empty string for all)
    Limit:            10,
    Offset:           0,
    OrderBy:          "date_created",
    OrderByDirection: "desc",
})

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

### AI-Powered Job Creation

Create job configurations from natural language prompts using AI:

```go
// Create job configurations from a natural language prompt
promptRequest := &scheduler0_go_client.PromptJobRequest{
    Prompt:     "Send weekly reports every Monday at 9 AM",
    Purposes:   []string{"reporting", "communication"},
    Events:     []string{"weekly_cycle"},
    Recipients: []string{"team@example.com", "manager@example.com"},
    Channels:   []string{"email"},
    Timezone:   "America/New_York",
}

// Generate job configurations from the prompt
// Note: This endpoint requires credits and validates credentials
jobConfigs, err := client.CreateJobFromPrompt(promptRequest)
if err != nil {
    log.Fatal(err)
}

// jobConfigs is an array of PromptJobResponse with generated configurations
for _, config := range jobConfigs {
    fmt.Printf("Kind: %s\n", config.Kind)
    fmt.Printf("Cron Expression: %s\n", config.CronExpression)
    if config.NextRunAt != nil {
        fmt.Printf("Next Run At: %s\n", *config.NextRunAt)
    }
    fmt.Printf("Recipients: %v\n", config.Recipients)
    
    // Use the generated configuration to create actual jobs
    job := &scheduler0_go_client.JobRequestBody{
        ProjectID:  123,
        Timezone:   config.Timezone,
        Spec:       config.CronExpression,
        CreatedBy:  "ai-prompt",
    }
    
    // Set optional fields if available
    if config.StartDate != nil {
        job.StartDate = *config.StartDate
    }
    if config.EndDate != nil {
        job.EndDate = *config.EndDate
    }
    if config.Subject != "" {
        job.Data = fmt.Sprintf(`{"subject": "%s", "recipients": %v}`, config.Subject, config.Recipients)
    }
    
    result, err := client.CreateJob(job)
    if err != nil {
        log.Printf("Failed to create job: %v", err)
        continue
    }
    
    fmt.Printf("Job created with request ID: %s\n", result.Data)
}
```

**Note**: The AI prompt endpoint requires:
- Valid API credentials (API Key + Secret)
- Account ID header
- Sufficient credits (1 credit per prompt execution)

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
- **Response Format**: Job creation returns `BatchJobResponse` with HTTP 202 Accepted status and a `Data` field containing the request ID (string) for async task tracking
- **Async Tracking**: Use the request ID with `GetAsyncTask()` to track job creation status

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
- `/api/v1/prompt` (AI prompt endpoint)

Account endpoints (`/api/v1/accounts/*`) and features (`/api/v1/features`) do not require account ID.

### Per-Request Account ID Override

You can override the Account ID set during client initialization on a per-request basis by including it in the params struct for list methods:

```go
// Override Account ID for a specific request
projects, err := client.ListProjects(scheduler0_go_client.ListProjectsParams{
    AccountID: 456,  // Overrides the client's default Account ID
    Limit:     10,
    Offset:    0,
})
```

For other methods, the Account ID can be set in the request body's `AccountID` field (which is excluded from JSON serialization but used for the `X-Account-ID` header).

## Credits and AI Features

The AI prompt endpoint (`/api/v1/prompt`) requires:
- **Credits**: 1 credit per prompt execution
- **Authentication**: Valid API Key + Secret credentials
- **Account ID**: Required header for credit deduction

Credits are automatically deducted when the prompt is successfully processed. If the prompt processing fails after credit deduction, credits are not refunded.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Development

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run tests with race detection
go test -v -race ./...

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### CI/CD

This project uses GitHub Actions for continuous integration. Tests are automatically run on:
- Push to `main`, `master`, or `develop` branches
- Pull requests to `main`, `master`, or `develop` branches

The CI pipeline tests against Go 1.23 and 1.24.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.