package analytics

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) RecordEvent(ctx context.Context, event Event) error {
	if err := event.Validate(); err != nil {
		return err
	}
	if event.ID == "" {
		event.ID = newID()
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now().UTC()
	}

	metadata := string(event.Metadata)
	if len(event.Metadata) == 0 {
		metadata = ""
	}

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO events (
			id, type, session_id, path, target, label, referrer, duration_ms, metadata_json, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		event.ID,
		event.Type,
		event.SessionID,
		event.Path,
		nullString(event.Target),
		nullString(event.Label),
		nullString(event.Referrer),
		event.DurationMS,
		nullString(metadata),
		event.CreatedAt.UTC().Format(time.RFC3339Nano),
	)
	return err
}

func (s *Store) StartSession(ctx context.Context, session Session) error {
	if session.ID == "" {
		return errors.New("session id is required")
	}
	if session.Path == "" {
		return errors.New("session path is required")
	}
	if session.StartedAt.IsZero() {
		session.StartedAt = time.Now().UTC()
	}

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO sessions (
			id, path, referrer, user_agent, started_at, last_seen_at
		) VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			path = excluded.path,
			referrer = excluded.referrer,
			user_agent = excluded.user_agent,
			last_seen_at = excluded.last_seen_at`,
		session.ID,
		session.Path,
		nullString(session.Referrer),
		nullString(session.UserAgent),
		session.StartedAt.UTC().Format(time.RFC3339Nano),
		session.StartedAt.UTC().Format(time.RFC3339Nano),
	)
	return err
}

func (s *Store) EndSession(ctx context.Context, sessionID string, durationMS int64, endedAt time.Time) error {
	if sessionID == "" {
		return errors.New("session id is required")
	}
	if endedAt.IsZero() {
		endedAt = time.Now().UTC()
	}

	res, err := s.db.ExecContext(ctx, `
		UPDATE sessions
		SET ended_at = ?, duration_ms = ?, last_seen_at = ?
		WHERE id = ?`,
		endedAt.UTC().Format(time.RFC3339Nano),
		durationMS,
		endedAt.UTC().Format(time.RFC3339Nano),
		sessionID,
	)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *Store) Health(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func nullString(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func newID() string {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return time.Now().UTC().Format("20060102T150405.000000000Z07:00")
	}
	return hex.EncodeToString(raw[:])
}
