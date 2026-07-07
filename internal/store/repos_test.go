package store

import (
	"context"
	"testing"

	"github.com/Lerma4/grimoire/internal/domain"
)

func TestTagRepo_CreateIdempotentAndJoins(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()
	tags := NewTagRepo(db)
	tasks := NewTaskRepo(db)
	notes := NewNoteRepo(db)

	t1, _ := tasks.Create(ctx, domain.Task{Title: "tag me"})
	n1, _ := notes.Create(ctx, domain.Note{Title: "tag me too"})

	first, err := tags.Create(ctx, domain.Tag{Name: "bug"})
	if err != nil {
		t.Fatalf("create tag: %v", err)
	}
	// creating the same name returns the existing tag (no duplicate).
	dup, err := tags.Create(ctx, domain.Tag{Name: "bug"})
	if err != nil || dup.ID != first.ID {
		t.Fatalf("expected idempotent create, got %+v %v", dup, err)
	}
	var count int
	db.QueryRow(`SELECT COUNT(*) FROM tags WHERE name='bug'`).Scan(&count)
	if count != 1 {
		t.Fatalf("expected 1 tag row, got %d", count)
	}

	if err := tags.SetTaskTag(ctx, t1.ID, first.ID); err != nil {
		t.Fatalf("set task tag: %v", err)
	}
	if err := tags.SetNoteTag(ctx, n1.ID, first.ID); err != nil {
		t.Fatalf("set note tag: %v", err)
	}

	tt, _ := tags.TaskTags(ctx, t1.ID)
	if len(tt) != 1 || tt[0].Name != "bug" {
		t.Fatalf("task tags = %+v", tt)
	}
	nt, _ := tags.NoteTags(ctx, n1.ID)
	if len(nt) != 1 {
		t.Fatalf("note tags = %+v", nt)
	}

	// double-set is a no-op.
	_ = tags.SetTaskTag(ctx, t1.ID, first.ID)
	tt, _ = tags.TaskTags(ctx, t1.ID)
	if len(tt) != 1 {
		t.Fatalf("double-set should not duplicate: %+v", tt)
	}

	if err := tags.ClearTaskTag(ctx, t1.ID, first.ID); err != nil {
		t.Fatalf("clear task tag: %v", err)
	}
	tt, _ = tags.TaskTags(ctx, t1.ID)
	if len(tt) != 0 {
		t.Fatalf("expected cleared, got %+v", tt)
	}
}

func TestLinkRepo_BothDirections(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()
	tasks := NewTaskRepo(db)
	notes := NewNoteRepo(db)
	links := NewLinkRepo(db)

	t1, _ := tasks.Create(ctx, domain.Task{Title: "t1"})
	n1, _ := notes.Create(ctx, domain.Note{Title: "n1"})
	n2, _ := notes.Create(ctx, domain.Note{Title: "n2"})

	if err := links.Link(ctx, t1.ID, n1.ID); err != nil {
		t.Fatalf("link: %v", err)
	}
	// idempotent
	if err := links.Link(ctx, t1.ID, n1.ID); err != nil {
		t.Fatalf("re-link: %v", err)
	}
	_ = links.Link(ctx, t1.ID, n2.ID)

	ok, _ := links.IsLinked(ctx, t1.ID, n1.ID)
	if !ok {
		t.Fatal("expected linked")
	}
	ok, _ = links.IsLinked(ctx, t1.ID, n2.ID+9999)
	if ok {
		t.Fatal("expected not linked")
	}

	notesForTask, _ := links.NotesForTask(ctx, t1.ID)
	if len(notesForTask) != 2 {
		t.Fatalf("notes for task = %d", len(notesForTask))
	}
	tasksForNote, _ := links.TasksForNote(ctx, n1.ID)
	if len(tasksForNote) != 1 || tasksForNote[0].ID != t1.ID {
		t.Fatalf("tasks for note = %+v", tasksForNote)
	}

	if err := links.Unlink(ctx, t1.ID, n1.ID); err != nil {
		t.Fatalf("unlink: %v", err)
	}
	notesForTask, _ = links.NotesForTask(ctx, t1.ID)
	if len(notesForTask) != 1 {
		t.Fatalf("after unlink notes = %d", len(notesForTask))
	}

	// cascading: deleting a task removes its links.
	_ = links.Link(ctx, t1.ID, n1.ID)
	_ = tasks.Delete(ctx, t1.ID)
	tasksForNote, _ = links.TasksForNote(ctx, n1.ID)
	if len(tasksForNote) != 0 {
		t.Fatalf("task delete should cascade links, got %d", len(tasksForNote))
	}
}
