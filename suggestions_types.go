package scheduler0_go_client

// SuggestionParticipant is a member of a conversation.
type SuggestionParticipant struct {
	ID          string `json:"id,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Timezone    string `json:"timezone,omitempty"`
}

// SuggestionMessage is a single conversation message. Speaker accepts either a
// SuggestionParticipant (full shape) or a plain string display name (minimal shape).
type SuggestionMessage struct {
	ID        string      `json:"id,omitempty"`
	Speaker   interface{} `json:"speaker"`
	Timestamp string      `json:"timestamp"`
	Message   string      `json:"message"`
}

// SuggestionOptions carries per-request analysis options.
//
// Locale is English only for the first release: it defaults to "en" and only
// accepts en* values. Any other locale is rejected by the API with
// UNSUPPORTED_LOCALE.
type SuggestionOptions struct {
	ReferenceTime              string   `json:"reference_time,omitempty"`
	Locale                     string   `json:"locale,omitempty"`
	DefaultTimezone            string   `json:"default_timezone,omitempty"`
	MinimumConfidence          *float64 `json:"minimum_confidence,omitempty"`
	IncludeLowConfidence       *bool    `json:"include_low_confidence,omitempty"`
	IncludeResolvedObligations *bool    `json:"include_resolved_obligations,omitempty"`
	DefaultDueTime             string   `json:"default_due_time,omitempty"`
	DefaultDeadlineTime        string   `json:"default_deadline_time,omitempty"`
}

// AnalyzeSuggestionsRequest is the request body for POST /suggestions/analyze.
type AnalyzeSuggestionsRequest struct {
	ConversationID string                  `json:"conversation_id,omitempty"`
	Messages       []SuggestionMessage     `json:"messages"`
	Participants   []SuggestionParticipant `json:"participants,omitempty"`
	Options        *SuggestionOptions      `json:"options,omitempty"`
}

// AnalyzeSuggestionsResult mirrors the analyzer's response. Individual
// suggestions and obligations carry a rich, evolving shape owned by the edge
// analyzer, so they are exposed as generic maps.
type AnalyzeSuggestionsResult struct {
	RequestID      string                   `json:"request_id,omitempty"`
	ConversationID string                   `json:"conversation_id,omitempty"`
	AnalyzedAt     string                   `json:"analyzed_at,omitempty"`
	Suggestions    []map[string]interface{} `json:"suggestions"`
	Obligations    []map[string]interface{} `json:"obligations"`
	Warnings       []map[string]interface{} `json:"warnings"`
	Engine         map[string]interface{}   `json:"engine,omitempty"`
}
