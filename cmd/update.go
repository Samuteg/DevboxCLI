package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Atualiza a devbox para a versão mais recente",
	Run: func(cmd *cobra.Command, args []string) {
		repo := viper.GetString("repo")

		if repo == "" {
			fmt.Println("❌ Erro: URL do repositório não configurada no .devbox.yaml")
			fmt.Println("Adicione: repo: github.com/seu-usuario/devbox")
			return
		}

		fmt.Printf("Checking for updates for %s...\n", repo)

		// O comando 'go install repo@latest' baixa a versão mais recente do branch principal
		// e substitui o binário atual no seu $GOPATH/bin
		updateProcess := exec.Command("go", "install", repo+"@latest")

		updateProcess.Stdout = os.Stdout
		updateProcess.Stderr = os.Stderr

		err := updateProcess.Run()
		if err != nil {
			fmt.Printf("❌ Erro ao atualizar: %v\n", err)
			return
		}

		fmt.Println("✅ devbox atualizada com sucesso! Reinicie seu terminal se necessário.")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
