package scheduler0_go_client

import "encoding/json"

// ClusterStatusResponse wraps the standard API envelope for cluster node-lifecycle
// operations (add/remove/promote/demote/rebuild/reset/transfer-leadership), which
// acknowledge with a short status message rather than returning a resource.
type ClusterStatusResponse struct {
	Success bool              `json:"success"`
	Data    map[string]string `json:"data"`
}

// Node describes a member of the Raft cluster, as returned by GET /cluster/list-nodes.
type Node struct {
	ClientAddress string `json:"clientAddress"`
	NodeAddress   string `json:"nodeAddress"`
	NodeId        uint64 `json:"nodeId"`
}

// ListNodesResponse wraps the standard API envelope for GET /cluster/list-nodes.
type ListNodesResponse struct {
	Success bool   `json:"success"`
	Data    []Node `json:"data"`
}

// ClusterDumpResponse wraps the standard API envelope for the /cluster/dump/* debug
// endpoints. Data is left as raw JSON: these expose internal server state (schedule
// queue, executions cache, job queues/versions) whose shape is an implementation
// detail, not a stable public contract. Unmarshal Data yourself into whatever shape
// you need for inspection/debugging.
type ClusterDumpResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
}
