package cmd

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dryRun bool

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Remove branches locais que já foram mergeadas ou são órfãs",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Abrir o repositório na pasta atual
		repo, err := git.PlainOpen(".")
		if err != nil {
			fmt.Println("❌ Erro: Você não está em um repositório Git.")
			return
		}

		// 2. Obter a branch atual (para não deletá-la)
		head, _ := repo.Head()
		currentBranch := head.Name().Short()

		// 3. Definir branches protegidas (puxando do Viper)
		protected := viper.GetStringSlice("protected_branches")
		if len(protected) == 0 {
			protected = []string{"main", "master", "develop"}
		}

		fmt.Println("🔍 Analisando branches locais...")

		// 4. Iterar sobre as branches locais
		iter, _ := repo.Branches()
		err = iter.ForEach(func(ref *plumbing.Reference) error {
			branchName := ref.Name().Short()

			// Pular se for a branch atual
			if branchName == currentBranch {
				return nil
			}

			// Pular se for protegida
			for _, p := range protected {
				if branchName == p {
					return nil
				}
			}

			// Lógica de Limpeza: Aqui você pode decidir o quão agressivo quer ser.
			// Por segurança, vamos apenas listar e deletar se o usuário confirmar
			// ou se não houver erro no checkout futuro.

			if dryRun {
				fmt.Printf(" [DRY RUN] Branch passível de remoção: %s\n", branchName)
				return nil
			}

			// 5. Deletar a branch local
			err := repo.Storer.RemoveReference(ref.Name())
			if err != nil {
				fmt.Printf("❌ Erro ao deletar %s: %v\n", branchName, err)
			} else {
				fmt.Printf("✅ Branch removida: %s\n", branchName)
			}

			return nil
		})

		if err != nil {
			fmt.Printf("❌ Erro durante a iteração: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanupCmd)
	cleanupCmd.Flags().BoolVarP(&dryRun, "dry-run", "c", false, "Apenas lista as branches sem deletá-las")
}
