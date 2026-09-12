package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateComponentRejectsTraversal(t *testing.T) {
	for _, evil := range []string{"../../evil", "/tmp/x", "-f", "a/b", "", "a"} {
		if _, err := CreateComponent("controller", evil, "tester"); err == nil {
			t.Fatalf("esperava erro para nome %q, veio nil", evil)
		}
	}
}

func TestCreateComponentRejectsUnknownType(t *testing.T) {
	if _, err := CreateComponent("nope", "valid-name", "tester"); err == nil {
		t.Fatal("esperava erro para tipo desconhecido")
	}
}

func TestDetectLanguage(t *testing.T) {
	tmpDir := t.TempDir()

	// Default when empty
	if lang := DetectLanguage(tmpDir); lang != LangGo {
		t.Errorf("esperado %s, veio %s", LangGo, lang)
	}

	// Node (package.json)
	os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte("{}"), 0644)
	if lang := DetectLanguage(tmpDir); lang != LangNodeJS {
		t.Errorf("esperado %s, veio %s", LangNodeJS, lang)
	}

	// Node TS (tsconfig.json)
	os.WriteFile(filepath.Join(tmpDir, "tsconfig.json"), []byte("{}"), 0644)
	if lang := DetectLanguage(tmpDir); lang != LangNodeTS {
		t.Errorf("esperado %s, veio %s", LangNodeTS, lang)
	}

	// Python
	tmpPy := t.TempDir()
	os.WriteFile(filepath.Join(tmpPy, "requirements.txt"), []byte("fastapi"), 0644)
	if lang := DetectLanguage(tmpPy); lang != LangPython {
		t.Errorf("esperado %s, veio %s", LangPython, lang)
	}

	// Ruby
	tmpRb := t.TempDir()
	os.WriteFile(filepath.Join(tmpRb, "Gemfile"), []byte("gem 'rails'"), 0644)
	if lang := DetectLanguage(tmpRb); lang != LangRuby {
		t.Errorf("esperado %s, veio %s", LangRuby, lang)
	}

	// Java
	tmpJava := t.TempDir()
	os.WriteFile(filepath.Join(tmpJava, "pom.xml"), []byte("<project></project>"), 0644)
	if lang := DetectLanguage(tmpJava); lang != LangJava {
		t.Errorf("esperado %s, veio %s", LangJava, lang)
	}
}

func TestCreateComponentPerLanguage(t *testing.T) {
	tmp := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmp)

	// In Node TS project
	os.WriteFile("tsconfig.json", []byte("{}"), 0644)
	path, err := CreateComponent("controller", "user", "tester")
	if err != nil {
		t.Fatalf("erro ao criar controller ts: %v", err)
	}
	if !strings.HasSuffix(path, ".ts") {
		t.Fatalf("esperava arquivo .ts, veio %s", path)
	}
	content, _ := os.ReadFile(path)
	if !strings.Contains(string(content), "UserController") {
		t.Fatalf("conteúdo inesperado: %s", string(content))
	}
}
