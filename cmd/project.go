package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// projectCmd representa o grupo de comandos relacionados a projetos
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Gerencia projetos (init, add, commit, cleanup)",
	Long: `Comandos focados no ciclo de vida de projetos.

Cria novos projetos (init), adiciona componentes (add),
auxilia na criação de commits convencionais (commit)
e mantém o projeto limpo (cleanup).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			MarginLeft(2).
			Render("📦  Devbox — Gestão de Projetos"))
		fmt.Println()

		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)
}
