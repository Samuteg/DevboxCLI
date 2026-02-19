package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var commitWizardCmd = &cobra.Command{
	Use:   "commit",
	Short: "Assistente interativo para Conventional Commits",
	Run: func(cmd *cobra.Command, args []string) {
		items := []string{
			"feat:     ✨ Nova funcionalidade",
			"fix:      🐛 Correção de bug",
			"docs:     📝 Documentação",
			"style:    🎨 Formatação/Estilo",
			"refactor: ♻️  Refatoração",
			"test:     ✅ Testes",
			"chore:    🔧 Manutenção",
		}

		promptType := promptui.Select{
			Label: lipgloss.NewStyle().Foreground(primaryColor).Render("Tipo de alteração"),
			Items: items,
			Size:  7,
		}

		_, result, err := promptType.Run()
		if err != nil {
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("  Commit cancelado."))
			return
		}

		// Extrai apenas o tipo (antes dos dois pontos)
		commitType := strings.TrimSpace(strings.Split(result, ":")[0])

		promptScope := promptui.Prompt{
			Label: "  🎯 Escopo (opcional)",
		}
		scope, _ := promptScope.Run()

		promptMsg := promptui.Prompt{
			Label: "  📝 Descrição curta",
			Validate: func(input string) error {
				if len(input) < 3 {
					return fmt.Errorf("a descrição precisa de pelo menos 3 caracteres")
				}
				return nil
			},
		}
		description, _ := promptMsg.Run()

		// Montar a mensagem final
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

		promptConfirm := promptui.Prompt{
			Label:     "  Confirmar commit?",
			IsConfirm: true,
		}

		if _, err := promptConfirm.Run(); err != nil {
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Render("  Commit cancelado pelo usuário."))
			return
		}

		// ---Executar Git com Feedback ---
		fmt.Println()
		printStep("active", "Preparando arquivos (git add .)")
		if err := exec.Command("git", "add", ".").Run(); err != nil {
			return
		}

		printStep("active", "Executando commit")
		cmdGit := exec.Command("git", "commit", "-m", finalMsg)

		// Se houver erro no commit (ex: nada para commitar)
		if output, err := cmdGit.CombinedOutput(); err != nil {
			printStep("todo", "Nada para commitar ou erro no Git.")
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).PaddingLeft(4).Render(string(output)))
		} else {
			printStep("done", "Commit registrado")
			// Chama a função de sucesso que você já tem
			showCommitSuccess(commitType, scope, description)
		}
	},
}

func showCommitSuccess(cType, scope, msg string) {
	// Define a cor baseada no tipo para o badge
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
	rootCmd.AddCommand(commitWizardCmd)
}
