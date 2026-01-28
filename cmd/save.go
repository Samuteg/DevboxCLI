package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Variável para a flag --force
var forceCommit bool

// Função auxiliar para verificar a branch atual
func checkBranchProtection() error {
	// 1. Obter branch atual via git
	out, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		return fmt.Errorf("não consegui ler a branch atual (é um repositório git?)")
	}
	currentBranch := strings.TrimSpace(string(out))

	// 2. Ler lista de bloqueio do Viper (ou usar padrão se vazio)
	protected := viper.GetStringSlice("protected_branches")
	if len(protected) == 0 {
		protected = []string{"main", "master"} // Fallback padrão
	}

	// 3. Verificar se está na lista negra
	for _, p := range protected {
		if currentBranch == p {
			if forceCommit {
				fmt.Printf("⚠️  ALERTA: Commitando na branch protegida '%s' (Bypass ativado).\n", currentBranch)
				return nil
			}
			return fmt.Errorf("🚫 AÇÃO BLOQUEADA: Você está na branch '%s'.\n   Não commite direto na main! Crie uma branch nova ou use --force.", currentBranch)
		}
	}

	return nil
}

// O comando SAVE atualizado
var saveCmd = &cobra.Command{
	Use:   "save [mensagem]",
	Short: "Commit seguro com validação de branch",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		message := args[0]

		// --- NOVO: CHECAGEM DE BRANCH ---
		if err := checkBranchProtection(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// (Aqui entra sua lógica de validação de regex feita anteriormente)
		// ...

		fmt.Println("🔄 Iniciando fluxo de sincronização...")

		exec.Command("git", "add", ".").Run()

		c := exec.Command("git", "commit", "-m", message)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			fmt.Println("❌ Erro no commit.")
			return
		}

		fmt.Println("🚀 Enviando...")
		exec.Command("git", "push").Run()
		fmt.Println("✅ Feito!")
	},
}

func init() {
	rootCmd.AddCommand(saveCmd)

	// Adiciona a flag --force ou -f
	saveCmd.Flags().BoolVarP(&forceCommit, "force", "f", false, "Ignora proteção de branch (main/master)")
}
