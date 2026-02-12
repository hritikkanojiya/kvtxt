package storage

import (
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func Open(path string) (*Storage, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	if err := applyPragmas(db); err != nil {
		return nil, err
	}

	if err := applySchema(db); err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) Ping() error {
	return s.db.Ping()
}

func IsUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}

func (s *Storage) DeleteExpired(now int64) (int64, error) {
	result, err := s.db.Exec(`
		DELETE FROM entries
		WHERE expires_at IS NOT NULL
		AND expires_at <= ?
	`, now)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
