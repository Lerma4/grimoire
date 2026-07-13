package service

import (
	"context"
	"strings"

	"github.com/Lerma4/grimoire/internal/domain"
	"github.com/Lerma4/grimoire/internal/store"
)

// NoteRepo is the subset of the store this service needs.
type NoteRepo interface {
	Create(context.Context, domain.Note) (domain.Note, error)
	Get(context.Context, int64) (domain.Note, error)
	List(context.Context, store.NoteFilter) ([]domain.Note, error)
	Update(context.Context, domain.Note) error
	Archive(context.Context, int64) error
	Unarchive(context.Context, int64) error
	Delete(context.Context, int64) error
	Count(context.Context) (int, error)
}

// NoteFilter is the TUI-facing filter, re-exported from store.
type NoteFilter = store.NoteFilter

// NoteService handles note lifecycle and validation.
type NoteService struct {
	repo NoteRepo
	tags *TagService
}

// NewNoteService builds a NoteService.
func NewNoteService(repo NoteRepo, tags *TagService) *NoteService {
	return &NoteService{repo: repo, tags: tags}
}

// NoteInput holds the editable fields of a note.
type NoteInput struct {
	Title string
	Body  string
}

// Create validates and inserts a note.
func (s *NoteService) Create(ctx context.Context, in NoteInput) (domain.Note, error) {
	in.Title = strings.TrimSpace(in.Title)
	if in.Title == "" {
		return domain.Note{}, errTitle
	}
	n := domain.Note{Title: in.Title, Body: in.Body}
	if err := n.Validate(); err != nil {
		return n, err
	}
	return s.repo.Create(ctx, n)
}

// Get returns a note by id.
func (s *NoteService) Get(ctx context.Context, id int64) (domain.Note, error) {
	return s.repo.Get(ctx, id)
}

// List returns notes matching the filter.
func (s *NoteService) List(ctx context.Context, f NoteFilter) ([]domain.Note, error) {
	return s.repo.List(ctx, store.NoteFilter(f))
}

// Update replaces a note's editable fields.
func (s *NoteService) Update(ctx context.Context, id int64, in NoteInput) (domain.Note, error) {
	n, err := s.repo.Get(ctx, id)
	if err != nil {
		return n, err
	}
	in.Title = strings.TrimSpace(in.Title)
	if in.Title == "" {
		return n, errTitle
	}
	n.Title = in.Title
	n.Body = in.Body
	if err := n.Validate(); err != nil {
		return n, err
	}
	return n, s.repo.Update(ctx, n)
}

// Archive archives a note.
func (s *NoteService) Archive(ctx context.Context, id int64) error { return s.repo.Archive(ctx, id) }

// Unarchive restores a note.
func (s *NoteService) Unarchive(ctx context.Context, id int64) error {
	return s.repo.Unarchive(ctx, id)
}

// Delete removes a note permanently.
func (s *NoteService) Delete(ctx context.Context, id int64) error { return s.repo.Delete(ctx, id) }

// TagByName attaches a tag to a note by name.
func (s *NoteService) TagByName(ctx context.Context, id int64, name string) error {
	_, err := s.tags.TagNote(ctx, id, name)
	return err
}

// Tags returns tags attached to a note.
func (s *NoteService) Tags(ctx context.Context, id int64) ([]domain.Tag, error) {
	return s.tags.NoteTags(ctx, id)
}
