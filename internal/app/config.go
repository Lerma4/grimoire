// Package app provides configuration and bootstrap helpers for Grimoire.
package app

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds resolved runtime configuration.
type Config struct {
	DBPath string
}

// DefaultConfig resolves configuration from the environment and OS-specific
// directories. The database path can be overridden with GRIMOIRE_DB.
func DefaultConfig() (Config, error) {
	dbPath, err := DefaultDBPath()
	if err != nil {
		return Config{}, err
	}
	return Config{DBPath: dbPath}, nil
}

// DefaultDBPath returns the path to the grimoire SQLite database, creating
// the parent directory if needed.
func DefaultDBPath() (string, error) {
	if v := os.Getenv("GRIMOIRE_DB"); v != "" {
		return v, nil
	}
	base, err := dataDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(base, 0o755); err != nil {
		return "", fmt.Errorf("create data dir %s: %w", base, err)
	}
	return filepath.Join(base, "grimoire.db"), nil
}

// dataDir returns the OS-appropriate per-user data directory for grimoire.
func dataDir() (string, error) {
	// Windows.
	if v := os.Getenv("LOCALAPPDATA"); v != "" {
		return filepath.Join(v, "grimoire"), nil
	}
	// macOS / Linux: follow XDG_DATA_HOME when set.
	if v := os.Getenv("XDG_DATA_HOME"); v != "" {
		return filepath.Join(v, "grimoire"), nil
	}
	// Fallback to ~/.local/share/grimoire (Linux) or ~/Library/Application Support/grimoire (macOS).
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	switch runtimeOS() {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "grimoire"), nil
	default:
		return filepath.Join(home, ".local", "share", "grimoire"), nil
	}
}
