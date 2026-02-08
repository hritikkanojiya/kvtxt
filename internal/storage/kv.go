package storage

import (
	"database/sql"
)

type Entry struct {
	Hash      string
	Payload   []byte
	CreatedAt int64
	ExpiresAt sql.NullInt64
}

func (s *Storage) Insert(e *Entry) error {
	const q = `
	INSERT INTO kv (hash, payload, created_at, expires_at)
	VALUES (?, ?, ?, ?)
	`

	_, err := s.db.Exec(
		q,
		e.Hash,
		e.Payload,
		e.CreatedAt,
		e.ExpiresAt,
	)

	return err
}
