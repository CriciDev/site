package analytics

import (
	"context"
	"os"
)

func (s *Store) Migrate(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx, `PRAGMA foreign_keys = ON;`); err != nil {
		return err
	}

	schema, err := os.ReadFile("db/schema.sql")
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, string(schema))
	return err
}
