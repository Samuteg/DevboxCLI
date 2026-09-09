package validation

import "testing"

func TestIsValidProjectName(t *testing.T) {
	valid := []string{"meu-api", "app1", "a1", "my_app-2"}
	for _, v := range valid {
		if !IsValidProjectName(v) {
			t.Fatalf("esperava válido: %q", v)
		}
	}
	invalid := []string{"", "a", "../evil", "/tmp/x", "-f", "a/b", "a b", ".", "..", "a..b/c"}
	for _, v := range invalid {
		if IsValidProjectName(v) {
			t.Fatalf("esperava inválido: %q", v)
		}
	}
}

func TestIsValidComponentName(t *testing.T) {
	if !IsValidComponentName("user-profile") {
		t.Fatal("esperava válido")
	}
	for _, v := range []string{"../../evil", "", "a", "-x", "/abs"} {
		if IsValidComponentName(v) {
			t.Fatalf("esperava inválido: %q", v)
		}
	}
}

func TestValidateDirInCWD(t *testing.T) {
	if err := ValidateDirInCWD(t.TempDir()); err == nil {
		t.Fatal("esperava erro para dir fora do cwd")
	}
}
