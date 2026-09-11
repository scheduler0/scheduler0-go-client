# Scheduler0 Go Client

<div align="center">
  <img src="https://raw.githubusercontent.com/scheduler0/scheduler0-go-client/main/logo.png" alt="Scheduler0 Logo" width="200">
</div>

The official Go client for the [Scheduler0](https://scheduler0.com) HTTP API (`/api/v1`). It is a thin, typed wrapper: one method per endpoint, request/response structs that mirror the server JSON, and no hidden behaviour (no retries, no caching, no background goroutines).

- API reference: <https://api-reference.scheduler0.com>
- Documentation: <https://docs.scheduler0.com>

## Installation

```bash
go get github.com/scheduler0/scheduler0-go-client/v2
```

The package name is `scheduler0_go_client`:

```go
import scheduler0_go_client "github.com/scheduler0/scheduler0-go-client/v2"
```

## Creating a client

`NewClient` takes the base URL, the API version segment (always `"v1"` today) and options. Requests are sent to `<baseURL>/api/<version>/<endpoint>`. There is no default base URL; the hosted API is `https://api.scheduler0.com`.

```go
client, err := scheduler0_go_client.NewClient(
    "https://api.scheduler0.com",
    "v1",
    scheduler0_go_client.WithAPIKey("your-api-key", "your-secret-key"),
    scheduler0_go_client.WithAccountID("123"),
)
```

Equivalent shorthands:

```go
// API key + secret + account ID (what almost every caller wants)
client, err := scheduler0_go_client.NewAPIClientWithAccount("https://api.scheduler0.com", "v1", "api-key", "secret-key", "123")

// API key + secret only; pass the account ID per call or via body AccountID fields
client, err = scheduler0_go_client.NewAPIClient("https://api.scheduler0.com", "v1", "api-key", "secret-key")

// Basic auth (self-hosted operator / peer path, sends X-Peer: cmd)
client, err = scheduler0_go_client.NewBasicAuthClient("http://127.0.0.1:9091", "v1", "admin", "admin")
```

### Authentication headers

With `WithAPIKey`, every request carries `X-API-Key` and `X-Secret-Key`. `X-Account-ID` is added when an account ID can be resolved, in this order of precedence:

1. an explicit per-call override (the trailing `accountIDOverride ...string` argument on methods that have one, or the `AccountID` field on `List*Params`),
2. a non-zero `AccountID` field on the request body struct (these fields are tagged `json:"-"` so they are used for the header only, not serialized),
3. the client default set with `WithAccountID`.

The server requires all three headers on every request except `GET /healthcheck`, including the `/accounts/*`, `/features` and `/cluster/*` routes. If `WithBasicAuth` is set it takes precedence over the API key and the request is sent with HTTP Basic auth plus `X-Peer: cmd`.

### HTTP client, timeouts, retries

`Client.HTTPClient` is a plain `&http.Client{}` with no timeout. Set your own:

```go
client.HTTPClient = &http.Client{Timeout: 30 * time.Second}
```

The client does not retry, back off or follow `Location` headers. Handle that in your application if you need it.

## Scopes

Each credential has a set of scopes. Calls made with a credential missing the required scope fail with `403`; expired credentials fail with `401`. `admin` satisfies every scope.

| Scope     | Grants                                                                                                                                                                            |
|-----------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `read`    | All `GET`s: jobs, projects, credentials, executors, executions (`/executions`, `/analytics`, `/totals`), async tasks, features, AI settings, AI models, AI prompt-request log, local-executor job pull |
| `write`   | Create/update/delete jobs, projects, credentials, executors; `PUT /ai/settings`; register local executors                                                                        |
| `execute` | `POST /ai/prompt`, `/ai/prompt/classify`, `/ai/schedule`, `/ai/suggestions/analyze`, `/ai/suggestions/time`, `/executions/cleanup-old-logs`, `/executors/{id}/test-invoke`, local-executor execution reports |
| `admin`   | `/accounts/*`, `/account/rotate-secret`, `/cluster/*` (also reachable with Basic auth when self-hosting)                                                                          |

Constants: `ScopeRead`, `ScopeWrite`, `ScopeExecute`, `ScopeAdmin`.

## Response envelope and errors

Every response is `{"success": bool, "data": ...}` and every response struct exposes it as `Success` and `Data`. `204` responses have no body and the corresponding methods return only `error`.

For any HTTP status `>= 400` the client returns `fmt.Errorf("API error: %s", body)` where `body` is the raw response (normally `{"success":false,"data":"<message>"}`). The status code is **not** included in the error string, and there are no typed errors except `*PromptSkippedError` (see [AI](#ai)). Inspect the message text if you need to branch:

```go
result, err := client.GetJob("42")
if err != nil {
    // err.Error() == `API error: {"success":false,"data":"job not found"}`
    log.Fatal(err)
}
```

A `2xx` response with `Success == false` is possible (see the credential-update note below), so check `Success` on responses you care about.

## Usage

### Projects

```go
// GET /projects
projects, err := client.ListProjects(scheduler0_go_client.ListProjectsParams{
    Limit:            10,             // optional; server default 10, max 100
    Offset:           0,
    OrderBy:          "date_created", // id | name | description | date_created | account_id
    OrderByDirection: "desc",         // asc | desc
})
for _, p := range projects.Data.Projects {
    fmt.Println(p.ID, p.Name)
}

// POST /projects -> 201. name, description and createdBy are required; name is unique per account.
created, err := client.CreateProject(&scheduler0_go_client.ProjectRequestBody{
    Name:        "billing",
    Description: "Invoice reminders",
    CreatedBy:   "user@example.com",
})

// GET /projects/{id}
project, err := client.GetProject(created.Data.ID)

// PUT /projects/{id}. Only description can change; name is immutable.
updated, err := client.UpdateProject(created.Data.ID, &scheduler0_go_client.ProjectUpdateRequestBody{
    Description: "Invoice reminders (EU)",
    ModifiedBy:  "user@example.com",
})

// DELETE /projects/{id} -> 204. Also deletes the project's jobs.
err = client.DeleteProject(created.Data.ID, &scheduler0_go_client.ProjectDeleteRequestBody{
    DeletedBy: "user@example.com",
})
```

### Jobs

`POST /jobs` always takes an **array** and returns `202 Accepted` with `Data` set to a request ID. The jobs are created asynchronously; call `GetAsyncTask(requestID)` to learn the outcome (`Output` holds the created jobs as JSON on success, or an error JSON on failure). `CreateJob` is a convenience wrapper that sends a one-element array.

```go
projectID := int64(12)
executorID := int64(7)

// POST /jobs -> 202. createdBy is required on every element.
resp, err := client.CreateJob(&scheduler0_go_client.JobRequestBody{
    ProjectID:  projectID,
    Timezone:   "UTC",
    ExecutorID: &executorID,
    Spec:       "0 30 * * * *",             // cron with seconds; empty spec = one-time job fired at StartDate
    Data:       `{"invoiceId":42}`,          // opaque payload delivered to the executor
    StartDate:  "2026-01-01T00:00:00Z",     // optional, RFC3339
    EndDate:    "2026-12-31T23:59:59Z",     // optional
    RetryMax:   3,                          // optional
    Status:     scheduler0_go_client.JobStatusActive, // active | inactive
    CreatedBy:  "user@example.com",
})
requestID := resp.Data

// Several jobs in one request
batch, err := client.BatchCreateJobs([]scheduler0_go_client.JobRequestBody{
    {ProjectID: projectID, Timezone: "UTC", Spec: "0 0 9 * * 1", CreatedBy: "user@example.com"},
    {ProjectID: projectID, Timezone: "UTC", Spec: "0 0 9 * * 5", CreatedBy: "user@example.com"},
})

// GET /async-tasks/{requestId}. Blocks server-side until the task finishes.
task, err := client.GetAsyncTask(requestID)
switch task.Data.State {
case scheduler0_go_client.AsyncTaskSuccess:
    fmt.Println("created:", task.Data.Output)
case scheduler0_go_client.AsyncTaskFailed:
    fmt.Println("failed:", task.Data.Output)
}

// GET /jobs
jobs, err := client.ListJobs(scheduler0_go_client.ListJobsParams{
    ProjectID:        "12",           // optional filter; "" = all projects
    Limit:            10,
    Offset:           0,
    OrderBy:          "date_created", // id | project_id | spec | date_created | timezone | account_id | date_modified | modified_by | deleted_by | executor_id | start_date | end_date | retry_max
    OrderByDirection: "desc",
})

// GET /jobs/{id}
job, err := client.GetJob("42")

// PUT /jobs/{id}. modifiedBy is required. timezone, timezoneOffset and executorId keep
// their current values when omitted.
updatedJob, err := client.UpdateJob("42", &scheduler0_go_client.JobUpdateRequestBody{
    Spec:       "0 0 * * * *",
    Status:     scheduler0_go_client.JobStatusInactive,
    ModifiedBy: "user@example.com",
})

// DELETE /jobs/{id} -> 204
err = client.DeleteJob("42", &scheduler0_go_client.JobDeleteRequestBody{DeletedBy: "user@example.com"})
```

`GetJob`, `UpdateJob`, `DeleteJob`, `CreateJob` and `BatchCreateJobs` accept an optional trailing account ID override, e.g. `client.GetJob("42", "456")`.

### Executors

Executor types (`ExecutorTypeCloudFunction`, `ExecutorTypeWebhookURL`, `ExecutorTypeLocal`):

| Type             | Required fields                                                     |
|------------------|---------------------------------------------------------------------|
| `webhook_url`    | `WebhookURL`, `WebhookMethod` (`GET`/`POST`/`PUT`/`DELETE`)          |
| `cloud_function` | `CloudResourceURL` (plus provider-specific `CloudProvider`, `Region`, `CloudAPIKey`, `CloudAPISecret`) |
| `local`          | `Command` (see [Local executors](#local-executors))                  |

`CloudAPIKey`, `CloudAPISecret` and `WebhookSecret` are returned **only** in the `CreateExecutor` response; they are empty on get, list and update.

```go
// GET /executors
executors, err := client.ListExecutors(scheduler0_go_client.ListExecutorsParams{
    Limit:            10,
    Offset:           0,
    OrderBy:          "date_created", // id | date_created | date_modified | created_by | modified_by | deleted_by
    OrderByDirection: "desc",
})

// POST /executors -> 201
webhook, err := client.CreateExecutor(&scheduler0_go_client.ExecutorRequestBody{
    Name:          "invoice-webhook",
    Description:   "Sends invoice reminder emails", // used by /ai/schedule to match executors
    Tags:          []string{"email", "billing"},
    Type:          scheduler0_go_client.ExecutorTypeWebhookURL,
    WebhookURL:    "https://example.com/hooks/invoice",
    WebhookMethod: "POST",
    WebhookSecret: "shared-secret",
    CreatedBy:     "user@example.com",
})

// GET /executors/{id}
executor, err := client.GetExecutor("7")

// PUT /executors/{id}. modifiedBy is required.
updatedExecutor, err := client.UpdateExecutor("7", &scheduler0_go_client.ExecutorUpdateRequestBody{
    Name:          "invoice-webhook",
    Type:          scheduler0_go_client.ExecutorTypeWebhookURL,
    WebhookURL:    "https://example.com/hooks/invoice-v2",
    WebhookMethod: "POST",
    ModifiedBy:    "user@example.com",
})

// DELETE /executors/{id} -> 204
err = client.DeleteExecutor("7", &scheduler0_go_client.ExecutorDeleteRequestBody{DeletedBy: "user@example.com"})

// POST /executors/{id}/test-invoke. Fires a synthetic job through the executor right now.
// No job, execution log or schedule is created. Body is optional (pass nil). Local
// executors cannot be test-invoked (400). HTTP 200 is returned even when the target
// fails; check Data.Success / Data.Error.
test, err := client.TestInvokeExecutor("7", &scheduler0_go_client.TestInvocationRequestBody{
    Job:           &scheduler0_go_client.Job{Spec: "0 0 2 * * *", Data: `{"dryRun":true}`, Timezone: "UTC"},
    Age:           "24h",                  // optional Go duration; how old the synthetic job looks
    ExecutionTime: "2026-01-15T02:00:00Z", // optional; defaults to now
})
fmt.Println(test.Data.Success, test.Data.DurationMs, test.Data.Error)
```

### Local executors

Local executors run a shell command on a machine you control. The `scheduler0` CLI registers one, pulls its jobs every minute and reports results; these are the three endpoints it uses.

```go
// POST /local-executors -> 201 {id}. name, command and createdBy are required.
reg, err := client.RegisterLocalExecutor(&scheduler0_go_client.LocalExecutorRegisterRequest{
    Name:       "build-box",
    Command:    "/usr/local/bin/run-job.sh",
    WorkingDir: "/srv/app",
    CreatedBy:  "user@example.com",
})
executorID := reg.Data.ID

// GET /local-executors/{id}/jobs -> active jobs assigned to this executor (also renews its lease)
pulled, err := client.PullLocalExecutorJobs(executorID)

// POST /local-executors/{id}/executions -> {committed}
report, err := client.ReportLocalExecutions(executorID, []scheduler0_go_client.LocalExecutionReport{
    {
        JobID:             pulled.Data[0].ID,
        UniqueID:          "exec-2026-01-01T00:00:00Z-1",
        State:             scheduler0_go_client.ExecutionStateSuccess, // 0 scheduled, 1 success, 2 failed
        LastExecutionTime: "2026-01-01T00:00:00Z",
        NextExecutionTime: "2026-01-02T00:00:00Z",
    },
})
fmt.Println(report.Data.Committed)
```

### Executions

```go
// GET /executions. Every parameter is optional.
executions, err := client.ListExecutions(scheduler0_go_client.ListExecutionsParams{
    StartDate:      "2026-01-01T00:00:00Z", // RFC3339
    EndDate:        "2026-01-31T23:59:59Z",
    ProjectID:      12,
    JobID:          42,
    State:          "failed",               // scheduled | success | failed
    OrderBy:        "dateCreated",          // dateCreated | lastExecutionDateTime | nextExecutionDateTime
    OrderDirection: "DESC",                 // ASC | DESC (note: not orderByDirection)
    Limit:          50,                     // server default 50
    Offset:         0,
})
for _, e := range executions.Data.Executions {
    fmt.Println(e.JobID, e.State, e.LastExecutionDatetime) // State: 0 scheduled, 1 success, 2 failed
}

// GET /executions/analytics?startDate=YYYY-MM-DD&startTime=HH:MM[:SS] -> per-minute buckets
analytics, err := client.GetDateRangeAnalytics(scheduler0_go_client.GetDateRangeAnalyticsParams{
    StartDate: "2026-01-01",
    StartTime: "00:00",
})
for _, p := range analytics.Data.Points {
    fmt.Println(p.Date, p.Time, p.Scheduled, p.Success, p.Failed)
}

// GET /executions/totals -> lifetime scheduled/success/failed counts. 0 uses the client's account ID.
totals, err := client.GetExecutionTotals(0)

// POST /executions/cleanup-old-logs (execute scope). accountId must equal X-Account-ID.
cleanup, err := client.CleanupOldExecutionLogs("123", 6) // delete logs older than 6 months
fmt.Println(cleanup.Data.Message)
```

### Credentials

```go
// GET /credentials
credentials, err := client.ListCredentials(scheduler0_go_client.ListCredentialsParams{
    Limit:            10,
    Offset:           0,
    OrderBy:          "date_created", // id | date_created | date_modified | created_by | modified_by | deleted_by | expires_at
    OrderByDirection: "desc",
})

// POST /credentials -> 201. Scopes is required (non-empty, no duplicates). Granting "admin"
// needs an admin credential or Basic auth. Expiry defaults to 90 days; ExpiresInSeconds
// can shorten it (the server clamps the value).
ttl := int64(8 * 60 * 60)
cred, err := client.CreateCredential(&scheduler0_go_client.CredentialCreateRequestBody{
    CreatedBy:        "user@example.com",
    Scopes:           []string{scheduler0_go_client.ScopeRead, scheduler0_go_client.ScopeWrite},
    ExpiresInSeconds: &ttl, // optional
})
// Data.PlaintextSecret is returned ONLY here. Store it now; it cannot be fetched again.
fmt.Println(cred.Data.APIKey, cred.Data.PlaintextSecret, cred.Data.Scopes)

// GET /credentials/{id} (no secret)
one, err := client.GetCredential("9")

// POST /credentials/{id}/archive -> 204. Disables the credential.
err = client.ArchiveCredential("9", "user@example.com")

// DELETE /credentials/{id} -> 204
err = client.DeleteCredential("9", &scheduler0_go_client.CredentialDeleteRequestBody{DeletedBy: "user@example.com"})
```

There is no rotate endpoint for credentials: create a new credential, switch your application over, then archive or delete the old one.

`UpdateCredential` (`PUT /credentials/{id}`) changes `Archived` and `ModifiedBy` only. `apiKey`, `apiSecret`, `scopes` and `expiresAt` are fixed at creation; the server rejects attempts to change the key or secret with `400`. `Archived` is `omitempty`, and the server treats an omitted `archived` as `false`, so `Archived: false` un-archives. Servers older than the credential-update fix answer every call with `200` and `success:false "api_key or api_secret cannot be empty"`; check `Success` if you target one.

### AI

All AI endpoints need the `execute` scope except the `GET`s (`read`). `POST /ai/prompt` and `POST /ai/schedule` count against the account's monthly prompt quota (`429` when exhausted) and, when using platform-hosted models, its AI credit balance (`402` when exhausted). `/ai/prompt/classify`, `/ai/suggestions/analyze` and `/ai/suggestions/time` do not invoke a model or consume credits.

#### Generate job configurations from a prompt

`POST /ai/prompt` runs the intent guardrail and then the configured model(s). If the guardrail rejects the prompt (or asks for clarification) the server returns `422` and the client returns a `*PromptSkippedError`.

```go
result, err := client.CreateJobFromPrompt(&scheduler0_go_client.PromptJobRequest{
    Prompt:     "Send the weekly sales report every Monday at 9am",
    Purposes:   []string{"reporting"},          // optional hints
    Recipients: []string{"sales@example.com"},
    Channels:   []string{"email"},
    Timezone:   "America/New_York",             // optional IANA zone; invalid -> 400. Defaults to UTC.
    Locale:     "en",                           // optional; the guardrail only runs for en* locales
})
if err != nil {
    var skipped *scheduler0_go_client.PromptSkippedError
    if errors.As(err, &skipped) { // or scheduler0_go_client.IsPromptSkippedError(err)
        fmt.Println("rejected:", skipped.Message, skipped.Classification.Decision, skipped.Classification.Reason)
        return
    }
    log.Fatal(err)
}
for _, p := range result.Providers {
    fmt.Println(p.Provider, p.Model, p.TotalTokens, p.DurationMs)
    for _, j := range p.Jobs {
        fmt.Println(j.Kind, j.CronExpression, j.Timezone, j.Subject) // Kind: FOLLOW_UP | REMINDER | DIGEST
    }
}
```

The result is only a suggestion; create the jobs yourself with `CreateJob`/`BatchCreateJobs`, or use `ScheduleFromPrompt` to have the server do it.

#### Schedule jobs from a prompt in one call

`POST /ai/schedule` -> `201`. Runs the same pipeline, resolves or creates a project, picks an executor (pinned `ExecutorID`, the account's only executor, or the best `Description`/`Tags` match) and creates the jobs synchronously. `409` when no executor/jobs can be resolved, `422` when the guardrail rejects the prompt (nothing is created).

```go
sched, err := client.ScheduleFromPrompt(&scheduler0_go_client.SchedulePromptRequest{
    Prompt:    "Remind the sales team every Monday at 9am to review the pipeline",
    Channels:  []string{"email"},
    CreatedBy: "user@example.com", // required
    // ProjectID: &projectID,                                                // reuse a project, or
    // Project:   &scheduler0_go_client.ScheduleProjectInput{Name: "Sales"}, // create-or-reuse by name
    // ExecutorID: &executorID,                                              // pin an executor
})
if err != nil {
    if scheduler0_go_client.IsPromptSkippedError(err) {
        log.Printf("rejected by guardrail: %v", err)
        return
    }
    log.Fatal(err)
}
fmt.Printf("project %d (created=%v), executor %d via %s, %d jobs\n",
    sched.Project.ID, sched.ProjectCreated, sched.Executor.ID, sched.ExecutorMatchedBy, len(sched.Jobs))
```

#### Classify a prompt only

`POST /ai/prompt/classify`. English (`en*`) only; other locales return `400`. `503` when the classifier is not configured.

```go
clf, err := client.ClassifyPrompt(&scheduler0_go_client.ClassifyPromptRequest{Prompt: "What is Kubernetes?"})
fmt.Println(clf.Decision, clf.Reason) // Decision: allow | clarify | reject
```

#### Analyze a conversation for follow-ups

`POST /ai/suggestions/analyze`. Deterministic, English only (`400 UNSUPPORTED_LOCALE` otherwise). Request and response use snake_case JSON; suggestions/obligations are returned as generic maps because their shape is owned by the analyzer.

```go
analysis, err := client.AnalyzeSuggestions(&scheduler0_go_client.AnalyzeSuggestionsRequest{
    ConversationID: "conv_123",
    Messages: []scheduler0_go_client.SuggestionMessage{
        {Speaker: "Victor", Timestamp: "2026-07-17T10:00:00-04:00", Message: "I'll send the proposal tomorrow."},
    },
    Options: &scheduler0_go_client.SuggestionOptions{Locale: "en", DefaultTimezone: "America/Toronto"},
})
for _, s := range analysis.Suggestions {
    fmt.Println(s["type"], s["reason"])
}
```

#### Recommend send times

`POST /ai/suggestions/time`. Deterministic time-zone math; nothing is sent or scheduled. `recipients[].timezone` is required. Validation errors come back as `400` with `{code, message, field}` in `data`.

```go
times, err := client.SendTimeSuggestions(&scheduler0_go_client.SendTimeSuggestionsRequest{
    Sender: &scheduler0_go_client.SendTimeParticipant{ID: "user_123", Timezone: "America/Toronto"},
    Recipients: []scheduler0_go_client.SendTimeParticipant{
        {ID: "user_456", Timezone: "America/Los_Angeles", Role: "primary"},
    },
    Message: &scheduler0_go_client.SendTimeMessage{Priority: "normal"},
})
for _, s := range times.Suggestions {
    fmt.Println(s["send_at"], s["score"], s["label"])
}
```

#### Prompt-request log and model catalog

```go
// GET /ai/prompt-requests (read). Fields are snake_case. limit defaults to 25, max 100.
logPage, err := client.ListPromptRequests(scheduler0_go_client.ListPromptRequestsParams{
    Provider:  "anthropic",            // optional
    Status:    "success",              // optional: success | failed | skipped_intent ...
    StartDate: "2026-01-01T00:00:00Z", // optional RFC3339 (query param "start")
    Order:     "DESC",                 // ASC | DESC
    Limit:     20,
})
for _, r := range logPage.Data.Requests {
    fmt.Println(r.DateCreated, r.Provider, r.Model, r.Status, r.TotalTokens, r.EstimatedCostUSD)
}

// GET /ai/models (read) -> map[provider][]ModelInfo. Only these models are accepted by PUT /ai/settings.
models, err := client.GetAIModels()
for provider, list := range models.Data {
    for _, m := range list {
        fmt.Println(provider, m.ID, m.DisplayName, m.Default)
    }
}
```

#### AI settings (bring your own keys)

`GET`/`PUT /ai/settings` operate on the account in `X-Account-ID`; the `accountID` argument is sent as that header. Keys are masked as `"•"` on read. `ActiveModels` is ordered (primary first, then fallbacks); each provider listed must have a key stored or supplied in the same request.

```go
settings, err := client.GetAccountAISettings("123")
fmt.Println(settings.Data.ActiveModels, settings.Data.AnthropicAPIKey) // key is "•" when set

saved, err := client.UpsertAccountAISettings("123", &scheduler0_go_client.AccountAISettings{
    ActiveModels: []scheduler0_go_client.ActiveModel{
        {Provider: "anthropic", Model: "claude-sonnet-4-5"},
        {Provider: "openai", Model: "gpt-4o"},
    },
    AnthropicAPIKey: "sk-ant-...",
    OpenAIAPIKey:    "sk-...",
})
fmt.Println(saved.Data.Message) // the PUT returns a message, not the saved settings
```

### Features

```go
// GET /features (read)
features, err := client.ListFeatures()
for _, f := range features.Data {
    fmt.Println(f.ID, f.Name)
}
```

### Accounts (self-hosting / admin scope)

These routes require an `admin` credential or Basic auth. For API-key callers the path `{id}` must equal `X-Account-ID` (`403` otherwise); the `accountID` argument is used for both.

```go
account, err := client.CreateAccount(&scheduler0_go_client.AccountCreateRequestBody{Name: "Acme"})          // POST /accounts -> 201
got, err := client.GetAccount("123")                                                                         // GET /accounts/{id}
renamed, err := client.UpdateAccount("123", &scheduler0_go_client.AccountUpdateRequestBody{Name: "Acme Inc"}) // PUT /accounts/{id}

// Feature flags
added, err := client.AddFeatureToAccount("123", &scheduler0_go_client.FeatureRequest{FeatureID: 1}) // PUT /accounts/{id}/feature -> 201
err = client.RemoveFeatureFromAccount("123", &scheduler0_go_client.FeatureRequest{FeatureID: 1})  // DELETE /accounts/{id}/feature -> 204
err = client.AddAllFeaturesToAccount("123")                                                        // PUT /accounts/{id}/features/all
err = client.RemoveAllFeaturesFromAccount("123")                                                   // DELETE /accounts/{id}/features/all

// Monthly execution counter
count, err := client.GetAccountExecutionCount("123")                 // GET /accounts/{id}/execution-count
bumped, err := client.IncreaseAccountExecutionCount("123", 1000)     // PUT /accounts/{id}/execution-count {count} -> {newExecutionCount}

// AI usage for the current period: prompt/classify limit-used-remaining and estimated cost
usage, err := client.GetAIUsage("123")                                // GET /accounts/{id}/ai/usage
fmt.Println(usage.Data.Prompt.Used, usage.Data.Prompt.Limit, usage.Data.EstimatedCostUSD)

// Platform token balance
tokens, err := client.GetAccountTokens("123")                         // GET /accounts/{id}/tokens -> {tokens}
balance, err := client.AddAccountTokens("123", 500)                   // PUT /accounts/{id}/tokens/add {amount} -> {newBalance}

// POST /account/rotate-secret: re-encrypt stored secrets after changing the server SecretKey.
// Update and reload the server's SecretKey first, then pass the previous key.
rotated, err := client.RotateSecret("<old-hex-secret-key>")
fmt.Println(rotated.Data.CredentialsRotated, rotated.Data.ExecutorsRotated, rotated.Data.AISettingsRotated)
```

### Cluster (self-hosting / admin scope or Basic auth + `X-Peer`)

```go
nodes, err := client.ListNodes()                                   // GET /cluster/list-nodes -> []Node{NodeId, NodeAddress, ClientAddress}
status, err := client.AddNode(2, "10.0.0.2:7071", "http://10.0.0.2:9091") // POST /cluster/add-node?nodeId=&nodeAddress=&clientAddress=
status, err = client.RemoveNode(2)                                 // POST /cluster/remove-node?nodeId=
status, err = client.PromoteNode(2)                                // POST /cluster/promote-node?nodeId=
status, err = client.DemoteNode(2)                                 // POST /cluster/demote-node?nodeId=
status, err = client.TransferLeadership()                          // POST /cluster/transfer-leadership
status, err = client.ForceRebuildCluster(1)                        // POST /cluster/force-rebuild?seedNodeId=
status, err = client.AddSelfToCluster()                            // POST /cluster/add-self
status, err = client.RemoveSelfFromCluster()                       // POST /cluster/remove-self
status, err = client.ResetRaftState()                              // POST /cluster/reset-raft (node exits afterwards)
fmt.Println(status.Data["status"])

// Debug dumps; Data is json.RawMessage because the shape is internal
queue, err := client.DumpScheduleQueue()                           // GET /cluster/dump/schedule-queue
cache, err := client.DumpJobExecutionsCache()                      // GET /cluster/dump/job-executions-cache
queues, err := client.DumpJobQueues()                              // GET /cluster/dump/job-queues
versions, err := client.DumpJobQueueVersions()                     // GET /cluster/dump/job-queue-versions

// Backup / restore -> 202 {status, requestId}
backup, err := client.BackupDatabase()                             // POST /cluster/backup
restore, err := client.RestoreDatabase("backup-2026-01-01.db")     // POST /cluster/restore {filePath}
fmt.Println(backup.Data["requestId"], restore.Data["status"])
```

### Health

```go
// GET /healthcheck. Sent without any auth headers.
health, err := client.Healthcheck()
fmt.Println(health.Data.LeaderAddress, health.Data.LeaderID, health.Data.RaftStats.State)
```

## Reference: enums

| Concept              | Values                                                                                   |
|----------------------|------------------------------------------------------------------------------------------|
| Credential scope     | `read`, `write`, `execute`, `admin` (`Scope*` constants)                                  |
| Executor type        | `cloud_function`, `webhook_url`, `local` (`ExecutorType*` constants)                     |
| Webhook method       | `GET`, `POST`, `PUT`, `DELETE`                                                            |
| Job status           | `active`, `inactive` (`JobStatus*` constants)                                             |
| Execution state      | `0` scheduled, `1` success, `2` failed (`ExecutionState*` constants); as a `/executions` filter: `scheduled`, `success`, `failed` |
| Async task state     | `0` not started, `1` in progress, `2` success, `3` failed (`AsyncTask*` constants)        |
| Intent decision      | `allow`, `clarify`, `reject`                                                              |
| Prompt job kind      | `FOLLOW_UP`, `REMINDER`, `DIGEST`                                                         |
| Executor matched by  | `pinned`, `only`, `llm` (`ScheduleResult.ExecutorMatchedBy`)                              |

## Development

```bash
go build ./...
go vet ./...
go test ./...
```

## License

MIT. See [LICENSE](LICENSE).
