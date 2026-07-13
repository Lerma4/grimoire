package domain

import "errors"

// Note is a markdown document tracked by Grimoire.
type Note struct {
	ID         int64
	Title      string
	Body       string
	CreatedAt  string
	UpdatedAt  string
	ArchivedAt string
}

// Validate enforces note invariants.
func (n Note) Validate() error {
	if n.Title == "" {
		return errors.New("note title is required")
	}
	return nil
}
