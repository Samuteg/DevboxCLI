package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [tipo] [nome]",
	Short: "Gera boilerplate para componentes (controller, usecase, etc)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		kind := strings.ToLower(args[0])
		name := args[1]

		if isNodeProject() {
			generateFromTemplate("node", kind, name)
		} else if isGoProject() {
			generateFromTemplate("go", kind, name)
		} else {
			fmt.Println("❌ Erro: Não detectei um projeto Node (package.json) ou Go (go.mod) nesta pasta.")
		}
	},
}

func generateFromTemplate(stack, kind, name string) {
	// 1. Define o caminho do template no embed e o destino no disco
	templatePath := fmt.Sprintf("templates/%s/%s.tmpl", stack, kind)

	// Define onde o arquivo será salvo baseado na stack
	var targetPath string
	capitalizedName := strings.Title(name)

	if stack == "node" {
		targetPath = filepath.Join("src", kind+"s", capitalizedName+strings.Title(kind)+".js")
	} else {
		targetPath = filepath.Join("internal", kind, name+".go")
	}

	// 2. Lê o template do embed
	content, err := templatesFS.ReadFile(templatePath)
	if err != nil {
		fmt.Printf("❌ Erro: O tipo '%s' não existe para a stack %s.\n", kind, stack)
		return
	}

	// 3. Substitui as variáveis
	output := strings.ReplaceAll(string(content), "{{NAME}}", capitalizedName)

	// 4. Cria a pasta se não existir e escreve o arquivo
	os.MkdirAll(filepath.Dir(targetPath), 0755)
	err = os.WriteFile(targetPath, []byte(output), 0644)
	if err != nil {
		fmt.Printf("❌ Erro ao criar arquivo: %v\n", err)
		return
	}

	fmt.Printf("✅ %s gerado com sucesso em: %s\n", strings.Title(kind), targetPath)
}

// Funções de detecção simples
func isNodeProject() bool { _, err := os.Stat("package.json"); return err == nil }
func isGoProject() bool   { _, err := os.Stat("go.mod"); return err == nil }

func init() {
	rootCmd.AddCommand(addCmd)
}
