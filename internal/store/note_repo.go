package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Lerma4/grimoire/internal/domain"
)

// NoteRepo persists notes.
type NoteRepo struct {
	db *sql.DB
}

// NewNoteRepo returns a NoteRepo backed by db.
func NewNoteRepo(db *sql.DB) *NoteRepo { return &NoteRepo{db: db} }

// NoteFilter narrows note queries.
type NoteFilter struct {
	TagID           int64
	Search          string
	IncludeArchived bool
}

// Create inserts a note and returns it with its ID.
func (r *NoteRepo) Create(ctx context.Context, n domain.Note) (domain.Note, error) {
	n.CreatedAt = domain.TimeStamp()
	n.UpdatedAt = n.CreatedAt
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO notes (title, body, created_at, updated_at)
		 VALUES (?,?,?,?)`,
		n.Title, n.Body, n.CreatedAt, n.UpdatedAt)
	if err != nil {
		return n, fmt.Errorf("insert note: %w", err)
	}
	n.ID, err = res.LastInsertId()
	return n, err
}

// Get returns a note by id.
func (r *NoteRepo) Get(ctx context.Context, id int64) (domain.Note, error) {
	row := r.db.QueryRowContext(ctx, noteCols+` WHERE n.id = ?`, id)
	return scanNote(row.Scan)
}

// List returns notes matching the filter.
func (r *NoteRepo) List(ctx context.Context, f NoteFilter) ([]domain.Note, error) {
	q := noteCols
	var args []any
	joins := ""
	clauses := []string{}
	if !f.IncludeArchived {
		clauses = append(clauses, "n.archived_at IS NULL")
	}
	if f.TagID != 0 {
		joins += " JOIN note_tags nt ON nt.note_id = n.id"
		clauses = append(clauses, "nt.tag_id = ?")
		args = append(args, f.TagID)
	}
	if f.Search != "" {
		clauses = append(clauses, "(n.title LIKE ? OR n.body LIKE ?)")
		pat := "%" + f.Search + "%"
		args = append(args, pat, pat)
	}
	q += joins
	if len(clauses) > 0 {
		q += " WHERE " + strings.Join(clauses, " AND ")
	}
	q += " ORDER BY n.updated_at DESC"

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("list notes: %w", err)
	}
	defer rows.Close()
	var out []domain.Note
	for rows.Next() {
		n, err := scanNote(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

// Update modifies a note.
func (r *NoteRepo) Update(ctx context.Context, n domain.Note) error {
	n.UpdatedAt = domain.TimeStamp()
	_, err := r.db.ExecContext(ctx,
		`UPDATE notes SET title=?, body=?, archived_at=?, updated_at=? WHERE id=?`,
		n.Title, n.Body, nullString(n.ArchivedAt), n.UpdatedAt, n.ID)
	if err != nil {
		return fmt.Errorf("update note %d: %w", n.ID, err)
	}
	return nil
}

// Archive stamps archived_at.
func (r *NoteRepo) Archive(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE notes SET archived_at=?, updated_at=? WHERE id=?`,
		domain.TimeStamp(), domain.TimeStamp(), id)
	return err
}

// Unarchive clears archived_at.
func (r *NoteRepo) Unarchive(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE notes SET archived_at=NULL, updated_at=? WHERE id=?`,
		domain.TimeStamp(), id)
	return err
}

// Delete removes a note by id.
func (r *NoteRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM notes WHERE id = ?`, id)
	return err
}

// Count returns the number of notes.
func (r *NoteRepo) Count(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM notes`).Scan(&n)
	return n, err
}

const noteCols = `SELECT n.id, n.title, COALESCE(n.body,''),
	n.created_at, n.updated_at, COALESCE(n.archived_at,'')
	FROM notes n`

func scanNote(scan scanner) (domain.Note, error) {
	var n domain.Note
	err := scan(&n.ID, &n.Title, &n.Body, &n.CreatedAt, &n.UpdatedAt, &n.ArchivedAt)
	if err != nil {
		return n, err
	}
	return n, nil
}
