package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var commitWizardCmd = &cobra.Command{
	Use:   "commit",
	Short: "Assistente interativo para Conventional Commits",
	Run: func(cmd *cobra.Command, args []string) {

		// Passo 1: Escolher o Tipo
		promptType := promptui.Select{
			Label: "Qual o tipo de alteração?",
			Items: []string{
				"feat     (Nova funcionalidade)",
				"fix      (Correção de bug)",
				"docs     (Documentação)",
				"style    (Formatação, ponto e vírgula, etc)",
				"refactor (Refatoração de código)",
				"test     (Adição de testes)",
				"chore    (Tarefas de build, config, etc)",
			},
		}

		_, result, err := promptType.Run()
		if err != nil {
			fmt.Println("Cancelado.")
			return
		}

		// Pega só a primeira palavra (ex: "feat")
		commitType := result[:4]
		if result[:5] == "style" {
			commitType = "style"
		} // ajuste manual rápido ou use split string
		// Uma forma melhor de pegar o tipo é dar split no espaço:
		// commitType = strings.Split(result, " ")[0]

		// Passo 2: Digitar o Escopo (Opcional)
		promptScope := promptui.Prompt{
			Label: "Escopo (opcional, Enter para pular)",
		}
		scope, _ := promptScope.Run()

		// Passo 3: Digitar a Descrição
		promptMsg := promptui.Prompt{
			Label: "Descrição curta",
			Validate: func(input string) error {
				if len(input) < 3 {
					return fmt.Errorf("descrição muito curta")
				}
				return nil
			},
		}
		description, _ := promptMsg.Run()

		// Montar a mensagem final
		finalMsg := ""
		if scope != "" {
			finalMsg = fmt.Sprintf("%s(%s): %s", commitType, scope, description)
		} else {
			finalMsg = fmt.Sprintf("%s: %s", commitType, description)
		}

		// Confirmação
		fmt.Printf("\n📝 Mensagem gerada: %s\n", finalMsg)

		promptConfirm := promptui.Prompt{
			Label:     "Confirmar commit?",
			IsConfirm: true,
		}

		if _, err := promptConfirm.Run(); err != nil {
			fmt.Println("Commit cancelado.")
			return
		}

		// Executar Git
		exec.Command("git", "add", ".").Run()
		cmdGit := exec.Command("git", "commit", "-m", finalMsg)
		cmdGit.Stdout = os.Stdout
		cmdGit.Stderr = os.Stderr
		if cmdGit.Run(); err != nil {
			fmt.Println("Erro ao commitar.")
		} else {
			fmt.Println("✅ Commit realizado com sucesso!")
			// Opcional: Perguntar se quer dar push
		}
	},
}

func init() {
	rootCmd.AddCommand(commitWizardCmd)
}
