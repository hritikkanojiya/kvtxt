package storage

func (s *Storage) DeleteExpired(now int64) (int64, error) {
	result, err := s.db.Exec(`
		DELETE FROM kv
		WHERE expires_at IS NOT NULL
		AND expires_at <= ?
	`, now)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
