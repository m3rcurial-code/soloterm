package snippet

import (
	"database/sql"
	"errors"
	"fmt"
	"soloterm/database"
)

// Repository handles database operations for snippets
type Repository struct {
	db *database.DBStore
}

// NewRepository creates a new Repository
func NewRepository(db *database.DBStore) *Repository {
	return &Repository{db: db}
}

// Save creates or updates a snippet
// Automatically manages created_at, and updated_at
// The snippet pointer is updated with the current values after save
func (r *Repository) Save(snippet *Snippet) error {
	if snippet.ID == 0 {
		// INSERT - new snippet
		return r.insert(snippet)
	} else {
		// UPDATE - existing snippet
		return r.update(snippet)
	}
}

// Delete removes a snippet by id
// Returns the number of rows deleted and an error if the id doesn't exist
func (r *Repository) Delete(id int64) (int64, error) {
	if id == 0 {
		return 0, errors.New("id cannot be empty")
	}

	query := `DELETE FROM snippets WHERE id = ?`

	result, err := r.db.Connection.Exec(query, id)

	if err != nil {
		return 0, err
	}

	// Check if a row was actually deleted
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if rows == 0 {
		return 0, fmt.Errorf("id '%d' not found", id)
	}

	return rows, nil
}

// GetByID retrieves a snippet by ID
func (r *Repository) GetByID(id int64) (*Snippet, error) {
	if id == 0 {
		return nil, errors.New("id cannot be zero")
	}

	var snippet Snippet
	err := r.db.Connection.Get(&snippet, "SELECT * FROM snippets WHERE id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("snippet not found")
		}
		return nil, err
	}

	return &snippet, nil
}

// GetAll retrieves all snippets ordered by position then name
func (r *Repository) GetAll() ([]*Snippet, error) {
	var snippets []*Snippet
	err := r.db.Connection.Select(&snippets, "SELECT * FROM snippets ORDER BY position, name")
	if err != nil {
		return nil, err
	}
	return snippets, nil
}

// GetGlobal retrieves snippets with no game scope, ordered by position then name
func (r *Repository) GetGlobal() ([]*Snippet, error) {
	var snippets []*Snippet
	err := r.db.Connection.Select(&snippets, "SELECT * FROM snippets WHERE game_id IS NULL ORDER BY position, name")
	return snippets, err
}

// GetByGameID retrieves snippets scoped to a specific game, ordered by position then name
func (r *Repository) GetByGameID(gameID int64) ([]*Snippet, error) {
	var snippets []*Snippet
	err := r.db.Connection.Select(&snippets, "SELECT * FROM snippets WHERE game_id = ? ORDER BY position, name", gameID)
	return snippets, err
}

func (r *Repository) GetByName(name string) (*Snippet, error) {
	var snippet Snippet
	err := r.db.Connection.Get(&snippet, "SELECT * FROM snippets WHERE lower(name) = lower(?) LIMIT 1", name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("snippet named '%s' not found.", name)
		}
		return nil, err
	}
	return &snippet, nil
}

// Inserts a new record using positions already set on the snippet struct
func (r *Repository) insert(snippet *Snippet) error {
	query := `
		INSERT INTO snippets (name, content, game_id, position, created_at, updated_at)
		VALUES (?, ?, ?, ?, datetime('now', 'subsec'), datetime('now', 'subsec'))
		RETURNING id, created_at, updated_at
	`

	return r.db.Connection.QueryRowx(query,
		snippet.Name,
		snippet.Content,
		snippet.GameID,
		snippet.Position,
	).StructScan(snippet)
}

// Updates an existing record including sort positions
func (r *Repository) update(snippet *Snippet) error {
	query := `
		UPDATE snippets SET name = ?, content = ?, game_id = ?, position = ?, updated_at = datetime('now','subsec')
		WHERE id = ?
		RETURNING created_at, updated_at
	`

	return r.db.Connection.QueryRowx(query,
		snippet.Name,
		snippet.Content,
		snippet.GameID,
		snippet.Position,
		snippet.ID,
	).StructScan(snippet)
}

// SwapPositions atomically swaps the position of two snippets
func (r *Repository) SwapPositions(idA int64, posA int, idB int64, posB int) error {
	query := `UPDATE snippets
	          SET position = CASE
	              WHEN id = ? THEN ?
	              WHEN id = ? THEN ?
	          END,
	          updated_at = datetime('now','subsec')
	          WHERE id IN (?, ?)`
	_, err := r.db.Connection.Exec(query, idA, posB, idB, posA, idA, idB)
	return err
}
