package cmd

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/Samuteg/DevboxCLI/internal/scaffold"
	"github.com/Samuteg/DevboxCLI/internal/system"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

//go:embed templates/*
var templatesFS embed.FS

var stacks = scaffold.DefaultStacks()

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Inicializa um novo projeto",
	Run:   runInit,
}

func runInit(cmd *cobra.Command, args []string) {
	printStep("active", "Configuração Inicial")
	projectName := promptInput("  📁 Nome do Projeto", "Nome muito curto", 2)
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
	printStep("active", "Gerando arquivos e diretórios...")
	spin := NewSpinner(info(" Escaneando templates..."))
	spin.Start()

	if err := os.MkdirAll(name, 0755); err != nil {
		spin.Stop()
		HandleError(err, "Criação da pasta do projeto")
		return
	}

	for _, d := range s.ExtraDirs {
		if err := os.MkdirAll(filepath.Join(name, d), 0755); err != nil {
			spin.Stop()
			HandleError(err, "Criação de diretórios adicionais")
			return
		}
	}

	walkErr := scaffold.MaterializeTemplates(templatesFS, s.Source, name)
	spin.Stop()
	if walkErr != nil {
		HandleError(walkErr, "Geração de templates")
		return
	}

	printStep("done", "Estrutura de arquivos finalizada")

	if s.RunInstall {
		packageJSONPath := filepath.Join(name, "package.json")
		if _, err := os.Stat(packageJSONPath); err != nil {
			printStep("todo", "package.json não encontrado; instalação automática ignorada")
		} else {
			installSpin := NewSpinner(info("Instalando dependências (npm install)..."))
			installSpin.Start()
			if err := system.ExecuteSilent("npm", []string{"install"}, name); err != nil {
				installSpin.Stop()
				LogWarning("Falha ao instalar dependências automaticamente. Rode 'npm install' manualmente.")
			} else {
				installSpin.Stop()
			}
		}
	}

	fmt.Println()
	fmt.Println(lipgloss.NewStyle().Bold(true).MarginLeft(2).Render("📦 Estrutura criada:"))
	renderMinimalTree(name, s)

	showSuccessBox(name, s.Name)
}

func handleFrontend(name string, s scaffold.Stack) {
	fmt.Printf("\n🎨 %s\n", info("Iniciando gerador oficial do "+s.Name))

	rawCmd := fmt.Sprintf(s.Source, name)
	parts := strings.Fields(rawCmd)
	if len(parts) == 0 {
		HandleError(errors.New("comando vazio para stack frontend"), "Configuração de Stack")
		return
	}
	commandName := parts[0]

	if runtime.GOOS == "windows" {
		if commandName == "npx" || commandName == "npm" {
			commandName += ".cmd"
		}
	}

	if err := system.Execute(commandName, parts[1:], ""); err != nil {
		HandleError(err, "Execução do gerador frontend")
		return
	}

	showSuccessBox(name, s.Name)
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
		HandleError(err, "Entrada do usuário")
		os.Exit(1)
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
		HandleError(err, "Seleção de opção")
		os.Exit(1)
	}
	return value
}

func init() {
	projectCmd.AddCommand(initCmd)
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
		HandleError(err, "Seleção de variante")
		os.Exit(1)
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
