package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/manifoldco/promptui" // Você vai precisar dessa lib para selects simples
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Assistente de Conventional Commits",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Selecionar o Tipo
		promptType := promptui.Select{
			Label: "Tipo de alteração",
			Items: []string{"feat", "fix", "docs", "style", "refactor", "test", "chore"},
		}
		_, typeStr, err := promptType.Run()
		if err != nil {
			return
		}

		// 2. Escopo (Opcional)
		promptScope := promptui.Prompt{
			Label: "Escopo (opcional, ex: ui, db)",
		}
		scope, _ := promptScope.Run()

		// 3. Mensagem
		promptMsg := promptui.Prompt{
			Label: "Mensagem do commit",
			Validate: func(input string) error {
				if len(input) < 3 {
					return fmt.Errorf("mensagem muito curta")
				}
				return nil
			},
		}
		msg, _ := promptMsg.Run()

		// Montar a string final
		commitMsg := typeStr
		if scope != "" {
			commitMsg += fmt.Sprintf("(%s)", scope)
		}
		commitMsg += ": " + msg

		// Mostrar o que será feito
		fmt.Println()
		LogInfo(fmt.Sprintf("Commitando: %s", commitMsg))

		// Executar o Git
		gitCmd := exec.Command("git", "commit", "-m", commitMsg)
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr

		if err := gitCmd.Run(); err != nil {
			LogError("Falha ao commitar. Você adicionou os arquivos (git add)?")
		} else {
			LogSuccess("Commit realizado com sucesso!")
		}
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
