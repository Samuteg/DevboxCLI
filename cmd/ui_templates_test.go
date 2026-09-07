package cmd

import (
	"strings"
	"testing"
	"text/template"

	"github.com/manifoldco/promptui"
)

// Os templates padrão do promptui v0.9.0 renderizam a linha confirmada com
// o atributo faint (SGR 2): Select usa `{{ . | faint }}` e Prompt usa
// `{{ . | faint }}` no Success. Isso deixava emojis e texto "apagados"
// após a escolha (ex: tipos de commit ✨🐛📝🎨🔁🧪🔧).
// Estes testes garantem que nossos templates não usam faint.

func renderWithFuncMap(t *testing.T, tplText string, data string) string {
	t.Helper()
	tpl, err := template.New("").Funcs(promptui.FuncMap).Parse(tplText)
	if err != nil {
		t.Fatalf("template inválido %q: %v", tplText, err)
	}
	var sb strings.Builder
	if err := tpl.Execute(&sb, data); err != nil {
		t.Fatalf("falha ao renderizar %q: %v", tplText, err)
	}
	return sb.String()
}

func TestSelectTemplateSemFaint(t *testing.T) {
	tpl := newSelectTemplates()
	if strings.Contains(tpl.Selected, "faint") {
		t.Fatalf("template Selected contém faint (ícone apagado): %q", tpl.Selected)
	}

	out := renderWithFuncMap(t, tpl.Selected, "feat:     ✨ Nova funcionalidade")
	if strings.Contains(out, "\x1b[2m") {
		t.Fatalf("linha selecionada renderiza com faint (apagada): %q", out)
	}
	if !strings.Contains(out, "✨") {
		t.Fatalf("emoji sumiu da linha selecionada: %q", out)
	}
}

func TestPromptTemplateSemFaint(t *testing.T) {
	tpl := newPromptTemplates()
	if strings.Contains(tpl.Success, "faint") {
		t.Fatalf("template Success contém faint (texto apagado): %q", tpl.Success)
	}

	out := renderWithFuncMap(t, tpl.Success, "meu-projeto 📁")
	if strings.Contains(out, "\x1b[2m") {
		t.Fatalf("linha confirmada renderiza com faint (apagada): %q", out)
	}
	if !strings.Contains(out, "meu-projeto") {
		t.Fatalf("texto sumiu da linha confirmada: %q", out)
	}
}
