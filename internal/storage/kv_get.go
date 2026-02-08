package storage

import (
	"database/sql"
)

func (s *Storage) Get(hash string) (*Entry, error) {
	const q = `
	SELECT hash, payload, content_type, created_at, expires_at
	FROM kv
	WHERE hash = ?
	`

	var e Entry
	err := s.db.QueryRow(q, hash).Scan(
		&e.Hash,
		&e.Payload,
		&e.ContentType,
		&e.CreatedAt,
		&e.ExpiresAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &e, nil
}
