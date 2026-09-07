# Devbox CLI — Plano de Refatoração Completa

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Corrigir bugs críticos, unificar identidade (nome/versão/idioma), migrar de `promptui` (unmaintained) para `huh`, adicionar CI e testes, e limpar dívida técnica.

**Architecture:** CLI Cobra com `cmd/` (comandos) e `internal/` (scaffold/system). Templates embutidos via `//go:embed` (`cmd/templates/*.zip`). UI via lipgloss + fatih/color. A migração promptui→huh preserva o mesmo fluxo interativo (Select + Input + Confirm). Cada task termina com build+vet+test verdes e um commit próprio.

**Tech Stack:** Go 1.25+, Cobra, lipgloss, huh (substitui promptui), Go embed, GitHub Actions.

**Spec:** Relatório de refatoração com 26 itens (apresentado em chat na sessão de análise). Mapeamento item→task: (1)→T5, (2)→T6, (3)→T7, (4)→T8, (5)→T5, (6)→T2, (7)→T3, (8)→T15, (9)→T25, (10)→T24, (11)→T19, (12)→T26+T1(teste base), (13)→T17, (14)→T16, (15)→T1, (16)→T21, (17)→T24, (18)→T18, (19)→T20, (20)→T4, (21)→T9–T14, (22)→T22, (23)→T23, (24)→T28, (25)→fora de escopo desta sequência (ver Ruling R5), (26)→higiene contínua, não é task.

## Global Constraints

- Nome oficial do binário é `devbox` (uma palavra, minúsculo) em help, README, exemplos e mensagens.
- Idioma único PT-BR em todas as mensagens, help e README.
- Go 1.25+ (go.mod); `gofmt -l` vazio; `go build ./...`, `go vet ./...`, `go test ./...` verdes ao fim de cada task.
- Um commit por task, mensagem em PT-BR ou conventional-commit curto.
- Não reintroduzir `promptui` após T14; não usar `faint` (SGR 2) em texto/emoji confirmado.
- Não quebrar o formato de saída final dos comandos (tree + success box).

---

## Mapa de Arquivos

| Arquivo | Responsabilidade |
|---------|-----------------|
| `cmd/colors.go` (novo, T1) | Tokens de cor centralizados |
| `cmd/version.go` (novo, T2) | Única fonte de versão (`Version`) |
| `cmd/ui_utils.go` | Funções de UI (spinner, log, tree, banner, steps) |
| `cmd/ui_templates_test.go` | Testes de templates (removido em T14) |
| `cmd/init.go` | Comando init — ordenar stacks (T7), huh (T10/T13), frontend (T22) |
| `cmd/commit_wizard.go` | Wizard — huh (T11), badges (T20), git add com erro (T25) |
| `cmd/cleanup.go` | Limpeza — confirmação + fix iconFail (T24) |
| `cmd/config.go` | Config — semântica warn (T18), texto Salvar (T4) |
| `cmd/doctor.go` | Diagnóstico — cores (T1), UTF-8 (T19) |
| `cmd/kill.go` | Kill — falso-sucesso + SIGTERM (T8), warn (T18) |
| `cmd/help.go` | Help — cores (T1), aliases/examples (T21) |
| `cmd/update.go` | Update — Version (T2), huh confirm (T12) |
| `cmd/root.go` | Root — defaults incondicionais (T6), texto PT-BR (T4) |
| `internal/scaffold/templates.go` | Extração zip — zip-slip + permissões (T5) |
| `internal/scaffold/templates_test.go` (novo, T5) | Testes zip-slip e permissão |
| `internal/scaffold/stacks.go` | Stacks — acento, case, source morto (T23) |
| `internal/system/command.go` | Executor — variante silenciosa (T16) |
| `winres/winres.json` | Metadata Windows (T28) |
| `README.md` | Nome devbox + stack Java (T3/T23/T27) |
| `.github/workflows/ci.yml` (novo, T26) | CI: build, vet, test, fmt |
| `go.mod`/`go.sum` | huh entra (T9), promptui sai (T14) |

---

## Fase 1 — Fundação

### Task 1: Tokens de cor centralizados

**Files:**
- Create: `cmd/colors.go`
- Modify: `cmd/ui_utils.go`, `cmd/doctor.go`, `cmd/help.go`
- Test: `go build ./... && go vet ./... && go test ./...`

**Interfaces:**
- Consumes: nada (primeira task).
- Produces: `ColorPrimary`, `ColorSecondary`, `ColorSuccess`, `ColorError`, `ColorWarning`, `ColorMuted` (=242), `ColorSubtle` (=244), `ColorWhite`, `ColorFeat`, `ColorFix`, `ColorDocs`, `ColorRefactor`, `ColorStyle`, `ColorTest`, `ColorChore` — usados por T18 e T20.

- [ ] **Step 1: Criar `cmd/colors.go`**

```go
package cmd

import "github.com/charmbracelet/lipgloss"

// Paleta centralizada do Devbox CLI.
var (
	ColorPrimary   = lipgloss.Color("#7D56F4")
	ColorSecondary = lipgloss.Color("#00ADD8")
	ColorSuccess   = lipgloss.Color("#27AE60")
	ColorError     = lipgloss.Color("#E74C3C")
	ColorWarning   = lipgloss.Color("#F1C40F")
	ColorMuted     = lipgloss.Color("242")
	ColorSubtle    = lipgloss.Color("244")
	ColorWhite     = lipgloss.Color("#FFF")

	ColorFeat     = lipgloss.Color("#A3BE8C")
	ColorFix      = lipgloss.Color("#BF616A")
	ColorDocs     = lipgloss.Color("#81A1C1")
	ColorRefactor = lipgloss.Color("#B48EAD")
	ColorStyle    = lipgloss.Color("#EBCB8B")
	ColorTest     = lipgloss.Color("#88C0D0")
	ColorChore    = lipgloss.Color("#616E88")
)
```

- [ ] **Step 2: Migrar `cmd/ui_utils.go`** — substituir `primaryColor`→`ColorPrimary`, `secondaryColor`→`ColorSecondary`, `successColor`→`ColorSuccess`, `errorColor`/`deleteColor`/`red`→`ColorError`, `neutralColor`→`ColorMuted`, `featColor`→`ColorFeat`, `fixColor`→`ColorFix`, `docsColor`→`ColorDocs`, `refactorColor`→`ColorRefactor`, `addComponentColor`→`ColorSecondary`, `addDirColor`→`ColorWarning`, `killColor`→`ColorFix`, `targetColor`→`ColorStyle`. Trocar `lipgloss.Color("240")` por `ColorSubtle` onde for texto/ícone apagado.
- [ ] **Step 3: Migrar `cmd/doctor.go`** — `headerColor`→`ColorPrimary`, `subtleColor`→`ColorMuted`, `failColor`→`ColorError`.
- [ ] **Step 4: Migrar `cmd/help.go`** — estilos de seção/comando/descrição/flag para os tokens.
- [ ] **Step 5: Verificar**

Run: `gofmt -l cmd/ internal/ ; go build ./... && go vet ./... && go test ./...`
Expected: sem saída do gofmt; build/vet/test OK.

- [ ] **Step 6: Commit**

```bash
git add cmd/colors.go cmd/ui_utils.go cmd/doctor.go cmd/help.go
git commit -m "refactor: centraliza paleta de cores em cmd/colors.go"
```

### Task 2: Versão unificada

**Files:**
- Create: `cmd/version.go`
- Modify: `cmd/update.go`, `cmd/ui_utils.go`
- Test: `go build ./... && go vet ./... && go test ./...`

**Interfaces:**
- Consumes: nada.
- Produces: `Version` (string, ex. `"1.0.1"`) usada por T28.

- [ ] **Step 1: Criar `cmd/version.go`**

```go
package cmd

// Version é a única fonte da versão do binário.
const Version = "1.0.1"
```

- [ ] **Step 2: Usar em `cmd/update.go`** — remover `const version = "1.0.1"` e trocar todos os usos de `version` por `Version`.
- [ ] **Step 3: Usar no banner `cmd/ui_utils.go`** — trocar `"v1.0.0 • Automation Tool"` por `fmt.Sprintf("v%s • Devbox CLI", Version)`.
- [ ] **Step 4: Verificar**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add cmd/version.go cmd/update.go cmd/ui_utils.go
git commit -m "refactor: unifica versão em cmd.Version"
```

### Task 3: Nome do binário unificado (`devbox`)

**Files:**
- Modify: `README.md`, `winres/winres.json`, `.goreleaser.yaml` (só se o binário gerado não se chamar `devbox`)
- Test: `grep -rn "DevboxCLI" README.md` deve retornar só o path de `go install` e o nome do repo.

**Interfaces:**
- Consumes: nada.
- Produces: convenção de nome usada por T27.

- [ ] **Step 1: README** — substituir ocorrências de `DevboxCLI` por `devbox` nos comandos e exemplos. Manter `github.com/Samuteg/DevboxCLI` no path de `go install` e no nome do repositório.
- [ ] **Step 2: winres.json (só identity/info)** — preencher `identity.name: "Devbox"` e `info` com `FileDescription: "Devbox CLI"`, `ProductName: "Devbox"` (o restante do arquivo fica para T28).
- [ ] **Step 3: Verificar binário do goreleaser** — ler `.goreleaser.yaml`; se o nome do binário não for `devbox`, ajustar o campo `name`.
- [ ] **Step 4: Verificar**

Run: `grep -rn "DevboxCLI" README.md cmd/root.go install.sh | grep -v "github.com/Samuteg/DevboxCLI"`
Expected: sem saída (só o path de import permanece).

- [ ] **Step 5: Commit**

```bash
git add README.md winres/winres.json .goreleaser.yaml
git commit -m "docs: unifica nome do binário para devbox"
```

### Task 4: Idioma unificado PT-BR

**Files:**
- Modify: `cmd/root.go`, `cmd/kill.go`, `cmd/config.go`
- Test: `grep -rni "ficheiro\|liberta\|Guardar" cmd/ | grep -v "_test.go"` vazio.

**Interfaces:**
- Consumes: nada.
- Produces: nada (texto apenas).

- [ ] **Step 1: `cmd/root.go`** — `"ficheiro de config (por omissão é $HOME/.devbox.yaml)"` → `"arquivo de configuração (padrão: $HOME/.devbox.yaml)"`.
- [ ] **Step 2: `cmd/kill.go`** — `"Porta libertada"` → `"Porta liberada"`.
- [ ] **Step 3: `cmd/config.go`** — `"Guardar Configuração"` → `"Salvar Configuração"`.
- [ ] **Step 4: Varredura final** — `grep -rni "ficheiro\|liberta\|Guardar\|por omissão" cmd/ internal/` e corrigir achados.
- [ ] **Step 5: Verificar**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add cmd/root.go cmd/kill.go cmd/config.go
git commit -m "refactor: unifica mensagens em PT-BR"
```

---

## Fase 2 — Bugs críticos

### Task 5: Zip-slip + permissões na extração

**Files:**
- Modify: `internal/scaffold/templates.go`
- Create: `internal/scaffold/templates_test.go`
- Test: `go test ./internal/scaffold/ -v`

**Interfaces:**
- Consumes: nada.
- Produces: `MaterializeTemplates` segura (contrato inalterado: `func(fsys fs.FS, sourceZip, targetRoot string) error`).

- [ ] **Step 1: Escrever o teste que falha (TDD)** — criar `internal/scaffold/templates_test.go`:

```go
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
```

- [ ] **Step 2: Rodar e ver falhar**

Run: `go test ./internal/scaffold/ -run "TestZipSlip|TestPermissao" -v`
Expected: FAIL (zip-slip passa batido; permissão sai 0644).

- [ ] **Step 3: Implementar** — em `internal/scaffold/templates.go`, adicionar imports `fmt` e `strings`, e no loop:

```go
for _, file := range zipReader.File {
	cleanName := filepath.Clean(file.Name)
	if cleanName == "." || strings.HasPrefix(cleanName, "..") || filepath.IsAbs(cleanName) {
		return fmt.Errorf("caminho inseguro no zip: %s", file.Name)
	}
	targetPath := filepath.Join(targetRoot, cleanName)
	// ... resto igual, exceto a criação do arquivo:
	mode := file.Mode()
	if mode == 0 {
		mode = 0644
	}
	outFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	// ...
}
```

- [ ] **Step 4: Rodar e ver passar**

Run: `go test ./internal/scaffold/ -v`
Expected: PASS nos dois testes.

- [ ] **Step 5: Verificar tudo**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add internal/scaffold/templates.go internal/scaffold/templates_test.go
git commit -m "fix: bloqueia zip-slip e preserva permissão na extração"
```

### Task 6: Defaults do viper incondicionais

**Files:**
- Modify: `cmd/root.go`
- Test: `go build ./... && go vet ./... && go test ./...`

**Interfaces:**
- Consumes: nada.
- Produces: `viper.GetString("author"/"default-port"/...)` nunca vazio (usado por T8 indiretamente).

- [ ] **Step 1: Mover `SetDefault` para fora do `else`** em `cmd/root.go:initConfig`:

```go
func initConfig() {
	viper.SetDefault("author", "Devbox User")
	viper.SetDefault("default-port", "8080")
	viper.SetDefault("update-channel", "stable")
	viper.SetDefault("template-style", "clean")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".devbox")
		viper.AutomaticEnv()
		viper.SetEnvPrefix("DEVBOX")
		configPath := filepath.Join(home, ".devbox.yaml")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			os.WriteFile(configPath, []byte(""), 0644)
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("DEVBOX")

	_ = viper.ReadInConfig()
}
```

- [ ] **Step 2: Verificar**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: PASS.

- [ ] **Step 3: Commit**

```bash
git add cmd/root.go
git commit -m "fix: defaults do viper valem com --config customizado"
```

### Task 7: Ordenar stacks no select

**Files:**
- Modify: `cmd/init.go`
- Test: `go build ./... && go vet ./... && go test ./...`

**Interfaces:**
- Consumes: nada.
- Produces: `options` ordenado (base para T10, que reescreve os helpers no mesmo arquivo).

- [ ] **Step 1: Ordenar** — em `cmd/init.go:runInit`, após montar `options`, adicionar `slices.Sort(options)` (importar `"slices"`).
- [ ] **Step 2: Verificar**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: PASS.

- [ ] **Step 3: Commit**

```bash
git add cmd/init.go
git commit -m "fix: ordena stacks no select do init"
```

### Task 8: Kill sem falso-sucesso + SIGTERM antes de SIGKILL

**Files:**
- Modify: `cmd/kill.go`
- Test: `go build ./... && go vet ./... && go test ./...`

**Interfaces:**
- Consumes: `printStep` (forma atual; T18 muda "todo"→"warn" depois — não fazer aqui).
- Produces: nada.

- [ ] **Step 1: Reescrever `killUnix`**:

```go
func killUnix(port string) {
	printStep("active", "Buscando PID via lsof...")

	cmdFind := exec.Command("lsof", "-t", "-i:"+port)
	out, err := cmdFind.Output()

	if err != nil || len(strings.TrimSpace(string(out))) == 0 {
		printStep("todo", "Nenhum processo ativo encontrado na porta")
		return
	}

	pid := strings.TrimSpace(string(out))

	validPid := regexp.MustCompile(`^[0-9]+$`)
	if !validPid.MatchString(pid) {
		HandleError(fmt.Errorf("PID retornado é suspeito: %q", pid), "Segurança")
		return
	}

	printStep("active", fmt.Sprintf("Encerrando processo %s", pidStyle.Render("("+pid+")")))

	killed := 0
	for _, p := range strings.Fields(pid) {
		if exec.Command("kill", p).Run() == nil {
			killed++
			continue
		}
		if err := exec.Command("kill", "-9", p).Run(); err != nil {
			HandleError(err, "Falha ao matar processo "+p)
			return
		}
		killed++
	}

	if killed == 0 {
		printStep("todo", "Nenhum processo foi terminado")
		return
	}

	printStep("done", fmt.Sprintf("Processo(s) terminado(s): %d", killed))
	showKillFinal(port)
}
```

- [ ] **Step 2: Verificar**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: PASS.

- [ ] **Step 3: Commit**

```bash
git add cmd/kill.go
git commit -m "fix: kill valida PID, tenta SIGTERM e evita falso-sucesso"
```

---

## Fase 3 — Migração promptui → huh

### Task 9: Adicionar dependência huh

**Files:**
- Modify: `go.mod`, `go.sum`
- Test: `go build ./...`

**Interfaces:**
- Consumes: nada.
- Produces: `github.com/charmbracelet/huh` disponível (usado por T10–T13, T15, T24).

- [ ] **Step 1: Adicionar**

Run: `go get github.com/charmbracelet/huh@latest && go mod tidy`
Expected: sem erro; `huh` em go.mod.

- [ ] **Step 2: Verificar**

Run: `go build ./...`
Expected: PASS.

- [ ] **Step 3: Commit**

```bash
git add go.mod go.sum
git commit -m "chore: adiciona dependência huh"
```

### Task 10: Migrar `cmd/init.go` para huh

**Files:**
- Modify: `cmd/init.go`
- Test: `go build ./... && go vet ./... && go test ./...`

**Interfaces:**
- Consumes: `huh` (T9), `options` ordenado (T7, manter o `slices.Sort`).
- Produces: helpers `promptInput`, `promptSelect`, `promptVariant` baseados em huh (usados por `add.go` sem mudança).

- [ ] **Step 1: Trocar imports** — remover `"github.com/manifoldco/promptui"`, adicionar `"github.com/charmbracelet/huh"`.
- [ ] **Step 2: Reescrever os 3 helpers**:

```go
func promptInput(label, errMsg string, minLen int) string {
	var value string
	_ = huh.NewInput().
		Title(label).
		Validate(func(s string) error {
			if len(s) < minLen {
				return errors.New(errMsg)
			}
			return nil
		}).
		Value(&value).
		Run()
	return value
}

func promptSelect(label string, items []string) string {
	var value string
	options := make([]huh.Option[string], len(items))
	for i, item := range items {
		options[i] = huh.NewOption(item, item)
	}
	_ = huh.NewSelect[string]().
		Title(label).
		Options(options...).
		Value(&value).
		Run()
	return value
}

func promptVariant(variants []scaffold.Variant) scaffold.Variant {
	var value string
	options := make([]huh.Option[string], len(variants))
	for i, v := range variants {
		options[i] = huh.NewOption(v.Name, v.Name)
	}
	_ = huh.NewSelect[string]().
		Title("Escolha uma variante").
		Options(options...).
		Value(&value).
		Run()
	for _, v := range variants {
		if v.Name == value {
			return v
		}
	}
	os.Exit(1)
	return scaffold.Variant{}
}
```

Manter os emojis dos labels (`📁`, `💻`, `🛠️`, `⚡`) e o `slices.Sort(options)` do `runInit`.

- [ ] **Step 3: Verificar**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: PASS.

- [ ] **Step 4: Commit**

```bash
git add cmd/init.go
git commit -m "refactor: migra prompts do init para huh"
```

### Task 11: Migrar `cmd/commit_wizard.go` para huh

**Files:**
- Modify: `cmd/commit_wizard.go`
- Test: `go build ./... && go vet ./... && go test ./...`

**Interfaces:**
- Consumes: `huh` (T9). Preserva emojis 🔁/🧪 e `showCommitSuccess`.
- Produces: nada.

- [ ] **Step 1: Trocar imports** — remover `promptui`, adicionar `huh`.
- [ ] **Step 2: Reescrever os 4 prompts** (tipo/select com os 7 items e valores `feat`/`fix`/…; escopo/input; descrição/input com validate ≥3; confirmação com `huh.NewConfirm`; cancelamento imprime mensagem e retorna). Manter o restante do fluxo (`git add .`, `git commit`, `showCommitSuccess`) idêntico.
- [ ] **Step 3: Verificar**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: PASS.

- [ ] **Step 4: Commit**

```bash
git add cmd/commit_wizard.go
git commit -m "refactor: migra commit wizard para huh"
```

### Task 12: Migrar `cmd/update.go` para huh

**Files:**
- Modify: `cmd/update.go`
- Test: `go build ./... && go vet ./... && go test ./...`

**Interfaces:**
- Consumes: `huh` (T9), `Version` (T2, manter).
- Produces: nada.

- [ ] **Step 1: Trocar imports** — remover `promptui`, adicionar `huh`.
- [ ] **Step 2: Trocar o confirm** por `huh.NewConfirm().Title("Deseja baixar e instalar agora?")`; se negado, imprime aviso e retorna (mesmo texto atual).
- [ ] **Step 3: Verificar**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: PASS.

- [ ] **Step 4: Commit**

```bash
git add cmd/update.go
git commit -m "refactor: migra confirmação do update para huh"
```

### Task 13: Confirmar migração total do init

**Files:** nenhum (verificação).

- [ ] **Step 1: Grep** — `grep -rn "promptui" cmd/init.go` deve retornar vazio. Se sobrar algo, migrar seguindo o padrão da T10.
- [ ] **Step 2: Verificar**

Run: `go build ./...`
Expected: PASS (sem commit se nada mudou; se mudou, commitar como `refactor: completa migração do init para huh`).

### Task 14: Remover promptui

**Files:**
- Modify: `cmd/ui_utils.go`, `go.mod`, `go.sum`
- Delete: `cmd/ui_templates_test.go` (testava os workarounds do promptui)
- Test: `go build ./... && go vet ./... && go test ./...`

**Interfaces:**
- Consumes: T10–T13 completas.
- Produces: zero referências a promptui (base para T15).

- [ ] **Step 1: Grep geral** — `grep -rn "promptui" cmd/ internal/` deve mostrar só `ui_utils.go` (helpers) e mais nada.
- [ ] **Step 2: Remover de `ui_utils.go`** — apagar `newSelectTemplates`, `newPromptTemplates` e o import do promptui. Apagar `cmd/ui_templates_test.go`.
- [ ] **Step 3: Limpar módulo**

Run: `go mod tidy && go build ./... && go vet ./... && go test ./...`
Expected: `grep -rn promptui go.mod` vazio; tudo PASS.

- [ ] **Step 4: Commit**

```bash
git add cmd/ui_utils.go cmd/ui_templates_test.go go.mod go.sum
git commit -m "chore: remove promptui (migração para huh concluída)"
```

---

## Fase 4 — UX e robustez

### Task 15: Ctrl+C limpo nos prompts huh

**Files:**
- Modify: `cmd/ui_utils.go`, `cmd/init.go`, `cmd/commit_wizard.go`, `cmd/update.go`, `cmd/cleanup.go` (se usar huh na T24)
- Test: `go build ./... && go vet ./... && go test ./...`

**Interfaces:**
- Consumes: prompts huh (T10–T12).
- Produces: helper de cancelamento reutilizável.

- [ ] **Step 1: Helper em `cmd/ui_utils.go`**:

```go
// promptAborted diz se o usuário cancelou (Ctrl+C) e já imprime aviso.
func promptAborted(err error) bool {
	if err == huh.ErrUserAborted {
		fmt.Println("  Operação cancelada.")
		return true
	}
	return false
}
```

(importar `huh` em ui_utils.)

- [ ] **Step 2: Usar nos helpers** — `promptInput`, `promptSelect`, `promptVariant` (`init.go`): se `Run()` retornar erro e `promptAborted(err)`, sair com `os.Exit(0)` após a mensagem (comportamento de CLI ao cancelar). Demais erros: `HandleError` + `os.Exit(1)`.
- [ ] **Step 3: Usar no wizard e update** — mesmo padrão nos `Run()` de `commit_wizard.go` e `update.go`.
- [ ] **Step 4: Verificar**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add cmd/ui_utils.go cmd/init.go cmd/commit_wizard.go cmd/update.go
git commit -m "feat: cancelamento limpo com Ctrl+C nos prompts"
```

### Task 16: Executor de comandos único

**Files:**
- Modify: `internal/system/command.go`, `cmd/init.go`, `cmd/ui_utils.go`
- Test: `go build ./... && go vet ./... && go test ./...`

**Interfaces:**
- Consumes: nada.
- Produces: `system.Execute` + `system.ExecuteSilent` (usado por T22).

- [ ] **Step 1: Adicionar `ExecuteSilent` em `internal/system/command.go`** (mesma validação de nome, `cmd.Run()` sem stdio herdado).
- [ ] **Step 2: Trocar uso em `cmd/init.go`** — `ExecuteCommandSilent` → `system.ExecuteSilent` (adicionar import do pacote system).
- [ ] **Step 3: Remover `ExecuteCommandSilent` e `NewSpinner`? Não** — remover só `ExecuteCommandSilent` de `ui_utils.go` (spinner fica).
- [ ] **Step 4: Verificar**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/system/command.go cmd/init.go cmd/ui_utils.go
git commit -m "refactor: unifica execução de comandos em system"
```

### Task 17: Funde as caixas de sucesso

**Files:**
- Modify: `cmd/ui_utils.go`, `cmd/init.go`, `cmd/add.go`, `cmd/update.go`
- Test: `go build ./... && go vet ./... && go test ./...`

**Interfaces:**
- Consumes: `Color*` (T1).
- Produces: só `ShowSuccessBox` existe.

- [ ] **Step 1: Escolher e fundir** — manter o corpo de `ShowSuccessBox` (versão simples com `IconStep`); apagar `showSuccessBox` (versão com borda lipgloss) e trocar suas chamadas em `init.go:handleBackend` para `ShowSuccessBox`.
- [ ] **Step 2: Grep** — `grep -rn "showSuccessBox" cmd/` deve mostrar só a definição de `ShowSuccessBox`.
- [ ] **Step 3: Verificar**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: PASS.

- [ ] **Step 4: Commit**

```bash
git add cmd/ui_utils.go cmd/init.go cmd/add.go cmd/update.go
git commit -m "refactor: funde caixas de sucesso duplicadas"
```

### Task 18: Semântica warn nos steps

**Files:**
- Modify: `cmd/ui_utils.go`, `cmd/kill.go`, `cmd/commit_wizard.go`, `cmd/config.go`
- Test: `go build ./... && go vet ./... && go test ./...`

**Interfaces:**
- Consumes: `ColorWarning` (T1).
- Produces: status `"warn"` em `printStep`.

- [ ] **Step 1: Novo case em `printStep`** (`ui_utils.go`): `case "warn"` com ícone `⚠` em `ColorWarning` e texto na mesma cor.
- [ ] **Step 2: Trocar usos de erro parcial**: `kill.go` (2 ocorrências de "Nenhum processo…"/"Porta…"), `commit_wizard.go` ("Nada para commitar…"), `config.go` ("Nenhuma configuração…"): `"todo"` → `"warn"`.
- [ ] **Step 3: Verificar**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: PASS.

- [ ] **Step 4: Commit**

```bash
git add cmd/ui_utils.go cmd/kill.go cmd/commit_wizard.go cmd/config.go
git commit -m "refactor: status warn para falhas parciais nos steps"
```

### Task 19: Truncamento por runas no doctor

**Files:**
- Modify: `cmd/doctor.go`
- Test: `go build ./... && go vet ./... && go test ./...`

**Interfaces:**
- Consumes: nada.
- Produces: nada.

- [ ] **Step 1: Corrigir `getVersion`**:

```go
v := strings.Split(string(out), "\n")[0]
if runes := []rune(v); len(runes) > 15 {
	return "v" + string(runes[:15]) + "..."
}
return strings.TrimSpace(v)
```

- [ ] **Step 2: Verificar**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: PASS.

- [ ] **Step 3: Commit**

```bash
git add cmd/doctor.go
git commit -m "fix: trunca versão por runas no doctor"
```

### Task 20: Badge colorido para todos os tipos de commit

**Files:**
- Modify: `cmd/commit_wizard.go`
- Test: `go build ./... && go vet ./... && go test ./...`

**Interfaces:**
- Consumes: `ColorStyle`, `ColorTest`, `ColorChore` (T1).
- Produces: nada.

- [ ] **Step 1: Completar o switch em `showCommitSuccess`**: `style`→`ColorStyle`, `test`→`ColorTest`, `chore`→`ColorChore`, `default`→cinza `ColorMuted`.
- [ ] **Step 2: Verificar**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: PASS.

- [ ] **Step 3: Commit**

```bash
git add cmd/commit_wizard.go
git commit -m "feat: badge colorido para todos os tipos de commit"
```

---

## Fase 5 — Arquitetura

### Task 21: Help completo (aliases, examples, path dinâmico)

**Files:**
- Modify: `cmd/help.go`
- Test: `go build ./... && go run . --help` mostra USAGE/COMMANDS/FLAGS sem `devbox` hardcodado.

**Interfaces:**
- Consumes: `Color*` (T1, já aplicado).
- Produces: nada.

- [ ] **Step 1: Path dinâmico** — trocar o prefixo hardcodado por `cmd.CommandPath()`.
- [ ] **Step 2: Seções EXAMPLES e ALIASES** — renderizar se `cmd.Example != ""` e se `len(cmd.Aliases) > 0`, no mesmo estilo das demais seções.
- [ ] **Step 3: Verificar**

Run: `go build ./... && go vet ./... && go run . --help && go run . project --help`
Expected: ajuda renderiza sem erro.

- [ ] **Step 4: Commit**

```bash
git add cmd/help.go
git commit -m "refactor: help exibe examples, aliases e path dinâmico"
```

### Task 22: Frontend sem `strings.Fields`

**Files:**
- Modify: `cmd/init.go`
- Test: `go build ./... && go vet ./... && go test ./...`

**Interfaces:**
- Consumes: `system.Execute` (T16).
- Produces: nada.

- [ ] **Step 1: Reescrever `handleFrontend`** — em vez de `Sprintf` + `Fields`, montar `commandName`/`args` por stack com o nome do projeto como argumento separado (sem quebrar em espaços). Manter o sufixo `.cmd` no Windows e o mesmo tratamento de erro. Não mudar comportamento para comandos simples.
- [ ] **Step 2: Verificar**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: PASS.

- [ ] **Step 3: Commit**

```bash
git add cmd/init.go
git commit -m "fix: frontend sem quebrar nomes com espaço"
```

### Task 23: Consistência das stacks + README

**Files:**
- Modify: `internal/scaffold/stacks.go`, `README.md`
- Test: `go build ./... && go vet ./... && go test ./...`

**Interfaces:**
- Consumes: nada.
- Produces: nada.

- [ ] **Step 1: `stacks.go`** — `"Simples (padrao)"`→`"Simples (padrão)"`; variante `"simple"` (ruby)→`"Simple"`; remover `Source: "templates/node/base"` morto do Node.js (as variantes já definem Source).
- [ ] **Step 2: `README.md`** — adicionar linha do Java na tabela Backend e corrigir nomes de zip divergentes (`Gin.zip`, `python.zip`, `ruby.zip`).
- [ ] **Step 3: Verificar**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: PASS.

- [ ] **Step 4: Commit**

```bash
git add internal/scaffold/stacks.go README.md
git commit -m "docs: consistência nas stacks e README"
```

---

## Fase 6 — Limpeza

### Task 24: Cleanup com confirmação + iconFail correto

**Files:**
- Modify: `cmd/cleanup.go`
- Test: `go build ./... && go vet ./... && go test ./...`

**Interfaces:**
- Consumes: huh (T9), `promptAborted` (T15, usar se existir; senão tratar `huh.ErrUserAborted` local).
- Produces: nada.

- [ ] **Step 1: Listar antes de apagar** — coletar `found` (dirs existentes), imprimir, confirmar com `huh.NewConfirm().Title("Confirmar limpeza?")`; se negado/cancelado, sair com mensagem (sem apagar nada).
- [ ] **Step 2: Fix iconFail** — trocar `iconFail.Render("!")` por `LogWarning(err.Error())` (o `Render("!")` descarta o ícone).
- [ ] **Step 3: Verificar**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: PASS.

- [ ] **Step 4: Commit**

```bash
git add cmd/cleanup.go
git commit -m "feat: cleanup pede confirmação antes de apagar"
```

### Task 25: Erro visível no `git add` do wizard

**Files:**
- Modify: `cmd/commit_wizard.go`
- Test: `go build ./... && go vet ./... && go test ./...`

**Interfaces:**
- Consumes: nada.
- Produces: nada.

- [ ] **Step 1: Tratar erro**:

```go
printStep("active", "Preparando arquivos (git add .)")
if err := exec.Command("git", "add", ".").Run(); err != nil {
	HandleError(err, "Falha ao adicionar arquivos")
	return
}
```

- [ ] **Step 2: Verificar**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: PASS.

- [ ] **Step 3: Commit**

```bash
git add cmd/commit_wizard.go
git commit -m "fix: erro do git add aparece no wizard"
```

---

## Fase 7 — CI e docs

### Task 26: GitHub Actions CI

**Files:**
- Create: `.github/workflows/ci.yml`
- Test: `cat .github/workflows/ci.yml` (YAML válido).

**Interfaces:**
- Consumes: nada.
- Produces: CI verde em push/PR.

- [ ] **Step 1: Criar `.github/workflows/ci.yml`**:

```yaml
name: CI
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.25"
      - run: go build ./...
      - run: go vet ./...
      - run: go test ./... -v
      - run: test -z "$(gofmt -l .)"
```

- [ ] **Step 2: Validar YAML**

Run: `python3 -c "import yaml,sys; yaml.safe_load(open('.github/workflows/ci.yml'))" 2>/dev/null || cat .github/workflows/ci.yml`
Expected: parse OK ou revisão visual.

- [ ] **Step 3: Commit**

```bash
git add .github/workflows/ci.yml
git commit -m "ci: build, vet, test e fmt no GitHub Actions"
```

### Task 27: README final

**Files:**
- Modify: `README.md`
- Test: revisão visual (`grep -n "DevboxCLI" README.md` só no install path).

**Interfaces:**
- Consumes: T3, T23 (não repetir o que já foi feito; só o que faltar).
- Produces: nada.

- [ ] **Step 1: Revisar** — garantir exemplos com `devbox`, tabela com Java, sem `DevboxCLI` fora do install path.
- [ ] **Step 2: Commit** (só se houve mudança)

```bash
git add README.md
git commit -m "docs: revisa README pós-refatoração"
```

### Task 28: Metadata do Windows + versão

**Files:**
- Modify: `winres/winres.json`
- Test: `python3 -m json.tool winres/winres.json` sem erro.

**Interfaces:**
- Consumes: `Version` (T2 — usar o mesmo valor, ex. `"1.0.1"`).
- Produces: nada.

- [ ] **Step 1: Preencher** — `identity.name: "Devbox"`, `identity.version` = Version, `description: "Devbox CLI - Acelere seu desenvolvimento"`, `info.FileDescription: "Devbox CLI"`, `CompanyName` (manter vazio ou "Samuteg"), `ProductName: "Devbox"`, `ProductVersion`/`FileVersion` = Version, `OriginalFilename: "devbox.exe"`.
- [ ] **Step 2: Validar JSON**

Run: `python3 -m json.tool winres/winres.json > /dev/null && echo OK`
Expected: OK.

- [ ] **Step 3: Commit**

```bash
git add winres/winres.json
git commit -m "chore: preenche metadata do executável Windows"
```

---

## Ordem de execução

| Fase | Tasks | Dependências |
|------|-------|--------------|
| 1 — Fundação | T1, T2, T3, T4 | independentes entre si |
| 2 — Bugs | T5, T6, T7, T8 | independentes entre si |
| 3 — Migração huh | T9→T10→T11→T12→T13→T14 | T9 primeiro; T14 por último |
| 4 — UX | T15–T20 | após T14 |
| 5 — Arquitetura | T21, T22, T23 | independentes entre si |
| 6 — Limpeza | T24, T25 | T24 após T9+T15 |
| 7 — CI/docs | T26, T27, T28 | independentes entre si |
