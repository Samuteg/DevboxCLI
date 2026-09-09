package scaffold

import "testing"

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
