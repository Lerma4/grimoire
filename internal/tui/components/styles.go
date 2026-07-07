// Package components holds the reusable visual widgets of the Grimoire TUI:
// the centralized palette/styles, the section/pane/mode enums, and the
// renderers for sidebar, list, detail, header and status bar.
package components

import "github.com/charmbracelet/lipgloss"

// Centralized palette. Dark mode by default; state is never communicated by
// color alone — textual glyphs accompany every status.
var (
	ColorBg       = lipgloss.Color("#1a1b26")
	colorPanel    = lipgloss.Color("#16161e")
	ColorText     = lipgloss.Color("#c0caf5")
	ColorMuted    = lipgloss.Color("#565f89")
	ColorAccent   = lipgloss.Color("#7aa2f7")
	ColorCyan     = lipgloss.Color("#7dcfff")
	ColorSuccess  = lipgloss.Color("#9ece6a")
	ColorWarning  = lipgloss.Color("#e0af68")
	ColorDanger   = lipgloss.Color("#f7768e")
	ColorBorder   = lipgloss.Color("#3b4261")
	ColorSelected = lipgloss.Color("#283457")
)

// Styles holds every lipgloss style used by the UI.
var Styles = struct {
	Header, StatusBar, Sidebar, List, Detail       lipgloss.Style
	Title, Muted, Accent, Success, Warning, Danger lipgloss.Style
	Selected, CursorRow, Row                       lipgloss.Style
	HelpKey, HelpDesc                              lipgloss.Style
	FocusedBorder, DimBorder                       lipgloss.Style
}{
	Header: lipgloss.NewStyle().
		Background(lipgloss.Color("#1f2335")).Foreground(ColorText).
		Padding(0, 1).Bold(true),
	StatusBar: lipgloss.NewStyle().
		Background(lipgloss.Color("#1f2335")).Foreground(ColorMuted).
		Padding(0, 1),
	Sidebar: lipgloss.NewStyle().Background(colorPanel).Foreground(ColorText).Padding(0, 1),
	List:    lipgloss.NewStyle().Background(ColorBg).Foreground(ColorText).Padding(0, 1),
	Detail:  lipgloss.NewStyle().Background(colorPanel).Foreground(ColorText).Padding(0, 1),
	Title:   lipgloss.NewStyle().Foreground(ColorCyan).Bold(true),
	Muted:   lipgloss.NewStyle().Foreground(ColorMuted),
	Accent:  lipgloss.NewStyle().Foreground(ColorAccent),
	Success: lipgloss.NewStyle().Foreground(ColorSuccess),
	Warning: lipgloss.NewStyle().Foreground(ColorWarning),
	Danger:  lipgloss.NewStyle().Foreground(ColorDanger),
	Selected: lipgloss.NewStyle().
		Background(ColorSelected).Foreground(ColorText).Bold(true),
	CursorRow:     lipgloss.NewStyle().Foreground(ColorAccent).Bold(true),
	Row:           lipgloss.NewStyle().Foreground(ColorText),
	HelpKey:       lipgloss.NewStyle().Foreground(ColorCyan).Bold(true),
	HelpDesc:      lipgloss.NewStyle().Foreground(ColorMuted),
	FocusedBorder: lipgloss.NewStyle().Foreground(ColorAccent),
	DimBorder:     lipgloss.NewStyle().Foreground(ColorBorder),
}
