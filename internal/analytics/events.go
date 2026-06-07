package analytics

import (
	"encoding/json"
	"errors"
	"time"
)

type EventType string

const (
	EventPageview     EventType = "pageview"
	EventClick        EventType = "click"
	EventSessionStart EventType = "session_start"
	EventSessionEnd   EventType = "session_end"
	EventHeartbeat    EventType = "heartbeat"
)

type Event struct {
	ID         string          `json:"id"`
	Type       EventType       `json:"type"`
	SessionID  string          `json:"session_id"`
	Path       string          `json:"path"`
	Target     string          `json:"target,omitempty"`
	Label      string          `json:"label,omitempty"`
	Referrer   string          `json:"referrer,omitempty"`
	DurationMS int64           `json:"duration_ms,omitempty"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`
	CreatedAt  time.Time       `json:"created_at,omitempty"`
}

type Session struct {
	ID         string    `json:"id"`
	Path       string    `json:"path"`
	Referrer   string    `json:"referrer,omitempty"`
	UserAgent  string    `json:"user_agent,omitempty"`
	StartedAt  time.Time `json:"started_at,omitempty"`
	EndedAt    time.Time `json:"ended_at,omitempty"`
	DurationMS int64     `json:"duration_ms,omitempty"`
}

func (e Event) Validate() error {
	if e.Type == "" {
		return errors.New("event type is required")
	}
	if e.SessionID == "" {
		return errors.New("session id is required")
	}
	if e.Path == "" {
		return errors.New("path is required")
	}
	return nil
}
