package cmd

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/Samuteg/DevboxCLI/internal/scaffold"
	"github.com/Samuteg/DevboxCLI/internal/system"
	"github.com/Samuteg/DevboxCLI/internal/validation"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

//go:embed templates/*
var templatesFS embed.FS

var stacks = scaffold.DefaultStacks()

var initCmd = &cobra.Command{
	Use:     "init",
	Short:   "Inicializa um novo projeto",
	Example: "  devbox init\n  devbox project init",
	Run:     runInit,
}

func runInit(cmd *cobra.Command, args []string) {
	printStep("active", "Configuração Inicial")
	projectName := promptInput("  Nome do Projeto", "Nome muito curto (min 2 caracteres)", 2)
	projectName = strings.TrimSpace(projectName)
	if !validation.IsValidProjectName(projectName) {
		HandleErrorAndExit(fmt.Errorf("nome %q inválido: use 2-64 caracteres alfanuméricos, '-' ou '_' (ex: meu-api)", projectName), "Validação de Entrada")
	}
	if _, err := os.Stat(projectName); err == nil {
		HandleErrorAndExit(fmt.Errorf("diretório %q já existe", projectName), "Validação de Entrada")
	}
	projectType := promptSelect("  💻 Tipo de Projeto", []string{"Backend", "Frontend"})

	var options []string
	for name, s := range stacks {
		if (projectType == "Backend" && s.IsBackend) || (projectType == "Frontend" && !s.IsBackend) {
			options = append(options, name)
		}
	}
	slices.Sort(options)

	stackName := promptSelect("🛠️  Escolha a Tech", options)
	selectedStack := stacks[stackName]

	if len(selectedStack.Variants) > 0 {
		selectedVariant := promptVariant(selectedStack.Variants)
		selectedStack.Source = selectedVariant.Source
		if len(selectedVariant.ExtraDirs) > 0 {
			selectedStack.ExtraDirs = selectedVariant.ExtraDirs
		}
	}

	printStep("done", fmt.Sprintf("Configurado: %s | %s", projectName, stackName))
	fmt.Println()

	if selectedStack.IsBackend {
		handleBackend(projectName, selectedStack)
	} else {
		handleFrontend(projectName, selectedStack)
	}
}

func handleBackend(name string, s scaffold.Stack) {
	if !jsonOutput {
		printStep("active", "Gerando arquivos e diretórios...")
	}
	var createdDirs []string
	err := withSpinner("Escaneando templates...", func() error {
		if err := os.MkdirAll(name, 0755); err != nil {
			return fmt.Errorf("criação da pasta do projeto: %w", err)
		}
		createdDirs = append(createdDirs, name)
		for _, d := range s.ExtraDirs {
			if err := os.MkdirAll(filepath.Join(name, d), 0755); err != nil {
				return fmt.Errorf("criação de diretórios adicionais: %w", err)
			}
			createdDirs = append(createdDirs, filepath.Join(name, d))
		}
		if err := scaffold.MaterializeTemplates(templatesFS, s.Source, name); err != nil {
			return err
		}
		createSnapshot(name, createdDirs)
		return nil
	})
	if err != nil {
		if jsonOutput {
			printJSON(JSONResult{Success: false, Command: "init", Message: err.Error()})
			return
		}
		HandleError(err, "Geração de templates")
		return
	}

	if jsonOutput {
		printJSON(JSONResult{
			Success: true,
			Command: "init",
			Message: "project created",
			Data:    map[string]string{"name": name, "path": name, "stack": s.Name, "type": "backend"},
		})
		return
	}

	printStep("done", "Estrutura de arquivos finalizada")

	if s.RunInstall {
		packageJSONPath := filepath.Join(name, "package.json")
		if _, err := os.Stat(packageJSONPath); err != nil {
			printStep("todo", "package.json não encontrado; instalação automática ignorada")
		} else if _, err := os.Stat(filepath.Join(name, "node_modules")); err == nil {
			printStep("todo", "node_modules já existe no template; instalação ignorada")
		} else {
			installMgr, installArgs := detectPackageManager(name)
			if err := withSpinner(fmt.Sprintf("Instalando dependências (%s, pode levar minutos)...", strings.Join(append([]string{installMgr}, installArgs...), " ")), func() error {
				return system.ExecuteSilentWithTimeout(installMgr, installArgs, name, 10*time.Minute)
			}); err != nil {
				LogWarning(fmt.Sprintf("Falha ao instalar dependências automaticamente. Rode '%s %s' manualmente em ./%s.", installMgr, strings.Join(installArgs, " "), name))
			}
		}
	}

	fmt.Println()
	fmt.Println(lipgloss.NewStyle().Bold(true).MarginLeft(2).Render("Estrutura criada:"))
	renderMinimalTree(name, s)

	ShowSuccessBox(name, s.Name)
}

// detectPackageManager prefere pnpm quando há pnpm-lock.yaml e pnpm instalado;
// caso contrário usa npm. Retorna (binário, args).
func detectPackageManager(projectDir string) (string, []string) {
	if _, err := os.Stat(filepath.Join(projectDir, "pnpm-lock.yaml")); err == nil {
		if _, err := exec.LookPath("pnpm"); err == nil {
			return "pnpm", []string{"install", "--prefer-offline"}
		}
	}
	return "npm", []string{"install", "--no-audit", "--no-fund"}
}

func handleFrontend(name string, s scaffold.Stack) {
	if !jsonOutput {
		fmt.Printf("\n🎨 %s\n", info("Iniciando gerador oficial do "+s.Name))
	}

	commandName, args := buildFrontendArgs(s.Source, name)
	if commandName == "" {
		if jsonOutput {
			printJSON(JSONResult{Success: false, Command: "init", Message: "empty frontend stack command"})
			return
		}
		HandleError(errors.New("comando vazio para stack frontend"), "Configuração de Stack")
		return
	}

	if runtime.GOOS == "windows" {
		if commandName == "npx" || commandName == "npm" {
			commandName += ".cmd"
		}
	}

	if err := system.Execute(commandName, args, ""); err != nil {
		if jsonOutput {
			printJSON(JSONResult{Success: false, Command: "init", Message: err.Error()})
			return
		}
		HandleError(err, "Execução do gerador frontend")
		return
	}

	if jsonOutput {
		printJSON(JSONResult{
			Success: true,
			Command: "init",
			Message: "project created",
			Data:    map[string]string{"name": name, "path": name, "stack": s.Name, "type": "frontend"},
		})
		return
	}

	ShowSuccessBox(name, s.Name)
}

// buildFrontendArgs divide o comando do gerador (ex: "pnpm create vite@latest %s")
// em binário + args, substituindo %s pelo nome do projeto. Partes sem %s
// são repassadas intactas (sem Sprintf: verbos ausentes gerariam %!(EXTRA...)).
func buildFrontendArgs(source, name string) (string, []string) {
	templateParts := strings.Fields(source)
	if len(templateParts) == 0 {
		return "", nil
	}
	commandName := templateParts[0]

	args := make([]string, 0, len(templateParts)-1)
	for _, p := range templateParts[1:] {
		args = append(args, strings.ReplaceAll(p, "%s", name))
	}
	return commandName, args
}

func promptInput(label, errMsg string, minLen int) string {
	var value string
	if err := huh.NewInput().
		Title(label).
		Validate(func(s string) error {
			if len(s) < minLen {
				return errors.New(errMsg)
			}
			return nil
		}).
		Value(&value).
		Run(); err != nil {
		if promptAborted(err) {
			os.Exit(0)
		}
		HandleErrorAndExit(err, "Entrada do usuário")
	}
	return value
}

func promptSelect(label string, items []string) string {
	var value string
	options := make([]huh.Option[string], len(items))
	for i, item := range items {
		options[i] = huh.NewOption(item, item)
	}
	if err := huh.NewSelect[string]().
		Title(label).
		Options(options...).
		Value(&value).
		Run(); err != nil {
		if promptAborted(err) {
			os.Exit(0)
		}
		HandleErrorAndExit(err, "Seleção de opção")
	}
	return value
}

func init() {
	registerDual(initCmd)
}

func promptVariant(variants []scaffold.Variant) scaffold.Variant {
	var value string
	options := make([]huh.Option[string], len(variants))
	for i, v := range variants {
		options[i] = huh.NewOption(v.Name, v.Name)
	}
	if err := huh.NewSelect[string]().
		Title("⚡ Escolha uma variante").
		Options(options...).
		Value(&value).
		Run(); err != nil {
		if promptAborted(err) {
			os.Exit(0)
		}
		HandleErrorAndExit(err, "Seleção de variante")
	}
	for _, v := range variants {
		if v.Name == value {
			return v
		}
	}
	os.Exit(1)
	return scaffold.Variant{}
}

func renderMinimalTree(projectName string, s scaffold.Stack) {
	branch := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("├──")
	lastBranch := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("└──")
	folder := lipgloss.NewStyle().Foreground(lipgloss.Color("#F4D03F")).Bold(true)

	fmt.Printf("  %s\n", folder.Render(projectName+"/"))

	limit := min(3, len(s.ExtraDirs))

	for i := 0; i < limit; i++ {
		char := branch
		if i == limit-1 && limit < 4 {
			char = lastBranch
		}
		fmt.Printf("  %s %s\n", char, folder.Render(s.ExtraDirs[i]))
	}

	if len(s.ExtraDirs) > limit {
		fmt.Printf("  %s %s\n", lastBranch, lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("242")).Render("... e mais diretórios"))
	}
}
