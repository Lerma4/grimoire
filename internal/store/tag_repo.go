package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Lerma4/grimoire/internal/domain"
)

// TagRepo persists tags.
type TagRepo struct {
	db *sql.DB
}

// NewTagRepo returns a TagRepo backed by db.
func NewTagRepo(db *sql.DB) *TagRepo { return &TagRepo{db: db} }

// Create inserts a tag, or returns the existing one if the name is taken.
func (r *TagRepo) Create(ctx context.Context, t domain.Tag) (domain.Tag, error) {
	t.CreatedAt = domain.TimeStamp()
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO tags (name, color, created_at) VALUES (?, ?, ?)`,
		t.Name, t.Color, t.CreatedAt)
	if err != nil {
		// Unique constraint: fetch existing.
		if existing, gerr := r.GetByName(ctx, t.Name); gerr == nil {
			return existing, nil
		}
		return t, fmt.Errorf("insert tag: %w", err)
	}
	t.ID, err = res.LastInsertId()
	return t, err
}

// Get returns a tag by id.
func (r *TagRepo) Get(ctx context.Context, id int64) (domain.Tag, error) {
	var t domain.Tag
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, COALESCE(color,''), created_at FROM tags WHERE id = ?`, id).
		Scan(&t.ID, &t.Name, &t.Color, &t.CreatedAt)
	if err != nil {
		return t, fmt.Errorf("get tag %d: %w", id, err)
	}
	return t, nil
}

// GetByName returns a tag by name.
func (r *TagRepo) GetByName(ctx context.Context, name string) (domain.Tag, error) {
	var t domain.Tag
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, COALESCE(color,''), created_at FROM tags WHERE name = ?`, name).
		Scan(&t.ID, &t.Name, &t.Color, &t.CreatedAt)
	if err != nil {
		return t, fmt.Errorf("get tag %q: %w", name, err)
	}
	return t, nil
}

// List returns all tags ordered by name.
func (r *TagRepo) List(ctx context.Context) ([]domain.Tag, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, COALESCE(color,''), created_at FROM tags ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}
	defer rows.Close()
	var out []domain.Tag
	for rows.Next() {
		var t domain.Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.Color, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

// Delete removes a tag by id. Associated task_tags/note_tags cascade.
func (r *TagRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM tags WHERE id = ?`, id)
	return err
}

// SetTaskTag ensures the tag is attached to the task.
func (r *TagRepo) SetTaskTag(ctx context.Context, taskID, tagID int64) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT OR IGNORE INTO task_tags (task_id, tag_id) VALUES (?, ?)`, taskID, tagID)
	return err
}

// ClearTaskTag removes the tag from the task.
func (r *TagRepo) ClearTaskTag(ctx context.Context, taskID, tagID int64) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM task_tags WHERE task_id = ? AND tag_id = ?`, taskID, tagID)
	return err
}

// TaskTags returns tags for a task.
func (r *TagRepo) TaskTags(ctx context.Context, taskID int64) ([]domain.Tag, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT t.id, t.name, COALESCE(t.color,''), t.created_at
		 FROM tags t JOIN task_tags tt ON tt.tag_id = t.id
		 WHERE tt.task_id = ? ORDER BY t.name`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTags(rows)
}

// NoteTags returns tags for a note.
func (r *TagRepo) NoteTags(ctx context.Context, noteID int64) ([]domain.Tag, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT t.id, t.name, COALESCE(t.color,''), t.created_at
		 FROM tags t JOIN note_tags nt ON nt.tag_id = t.id
		 WHERE nt.note_id = ? ORDER BY t.name`, noteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTags(rows)
}

// SetNoteTag ensures the tag is attached to the note.
func (r *TagRepo) SetNoteTag(ctx context.Context, noteID, tagID int64) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT OR IGNORE INTO note_tags (note_id, tag_id) VALUES (?, ?)`, noteID, tagID)
	return err
}

// ClearNoteTag removes the tag from the note.
func (r *TagRepo) ClearNoteTag(ctx context.Context, noteID, tagID int64) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM note_tags WHERE note_id = ? AND tag_id = ?`, noteID, tagID)
	return err
}

func scanTags(rows *sql.Rows) ([]domain.Tag, error) {
	var out []domain.Tag
	for rows.Next() {
		var t domain.Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.Color, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}
