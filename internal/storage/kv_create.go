package storage

import (
	"database/sql"
)

type Entry struct {
	Hash        string
	Payload     []byte
	ContentType string
	CreatedAt   int64
	ExpiresAt   sql.NullInt64
}

func (s *Storage) Insert(e *Entry) error {
	const q = `
	INSERT INTO kv (hash, payload, content_type, created_at, expires_at)
	VALUES (?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(
		q,
		e.Hash,
		e.Payload,
		e.ContentType,
		e.CreatedAt,
		e.ExpiresAt,
	)

	return err
}

func (e *Entry) ExpiresAtPtr() *int64 {
	if e.ExpiresAt.Valid {
		return &e.ExpiresAt.Int64
	}
	return nil
}
