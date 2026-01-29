package cmd

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// 1. Constantes para garantir que o texto do menu seja idêntico ao do switch
const (
	OptGo   = "Go (Clean Architecture)"
	OptNode = "Node.js (Express + TS)"
	OptPy   = "Python (FastAPI)"
)

// 2. Embutindo a pasta de templates que está em cmd/templates
//
//go:embed templates/*
var templatesFS embed.FS

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Cria um novo projeto a partir de um template interativo",
	Run: func(cmd *cobra.Command, args []string) {
		projectName := promptProjectName()
		projectType := promptType()

		if projectType == "Backend" {
			language := promptBackendLanguage()
			generateProject(projectName, language)
		} else {
			fmt.Println("🚧 Frontend vindo em breve no DevBox!")
		}
	},
}

func generateProject(projectName, techOption string) {
	fmt.Printf("\n🚀 Montando setup para '%s'...\n", projectName)

	var templateSourceDir string
	var extraDirs []string // Lista de pastas que DEVEM ser criadas
	isNode := false

	switch techOption {
	case OptGo:
		templateSourceDir = "templates/go"
		// Pastas padrão para Clean Arch
		extraDirs = []string{"internal/entity", "internal/usecase", "internal/infra/repository", "internal/infra/web"}
	case OptNode:
		templateSourceDir = "templates/node"
		isNode = true
		// Pastas que você solicitou para o Node
		extraDirs = []string{"src/models", "src/routes", "src/controllers", "src/middleware", "src/libs"}
	case OptPy:
		templateSourceDir = "templates/python"
	default:
		fmt.Printf("❌ Erro: Template [%s] não reconhecido.\n", techOption)
		return
	}

	// --- PASSO 1: Criar a pasta raiz ---
	os.MkdirAll(projectName, 0755)

	// --- PASSO 2: Criar as pastas extras (Garante que existam mesmo se vazias) ---
	for _, dir := range extraDirs {
		targetDir := filepath.Join(projectName, dir)
		os.MkdirAll(targetDir, 0755)
	}

	// --- PASSO 3: Varrer e copiar arquivos do EMBED ---
	err := fs.WalkDir(templatesFS, templateSourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(templateSourceDir, path)
		if relPath == "." || d.IsDir() {
			return nil
		}

		targetName := strings.TrimSuffix(relPath, ".tmpl")
		targetPath := filepath.Join(projectName, targetName)

		// Garante a pasta pai do arquivo
		os.MkdirAll(filepath.Dir(targetPath), 0755)

		// Lê e processa o conteúdo
		content, _ := templatesFS.ReadFile(path)
		code := strings.ReplaceAll(string(content), "{{PROJECT_NAME}}", projectName)

		fmt.Printf("   📄 Criando arquivo: %s\n", targetName)
		return os.WriteFile(targetPath, []byte(code), 0644)
	})

	if err != nil {
		fmt.Printf("❌ Erro ao gerar arquivos: %v\n", err)
		return
	}

	if isNode {
		runNpmInstall(projectName)
	}

	fmt.Printf("\n✨ Projeto %s finalizado!\n", projectName)
}

// --- Funções Auxiliares de UI e Execução ---

func runNpmInstall(projectDir string) {
	fmt.Println("📦 Instalando dependências Node...")
	cmd := exec.Command("npm", "install")
	cmd.Dir = projectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("⚠️  Aviso: npm install falhou. Rode manualmente em ./%s\n", projectDir)
	}
}

func promptProjectName() string {
	validate := func(input string) error {
		if len(input) < 2 {
			return errors.New("nome muito curto")
		}
		return nil
	}
	prompt := promptui.Prompt{Label: "📁 Nome do Projeto", Validate: validate}
	res, _ := prompt.Run()
	return res
}

func promptType() string {
	prompt := promptui.Select{
		Label: "💻 Tipo de Projeto",
		Items: []string{"Backend", "Frontend"},
	}
	_, res, _ := prompt.Run()
	return res
}

func promptBackendLanguage() string {
	prompt := promptui.Select{
		Label: "🛠️  Escolha a Stack",
		Items: []string{OptGo, OptNode, OptPy},
	}
	_, res, _ := prompt.Run()
	return res
}

func init() {
	rootCmd.AddCommand(initCmd)
}
