package cmd

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

//go:embed templates/*
var templatesFS embed.FS

// Definição das Stacks para evitar strings espalhadas
type Stack struct {
	Name       string
	IsBackend  bool
	Source     string // Nome da pasta no templates/ ou comando npx
	ExtraDirs  []string
	RunInstall bool
}

var stacks = map[string]Stack{
	"Go (Clean Arch)": {Name: "Go", IsBackend: true, Source: "templates/go", ExtraDirs: []string{"internal/entity", "internal/usecase"}},
	"Node (Express)":  {Name: "Node", IsBackend: true, Source: "templates/node", RunInstall: true, ExtraDirs: []string{"src/controllers", "src/models"}},
	"Next.js":         {Name: "Next", IsBackend: false, Source: "pnpm create next-app@latest %s"},
	"Vite":            {Name: "Vite", IsBackend: false, Source: "pnpm create vite@latest %s"},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Inicializa um novo projeto",
	Run:   runInit,
}

func runInit(cmd *cobra.Command, args []string) {
	projectName := promptInput("📁 Nome do Projeto", "nome muito curto", 2)
	projectType := promptSelect("💻 Tipo de Projeto", []string{"Backend", "Frontend"})

	// Filtra stacks pelo tipo escolhido
	var options []string
	for name, s := range stacks {
		if (projectType == "Backend" && s.IsBackend) || (projectType == "Frontend" && !s.IsBackend) {
			options = append(options, name)
		}
	}

	selectedStack := stacks[promptSelect("🛠️  Escolha a Stack", options)]

	if selectedStack.IsBackend {
		handleBackend(projectName, selectedStack)
	} else {
		handleFrontend(projectName, selectedStack)
	}
}

func handleBackend(name string, s Stack) {
	// Iniciamos o spinner para dar feedback visual
	spin := NewSpinner(info("Construindo a estrutura do backend..."))
	spin.Start()

	os.MkdirAll(name, 0755)
	for _, d := range s.ExtraDirs {
		os.MkdirAll(filepath.Join(name, d), 0755)
	}

	walkErr := fs.WalkDir(templatesFS, s.Source, func(path string, d fs.DirEntry, err error) error {
		// ... lógica ...
		return nil
	})

	spin.Stop()

	if walkErr != nil {
		// Agora 'err' aqui se refere à função de cor do ui_utils.go
		fmt.Printf("%s %v\n", err("❌ Erro ao gerar templates:"), walkErr)
		return
	}

	if s.RunInstall {
		installSpin := NewSpinner(info("Instalando dependências (npm install)..."))
		installSpin.Start()
		executeCommandSilent("npm", []string{"install"}, name) // Versão "silenciosa" para não quebrar o spinner
		installSpin.Stop()
	}

	showSuccessBox(name, s.Name)
}

func handleFrontend(name string, s Stack) {
	fmt.Printf("\n🎨 %s\n", info("Iniciando gerador oficial do "+s.Name))

	rawCmd := fmt.Sprintf(s.Source, name)
	parts := strings.Fields(rawCmd)
	commandName := parts[0]

	// Ajuste para Windows: se for npx ou npm, adiciona .cmd
	if runtime.GOOS == "windows" {
		if commandName == "npx" || commandName == "npm" {
			commandName = commandName + ".cmd"
		}
	}

	executeCommand(commandName, parts[1:], "")

	showSuccessBox(name, s.Name)
}

// --- Utilitários de Baixo Nível ---

func executeCommand(name string, args interface{}, dir string) {
	var cmd *exec.Cmd
	switch v := args.(type) {
	case string:
		cmd = exec.Command(name, v)
	case []string:
		cmd = exec.Command(name, v...)
	}

	cmd.Dir = dir
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Printf("⚠️ Erro ao executar %s: %v\n", name, err)
	}
}

// --- Helpers de UI ---

func promptInput(label, errMsg string, minLen int) string {
	p := promptui.Prompt{
		Label: label,
		Validate: func(s string) error {
			if len(s) < minLen {
				return fmt.Errorf(errMsg)
			}
			return nil
		},
	}
	res, _ := p.Run()
	return res
}

func promptSelect(label string, items []string) string {
	p := promptui.Select{Label: label, Items: items}
	_, res, _ := p.Run()
	return res
}

func init() {
	rootCmd.AddCommand(initCmd)
}
