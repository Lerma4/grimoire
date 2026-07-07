package components

import (
	"strings"

	"github.com/charmbracelet/glamour"
)

// RenderMarkdown renders body as styled markdown for the terminal. On any
// rendering error it falls back to the plain text so the UI never breaks.
func RenderMarkdown(body string, width int) string {
	body = strings.TrimSpace(body)
	if body == "" {
		return Styles.Muted.Render("(empty)")
	}
	if width < 20 {
		return body
	}
	r, err := glamour.NewTermRenderer(
		// ponytail: avoid WithAutoStyle; it queries the terminal from View and
		// races Bubble Tea's stdin reader when the notes pane renders markdown.
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(width),
		glamour.WithEmoji(),
	)
	if err != nil {
		return body
	}
	out, err := r.Render(body)
	if err != nil {
		return body
	}
	return strings.TrimRight(out, "\n")
}
