package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var commitWizardCmd = &cobra.Command{
	Use:   "commit",
	Short: "Assistente interativo para Conventional Commits",
	Run: func(cmd *cobra.Command, args []string) {
		var commitType string
		err := huh.NewSelect[string]().
			Title("Tipo de alteração").
			Options(
				huh.NewOption("feat:     ✨ Nova funcionalidade", "feat"),
				huh.NewOption("fix:      🐛 Correção de bug", "fix"),
				huh.NewOption("docs:     📝 Documentação", "docs"),
				huh.NewOption("style:    🎨 Formatação/Estilo", "style"),
				huh.NewOption("refactor: 🔁 Refatoração", "refactor"),
				huh.NewOption("test:     🧪 Testes", "test"),
				huh.NewOption("chore:    🔧 Manutenção", "chore"),
			).
			Value(&commitType).
			Run()
		if err != nil {
			if promptAborted(err) {
				os.Exit(0)
			}
			HandleError(err, "Tipo de alteração")
			os.Exit(1)
		}

		var scope string
		if err := huh.NewInput().
			Title("  🎯 Escopo (opcional)").
			Value(&scope).
			Run(); err != nil {
			if promptAborted(err) {
				os.Exit(0)
			}
			HandleError(err, "Escopo do commit")
			os.Exit(1)
		}

		var description string
		if err := huh.NewInput().
			Title("  📝 Descrição curta").
			Validate(func(s string) error {
				if len(s) < 3 {
					return fmt.Errorf("a descrição precisa de pelo menos 3 caracteres")
				}
				return nil
			}).
			Value(&description).
			Run(); err != nil {
			if promptAborted(err) {
				os.Exit(0)
			}
			HandleError(err, "Descrição do commit")
			os.Exit(1)
		}

		finalMsg := commitType
		if scope != "" {
			finalMsg = fmt.Sprintf("%s(%s): %s", commitType, scope, description)
		} else {
			finalMsg = fmt.Sprintf("%s: %s", commitType, description)
		}

		fmt.Println()
		fmt.Printf("  %s %s\n",
			lipgloss.NewStyle().Bold(true).Render("Mensagem gerada:"),
			lipgloss.NewStyle().Foreground(secondaryColor).Render(finalMsg),
		)

		var confirmed bool
		if err := huh.NewConfirm().
			Title("  Confirmar commit?").
			Affirmative("Sim").
			Negative("Não").
			Value(&confirmed).
			Run(); err != nil {
			if promptAborted(err) {
				os.Exit(0)
			}
			HandleError(err, "Confirmação do commit")
			os.Exit(1)
		}
		if !confirmed {
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Render("  Commit cancelado pelo usuário."))
			return
		}

		fmt.Println()
		printStep("active", "Preparando arquivos (git add .)")
		if exec.Command("git", "add", ".").Run() != nil {
			return
		}

		printStep("active", "Executando commit")
		cmdGit := exec.Command("git", "commit", "-m", finalMsg)

		if output, err := cmdGit.CombinedOutput(); err != nil {
			printStep("todo", "Nada para commitar ou erro no Git.")
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).PaddingLeft(4).Render(string(output)))
		} else {
			printStep("done", "Commit registrado")
			showCommitSuccess(commitType, scope, description)
		}
	},
}

func showCommitSuccess(cType, scope, msg string) {
	bgColor := lipgloss.Color("240")
	switch cType {
	case "feat":
		bgColor = featColor
	case "fix":
		bgColor = fixColor
	case "docs":
		bgColor = docsColor
	case "refactor":
		bgColor = refactorColor
	}

	badge := commitTypeStyle.Background(bgColor).Render(strings.ToUpper(cType))

	formattedScope := ""
	if scope != "" {
		formattedScope = commitScopeStyle.Render("(" + scope + ")")
	}

	content := fmt.Sprintf(
		"📦 %s%s %s\n\n%s",
		badge,
		formattedScope,
		lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Render(msg),
		lipgloss.NewStyle().Italic(true).Foreground(successColor).Render("Histórico do Git atualizado!"),
	)

	fmt.Println(lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(successColor).
		Padding(1, 2).
		MarginTop(1).
		Render(content))
}

func init() {
	projectCmd.AddCommand(commitWizardCmd)
}
