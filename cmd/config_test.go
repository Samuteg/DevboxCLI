package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestIsAllowedConfigKey(t *testing.T) {
	for _, k := range []string{"author", "default-port", "update-channel", "template-style"} {
		if !isAllowedConfigKey(k) {
			t.Fatalf("esperava chave permitida: %q", k)
		}
	}
	if isAllowedConfigKey("hacker") {
		t.Fatal("deveria rejeitar chave desconhecida")
	}
}

func TestDetectPackageManager(t *testing.T) {
	plain := t.TempDir()
	mgr, _ := detectPackageManager(plain)
	if mgr != "npm" {
		t.Fatalf("sem lockfile esperava npm, veio %q", mgr)
	}

	withLock := t.TempDir()
	if err := os.WriteFile(filepath.Join(withLock, "pnpm-lock.yaml"), []byte("lockfileVersion: 9"), 0644); err != nil {
		t.Fatal(err)
	}
	mgr, _ = detectPackageManager(withLock)
	// pnpm está instalado neste ambiente (doctor confirma); se não estiver, cai para npm.
	if _, err := exec.LookPath("pnpm"); err == nil && mgr != "pnpm" {
		t.Fatalf("com pnpm-lock + pnpm instalado esperava pnpm, veio %q", mgr)
	}
	if err := os.WriteFile(filepath.Join(withLock, "package.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
}
