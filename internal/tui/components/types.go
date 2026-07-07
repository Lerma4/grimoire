package components

// Section identifies a sidebar section / list source.
type Section int

const (
	SectionTasks Section = iota
	SectionNotes
	SectionToday
	SectionProjects
	SectionTags
	SectionLinks
	SectionArchive
)

func (s Section) String() string {
	switch s {
	case SectionTasks:
		return "Tasks"
	case SectionNotes:
		return "Notes"
	case SectionToday:
		return "Today"
	case SectionProjects:
		return "Projects"
	case SectionTags:
		return "Tags"
	case SectionLinks:
		return "Links"
	case SectionArchive:
		return "Archive"
	}
	return "?"
}

func (s Section) Glyph() string {
	switch s {
	case SectionTasks:
		return "○"
	case SectionNotes:
		return "✎"
	case SectionToday:
		return "★"
	case SectionProjects:
		return "▣"
	case SectionTags:
		return "#"
	case SectionLinks:
		return "↔"
	case SectionArchive:
		return "▽"
	}
	return "?"
}

// AllSections is the sidebar order.
var AllSections = []Section{
	SectionTasks,
	SectionNotes,
	SectionToday,
	SectionProjects,
	SectionTags,
	SectionLinks,
	SectionArchive,
}

// Pane identifies which column holds focus.
type Pane int

const (
	PaneSidebar Pane = iota
	PaneList
	PaneDetail
)

// Mode is the current editor mode (vim-style).
type Mode int

const (
	ModeNormal Mode = iota
	ModeInsert
	ModeCommand
	ModeSearch
	ModeHelp
)

func (m Mode) String() string {
	switch m {
	case ModeNormal:
		return "NORMAL"
	case ModeInsert:
		return "INSERT"
	case ModeCommand:
		return "COMMAND"
	case ModeSearch:
		return "SEARCH"
	case ModeHelp:
		return "HELP"
	}
	return "?"
}
