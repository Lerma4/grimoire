package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Lerma4/grimoire/internal/domain"
	"github.com/Lerma4/grimoire/internal/store"
)

// TaskRepo is the subset of the store this service needs.
type TaskRepo interface {
	Create(context.Context, domain.Task) (domain.Task, error)
	Get(context.Context, int64) (domain.Task, error)
	List(context.Context, store.TaskFilter) ([]domain.Task, error)
	Update(context.Context, domain.Task) error
	SetStatus(context.Context, int64, string) error
	Delete(context.Context, int64) error
	Count(context.Context) (int, error)
}

// TaskFilter is the TUI-facing filter, re-exported from store.
type TaskFilter = store.TaskFilter

// TaskService handles task lifecycle and validation.
type TaskService struct {
	repo TaskRepo
	tags *TagService
}

// NewTaskService builds a TaskService.
func NewTaskService(repo TaskRepo, tags *TagService) *TaskService {
	return &TaskService{repo: repo, tags: tags}
}

// TaskInput holds the editable fields of a task.
type TaskInput struct {
	Title       string
	Description string
	Priority    string
	DueDate     string // RFC3339 or empty
}

// Create validates and inserts a task.
func (s *TaskService) Create(ctx context.Context, in TaskInput) (domain.Task, error) {
	in.Title = strings.TrimSpace(in.Title)
	if in.Title == "" {
		return domain.Task{}, errTitle
	}
	if in.Priority == "" {
		in.Priority = domain.PriorityMedium
	}
	t := domain.Task{
		Title:       in.Title,
		Description: strings.TrimSpace(in.Description),
		Priority:    in.Priority,
		DueDate:     in.DueDate,
		Status:      domain.StatusTodo,
	}
	if err := t.Validate(); err != nil {
		return t, err
	}
	return s.repo.Create(ctx, t)
}

// Get returns a task by id.
func (s *TaskService) Get(ctx context.Context, id int64) (domain.Task, error) {
	return s.repo.Get(ctx, id)
}

// List returns tasks matching the filter.
func (s *TaskService) List(ctx context.Context, f TaskFilter) ([]domain.Task, error) {
	return s.repo.List(ctx, store.TaskFilter(f))
}

// Update replaces a task's editable fields.
func (s *TaskService) Update(ctx context.Context, id int64, in TaskInput) (domain.Task, error) {
	t, err := s.repo.Get(ctx, id)
	if err != nil {
		return t, err
	}
	in.Title = strings.TrimSpace(in.Title)
	if in.Title == "" {
		return t, errTitle
	}
	t.Title = in.Title
	t.Description = strings.TrimSpace(in.Description)
	if in.Priority != "" {
		t.Priority = in.Priority
	}
	t.DueDate = in.DueDate
	if err := t.Validate(); err != nil {
		return t, err
	}
	return t, s.repo.Update(ctx, t)
}

// SetStatus transitions a task and stamps the relevant timestamp.
func (s *TaskService) SetStatus(ctx context.Context, id int64, status string) error {
	if !validStatus(status) {
		return errors.New("invalid status: " + status)
	}
	return s.repo.SetStatus(ctx, id, status)
}

// ToggleDone flips a task between todo/doing and done.
func (s *TaskService) ToggleDone(ctx context.Context, id int64) (domain.Task, error) {
	t, err := s.repo.Get(ctx, id)
	if err != nil {
		return t, err
	}
	if t.Status == domain.StatusDone {
		err = s.repo.SetStatus(ctx, id, domain.StatusTodo)
		t.Status = domain.StatusTodo
	} else {
		err = s.repo.SetStatus(ctx, id, domain.StatusDone)
		t.Status = domain.StatusDone
	}
	return t, err
}

// Archive archives a task.
func (s *TaskService) Archive(ctx context.Context, id int64) error {
	return s.repo.SetStatus(ctx, id, domain.StatusArchived)
}

// Delete removes a task permanently.
func (s *TaskService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// TagByName attaches a tag to a task by name.
func (s *TaskService) TagByName(ctx context.Context, id int64, name string) error {
	_, err := s.tags.TagTask(ctx, id, name)
	return err
}

// Tags returns tags attached to a task.
func (s *TaskService) Tags(ctx context.Context, id int64) ([]domain.Tag, error) {
	return s.tags.TaskTags(ctx, id)
}

// ParseDueDate accepts "YYYY-MM-DD" or RFC3339 and normalizes to RFC3339.
func ParseDueDate(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", nil
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC().Format(time.RFC3339), nil
		}
	}
	return "", errors.New("invalid due date; use YYYY-MM-DD or RFC3339")
}

func validStatus(s string) bool {
	for _, v := range domain.ValidStatuses {
		if v == s {
			return true
		}
	}
	return false
}
