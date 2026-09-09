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
	Use:     "commit",
	Short:   "Assistente interativo para Conventional Commits",
	Example: "  devbox commit",
	Run:     runCommitWizard,
}

func runCommitWizard(cmd *cobra.Command, args []string) {
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
		Title("Escopo (opcional, ex: api)").
		Value(&scope).
		Run(); err != nil {
		if promptAborted(err) {
			os.Exit(0)
		}
		HandleError(err, "Escopo do commit")
		os.Exit(1)
	}
	scope = strings.TrimSpace(scope)

	var description string
	if err := huh.NewInput().
		Title("Descrição curta (max 72 caracteres)").
		Validate(func(s string) error {
			s = strings.TrimSpace(s)
			if len(s) < 3 {
				return fmt.Errorf("a descrição precisa de pelo menos 3 caracteres")
			}
			if len(s) > 72 {
				return fmt.Errorf("descrição muito longa (%d/72): resuma em uma linha", len(s))
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
	description = strings.TrimSpace(description)

	finalMsg := commitType
	if scope != "" {
		finalMsg = fmt.Sprintf("%s(%s): %s", commitType, scope, description)
	} else {
		finalMsg = fmt.Sprintf("%s: %s", commitType, description)
	}

	fmt.Println()
	fmt.Printf("  %s %s\n",
		lipgloss.NewStyle().Bold(true).Render("Mensagem gerada:"),
		lipgloss.NewStyle().Foreground(ColorSecondary).Render(finalMsg),
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
	if out, err := exec.Command("git", "status", "--short").Output(); err != nil {
		HandleError(fmt.Errorf("não é um repositório git ou git indisponível: %w", err), "Git status")
		return
	} else if len(strings.TrimSpace(string(out))) == 0 {
		printStep("warn", "Nada para commitar: working tree limpo.")
		fmt.Println("  Próximo passo: edite arquivos e repita devbox commit")
		return
	} else {
		fmt.Println(lipgloss.NewStyle().Foreground(ColorSubtle).Render("  Arquivos com alteração:"))
		fmt.Println(lipgloss.NewStyle().Foreground(ColorSubtle).PaddingLeft(4).Render(strings.TrimSpace(string(out))))
		fmt.Println()
	}
	printStep("active", "Preparando arquivos (git add .)")
	fmt.Println(lipgloss.NewStyle().Foreground(ColorWarning).Render("  Atenção: 'git add .' inclui tudo, inclusive .env. Confira a lista acima."))
	if err := exec.Command("git", "add", ".").Run(); err != nil {
		HandleError(err, "Falha ao adicionar arquivos")
		return
	}

	printStep("active", "Executando commit")
	cmdGit := exec.Command("git", "commit", "-m", finalMsg)

	if output, err := cmdGit.CombinedOutput(); err != nil {
		outStr := strings.TrimSpace(string(output))
		if strings.Contains(outStr, "nothing to commit") {
			printStep("warn", "Nada para commitar.")
		} else {
			HandleError(fmt.Errorf("git commit falhou: %s", outStr), "Git commit")
			fmt.Println("  Próximo passo: confira 'git status' e tente de novo")
		}
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).PaddingLeft(4).Render(string(output)))
	} else {
		printStep("done", "Commit registrado")
		showCommitSuccess(commitType, scope, description)
	}
}

func showCommitSuccess(cType, scope, msg string) {
	bgColor := lipgloss.Color("240")
	switch cType {
	case "feat":
		bgColor = ColorFeat
	case "fix":
		bgColor = ColorFix
	case "docs":
		bgColor = ColorDocs
	case "refactor":
		bgColor = ColorRefactor
	case "style":
		bgColor = ColorStyle
	case "test":
		bgColor = ColorTest
	case "chore":
		bgColor = ColorChore
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
		lipgloss.NewStyle().Italic(true).Foreground(ColorSuccess).Render("Histórico do Git atualizado!"),
	)

	fmt.Println(lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorSuccess).
		Padding(1, 2).
		MarginTop(1).
		Render(content))
}

func init() {
	projectCmd.AddCommand(commitWizardCmd)
	rootCmd.AddCommand(&cobra.Command{
		Use:     "commit",
		Short:   "Assistente interativo para Conventional Commits",
		Example: "  devbox commit",
		Run:     runCommitWizard,
	})
}
