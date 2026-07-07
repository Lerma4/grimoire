package domain

// Tag is a label that can be attached to tasks and/or notes.
type Tag struct {
	ID        int64
	Name      string
	Color     string
	CreatedAt string
}
