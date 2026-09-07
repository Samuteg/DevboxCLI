package scaffold

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func zipBytes(t *testing.T, name string, mode os.FileMode, content string) []byte {
	t.Helper()
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	h := &zip.FileHeader{Name: name}
	h.SetMode(mode)
	entry, err := w.CreateHeader(h)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := entry.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestZipSlipBloqueado(t *testing.T) {
	fsys := fstest.MapFS{
		"evil.zip": {Data: zipBytes(t, "../../../tmp/evil.txt", 0644, "evil")},
	}
	err := MaterializeTemplates(fsys, "evil.zip", t.TempDir())
	if err == nil {
		t.Fatal("esperava erro para caminho com ../, veio nil")
	}
}

func TestPermissaoPreservada(t *testing.T) {
	fsys := fstest.MapFS{
		"ok.zip": {Data: zipBytes(t, "script.sh", 0755, "#!/bin/sh\necho hi")},
	}
	target := t.TempDir()
	if err := MaterializeTemplates(fsys, "ok.zip", target); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	info, err := os.Stat(filepath.Join(target, "script.sh"))
	if err != nil {
		t.Fatalf("arquivo não extraído: %v", err)
	}
	if info.Mode().Perm() != 0755 {
		t.Fatalf("permissão = %o, esperava 755", info.Mode().Perm())
	}
}
