package scheduler0_go_client

// RaftStats represents the Raft cluster statistics
type RaftStats struct {
	AppliedIndex        string `json:"applied_index"`
	CommitIndex         string `json:"commit_index"`
	FSMPending          string `json:"fsm_pending"`
	LastContact         string `json:"last_contact"`
	LastLogIndex        string `json:"last_log_index"`
	LastLogTerm         string `json:"last_log_term"`
	LastSnapshotIndex   string `json:"last_snapshot_index"`
	LastSnapshotTerm    string `json:"last_snapshot_term"`
	LatestConfiguration string `json:"latest_configuration"`
	LatestConfigIndex   string `json:"latest_configuration_index"`
	NumPeers            string `json:"num_peers"`
	ProtocolVersion     string `json:"protocol_version"`
	ProtocolVersionMax  string `json:"protocol_version_max"`
	ProtocolVersionMin  string `json:"protocol_version_min"`
	SnapshotVersionMax  string `json:"snapshot_version_max"`
	SnapshotVersionMin  string `json:"snapshot_version_min"`
	State               string `json:"state"`
	Term                string `json:"term"`
}

// HealthcheckData represents the healthcheck response data
type HealthcheckData struct {
	LeaderAddress string    `json:"leaderAddress"`
	LeaderID      string    `json:"leaderId"`
	RaftStats     RaftStats `json:"raftStats"`
}

// HealthcheckResponse represents the healthcheck response
type HealthcheckResponse struct {
	Success bool            `json:"success"`
	Data    HealthcheckData `json:"data"`
}

