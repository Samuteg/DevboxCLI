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
	golangSubDirs = []string{
		"cmd/apps",
		"docs/images",
		"pkg/conf",
		"pkg/logger",
		"pkg/server/auth/acl",
		"pkg/server/controller",
		"pkg/server/model",
		"pkg/server/router/middleware",
		"pkg/server/storage/cache/local",
		"pkg/server/storage/cache/redis",
		"pkg/server/storage/db/mysql",
	}
	rubySubDirs = []string{
		".github/workflows",
		"lib/generators/rails/templates",
		"lib/jb",
		"test/dummy_app/app/assents/config",
		"test/dummy_app/app/assents/javascripts",
		"test/dummy_app/app/assents/stylesheets",
		"test/dummy_app/app/controllers",
		"test/dummy_app/app/helpers",
		"test/dummy_app/app/mailers",
		"test/dummy_app/app/models",
		"test/dummy_app/app/views",
		"test/dummy_app/bin",
		"test/dummy_app/config/environments",
		"test/dummy_app/config/initializers",
		"test/dummy_app/config/features",
	}
)

type Variant struct {
	Name      string
	Source    string
	ExtraDirs []string
}

type Stack struct {
	Name       string
	IsBackend  bool
	Source     string
	ExtraDirs  []string
	RunInstall bool
	Variants   []Variant
}

var stacks = map[string]Stack{
	"Node.js": {
		Name:      "Node",
		IsBackend: true,
		Source:    "templates/node/base",
		Variants: []Variant{
			{
				Name:      "TypeScript (Recomendado)",
				Source:    "templates/node/ts",
				ExtraDirs: []string{"src/types", "src/controllers"},
			},
			{
				Name:      "JavaScript",
				Source:    "templates/node/js",
				ExtraDirs: []string{"src/controllers", "src/libs", "src/middleware", "src/routes", "src/models"},
			},
		},
		RunInstall: true,
	},
	"Go": {
		Name:      "Go",
		IsBackend: true,
		Source:    "templates/go",
		Variants: []Variant{
			{Name: "Simples (padrao)", Source: "templates/golang/simple", ExtraDirs: []string{"cmd/api", "internal/entity", "internal/infra/repository", "internal/infra/web", "internal/usecase"}},
			{Name: "Gin", Source: "templates/golang/Gin", ExtraDirs: golangSubDirs},
		},
	},
	"Python": {
		Name:      "Python",
		IsBackend: true,
		Source:    "templates/python",
		Variants: []Variant{
			{
				Name:   "Simples (Recomendado)",
				Source: "templates/python/simple",
			},
			{
				Name:      "FastAPI",
				Source:    "templates/python/FastAPI",
				ExtraDirs: prefixPaths("backend", pythonSubDirs),
			},
		},
	},
	"Ruby": {
		Name:      "Ruby",
		IsBackend: true,
		Source:    "templates/ruby",
		Variants: []Variant{
			{
				Name:   "simple",
				Source: "templates/ruby/simple",
			},
			{
				Name:      "Rails",
				Source:    "templates/ruby/On_rails",
				ExtraDirs: rubySubDirs,
			},
		},
	},
	"Next.js": {Name: "Next", IsBackend: false, Source: "pnpm create next-app@latest %s"},
	"Vite":    {Name: "Vite", IsBackend: false, Source: "pnpm create vite@latest %s"},
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

	// 1. Escolhe a Stack Principal (ex: Node, Go)
	stackName := promptSelect("🛠️  Escolha a Tech", options)
	selectedStack := stacks[stackName]

	// 2. LÓGICA NOVA: Verifica se tem Variantes
	if len(selectedStack.Variants) > 0 {
		// Mostra o menu de variantes
		selectedVariant := promptVariant(selectedStack.Variants)

		// Sobrescreve as configurações da Stack com as da Variante
		selectedStack.Source = selectedVariant.Source

		// Se a variante tiver diretorios extras definidos, usa eles
		if len(selectedVariant.ExtraDirs) > 0 {
			selectedStack.ExtraDirs = selectedVariant.ExtraDirs
		}
	}

	// 3. Continua para a criação
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

func promptVariant(variants []Variant) Variant {
	var items []string
	for _, v := range variants {
		items = append(items, v.Name)
	}

	prompt := promptui.Select{
		Label: "⚡ Escolha uma variante",
		Items: items,
		Size:  5,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		os.Exit(1)
	}

	return variants[idx]
}
