package cmd

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

//go:embed templates/*
var templatesFS embed.FS

var (
	// Definimos as subpastas separadamente para não poluir o mapa
	pythonSubDirs = []string{
		"src/api/routes",
		"src/api/dependencies",
		"src/config/settings",
		"src/models/schemas",
		"src/models/db",
		"src/repository/crud",
		"src/repository/migrations/versions",
		"src/securities/authorizations",
		"src/securities/hashing",
		"src/securities/verifications",
		"src/utilities/exceptions/http",
		"src/utilities/formatters",
		"src/utilities/menssages/exception/http",
		"tests/end_to_end_tests",
		"tests/unit_tests",
		"tests/integration_tests",
		"tests/security_tests",
	}
)

// Definição das Stacks para evitar strings espalhadas
type Stack struct {
	Name       string
	IsBackend  bool
	Source     string // Nome da pasta no templates/ ou comando npx
	ExtraDirs  []string
	RunInstall bool
}

var stacks = map[string]Stack{
	"Go (Clean Arch)": {
		Name:      "Go",
		IsBackend: true,
		Source:    "templates/go",
		ExtraDirs: []string{"internal/entity", "internal/usecase"},
	},
	"Python (FastAPI)": {
		Name:      "Python",
		IsBackend: true,
		Source:    "templates/python-fastapi",
		// Aqui usamos o helper para injetar "backend" na frente de tudo
		ExtraDirs: prefixPaths("backend", pythonSubDirs),
	},
	"Node (Js)": {Name: "Node", IsBackend: true, Source: "templates/node", RunInstall: true, ExtraDirs: []string{"src/controllers", "src/models"}},
	"Next.js":   {Name: "Next", IsBackend: false, Source: "pnpm create next-app@latest %s"},
	"Vite":      {Name: "Vite", IsBackend: false, Source: "pnpm create vite@latest %s"},
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

	// Cria a pasta raiz do projeto
	os.MkdirAll(name, 0755)

	// Cria diretórios extras definidos na struct (opcional, pois o WalkDir cria tb)
	for _, d := range s.ExtraDirs {
		os.MkdirAll(filepath.Join(name, d), 0755)
	}

	// --- AQUI ESTA A CORREÇÃO ---
	walkErr := fs.WalkDir(templatesFS, s.Source, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 1. Calcula o caminho relativo (remove 'templates/go' do caminho)
		// Ex: "templates/go/cmd/main.go.tmpl" vira "cmd/main.go.tmpl"
		relPath, err := filepath.Rel(s.Source, path)
		if err != nil {
			return err
		}

		// Ignora o próprio diretório raiz (".")
		if relPath == "." {
			return nil
		}

		// 2. Define o caminho final no projeto do usuário
		targetPath := filepath.Join(name, relPath)

		// 3. Se for diretório, cria no disco e retorna
		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// 4. Se for arquivo:
		// A mágica acontece aqui: Remove o .tmpl do nome final se existir
		finalPath := strings.TrimSuffix(targetPath, ".tmpl")

		// Lê o conteúdo do template embarcado
		content, err := templatesFS.ReadFile(path)
		if err != nil {
			return err
		}

		// Escreve o arquivo no disco com o nome correto (sem .tmpl)
		return os.WriteFile(finalPath, content, 0644)
	})
	spin.Stop()

	if walkErr != nil {
		fmt.Printf("%s %v\n", errColor("❌ Erro ao gerar templates:"), walkErr)
		return
	}

	if s.RunInstall {
		installSpin := NewSpinner(info("Instalando dependências (npm install)..."))
		installSpin.Start()
		ExecuteCommandSilent("npm", []string{"install"}, name)
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

	ExecuteCommand(commandName, parts[1:], "")

	showSuccessBox(name, s.Name)
}

// --- Utilitários de Baixo Nível ---

func ExecuteCommand(name string, args interface{}, dir string) {
	validInput := regexp.MustCompile(`^[a-zA-Z0-9_\-\./\\]+$`)
	if !validInput.MatchString(name) {
		fmt.Printf("⚠️ Invalid command name\n")
		return
	}
	var cmd *exec.Cmd
	switch v := args.(type) {
	case string:
		if !validInput.MatchString(v) {
			fmt.Printf("⚠️ Invalid command argument\n")
			return
		}
		parts := strings.Fields(v)
		cmd = exec.Command(name, parts...)
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

// prefixPaths adiciona um prefixo a uma lista de subpastas
func prefixPaths(prefix string, subs []string) []string {
	fullPaths := make([]string, len(subs))
	for i, sub := range subs {
		fullPaths[i] = filepath.Join(prefix, sub)
	}
	return fullPaths
}

func init() {
	rootCmd.AddCommand(initCmd)
}
