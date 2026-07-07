// Package service contains the application logic of Grimoire. Services wrap
// repositories and enforce invariants and validation; they are the only layer
// the TUI talks to.
package service

import (
	"context"
	"database/sql"

	"github.com/Lerma4/grimoire/internal/store"
)

// Ctx is the default background context used by service calls.
var Ctx = context.Background()

// Services is the dependency container wiring all repositories into services.
// Construct it once at startup and hand it to the TUI.
type Services struct {
	Projects *ProjectService
	Tags     *TagService
	Tasks    *TaskService
	Notes    *NoteService
	Links    *LinkService
	Search   *SearchService
}

// NewServices builds the whole service layer over a single database handle.
func NewServices(db *sql.DB) *Services {
	projectRepo := store.NewProjectRepo(db)
	tagRepo := store.NewTagRepo(db)
	projects := NewProjectService(projectRepo)
	tags := NewTagService(tagRepo)
	tasks := NewTaskService(store.NewTaskRepo(db), projects, tags)
	notes := NewNoteService(store.NewNoteRepo(db), projects, tags)
	links := NewLinkService(store.NewLinkRepo(db), tasks, notes)
	return &Services{
		Projects: projects,
		Tags:     tags,
		Tasks:    tasks,
		Notes:    notes,
		Links:    links,
		Search:   NewSearchService(tasks, notes),
	}
}
