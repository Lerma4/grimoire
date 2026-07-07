package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Lerma4/grimoire/internal/domain"
)

// LinkRepo persists task↔note associations.
type LinkRepo struct {
	db *sql.DB
}

// NewLinkRepo returns a LinkRepo backed by db.
func NewLinkRepo(db *sql.DB) *LinkRepo { return &LinkRepo{db: db} }

// Link associates a task with a note. It is idempotent.
func (r *LinkRepo) Link(ctx context.Context, taskID, noteID int64) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT OR IGNORE INTO task_notes (task_id, note_id, relation_type, created_at)
		 VALUES (?, ?, ?, ?)`,
		taskID, noteID, domain.RelationReference, domain.TimeStamp())
	if err != nil {
		return fmt.Errorf("link %d->%d: %w", taskID, noteID, err)
	}
	return nil
}

// Unlink removes a task↔note association.
func (r *LinkRepo) Unlink(ctx context.Context, taskID, noteID int64) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM task_notes WHERE task_id = ? AND note_id = ?`, taskID, noteID)
	return err
}

// IsLinked reports whether a task and note are associated.
func (r *LinkRepo) IsLinked(ctx context.Context, taskID, noteID int64) (bool, error) {
	var ok bool
	err := r.db.QueryRowContext(ctx,
		`SELECT 1 FROM task_notes WHERE task_id = ? AND note_id = ?`, taskID, noteID).Scan(&ok)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return ok, err
}

// NotesForTask returns notes linked to a task.
func (r *LinkRepo) NotesForTask(ctx context.Context, taskID int64) ([]domain.Note, error) {
	q := noteCols + ` JOIN task_notes tn ON tn.note_id = n.id WHERE tn.task_id = ? ORDER BY n.updated_at DESC`
	rows, err := r.db.QueryContext(ctx, q, taskID)
	if err != nil {
		return nil, fmt.Errorf("notes for task %d: %w", taskID, err)
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

// TasksForNote returns tasks linked to a note.
func (r *LinkRepo) TasksForNote(ctx context.Context, noteID int64) ([]domain.Task, error) {
	q := taskCols + ` JOIN task_notes tn ON tn.task_id = t.id WHERE tn.note_id = ? ORDER BY t.updated_at DESC`
	rows, err := r.db.QueryContext(ctx, q, noteID)
	if err != nil {
		return nil, fmt.Errorf("tasks for note %d: %w", noteID, err)
	}
	defer rows.Close()
	var out []domain.Task
	for rows.Next() {
		t, err := scanTask(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}
