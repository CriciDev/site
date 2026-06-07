PRAGMA journal_mode = WAL;

CREATE TABLE IF NOT EXISTS sessions (
  id TEXT PRIMARY KEY,
  path TEXT NOT NULL,
  referrer TEXT,
  user_agent TEXT,
  started_at TEXT NOT NULL,
  last_seen_at TEXT NOT NULL,
  ended_at TEXT,
  duration_ms INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS events (
  id TEXT PRIMARY KEY,
  type TEXT NOT NULL,
  session_id TEXT NOT NULL,
  path TEXT NOT NULL,
  target TEXT,
  label TEXT,
  referrer TEXT,
  duration_ms INTEGER NOT NULL DEFAULT 0,
  metadata_json TEXT,
  created_at TEXT NOT NULL,
  FOREIGN KEY(session_id) REFERENCES sessions(id)
);

CREATE INDEX IF NOT EXISTS idx_events_session_id ON events(session_id);
CREATE INDEX IF NOT EXISTS idx_events_type ON events(type);
CREATE INDEX IF NOT EXISTS idx_sessions_started_at ON sessions(started_at);
