package scheduler0_go_client

// Send-time suggestion types for POST /ai/suggestions/time. The engine is
// deterministic scheduling/time-zone math (no LLM). All JSON keys are snake_case
// to match the API contract.

// SendTimeWorkingHours is a single continuous local working interval.
type SendTimeWorkingHours struct {
	Days  []string `json:"days,omitempty"`
	Start string   `json:"start,omitempty"`
	End   string   `json:"end,omitempty"`
}

// SendTimeQuietHours is a local "HH:MM" interval that may cross midnight.
type SendTimeQuietHours struct {
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
}

// SendTimeParticipant is a sender or recipient. Timezone is a required IANA id.
type SendTimeParticipant struct {
	ID           string                `json:"id,omitempty"`
	DisplayName  string                `json:"display_name,omitempty"`
	Timezone     string                `json:"timezone"`
	Role         string                `json:"role,omitempty"`
	WorkingHours *SendTimeWorkingHours `json:"working_hours,omitempty"`
	QuietHours   *SendTimeQuietHours   `json:"quiet_hours,omitempty"`
}

// SendTimeMessage carries timing-relevant metadata. Text is optional.
type SendTimeMessage struct {
	Channel            string `json:"channel,omitempty"`
	Priority           string `json:"priority,omitempty"`
	Intent             string `json:"intent,omitempty"`
	Text               string `json:"text,omitempty"`
	EstimatedAttention string `json:"estimated_attention,omitempty"`
}

// SendTimeConstraints are hard rules; violating candidates are rejected.
type SendTimeConstraints struct {
	EarliestSendAt      string `json:"earliest_send_at,omitempty"`
	LatestSendAt        string `json:"latest_send_at,omitempty"`
	MinimumDelaySeconds *int64 `json:"minimum_delay_seconds,omitempty"`
	WorkingHoursOnly    *bool  `json:"working_hours_only,omitempty"`
	AvoidWeekends       *bool  `json:"avoid_weekends,omitempty"`
	AvoidHolidays       *bool  `json:"avoid_holidays,omitempty"`
	RespectQuietHours   *bool  `json:"respect_quiet_hours,omitempty"`
	RequireCalendarFree *bool  `json:"require_calendar_free,omitempty"`
}

// SendTimeWindow is a local "HH:MM" preference window.
type SendTimeWindow struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// SendTimePreferences influence scoring but do not invalidate candidates.
type SendTimePreferences struct {
	PreferredRecipientWindows    []SendTimeWindow `json:"preferred_recipient_windows,omitempty"`
	AvoidRecipientWindows        []SendTimeWindow `json:"avoid_recipient_windows,omitempty"`
	PreferSenderRecipientOverlap *bool            `json:"prefer_sender_recipient_overlap,omitempty"`
}

// SendTimeGroupPolicy controls multi-recipient coverage.
type SendTimeGroupPolicy struct {
	Strategy                 string   `json:"strategy,omitempty"`
	MinimumRecipientCoverage *float64 `json:"minimum_recipient_coverage,omitempty"`
}

// SendTimeOptions tune candidate generation and the response shape.
type SendTimeOptions struct {
	ReferenceTime                   string `json:"reference_time,omitempty"`
	SuggestionCount                 *int   `json:"suggestion_count,omitempty"`
	CandidateIntervalMinutes        *int   `json:"candidate_interval_minutes,omitempty"`
	Locale                          string `json:"locale,omitempty"`
	IncludeScoreBreakdown           *bool  `json:"include_score_breakdown,omitempty"`
	IncludeRejectedSummary          *bool  `json:"include_rejected_summary,omitempty"`
	SearchHorizonDays               *int   `json:"search_horizon_days,omitempty"`
	EvaluateSendNow                 *bool  `json:"evaluate_send_now,omitempty"`
	DiversifySuggestions            *bool  `json:"diversify_suggestions,omitempty"`
	MinimumSuggestionSpacingMinutes *int   `json:"minimum_suggestion_spacing_minutes,omitempty"`
}

// SendTimeHolidayPolicy supplies holiday information. Only Dates are honored in
// the first release.
type SendTimeHolidayPolicy struct {
	Country    string   `json:"country,omitempty"`
	Region     string   `json:"region,omitempty"`
	CalendarID string   `json:"calendar_id,omitempty"`
	Dates      []string `json:"dates,omitempty"`
}

// SendTimeBusyInterval is an absolute (RFC3339) busy period.
type SendTimeBusyInterval struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// SendTimeAvailability is a participant's calendar busy intervals.
type SendTimeAvailability struct {
	ParticipantID string                 `json:"participant_id"`
	BusyIntervals []SendTimeBusyInterval `json:"busy_intervals,omitempty"`
}

// SendTimeSuggestionsRequest is the request body for POST /ai/suggestions/time.
type SendTimeSuggestionsRequest struct {
	Sender        *SendTimeParticipant   `json:"sender,omitempty"`
	Recipients    []SendTimeParticipant  `json:"recipients"`
	Message       *SendTimeMessage       `json:"message,omitempty"`
	Constraints   *SendTimeConstraints   `json:"constraints,omitempty"`
	Preferences   *SendTimePreferences   `json:"preferences,omitempty"`
	GroupPolicy   *SendTimeGroupPolicy   `json:"group_policy,omitempty"`
	Options       *SendTimeOptions       `json:"options,omitempty"`
	HolidayPolicy *SendTimeHolidayPolicy `json:"holiday_policy,omitempty"`
	Availability  []SendTimeAvailability `json:"availability,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// SendTimeSuggestionsResult mirrors the response. The top-level fields are typed;
// individual suggestions and the search/send_now blocks carry a rich, evolving
// shape owned by the engine, so they are exposed as generic maps.
type SendTimeSuggestionsResult struct {
	RequestID       string                   `json:"request_id,omitempty"`
	ReferenceTime   string                   `json:"reference_time,omitempty"`
	Policy          map[string]interface{}   `json:"policy,omitempty"`
	Engine          map[string]interface{}   `json:"engine,omitempty"`
	Suggestions     []map[string]interface{} `json:"suggestions"`
	Search          map[string]interface{}   `json:"search,omitempty"`
	RejectedSummary map[string]int           `json:"rejected_summary,omitempty"`
	NoSuggestion    map[string]interface{}   `json:"no_suggestion,omitempty"`
	SendNow         map[string]interface{}   `json:"send_now,omitempty"`
	Warnings        []map[string]interface{} `json:"warnings"`
	Metadata        map[string]interface{}   `json:"metadata,omitempty"`
}
