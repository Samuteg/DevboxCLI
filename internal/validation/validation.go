package validation

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var nameRe = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-_]{1,63}$`)

// IsValidProjectName aceita 2-64 chars, alfanumérico + - _, sem path/flag injection.
func IsValidProjectName(s string) bool {
	if !nameRe.MatchString(s) {
		return false
	}
	if strings.Contains(s, "..") || strings.Contains(s, "/") || strings.Contains(s, "\\") {
		return false
	}
	return true
}

// IsValidComponentName usa a mesma regra de projeto.
func IsValidComponentName(s string) bool {
	return IsValidProjectName(s)
}

// ValidateDirInCWD garante que dir resolve para dentro do cwd (anti-traversal).
func ValidateDirInCWD(dir string) error {
	if dir == "" {
		return errInvalid("diretório vazio")
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(cwd, abs)
	if err != nil {
		return err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(dir) && !strings.HasPrefix(abs, cwd+string(filepath.Separator)) && abs != cwd {
		// Permite subdir do cwd; bloqueia /tmp, .., absoluto fora
		if abs != cwd && !strings.HasPrefix(abs, cwd+string(filepath.Separator)) {
			return errInvalid("diretório fora do projeto: " + dir)
		}
	}
	return nil
}

type validationError string

func errInvalid(s string) error { return validationError(s) }

func (e validationError) Error() string { return string(e) }
