package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	cleanupDryRun bool
	cleanupYes    bool
)

var cleanupCmd = &cobra.Command{
	Use:     "cleanup",
	Short:   "Remove arquivos temporários e dependências (node_modules, dist, etc)",
	Example: "  devbox cleanup\n  devbox cleanup --dry-run\n  devbox cleanup --yes",
	Run:     runCleanup,
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

	// Trava de projeto: não rodar em $HOME ou fora de projeto reconhecível.
	if !hasProjectMarker() && !cleanupYes {
		HandleError(fmt.Errorf("nenhum marcador de projeto (go.mod, package.json, .git) no diretório atual"), "Segurança")
		fmt.Println("  Rode dentro de um projeto ou use --yes para forçar.")
		return
	}

	var foundPaths []string
	var foundDisplay []string
	for _, target := range targets {
		if _, err := os.Stat(target); err == nil {
			abs, _ := filepath.Abs(target)
			foundPaths = append(foundPaths, target)
			foundDisplay = append(foundDisplay, abs)
		}
	}

	if len(foundPaths) == 0 {
		fmt.Println(lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true).
			MarginLeft(2).
			Render("\n✨ Nada para limpar! Seu ambiente já está brilhando."))
		return
	}

	for _, f := range foundDisplay {
		printStep("todo", fmt.Sprintf("%s encontrado", f))
	}
	fmt.Println()

	if cleanupDryRun {
		printStep("done", "Dry-run: nada foi removido.")
		fmt.Println("  Para remover de verdade, rode sem --dry-run.")
		return
	}

	var confirmed bool
	if cleanupYes {
		confirmed = true
	} else if err := huh.NewConfirm().
		Title("Confirmar limpeza?").
		Description(fmt.Sprintf("Serão removidos: %v", foundDisplay)).
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
	for _, target := range foundPaths {
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

func hasProjectMarker() bool {
	for _, m := range []string{"go.mod", "package.json", ".git", "Cargo.toml", "pyproject.toml", "Gemfile"} {
		if _, err := os.Stat(m); err == nil {
			return true
		}
	}
	return false
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
	cleanupCmd.Flags().BoolVar(&cleanupDryRun, "dry-run", false, "mostra o que seria removido sem remover")
	cleanupCmd.Flags().BoolVar(&cleanupYes, "yes", false, "pula confirmação (perigoso fora de projeto)")
	projectCmd.AddCommand(cleanupCmd)
}
