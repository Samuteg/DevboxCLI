package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Remove arquivos temporários e dependências (node_modules, dist, etc)",
	Run:   runCleanup,
}

func runCleanup(cmd *cobra.Command, args []string) {
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Render("  Iniciando limpeza profunda do ambiente...\n"))

	targets := []string{
		"node_modules",
		"dist",
		"build",
		"bin",
		".next",
		"vendor",
	}

	filesRemoved := 0
	var removedDirs []string

	for _, target := range targets {
		// Verifica se o diretório existe
		if _, err := os.Stat(target); err == nil {
			printStep("active", fmt.Sprintf("Removendo %s...", target))

			err := os.RemoveAll(target)
			if err != nil {
				fmt.Printf("    %s %v\n", iconFail.Render("!"), err)
			} else {
				filesRemoved++
				removedDirs = append(removedDirs, target)
				// Feedback visual de item deletado
				fmt.Printf("    %s %s\n", delStyle.Render("🗑"), pathStyle.Render(target+" removido"))
			}
		}
	}

	if filesRemoved > 0 {
		showCleanupSummary(removedDirs)
	} else {
		fmt.Println(lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true).
			MarginLeft(2).
			Render("\n✨ Nada para limpar! Seu ambiente já está brilhando."))
	}
}

func showCleanupSummary(dirs []string) {
	title := lipgloss.NewStyle().Bold(true).Foreground(secondaryColor).Render("LIMPEZA CONCLUÍDA")

	content := fmt.Sprintf(
		"%s\n\nDiretórios limpos: %v\nStatus: %s",
		title,
		len(dirs),
		lipgloss.NewStyle().Foreground(successColor).Render("Ambiente Otimizado 🚀"),
	)

	fmt.Println(summaryBox.Render(content))
}

func init() {
	rootCmd.AddCommand(cleanupCmd)
}
