package service

import (
	"context"

	"github.com/Lerma4/grimoire/internal/domain"
)

// LinkRepo is the subset of the store this service needs.
type LinkRepo interface {
	Link(context.Context, int64, int64) error
	Unlink(context.Context, int64, int64) error
	IsLinked(context.Context, int64, int64) (bool, error)
	NotesForTask(context.Context, int64) ([]domain.Note, error)
	TasksForNote(context.Context, int64) ([]domain.Task, error)
}

// LinkService owns task↔note associations and both directions of lookup.
type LinkService struct {
	repo  LinkRepo
	tasks taskGetter
	notes noteGetter
}

// taskGetter is the subset of TaskService LinkService needs.
type taskGetter interface {
	Get(context.Context, int64) (domain.Task, error)
}

// noteGetter is the subset of NoteService LinkService needs.
type noteGetter interface {
	Get(context.Context, int64) (domain.Note, error)
}

// NewLinkService builds a LinkService.
func NewLinkService(repo LinkRepo, tasks taskGetter, notes noteGetter) *LinkService {
	return &LinkService{repo: repo, tasks: tasks, notes: notes}
}

// LinkTaskNote associates a task with a note (idempotent).
func (s *LinkService) LinkTaskNote(ctx context.Context, taskID, noteID int64) error {
	if _, err := s.tasks.Get(ctx, taskID); err != nil {
		return err
	}
	if _, err := s.notes.Get(ctx, noteID); err != nil {
		return err
	}
	return s.repo.Link(ctx, taskID, noteID)
}

// UnlinkTaskNote removes the association.
func (s *LinkService) UnlinkTaskNote(ctx context.Context, taskID, noteID int64) error {
	return s.repo.Unlink(ctx, taskID, noteID)
}

// IsLinked reports whether the two are associated.
func (s *LinkService) IsLinked(ctx context.Context, taskID, noteID int64) (bool, error) {
	return s.repo.IsLinked(ctx, taskID, noteID)
}

// NotesForTask returns notes linked to a task.
func (s *LinkService) NotesForTask(ctx context.Context, taskID int64) ([]domain.Note, error) {
	return s.repo.NotesForTask(ctx, taskID)
}

// TasksForNote returns tasks linked to a note.
func (s *LinkService) TasksForNote(ctx context.Context, noteID int64) ([]domain.Task, error) {
	return s.repo.TasksForNote(ctx, noteID)
}
