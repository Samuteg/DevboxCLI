package scaffold

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func MaterializeTemplates(fsys fs.FS, sourceDir, targetRoot string) error {
	return fs.WalkDir(fsys, sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}

		targetPath := filepath.Join(targetRoot, relPath)
		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		finalPath := strings.TrimSuffix(targetPath, ".tmpl")
		if err := os.MkdirAll(filepath.Dir(finalPath), 0755); err != nil {
			return err
		}

		content, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}

		return os.WriteFile(finalPath, content, 0644)
	})
}
