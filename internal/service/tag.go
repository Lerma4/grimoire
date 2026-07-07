package service

import (
	"context"
	"strings"

	"github.com/Lerma4/grimoire/internal/domain"
)

// TagRepo is the subset of the store this service needs.
type TagRepo interface {
	Create(context.Context, domain.Tag) (domain.Tag, error)
	Get(context.Context, int64) (domain.Tag, error)
	GetByName(context.Context, string) (domain.Tag, error)
	List(context.Context) ([]domain.Tag, error)
	Delete(context.Context, int64) error
	SetTaskTag(context.Context, int64, int64) error
	ClearTaskTag(context.Context, int64, int64) error
	TaskTags(context.Context, int64) ([]domain.Tag, error)
	SetNoteTag(context.Context, int64, int64) error
	ClearNoteTag(context.Context, int64, int64) error
	NoteTags(context.Context, int64) ([]domain.Tag, error)
}

// TagService applies validation around tag persistence.
type TagService struct {
	repo TagRepo
}

// NewTagService builds a TagService over repo.
func NewTagService(repo TagRepo) *TagService {
	return &TagService{repo: repo}
}

// Create inserts a tag (or returns the existing one with the same name).
func (s *TagService) Create(ctx context.Context, name, color string) (domain.Tag, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.Tag{}, errTagName
	}
	return s.repo.Create(ctx, domain.Tag{Name: name, Color: strings.TrimSpace(color)})
}

// FindOrCreate returns a tag by name, creating it if missing.
func (s *TagService) FindOrCreate(ctx context.Context, name string) (domain.Tag, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.Tag{}, errTagName
	}
	if t, err := s.repo.GetByName(ctx, name); err == nil {
		return t, nil
	}
	return s.repo.Create(ctx, domain.Tag{Name: name})
}

// List returns all tags.
func (s *TagService) List(ctx context.Context) ([]domain.Tag, error) { return s.repo.List(ctx) }

// TagTask attaches a tag to a task by tag name.
func (s *TagService) TagTask(ctx context.Context, taskID int64, name string) (domain.Tag, error) {
	t, err := s.FindOrCreate(ctx, name)
	if err != nil {
		return t, err
	}
	return t, s.repo.SetTaskTag(ctx, taskID, t.ID)
}

// UntagTask removes a tag from a task.
func (s *TagService) UntagTask(ctx context.Context, taskID, tagID int64) error {
	return s.repo.ClearTaskTag(ctx, taskID, tagID)
}

// TaskTags returns tags attached to a task.
func (s *TagService) TaskTags(ctx context.Context, taskID int64) ([]domain.Tag, error) {
	return s.repo.TaskTags(ctx, taskID)
}

// TagNote attaches a tag to a note by tag name.
func (s *TagService) TagNote(ctx context.Context, noteID int64, name string) (domain.Tag, error) {
	t, err := s.FindOrCreate(ctx, name)
	if err != nil {
		return t, err
	}
	return t, s.repo.SetNoteTag(ctx, noteID, t.ID)
}

// UntagNote removes a tag from a note.
func (s *TagService) UntagNote(ctx context.Context, noteID, tagID int64) error {
	return s.repo.ClearNoteTag(ctx, noteID, tagID)
}

// NoteTags returns tags attached to a note.
func (s *TagService) NoteTags(ctx context.Context, noteID int64) ([]domain.Tag, error) {
	return s.repo.NoteTags(ctx, noteID)
}
