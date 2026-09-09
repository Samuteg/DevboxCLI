package cmd

import "testing"

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
