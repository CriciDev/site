package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"cricidev/site/internal/analytics"
	"cricidev/site/internal/server"
	_ "modernc.org/sqlite"
)

func main() {
	addr := env("ADDR", ":8080")
	dbPath := env("DB_PATH", "/data/cricidev.db")

	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(1)
	db.SetConnMaxLifetime(0)

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	store := analytics.NewStore(db)
	if err := store.Migrate(context.Background()); err != nil {
		log.Fatal(err)
	}

	srv := server.New(store)
	httpServer := &http.Server{
		Addr:              addr,
		Handler:           srv.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("cricidev site listening on %s", addr)
	log.Fatal(httpServer.ListenAndServe())
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
