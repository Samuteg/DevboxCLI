package system

import (
	"testing"
	"time"
)

func TestExecuteRejectsInvalidName(t *testing.T) {
	if err := Execute("bad name!", nil, ""); err == nil {
		t.Fatal("esperava erro para nome inválido")
	}
	if err := ExecuteSilentWithTimeout("bad;name", nil, "", time.Second); err == nil {
		t.Fatal("esperava erro para nome inválido")
	}
}

func TestExecuteSilentTimeout(t *testing.T) {
	// sleep 5 com timeout de 300ms deve falhar rápido (não pendurar).
	if err := ExecuteSilentWithTimeout("sleep", []string{"5"}, "", 300*time.Millisecond); err == nil {
		t.Fatal("esperava timeout, veio nil")
	}
}
