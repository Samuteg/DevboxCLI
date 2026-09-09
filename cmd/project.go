package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// projectCmd representa o grupo LEGADO de comandos de projeto.
// Prefira os comandos diretos: `devbox init`, `devbox add`, `devbox commit`, `devbox cleanup`.
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Gerencia projetos (legado: prefira devbox init/add/commit/cleanup)",
	Long: `Grupo legado mantido por compatibilidade.

Prefira os comandos diretos:
  devbox init
  devbox add [tipo] [nome]
  devbox commit
  devbox cleanup`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginLeft(2).
			Render("📦  Devbox — Gestão de Projetos"))
		fmt.Println()

		_ = cmd.Help()
	},
}

// registerDual registra o comando em projectCmd (legado) E no rootCmd (direto).
// Evita duplicar o Run e mantém `devbox init` e `devbox project init` funcionando.
func registerDual(c *cobra.Command) {
	projectCmd.AddCommand(c)
	rootCmd.AddCommand(c)
}

func init() {
	rootCmd.AddCommand(projectCmd)
}
