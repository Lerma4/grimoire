package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Lerma4/grimoire/internal/domain"
)

// ProjectRepo persists projects.
type ProjectRepo struct {
	db *sql.DB
}

// NewProjectRepo returns a ProjectRepo backed by db.
func NewProjectRepo(db *sql.DB) *ProjectRepo { return &ProjectRepo{db: db} }

// Create inserts a new project and returns it with its ID.
func (r *ProjectRepo) Create(ctx context.Context, p domain.Project) (domain.Project, error) {
	p.CreatedAt = domain.TimeStamp()
	p.UpdatedAt = p.CreatedAt
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO projects (name, description, created_at, updated_at) VALUES (?, ?, ?, ?)`,
		p.Name, p.Description, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return p, fmt.Errorf("insert project: %w", err)
	}
	p.ID, err = res.LastInsertId()
	return p, err
}

// Get returns a project by id.
func (r *ProjectRepo) Get(ctx context.Context, id int64) (domain.Project, error) {
	var p domain.Project
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, COALESCE(description,''), created_at, updated_at
		 FROM projects WHERE id = ?`, id).
		Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return p, fmt.Errorf("get project %d: %w", id, err)
	}
	return p, nil
}

// GetByName returns a project by name.
func (r *ProjectRepo) GetByName(ctx context.Context, name string) (domain.Project, error) {
	var p domain.Project
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, COALESCE(description,''), created_at, updated_at
		 FROM projects WHERE name = ?`, name).
		Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return p, fmt.Errorf("get project %q: %w", name, err)
	}
	return p, nil
}

// List returns all projects ordered by name.
func (r *ProjectRepo) List(ctx context.Context) ([]domain.Project, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, COALESCE(description,''), created_at, updated_at
		 FROM projects ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	defer rows.Close()
	var out []domain.Project
	for rows.Next() {
		var p domain.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// Update modifies an existing project.
func (r *ProjectRepo) Update(ctx context.Context, p domain.Project) error {
	p.UpdatedAt = domain.TimeStamp()
	_, err := r.db.ExecContext(ctx,
		`UPDATE projects SET name=?, description=?, updated_at=? WHERE id=?`,
		p.Name, p.Description, p.UpdatedAt, p.ID)
	return err
}

// Delete removes a project by id.
func (r *ProjectRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM projects WHERE id = ?`, id)
	return err
}
