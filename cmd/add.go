package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Samuteg/DevboxCLI/internal/scaffold"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var validTypes = scaffold.ValidComponentTypes

var addCmd = &cobra.Command{
	Use:   "add [tipo] [nome]",
	Short: "Adiciona um novo componente ao projeto",
	Args:  cobra.MaximumNArgs(2),
	Run:   runAdd,
}

func runAdd(cmd *cobra.Command, args []string) {
	var resourceType, resourceName string

	if len(args) > 0 {
		resourceType = strings.ToLower(args[0])
	}
	if !scaffold.IsValidComponentType(resourceType) {
		if resourceType != "" {
			fmt.Printf("  %s \n", lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Render("Tipo '"+resourceType+"' desconhecido."))
		}
		resourceType = promptSelect("  O que você quer criar?", validTypes)
	}

	if len(args) > 1 {
		resourceName = args[1]
	} else {
		resourceName = promptInput("  Qual o nome do componente?", "O nome é obrigatório", 2)
	}

	printStep("active", fmt.Sprintf("Gerando %s: %s", resourceType, resourceName))

	path, err := scaffold.CreateComponent(resourceType, resourceName, viper.GetString("author"))
	if err != nil {
		HandleError(err, "Criação de Componente")
		os.Exit(1)
	}

	printStep("done", "Componente criado com sucesso!")

	fmt.Println()
	fmt.Println(lipgloss.NewStyle().Bold(true).MarginLeft(2).Render("📂 Arquivo gerado:"))
	renderDynamicTree(path)

	showSuccessBox(resourceName, toTitle(resourceType))
}

func toTitle(raw string) string {
	if raw == "" {
		return raw
	}
	return strings.ToUpper(raw[:1]) + raw[1:]
}

// --- UI: Árvore Dinâmica ---

func renderDynamicTree(path string) {
	parts := strings.Split(filepath.ToSlash(path), "/")

	folderStyle := lipgloss.NewStyle().Foreground(addDirColor).Bold(true)
	fileStyle := lipgloss.NewStyle().Foreground(addComponentColor)
	indent := "  "

	for i, part := range parts {
		isLast := i == len(parts)-1

		if i == 0 {
			fmt.Printf("%s%s\n", indent, folderStyle.Render(part+"/"))
		} else if !isLast {
			fmt.Printf("%s%s %s\n", indent, treeBranch, folderStyle.Render(part+"/"))
			indent += "│   "
		} else {
			fmt.Printf("%s%s %s\n", indent, treeLast, fileStyle.Render(part))
		}
	}
	fmt.Println()
}

func init() {
	rootCmd.AddCommand(addCmd)
}
