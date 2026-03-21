package system

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

var validCommandName = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)

func Execute(name string, args []string, dir string) error {
	if !validCommandName.MatchString(name) {
		return fmt.Errorf("nome de comando inválido: %s", name)
	}

	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("erro ao executar %s: %w", name, err)
	}

	return nil
}
