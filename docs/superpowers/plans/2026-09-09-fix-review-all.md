# Fix Review All Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Corrigir todos os achados P0/P1/P2 do review (segurança, UX, performance, arquitetura) no worktree fix/review-all sem quebrar `go vet/test/build`.

**Architecture:** Fixes locais por arquivo, sem reestruturação grande: novo `internal/validation` + `internal/term` para nomes/cores, `system.Execute` com context/timeout, `doctor` paralelo com timeout, `templates.go` com limites streaming, aliases Cobra para hierarquia, config whitelist, versão via ldflags com fallback.

**Tech Stack:** Go 1.25, Cobra, Viper, huh, lipgloss, fatih/color (remover uso), spinner, selfupdate, goreleaser.

**Spec:** Review consolidado da sessão (UX/P0 sem NO_COLOR, hierarquia project vs root, kill sem confirm, traversal via nome, install.sh sem checksum, templates 153M/20k arquivos, doctor serial, config frágil, versão hardcoded).

## Global Constraints

- `go vet ./...` deve passar com exit 0.
- `go test ./...` deve passar.
- `gofmt -l` limpo nos arquivos tocados.
- Não reescrever zips de 153M nesta sessão; prover `scripts/slim-templates.sh` como follow-up.
- Mensagens CLI em PT, sem misturar EN em UI nova.

---

### Task 1: Segurança — validação de nomes + traversal

**Files:**
- Create: `internal/validation/validation.go`
- Create: `internal/validation/validation_test.go`
- Modify: `cmd/init.go` (validar projectName, checar dir, bloquear flag injection)
- Modify: `internal/scaffold/components.go` (validar componentName)
- Modify: `cmd/add.go` (validar antes de chamar CreateComponent)
- Modify: `internal/system/command.go` (validar dir, bloquear args com NUL, checar prefixo no caller)

**Interfaces:**
- Consumes: nada anterior.
- Produces: `validation.IsValidProjectName(s string) bool`, `validation.IsValidComponentName(s string) bool`, `validation.ValidateDirInCWD(dir string) error`.

- [ ] **Step 1: Write the failing test**

```go
package validation

import "testing"

func TestIsValidProjectName(t *testing.T) {
    valid := []string{"meu-api", "app1", "a1"}
    for _, v := range valid {
        if !IsValidProjectName(v) {
            t.Fatalf("esperava válido: %q", v)
        }
    }
    invalid := []string{"", "a", "../evil", "/tmp/x", "-f", "a/b", "a b", ".", ".."}
    for _, v := range invalid {
        if IsValidProjectName(v) {
            t.Fatalf("esperava inválido: %q", v)
        }
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/validation/ -run TestIsValidProjectName -v`
Expected: FAIL with "undefined: IsValidProjectName"

- [ ] **Step 3: Write minimal implementation**

```go
package validation

import "regexp"

var projectRe = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-_]{1,63}$`)

func IsValidProjectName(s string) bool { return projectRe.MatchString(s) }
func IsValidComponentName(s string) bool { return projectRe.MatchString(s) }
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/validation/ -v`
Expected: PASS

- [ ] **Step 5: Wire into cmd/init.go, cmd/add.go, components.go + test components traversal**

```go
// components_test.go
func TestCreateComponentRejectsTraversal(t *testing.T) {
    if _, err := CreateComponent("controller", "../../evil", ""); err == nil {
        t.Fatal("esperava erro para traversal")
    }
}
```

Run: `go test ./internal/scaffold/ -v` Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/validation cmd/init.go cmd/add.go internal/scaffold/components.go
git commit -m "fix: validar nomes de projeto/componente contra traversal"
```

### Task 2: kill seguro + commit seguro + cleanup com trava

**Files:**
- Modify: `cmd/kill.go` (validar range 1-65535, flag --force, confirm, SIGTERM antes de -9, timeout lsof)
- Modify: `cmd/commit_wizard.go` (git status antes de add, separar nothing-to-commit, limite 72)
- Modify: `cmd/cleanup.go` (dry-run, exigir marcador de projeto, Abs, tamanhos)

**Interfaces:**
- Consumes: validation (reuso de regex porta).
- Produces: flags `--force` (kill), `--dry-run --yes` (cleanup).

- [ ] **Step 1: Write failing test para porta**

```go
func TestParsePortRange(t *testing.T) {
    if parsePort("99999999") == "" { t.Fatal("deveria rejeitar") }
}
```

Adaptar: se parsePort não existir, criar `parsePort` em kill.go primeiro como função pura, teste em `cmd/kill_test.go` (package cmd).
- [ ] **Step 2: Run** `go test ./cmd/ -run TestParsePortRange -v` Expected FAIL.
- [ ] **Step 3: Implementar parsePort + --force + confirm + timeout.**
- [ ] **Step 4: Run** `go test ./cmd/ -v` PASS + `go vet ./...` PASS.
- [ ] **Step 5: Commit** `fix: kill com confirmacao e force, commit/cleanup seguros`

### Task 3: UX — NO_COLOR/TTY, cores unificadas, hierarquia, help

**Files:**
- Create: `internal/term/term.go` (Plain(), NoEmoji())
- Modify: `cmd/ui_utils.go` (usar term, remover fatih/color, defer spinner, withSpinner, erros por contexto)
- Modify: `cmd/colors.go` (manter tokens, remover hex cru nos callers)
- Modify: `cmd/project.go`, `cmd/init.go`, `cmd/add.go`, `cmd/commit_wizard.go`, `cmd/cleanup.go`, `cmd/config.go` (aliases root + Example + PT único)
- Modify: `cmd/help.go` (sem largura fixa, sem fundo #333, mostra path completo)
- Modify: `cmd/devbox/main.go` (banner só interativo/verbose)

**Interfaces:**
- Produces: `term.Plain() bool` (true se NO_COLOR/CI ou !isatty), `term.NoEmoji() bool`.

- [ ] **Step 1: Teste failing**

```go
func TestPlainRespectsNoColor(t *testing.T) {
    t.Setenv("NO_COLOR", "1")
    if !Plain() { t.Fatal("esperava plain com NO_COLOR=1") }
}
```

Run: `go test ./internal/term/ -v` FAIL antes de criar.
- [ ] **Step 2-4: Implementar term.go, trocar callers, adicionar Aliases + Example.**
Exemplo alias: `initCmd.Aliases = []string{"new"}` + registrar também no root com `project` como grupo legado (manter ambos, documentar).
- [ ] **Step 5: Commit** `fix: NO_COLOR/TTY, hierarquia e help acionavel`

### Task 4: Performance — templates limites, doctor paralelo, timeouts, startup

**Files:**
- Modify: `internal/scaffold/templates.go` (LimitReader 100MB, rejeitar symlink, checar prefixo Abs, máscara mode)
- Modify: `internal/system/command.go` (ExecuteContext com timeout)
- Modify: `cmd/doctor.go` (errgroup paralelo, timeout 3s, sem truncar versão, required vs optional)
- Modify: `cmd/root.go` (lazy initConfig, SafeWriteConfig, sem duplicata AutomaticEnv)
- Modify: `cmd/init.go` (pular npm install se node_modules existe, timeout 10min)

**Interfaces:**
- Produces: `system.ExecuteContext(ctx, name, args, dir)`, `system.ExecuteWithTimeout(...)`.

- [ ] **Step 1: Teste failing**

```go
func TestMaterializeRejectsSymlink(t *testing.T) {
    // zip com symlink deve retornar erro
}
func TestMaterializeRejectsOversize(t *testing.T) {
    // LimitReader deve barrar > limite
}
```

Run: `go test ./internal/scaffold/ -run TestMaterializeRejects -v` FAIL.
- [ ] **Step 2-4: Implementar limites + ExecuteContext + doctor paralelo.**
- [ ] **Step 5: Commit** `perf: limites templates, doctor paralelo, timeouts`

### Task 5: Arquitetura — config, versão, goreleaser, install.sh, templates script

**Files:**
- Modify: `cmd/config.go` (whitelist author/default-port/update-channel/template-style)
- Modify: `cmd/root.go` (struct config, checar erros)
- Modify: `cmd/version.go` (var Version com ldflags, fallback 1.0.1-dev)
- Modify: `cmd/update.go` (tratar semver.Make erro)
- Modify: `.goreleaser.yaml` (ldflags -s -w -trimpath + -X Version, checksum, changelog desc)
- Modify: `install.sh` (curl -fsSL --proto, checksum sha256, --no-same-owner)
- Create: `scripts/slim-templates.sh` (re zip sem node_modules/.db)

- [ ] **Step 1: Teste failing config whitelist**

```go
func TestConfigSetRejectsUnknownKey(t *testing.T) {
    if isAllowedConfigKey("hacker") { t.Fatal("deveria rejeitar") }
}
```

- [ ] **Step 2-4: Implementar + `go build -ldflags "-X ...Version=9.9.9" ./cmd/devbox && ./devbox --version` confere.**
- [ ] **Step 5: Commit** `chore: config tipada, versao ldflags, install endurecido`
