package analytics

import (
	"encoding/json"
	"errors"
	"strings"
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

	switch e.Type {
	case EventPageview, EventClick, EventSessionStart, EventSessionEnd, EventHeartbeat:
	default:
		return errors.New("invalid event type")
	}

	if strings.TrimSpace(e.SessionID) == "" {
		return errors.New("session id is required")
	}
	if len(e.SessionID) > 128 {
		return errors.New("session id is too long")
	}

	if strings.TrimSpace(e.Path) == "" {
		return errors.New("path is required")
	}
	if len(e.Path) > 2048 {
		return errors.New("path is too long")
	}

	if e.DurationMS < 0 {
		return errors.New("duration must be positive")
	}
	if e.DurationMS > int64(24*time.Hour/time.Millisecond) {
		return errors.New("duration is too long")
	}

	if len(e.Target) > 256 {
		return errors.New("target is too long")
	}
	if len(e.Label) > 256 {
		return errors.New("label is too long")
	}
	if len(e.Referrer) > 2048 {
		return errors.New("referrer is too long")
	}
	if len(e.Metadata) > 8192 {
		return errors.New("metadata is too large")
	}

	return nil
}
