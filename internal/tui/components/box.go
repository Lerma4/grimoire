package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// boxed renders body inside a titled border of the given width. When focused
// the border uses the accent color, otherwise a dim gray. The title appears as
// the first line of content (lipgloss v1 has no border-title method).
func boxed(title string, focused bool, width int, body string) string {
	border := lipgloss.NormalBorder()
	style := lipgloss.NewStyle().Border(border, true).
		BorderBackground(ColorBg).Padding(0, 1)
	if focused {
		style = style.BorderForeground(ColorAccent).Foreground(ColorText)
	} else {
		style = style.BorderForeground(ColorBorder)
	}
	style = style.Width(clampWidth(width))
	if title != "" {
		head := Styles.Muted.Render(title)
		body = head + "\n" + Styles.Muted.Render(strings.Repeat("─", clampWidth(width)-4)) + "\n" + body
	}
	return style.Render(body)
}

// clampWidth keeps column widths sane on tiny terminals.
func clampWidth(w int) int {
	if w < 8 {
		return 8
	}
	return w
}

// Box is the exported wrapper around boxed for callers outside the package.
func Box(title string, focused bool, width int, body string) string {
	return boxed(title, focused, width, body)
}

// BoxCenter renders body centered in a titled box spanning width.
func BoxCenter(title string, width int, body string) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true).
		BorderForeground(ColorAccent).
		BorderBackground(ColorBg).
		Padding(0, 2).
		Width(clampWidth(width))
	return style.Render(body)
}
