package domain

// TaskNoteLink is a many-to-many association between a task and a note.
type TaskNoteLink struct {
	TaskID       int64
	NoteID       int64
	RelationType string
	CreatedAt    string
}

// RelationReference is the default MVP relation type.
const RelationReference = "reference"
