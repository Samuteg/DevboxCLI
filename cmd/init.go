package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

var (
	isFrontend bool
	isBackend  bool // Nova flag
)

var initCmd = &cobra.Command{
	Use:   "init [nome-do-projeto]",
	Short: "Inicializa um novo projeto",
	Long:  `Cria estrutura de pastas, git, .env e templates para Frontend (Vite) ou Backend (Go Standard Layout).`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		// 1. Criar Diretório Raiz
		if err := os.MkdirAll(projectName, 0755); err != nil {
			fmt.Printf("❌ Erro ao criar pasta raiz: %v\n", err)
			return
		}
		fmt.Printf("📁 Pasta '%s' criada.\n", projectName)

		// 2. Inicializar Git
		_, err := git.PlainInit(projectName, false)
		if err != nil {
			// Ignora erro se já existir git, senão avisa
			if err.Error() != "repository already exists" {
				fmt.Printf("⚠️  Aviso Git: %v\n", err)
			}
		} else {
			fmt.Println("g  Git inicializado.")
		}

		// 3. Criar .env Básico
		createFile(filepath.Join(projectName, ".env"), "ENV=development\nPORT=8080\nDB_URL=postgres://user:pass@localhost:5432/db")

		// --- LÓGICA DO BACKEND ---
		if isBackend {
			setupBackend(projectName)
		}

		// --- LÓGICA DO FRONTEND ---
		if isFrontend {
			setupFrontend(projectName)
		}

		fmt.Printf("\n✨ Projeto %s finalizado!\n", projectName)
		if isBackend {
			fmt.Printf("   Run: cd %s && go run cmd/api/main.go\n", projectName)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&isFrontend, "frontend", "f", false, "Gera template Frontend (Vite)")
	initCmd.Flags().BoolVarP(&isBackend, "backend", "b", false, "Gera template Backend (Go Clean Arch)")
}

// --- FUNÇÕES AUXILIARES ---

func setupBackend(root string) {
	fmt.Println("\n🔨 Configurando Backend (Go + Clean Arch)...")

	// 1. Criar Estrutura de Pastas
	dirs := []string{
		"cmd/api",
		"internal/handler",
		"internal/service",
		"internal/repository",
		"pkg/util",
	}

	for _, dir := range dirs {
		path := filepath.Join(root, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			fmt.Printf("❌ Erro criando pasta %s: %v\n", dir, err)
		}
	}

	// 2. Inicializar go.mod
	// Supomos que o nome do módulo é o nome do projeto para facilitar
	goCmd := exec.Command("go", "mod", "init", root) // root aqui serve como module name provisório
	goCmd.Dir = root
	if out, err := goCmd.CombinedOutput(); err != nil {
		fmt.Printf("❌ Erro no go mod init: %s\n", string(out))
	} else {
		fmt.Println("📦 go.mod criado.")
	}

	// 3. Criar main.go
	mainContent := `package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("🚀 Servidor rodando na porta 8080...")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("DevBox Backend Running!"))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
`
	createFile(filepath.Join(root, "cmd/api", "main.go"), mainContent)

	// 4. Criar .gitignore específico para Go
	gitIgnore := `# Binaries
bin/
*.exe
*.exe~
*.dll

# Go
vendor/
go.work

# Environment
.env

# IDE
.vscode/
.idea/
`
	createFile(filepath.Join(root, ".gitignore"), gitIgnore)

	// 5. Criar docker-compose.yml
	dockerCompose := `version: '3.8'

services:
  db:
    image: postgres:15-alpine
    restart: always
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
      POSTGRES_DB: dbname
    ports:
      - "5432:5432"
    volumes:
      - db_data:/var/lib/postgresql/data

volumes:
  db_data:
`
	createFile(filepath.Join(root, "docker-compose.yml"), dockerCompose)

	fmt.Println("✅ Estrutura Backend criada com sucesso.")
}

func setupFrontend(root string) {
	fmt.Println("\n🎨 Configurando Frontend (Vite)...")

	npmBin := "npm"
	if runtime.GOOS == "windows" {
		npmBin = "npm.cmd"
	}

	// Verifica se npm existe
	if _, err := exec.LookPath(npmBin); err != nil {
		fmt.Println("❌ Erro: NPM não encontrado. Pulei a etapa do frontend.")
		return
	}

	// Executa create vite
	cmd := exec.Command(npmBin, "create", "vite@latest", ".", "--", "--yes")
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Erro no NPM: %v\n", err)
	} else {
		fmt.Println("✅ Frontend configurado.")
	}
}

// Função utilitária para escrever arquivos rapidamente
func createFile(path, content string) {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		fmt.Printf("⚠️ Erro ao escrever %s: %v\n", path, err)
	} else {
		// fmt.Printf("📄 %s criado.\n", filepath.Base(path)) // Opcional: logar cada arquivo
	}
}
