package server

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cricidev/site/internal/analytics"
	_ "modernc.org/sqlite"
)

func TestEventsHandlerAcceptsPageviewBeforeSessionStart(t *testing.T) {
	store, db := newTestStore(t)
	defer db.Close()

	srv := New(store)

	sessionID := "session-123"

	assertStatusCode(t, srv.Handler(), `{"type":"pageview","session_id":"`+sessionID+`","path":"/"}`, http.StatusAccepted)
	assertStatusCode(t, srv.Handler(), `{"type":"session_start","session_id":"`+sessionID+`","path":"/"}`, http.StatusAccepted)

	var events int
	if err := db.QueryRow(`SELECT count(*) FROM events WHERE session_id = ?`, sessionID).Scan(&events); err != nil {
		t.Fatalf("count events: %v", err)
	}
	if events != 2 {
		t.Fatalf("expected 2 events, got %d", events)
	}

	var sessions int
	if err := db.QueryRow(`SELECT count(*) FROM sessions WHERE id = ?`, sessionID).Scan(&sessions); err != nil {
		t.Fatalf("count sessions: %v", err)
	}
	if sessions != 1 {
		t.Fatalf("expected 1 session, got %d", sessions)
	}
}

func TestEventsHandlerClosesSessionAfterOutOfOrderEvents(t *testing.T) {
	store, db := newTestStore(t)
	defer db.Close()

	srv := New(store)

	sessionID := "session-456"

	assertStatusCode(t, srv.Handler(), `{"type":"pageview","session_id":"`+sessionID+`","path":"/"}`, http.StatusAccepted)
	assertStatusCode(t, srv.Handler(), `{"type":"session_end","session_id":"`+sessionID+`","path":"/","duration_ms":123}`, http.StatusAccepted)

	var duration int64
	var endedAt sql.NullString
	if err := db.QueryRow(`SELECT duration_ms, ended_at FROM sessions WHERE id = ?`, sessionID).Scan(&duration, &endedAt); err != nil {
		t.Fatalf("query session: %v", err)
	}

	if duration != 123 {
		t.Fatalf("expected duration 123, got %d", duration)
	}
	if !endedAt.Valid {
		t.Fatal("expected ended_at to be set")
	}
}

func TestHomePageDoesNotEmbedContribRocks(t *testing.T) {
	store, db := newTestStore(t)
	defer db.Close()

	srv := New(store)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if strings.Contains(body, "contrib.rocks/embed") {
		t.Fatal("home page still embeds contrib.rocks")
	}
	if !strings.Contains(body, `src="https://contrib.rocks/image?repo=CriciDev/site"`) {
		t.Fatal("home page is missing the contrib.rocks image widget")
	}
}

func newTestStore(t *testing.T) (*analytics.Store, *sql.DB) {
	t.Helper()

	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir root: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(cwd)
	})

	dbPath := filepath.Join(t.TempDir(), "analytics.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	db.SetMaxOpenConns(1)
	db.SetConnMaxLifetime(0)

	if err := db.PingContext(context.Background()); err != nil {
		t.Fatalf("ping db: %v", err)
	}

	store := analytics.NewStore(db)
	if err := store.Migrate(context.Background()); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	return store, db
}

func assertStatusCode(t *testing.T, handler http.Handler, body string, want int) {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/api/events", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != want {
		t.Fatalf("expected status %d, got %d; body=%s", want, rec.Code, rec.Body.String())
	}
}
