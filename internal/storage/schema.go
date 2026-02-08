package storage

import "database/sql"

func applyPragmas(db *sql.DB) error {
	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA foreign_keys = ON;",
	}

	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			return err
		}
	}

	return nil
}

func applySchema(db *sql.DB) error {
	const schema = `
	CREATE TABLE IF NOT EXISTS kv (
		hash TEXT PRIMARY KEY,
		payload BLOB NOT NULL,
		content_type TEXT NOT NULL,
		created_at INTEGER NOT NULL,
		expires_at INTEGER
	);
	`

	_, err := db.Exec(schema)
	return err
}
