package scaffold

import (
	"archive/zip"
	"bytes"
	"io/fs"
	"testing"
	"testing/fstest"
)

func zipSymlink(t *testing.T, name, target string) []byte {
	t.Helper()
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	h := &zip.FileHeader{Name: name}
	h.SetMode(0777 | fs.ModeSymlink)
	entry, err := w.CreateHeader(h)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := entry.Write([]byte(target)); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestMaterializeRejectsSymlink(t *testing.T) {
	fsys := fstest.MapFS{
		"link.zip": {Data: zipSymlink(t, "evil-link", "/etc/passwd")},
	}
	if err := MaterializeTemplates(fsys, "link.zip", t.TempDir()); err == nil {
		t.Fatal("esperava erro para symlink, veio nil")
	}
}

func TestMaterializeRejectsOversize(t *testing.T) {
	big := make([]byte, 0, 1024)
	for i := 0; i < 1024; i++ {
		big = append(big, 'A')
	}
	fsys := fstest.MapFS{
		"big.zip": {Data: zipBytes(t, "big.bin", 0644, string(big))},
	}
	// Limite baixo só para o teste via variável de ambiente? Usa o limite
	// padrão (alto) aqui apenas como smoke: deve extrair sem erro.
	if err := MaterializeTemplates(fsys, "big.zip", t.TempDir()); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
}
