package cmd

import (
	"fmt"
	"os"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

// Defina a versão atual da sua CLI aqui
const version = "0.0.2"

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Atualiza a devbox para a versão mais recente",
	Run: func(cmd *cobra.Command, args []string) {
		updateCLI()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func updateCLI() {
	// 1. Inicia o Spinner
	s := NewSpinner("Verificando atualizações no GitHub...")
	s.Start()
	defer s.Stop()

	// 2. Verifica a versão mais recente
	// IMPORTANTE: Troque "seu-usuario/devbox" pelo seu repositório real
	latest, found, err := selfupdate.DetectLatest("samuteg/DevboxCLI")
	if err != nil {
		s.Stop()
		LogError(fmt.Sprintf("Erro ao verificar atualizações: %v", err))
		return
	}

	// 3. Valida se encontrou algo
	if !found {
		s.Stop()
		LogWarning("Nenhuma release encontrada no repositório.")
		return
	}

	// 4. Compara as versões (SemVer)
	vCurrent, _ := semver.Make(version)
	if latest.Version.LTE(vCurrent) {
		s.Stop()
		LogSuccess(fmt.Sprintf("Sua Devbox já está atualizada! (v%s)", version))
		return
	}

	// 5. Se houver atualização, pergunta ou prossegue
	s.Suffix = fmt.Sprintf(" Nova versão encontrada: %s. Baixando...", latest.Version)

	// Pega o caminho do executável atual para substituir
	exe, err := os.Executable()
	if err != nil {
		s.Stop()
		LogError("Não foi possível localizar o executável atual.")
		return
	}

	// 6. Realiza o Update (Download + Substituição do Binário)
	if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
		s.Stop()
		LogError(fmt.Sprintf("Erro ao realizar update: %v", err))
		return
	}

	s.Stop()
	LogSuccess(fmt.Sprintf("Atualizado com sucesso para v%s!", latest.Version))
	LogInfo("Por favor, reinicie seu terminal para as alterações surtirem efeito.")
}
