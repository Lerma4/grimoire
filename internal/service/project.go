package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Lerma4/grimoire/internal/domain"
)

// ProjectRepo is the subset of the store this service needs.
type ProjectRepo interface {
	Create(context.Context, domain.Project) (domain.Project, error)
	Get(context.Context, int64) (domain.Project, error)
	GetByName(context.Context, string) (domain.Project, error)
	List(context.Context) ([]domain.Project, error)
	Update(context.Context, domain.Project) error
	Delete(context.Context, int64) error
}

// ProjectService applies validation around project persistence.
type ProjectService struct {
	repo ProjectRepo
}

// NewProjectService builds a ProjectService over repo.
func NewProjectService(repo ProjectRepo) *ProjectService {
	return &ProjectService{repo: repo}
}

// ErrValidation is returned when input fails validation.
var ErrValidation = errors.New("validation error")

// Create validates and inserts a project.
func (s *ProjectService) Create(ctx context.Context, name, description string) (domain.Project, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.Project{}, errors.New("project name is required")
	}
	return s.repo.Create(ctx, domain.Project{Name: name, Description: strings.TrimSpace(description)})
}

// Get returns a project by id.
func (s *ProjectService) Get(ctx context.Context, id int64) (domain.Project, error) {
	return s.repo.Get(ctx, id)
}

// FindOrCreate returns a project by name, creating it if missing.
func (s *ProjectService) FindOrCreate(ctx context.Context, name string) (domain.Project, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.Project{}, errors.New("project name is required")
	}
	if p, err := s.repo.GetByName(ctx, name); err == nil {
		return p, nil
	}
	return s.repo.Create(ctx, domain.Project{Name: name})
}

// List returns all projects.
func (s *ProjectService) List(ctx context.Context) ([]domain.Project, error) {
	return s.repo.List(ctx)
}

// Update modifies a project.
func (s *ProjectService) Update(ctx context.Context, p domain.Project) (domain.Project, error) {
	if strings.TrimSpace(p.Name) == "" {
		return p, errors.New("project name is required")
	}
	return p, s.repo.Update(ctx, p)
}

// Delete removes a project.
func (s *ProjectService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
