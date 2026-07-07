package store

// nullString returns nil for empty strings so they are stored as SQL NULL.
func nullString(s string) any {
	if s == "" {
		return nil
	}
	return s
}

// nullID returns nil for zero ids so they are stored as SQL NULL.
func nullID(id int64) any {
	if id == 0 {
		return nil
	}
	return id
}
