package domain

import (
	"errors"
	"time"
)

// Task status values.
const (
	StatusTodo     = "todo"
	StatusDoing    = "doing"
	StatusDone     = "done"
	StatusArchived = "archived"
)

// Priority values.
const (
	PriorityLow    = "low"
	PriorityMedium = "medium"
	PriorityHigh   = "high"
	PriorityUrgent = "urgent"
)

// Task is a unit of work tracked by Grimoire.
type Task struct {
	ID          int64
	Title       string
	Description string
	Status      string
	Priority    string
	DueDate     string
	ProjectID   int64
	CreatedAt   string
	UpdatedAt   string
	CompletedAt string
	ArchivedAt  string
}

// ValidStatuses is the set of allowed task statuses.
var ValidStatuses = []string{StatusTodo, StatusDoing, StatusDone, StatusArchived}

// ValidPriorities is the set of allowed task priorities.
var ValidPriorities = []string{PriorityLow, PriorityMedium, PriorityHigh, PriorityUrgent}

// Validate enforces task invariants.
func (t Task) Validate() error {
	if t.Title == "" {
		return errors.New("task title is required")
	}
	if t.Status == "" {
		return errors.New("task status is required")
	}
	if !contains(ValidStatuses, t.Status) {
		return errors.New("invalid task status: " + t.Status)
	}
	if t.Priority != "" && !contains(ValidPriorities, t.Priority) {
		return errors.New("invalid task priority: " + t.Priority)
	}
	return nil
}

// StatusGlyph returns the textual symbol for a status, so state is never
// communicated by color alone.
func StatusGlyph(status string) string {
	switch status {
	case StatusDone:
		return "●"
	case StatusDoing:
		return "◐"
	case StatusArchived:
		return "□"
	default:
		return "○"
	}
}

// IsOverdue reports whether a task's due date is in the past.
func (t Task) IsOverdue() bool {
	if t.DueDate == "" || t.Status == StatusDone || t.Status == StatusArchived {
		return false
	}
	due, err := time.Parse(time.RFC3339, t.DueDate)
	if err != nil {
		return false
	}
	return due.Before(time.Now())
}

// OverdueGlyph returns "!" if overdue and not done, else the status glyph.
func (t Task) OverdueGlyph() string {
	if t.IsOverdue() {
		return "!"
	}
	return StatusGlyph(t.Status)
}

func contains(haystack []string, needle string) bool {
	for _, h := range haystack {
		if h == needle {
			return true
		}
	}
	return false
}
