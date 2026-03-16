package snippet

import (
	"soloterm/database"
)

func init() {
	// Register this package's migrations with the database package
	database.RegisterMigration(Migrate)
}

// Migrate runs all migrations for the snippet domain
func Migrate(dbStore *database.DBStore) error {
	// Migration: Create snippets table (previously saved_rolls)
	if err := createSnippetsTable(dbStore); err != nil {
		return err
	}

	return nil
}

// createSnippetsTable creates the initial snippets table and index
func createSnippetsTable(dbStore *database.DBStore) error {
	schema := `
		CREATE TABLE IF NOT EXISTS snippets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			game_id INTEGER,
			name STRING NOT NULL,
			content STRING NOT NULL,
			position INTEGER NOT NULL default 0,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_snippets_by_name ON snippets (name);
		CREATE INDEX IF NOT EXISTS idx_snippets_by_game_id ON snippets (game_id);
	`
	_, err := dbStore.Connection.Exec(schema)
	return err
}
