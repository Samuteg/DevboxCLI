package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Samuteg/DevboxCLI/internal/scaffold"
)

// Garante que os templates Node embutidos continuam magros e válidos:
// sem node_modules, com package.json + src. Evita regressão aos 153M.
func TestNodeTemplatesSlim(t *testing.T) {
	for _, zip := range []string{"templates/node/ts.zip", "templates/node/js.zip"} {
		target := t.TempDir()
		if err := scaffold.MaterializeTemplates(templatesFS, zip, target); err != nil {
			t.Fatalf("%s: erro ao extrair: %v", zip, err)
		}
		root := filepath.Join(target, strings.TrimSuffix(filepath.Base(zip), ".zip")) // ts/ ou js/
		if _, err := os.Stat(filepath.Join(root, "package.json")); err != nil {
			t.Fatalf("%s: package.json ausente: %v", zip, err)
		}
		if _, err := os.Stat(filepath.Join(root, "src")); err != nil {
			t.Fatalf("%s: src/ ausente: %v", zip, err)
		}
		if _, err := os.Stat(filepath.Join(root, "node_modules")); err == nil {
			t.Fatalf("%s: node_modules não deveria estar no template", zip)
		}
	}
}
