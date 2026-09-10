package cmd

import (
	"reflect"
	"testing"
)

// Reproduz https:// — devbox init (Vite/Next) falhava com
// ERR_PNPM_NO_IMPORTER_MANIFEST_FOUND porque fmt.Sprintf anexava
// "%!(EXTRA string=<nome>)" às partes do comando sem verbo.
func TestBuildFrontendArgs(t *testing.T) {
	cases := []struct {
		name    string
		source  string
		project string
		wantCmd string
		wantArg []string
	}{
		{
			name:    "vite",
			source:  "pnpm create vite@latest %s",
			project: "advogacia",
			wantCmd: "pnpm",
			wantArg: []string{"create", "vite@latest", "advogacia"},
		},
		{
			name:    "next",
			source:  "pnpm create next-app@latest %s",
			project: "minha-loja",
			wantCmd: "pnpm",
			wantArg: []string{"create", "next-app@latest", "minha-loja"},
		},
		{
			name:    "sem placeholder",
			source:  "npm init -y",
			project: "qualquer",
			wantCmd: "npm",
			wantArg: []string{"init", "-y"},
		},
		{
			name:    "vazio",
			source:  "",
			project: "x",
			wantCmd: "",
			wantArg: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotCmd, gotArg := buildFrontendArgs(tc.source, tc.project)
			if gotCmd != tc.wantCmd {
				t.Fatalf("comando: esperado %q, veio %q", tc.wantCmd, gotCmd)
			}
			if !reflect.DeepEqual(gotArg, tc.wantArg) {
				t.Fatalf("args: esperado %q, veio %q", tc.wantArg, gotArg)
			}
		})
	}
}
