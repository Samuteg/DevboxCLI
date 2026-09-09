package system

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"time"
)

var validCommandName = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)

func Execute(name string, args []string, dir string) error {
	return ExecuteWithTimeout(name, args, dir, 10*time.Minute)
}

// ExecuteWithTimeout executa com timeout e valida dir (sem shell, sem injeção).
func ExecuteWithTimeout(name string, args []string, dir string, timeout time.Duration) error {
	if !validCommandName.MatchString(name) {
		return fmt.Errorf("nome de comando inválido: %s", name)
	}
	for _, a := range args {
		if len(a) > 0 && a[0] == 0 {
			return fmt.Errorf("argumento inválido")
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return ExecuteContext(ctx, name, args, dir)
}

// ExecuteContext é a base com ctx (usada por doctor/kill com timeout curto).
func ExecuteContext(ctx context.Context, name string, args []string, dir string) error {
	if !validCommandName.MatchString(name) {
		return fmt.Errorf("nome de comando inválido: %s", name)
	}
	cmd := exec.CommandContext(ctx, name, args...)
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

func ExecuteSilent(name string, args []string, dir string) error {
	return ExecuteSilentWithTimeout(name, args, dir, 10*time.Minute)
}

// ExecuteSilentWithTimeout captura output e devolve tail no erro (sem pendurar CLI).
func ExecuteSilentWithTimeout(name string, args []string, dir string, timeout time.Duration) error {
	if !validCommandName.MatchString(name) {
		return fmt.Errorf("nome de comando inválido: %s", name)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	if err := cmd.Run(); err != nil {
		out := buf.String()
		if len(out) > 2000 {
			out = out[len(out)-2000:]
		}
		return fmt.Errorf("erro ao executar %s: %w\n%s", name, err, out)
	}
	return nil
}
