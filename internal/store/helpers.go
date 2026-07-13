package store

// nullString returns nil for empty strings so they are stored as SQL NULL.
func nullString(s string) any {
	if s == "" {
		return nil
	}
	return s
}
