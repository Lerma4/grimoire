package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Lerma4/grimoire/internal/domain"
)

// TaskRepo persists tasks.
type TaskRepo struct {
	db *sql.DB
}

// NewTaskRepo returns a TaskRepo backed by db.
func NewTaskRepo(db *sql.DB) *TaskRepo { return &TaskRepo{db: db} }

// TaskFilter narrows task queries. Zero-value fields are ignored.
type TaskFilter struct {
	Status          string
	ProjectID       int64
	TagID           int64
	Overdue         bool
	Search          string
	IncludeArchived bool
}

// Create inserts a task and returns it with its ID.
func (r *TaskRepo) Create(ctx context.Context, t domain.Task) (domain.Task, error) {
	if t.Status == "" {
		t.Status = domain.StatusTodo
	}
	if t.Priority == "" {
		t.Priority = domain.PriorityMedium
	}
	t.CreatedAt = domain.TimeStamp()
	t.UpdatedAt = t.CreatedAt
	res, err := r.db.ExecContext(ctx, `INSERT INTO tasks
		(title, description, status, priority, due_date, project_id, created_at, updated_at)
		VALUES (?,?,?,?,?,?,?,?)`,
		t.Title, nullString(t.Description), t.Status, t.Priority,
		nullString(t.DueDate), nullID(t.ProjectID), t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return t, fmt.Errorf("insert task: %w", err)
	}
	t.ID, err = res.LastInsertId()
	return t, err
}

// Get returns a task by id.
func (r *TaskRepo) Get(ctx context.Context, id int64) (domain.Task, error) {
	row := r.db.QueryRowContext(ctx, taskCols+` WHERE t.id = ?`, id)
	return scanTask(row.Scan)
}

// List returns tasks matching the filter, most recently updated first.
func (r *TaskRepo) List(ctx context.Context, f TaskFilter) ([]domain.Task, error) {
	q := taskCols
	var args []any
	joins := ""
	clauses := []string{}
	if !f.IncludeArchived {
		clauses = append(clauses, "t.status <> 'archived'")
	}
	if f.Status != "" {
		clauses = append(clauses, "t.status = ?")
		args = append(args, f.Status)
	}
	if f.ProjectID != 0 {
		clauses = append(clauses, "t.project_id = ?")
		args = append(args, f.ProjectID)
	}
	if f.TagID != 0 {
		joins += " JOIN task_tags tt ON tt.task_id = t.id"
		clauses = append(clauses, "tt.tag_id = ?")
		args = append(args, f.TagID)
	}
	if f.Overdue {
		clauses = append(clauses, "t.due_date IS NOT NULL AND t.due_date < ? AND t.status NOT IN ('done','archived')")
		args = append(args, domain.TimeStamp())
	}
	if f.Search != "" {
		clauses = append(clauses, "(t.title LIKE ? OR t.description LIKE ?)")
		pat := "%" + f.Search + "%"
		args = append(args, pat, pat)
	}
	q += joins
	if len(clauses) > 0 {
		q += " WHERE " + strings.Join(clauses, " AND ")
	}
	q += " ORDER BY t.updated_at DESC"

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
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

// Update modifies a task. Status/timestamp transitions are handled by the service.
func (r *TaskRepo) Update(ctx context.Context, t domain.Task) error {
	t.UpdatedAt = domain.TimeStamp()
	_, err := r.db.ExecContext(ctx, `UPDATE tasks SET
		title=?, description=?, status=?, priority=?, due_date=?, project_id=?,
		completed_at=?, archived_at=?, updated_at=? WHERE id=?`,
		t.Title, nullString(t.Description), t.Status, t.Priority,
		nullString(t.DueDate), nullID(t.ProjectID),
		nullString(t.CompletedAt), nullString(t.ArchivedAt), t.UpdatedAt, t.ID)
	if err != nil {
		return fmt.Errorf("update task %d: %w", t.ID, err)
	}
	return nil
}

// SetStatus updates only status and the relevant timestamp.
func (r *TaskRepo) SetStatus(ctx context.Context, id int64, status string) error {
	ts := domain.TimeStamp()
	q := "UPDATE tasks SET status=?, updated_at=?"
	args := []any{status, ts}
	switch status {
	case domain.StatusDone:
		q += ", completed_at=?"
		args = append(args, ts)
	case domain.StatusArchived:
		q += ", archived_at=?"
		args = append(args, ts)
	case domain.StatusTodo, domain.StatusDoing:
		q += ", completed_at=NULL"
	}
	q += " WHERE id=?"
	args = append(args, id)
	_, err := r.db.ExecContext(ctx, q, args...)
	return err
}

// Delete removes a task by id.
func (r *TaskRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = ?`, id)
	return err
}

// Count returns the number of tasks.
func (r *TaskRepo) Count(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM tasks`).Scan(&n)
	return n, err
}

const taskCols = `SELECT t.id, t.title, COALESCE(t.description,''), t.status, t.priority,
	COALESCE(t.due_date,''), COALESCE(t.project_id,0),
	t.created_at, t.updated_at, COALESCE(t.completed_at,''), COALESCE(t.archived_at,'')
	FROM tasks t`

type scanner func(dest ...any) error

func scanTask(scan scanner) (domain.Task, error) {
	var t domain.Task
	err := scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority,
		&t.DueDate, &t.ProjectID, &t.CreatedAt, &t.UpdatedAt, &t.CompletedAt, &t.ArchivedAt)
	if err != nil {
		return t, err
	}
	return t, nil
}
