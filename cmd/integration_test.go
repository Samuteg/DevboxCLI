package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Samuteg/DevboxCLI/internal/scaffold"
)

func TestInitCreatesProject(t *testing.T) {
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	stack := scaffold.Stack{
		Name:       "Go",
		IsBackend:  true,
		Source:     "templates/golang/simple.zip",
		ExtraDirs:  []string{"cmd/api", "internal/entity"},
		RunInstall: false,
	}

	handleBackend("testproject", stack)

	// Project directory created
	if info, err := os.Stat("testproject"); err != nil || !info.IsDir() {
		t.Fatal("diretório do projeto não foi criado")
	}

	// ExtraDirs created
	for _, d := range []string{"cmd/api", "internal/entity"} {
		if info, err := os.Stat(filepath.Join("testproject", d)); err != nil || !info.IsDir() {
			t.Fatalf("ExtraDir %q não foi criado", d)
		}
	}

	// Manifest created with file list
	undoDirs, err := os.ReadDir(filepath.Join("testproject", ".devbox", "undo"))
	if err != nil || len(undoDirs) == 0 {
		t.Fatal("snapshot de undo não foi criado")
	}
	manifestPath := filepath.Join("testproject", ".devbox", "undo", undoDirs[0].Name(), "manifest.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("manifest.json não encontrado: %v", err)
	}
	var paths []string
	if err := json.Unmarshal(data, &paths); err != nil {
		t.Fatalf("manifest.json inválido: %v", err)
	}
	if len(paths) == 0 {
		t.Fatal("manifest.json vazio")
	}
}

func TestInitUndoRestoresState(t *testing.T) {
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	stack := scaffold.Stack{
		Name:      "Go",
		IsBackend: true,
		Source:    "templates/golang/simple.zip",
		ExtraDirs: []string{"cmd/api"},
	}
	handleBackend("undotest", stack)

	// Verify project exists before undo
	if _, err := os.Stat("undotest"); err != nil {
		t.Fatal("projeto não foi criado antes do undo")
	}

	// Simulate undo: read manifest and remove listed paths
	undoDirs, err := os.ReadDir(filepath.Join("undotest", ".devbox", "undo"))
	if err != nil || len(undoDirs) == 0 {
		t.Fatal("snapshot não encontrado")
	}
	manifestPath := filepath.Join("undotest", ".devbox", "undo", undoDirs[0].Name(), "manifest.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	var paths []string
	if err := json.Unmarshal(data, &paths); err != nil {
		t.Fatal(err)
	}
	// Remove in reverse order (deepest first)
	for i := len(paths) - 1; i >= 0; i-- {
		os.RemoveAll(paths[i])
	}
	os.RemoveAll(filepath.Join("undotest", ".devbox", "undo", undoDirs[0].Name()))

	// Project directory removed
	if _, err := os.Stat("undotest"); !os.IsNotExist(err) {
		t.Fatal("diretório do projeto ainda existe após undo")
	}
}

func TestAddCreatesComponent(t *testing.T) {
	tmp := t.TempDir()
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	for _, compType := range []string{"controller", "usecase", "repository", "handler"} {
		name := "user-" + compType
		path, err := scaffold.CreateComponent(compType, name, "tester")
		if err != nil {
			t.Fatalf("CreateComponent(%q) erro: %v", compType, err)
		}

		if _, err := os.Stat(path); err != nil {
			t.Fatalf("arquivo %q não foi criado", path)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(content), "User") {
			t.Fatalf("conteúdo não contém nome do componente: %s", path)
		}
	}
}

func TestCollectSafeFiles(t *testing.T) {
	tmp := t.TempDir()
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	// Init git repo
	for _, cmd := range [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@test.com"},
		{"git", "config", "user.name", "test"},
	} {
		if out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput(); err != nil {
			t.Fatalf("git init falhou: %s", out)
		}
	}

	// Create files
	os.WriteFile("main.go", []byte("package main"), 0644)
	os.WriteFile(".env", []byte("SECRET=x"), 0644)
	os.WriteFile(".env.local", []byte("KEY=y"), 0644)
	os.WriteFile("README.md", []byte("# hi"), 0644)
	os.MkdirAll("internal/secret_stuff", 0755)
	os.WriteFile("internal/secret_stuff/data.go", []byte("package secret_stuff"), 0644)

	safe, skipped := collectSafeFiles()

	// .env files excluded
	for _, s := range skipped {
		if !strings.HasPrefix(s, ".env") {
			t.Fatalf("arquivo sensível %q não deveria estar em skipped", s)
		}
	}

	// Non-sensitive files included
	safeMap := make(map[string]bool)
	for _, s := range safe {
		safeMap[s] = true
	}
	if !safeMap["main.go"] {
		t.Fatal("main.go deveria estar em safe")
	}
	if !safeMap["README.md"] {
		t.Fatal("README.md deveria estar em safe")
	}

	// .env should NOT be in safe
	if safeMap[".env"] || safeMap[".env.local"] {
		t.Fatal("arquivos .env não deveriam estar em safe")
	}
}

func TestJSONOutput(t *testing.T) {
	origJSON := jsonOutput
	jsonOutput = true
	defer func() { jsonOutput = origJSON }()

	result := JSONResult{
		Success: true,
		Command: "test",
		Message: "ok",
		Data:    map[string]string{"key": "value"},
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printJSON(result)

	w.Close()
	os.Stdout = old

	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	r.Close()
	output := strings.TrimSpace(string(buf[:n]))

	// Output must be valid JSON
	var parsed JSONResult
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("saída JSON inválida: %v\nraw: %s", err, output)
	}
	if !parsed.Success {
		t.Fatal("success deveria ser true")
	}
	if parsed.Command != "test" {
		t.Fatalf("command deveria ser 'test', veio %q", parsed.Command)
	}
}
