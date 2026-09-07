package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
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

	var found []string
	for _, target := range targets {
		if _, err := os.Stat(target); err == nil {
			found = append(found, target)
		}
	}

	if len(found) == 0 {
		fmt.Println(lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true).
			MarginLeft(2).
			Render("\n✨ Nada para limpar! Seu ambiente já está brilhando."))
		return
	}

	for _, f := range found {
		printStep("todo", fmt.Sprintf("%s encontrado", f))
	}
	fmt.Println()

	var confirmed bool
	if err := huh.NewConfirm().
		Title("Confirmar limpeza?").
		Description(fmt.Sprintf("Serão removidos: %v", found)).
		Affirmative("Sim").
		Negative("Não").
		Value(&confirmed).
		Run(); err != nil {
		if promptAborted(err) {
			os.Exit(0)
		}
		HandleError(err, "Confirmação de limpeza")
		os.Exit(1)
	}
	if !confirmed {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Render("  Limpeza cancelada pelo usuário."))
		return
	}

	fmt.Println()
	var removedDirs []string
	for _, target := range found {
		printStep("active", fmt.Sprintf("Removendo %s...", target))

		err := os.RemoveAll(target)
		if err != nil {
			LogWarning(err.Error())
		} else {
			removedDirs = append(removedDirs, target)
			fmt.Printf("    %s %s\n", delStyle.Render("🗑"), pathStyle.Render(target+" removido"))
		}
	}

	if len(removedDirs) > 0 {
		showCleanupSummary(removedDirs)
	}
}

func showCleanupSummary(dirs []string) {
	title := lipgloss.NewStyle().Bold(true).Foreground(ColorSecondary).Render("LIMPEZA CONCLUÍDA")

	content := fmt.Sprintf(
		"%s\n\nDiretórios limpos: %v\nStatus: %s",
		title,
		len(dirs),
		lipgloss.NewStyle().Foreground(ColorSuccess).Render("Ambiente Otimizado 🚀"),
	)

	fmt.Println(summaryBox.Render(content))
}

func init() {
	projectCmd.AddCommand(cleanupCmd)
}
