package scaffold

import (
	"io/fs"
	"os"
	"path/filepath"
)

// overrideDir returns ~/.devbox/templates/ path.
func overrideDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".devbox", "templates")
}

// LoadOverrideFS returns an OS filesystem rooted at ~/.devbox/templates/
// if the directory exists, nil otherwise.
func LoadOverrideFS() fs.FS {
	dir := overrideDir()
	if dir == "" {
		return nil
	}
	if _, err := os.Stat(dir); err != nil {
		return nil
	}
	return os.DirFS(dir)
}

// HasOverride checks if a specific template zip exists in the override directory.
func HasOverride(name string) bool {
	dir := overrideDir()
	if dir == "" {
		return false
	}
	path := filepath.Join(dir, name)
	_, err := os.Stat(path)
	return err == nil
}
