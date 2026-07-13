package components

import (
	"strings"
)

// Header renders the top bar: section, filter and db status.
func Header(section, filter, dbStatus string, width int) string {
	parts := []string{
		Styles.Accent.Render("✦ Grimoire"),
		Styles.Muted.Render("│"),
		"section:", section,
	}
	if filter != "" {
		parts = append(parts, Styles.Muted.Render("│"), "filter:", Styles.Warning.Render(filter))
	}
	parts = append(parts, Styles.Muted.Render("│"), dbStatus)
	return Styles.Header.Width(width).Render(joinInline(parts))
}

// StatusBar renders the bottom bar: current mode, contextual hint, message.
func StatusBar(mode string, hint, msg string, msgErr bool, width int) string {
	modePart := Styles.Accent.Bold(true).Render("[" + mode + "]")
	msgPart := ""
	if msg != "" {
		if msgErr {
			msgPart = Styles.Danger.Render(msg)
		} else {
			msgPart = Styles.Success.Render(msg)
		}
	}
	hintPart := Styles.Muted.Render(hint)
	content := joinInline([]string{modePart, hintPart, msgPart})
	return Styles.StatusBar.Width(width).Render(content)
}

// SearchBar renders the '/' search prompt with the current query.
func SearchBar(query string, width int) string {
	prompt := Styles.Accent.Render("/")
	return Styles.StatusBar.Width(width).Render(prompt + " " + query + block())
}

// CommandBar renders the ':' command prompt with the current input.
func CommandBar(input string, width int) string {
	prompt := Styles.Accent.Render(":")
	return Styles.StatusBar.Width(width).Render(prompt + " " + input + block())
}

// InputBar renders a generic inline prompt (used for quick create/link).
func InputBar(prompt, input string, width int) string {
	return Styles.StatusBar.Width(width).Render(prompt + " " + input + block())
}

func block() string { return Styles.Muted.Render("  ⏎ enter · esc cancel") }

func joinInline(parts []string) string {
	var b strings.Builder
	for i, p := range parts {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(p)
	}
	return b.String()
}
