package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"cricidev/site/internal/analytics"
	home "cricidev/site/pages/home"
)

type Server struct {
	mux *http.ServeMux
}

func New(store *analytics.Store) *Server {
	mux := http.NewServeMux()

	mux.Handle("/", home.NewHandler())
	mux.Handle("/assets/", assetServer())
	mux.HandleFunc("/api/events", eventsHandler(store))
	mux.HandleFunc("/health", healthHandler(store))

	return &Server{mux: mux}
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func assetServer() http.Handler {
	return http.StripPrefix("/assets/", http.FileServer(http.Dir("assets")))
}

func healthHandler(store *analytics.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := store.Health(r.Context()); err != nil {
			http.Error(w, "unhealthy", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":    "ok",
			"timestamp": time.Now().UTC(),
		})
	}
}

func eventsHandler(store *analytics.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var payload analytics.Event
			decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10))
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(&payload); err != nil {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}
			if err := payload.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if err := store.RecordEvent(r.Context(), payload); err != nil {
				http.Error(w, "failed to store event", http.StatusInternalServerError)
				return
			}

			switch payload.Type {
			case analytics.EventSessionStart:
				if err := store.StartSession(r.Context(), analytics.Session{
					ID:        payload.SessionID,
					Path:      payload.Path,
					Referrer:  payload.Referrer,
					UserAgent: r.UserAgent(),
					StartedAt: payload.CreatedAt,
				}); err != nil {
					http.Error(w, "failed to store session", http.StatusInternalServerError)
					return
				}
			case analytics.EventSessionEnd:
				if err := store.EndSession(r.Context(), payload.SessionID, payload.DurationMS, payload.CreatedAt); err != nil && !errors.Is(err, sql.ErrNoRows) {
					http.Error(w, "failed to close session", http.StatusInternalServerError)
					return
				}
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)
			_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
		default:
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
