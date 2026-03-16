package snippet

import (
	"fmt"
)

// Service handles snippet business logic
type Service struct {
	repo *Repository
}

// NewService creates a new snippet service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Save validates and saves a snippet (create or update)
func (s *Service) Save(snippet *Snippet) (*Snippet, error) {
	// Validate
	validator := snippet.Validate()
	if validator.HasErrors() {
		return nil, validator
	}

	// Save to database
	err := s.repo.Save(snippet)
	if err != nil {
		return nil, err
	}

	return snippet, nil
}

// Delete removes a snippet by ID
func (s *Service) Delete(id int64) error {
	_, err := s.repo.Delete(id)
	return err
}

// GetByID retrieves a snippet by ID
func (s *Service) GetByID(id int64) (*Snippet, error) {
	return s.repo.GetByID(id)
}

// GetAll retrieves all snippets
func (s *Service) GetAll() ([]*Snippet, error) {
	return s.repo.GetAll()
}

// GetGlobal retrieves snippets with no game scope
func (s *Service) GetGlobal() ([]*Snippet, error) {
	return s.repo.GetGlobal()
}

// GetByGameID retrieves snippets scoped to a specific game
func (s *Service) GetByGameID(gameID int64) ([]*Snippet, error) {
	return s.repo.GetByGameID(gameID)
}

func (s *Service) Reorder(snippetID int64, direction int) (int64, error) {
	// Find the snippet first so we know its scope (game or global)
	curr, err := s.repo.GetByID(snippetID)
	if err != nil {
		return 0, fmt.Errorf("snippet %d not found", snippetID)
	}

	// Operate only within the same scope so game snippets never swap with global ones
	var snippets []*Snippet
	if curr.GameID != nil {
		snippets, err = s.repo.GetByGameID(*curr.GameID)
	} else {
		snippets, err = s.repo.GetGlobal()
	}
	if err != nil {
		return 0, err
	}

	idx := -1
	for i, s := range snippets {
		if s.ID == snippetID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return 0, fmt.Errorf("snippet %d not found in its scope", snippetID)
	}

	neighborIdx := idx + direction
	if neighborIdx < 0 || neighborIdx >= len(snippets) {
		return 0, nil // List boundary
	}

	// Swap using list indices as positions so duplicates are resolved
	if err := s.repo.SwapPositions(curr.ID, idx, snippets[neighborIdx].ID, neighborIdx); err != nil {
		return 0, err
	}
	return curr.ID, nil
}
