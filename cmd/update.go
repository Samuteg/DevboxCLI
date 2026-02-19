package cmd

import (
	"fmt"
	"os"

	"github.com/blang/semver"
	"github.com/charmbracelet/lipgloss"
	"github.com/manifoldco/promptui"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

// Versão atual da CLI
const version = "1.0.0"

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Atualiza a devbox para a versão mais recente",
	Run: func(cmd *cobra.Command, args []string) {
		updateCLI()
	},
}

func updateCLI() {
	printStep("active", "Buscando atualizações no GitHub...")

	latest, found, err := selfupdate.DetectLatest("Samuteg/DevboxCLI")
	if err != nil {
		HandleError(err, "Falha na conexão com GitHub")
		return
	}

	if !found {
		fmt.Printf("\n  %s \n", lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Render("Nenhuma release encontrada no repositório."))
		return
	}

	// Compara as versões (SemVer)
	vCurrent, _ := semver.Make(version)
	if latest.Version.LTE(vCurrent) {
		printStep("done", "Você já está na última versão!")

		fmt.Println(lipgloss.NewStyle().
			MarginLeft(4).
			Foreground(lipgloss.Color("240")).
			Render(fmt.Sprintf("Versão atual: v%s", version)))
		return
	}

	printStep("done", "Nova versão disponível!")
	fmt.Println()

	// Box de comparação de versões
	compareBox := lipgloss.JoinHorizontal(lipgloss.Center,
		lipgloss.NewStyle().Foreground(lipgloss.Color("246")).Render("v"+version),
		lipgloss.NewStyle().Padding(0, 2).Foreground(lipgloss.Color("240")).Render("→"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Bold(true).Render("v"+latest.Version.String()),
	)

	banner := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFF")).
		Background(primaryColor).
		Padding(0, 1).
		Bold(true).
		Render(" UPDATE DISPONÍVEL ")

	mainBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1, 3).
		Align(lipgloss.Center).
		Render(fmt.Sprintf("%s\n\n%s\n\nNovas melhorias e correções esperam por você.", banner, compareBox))

	fmt.Println(lipgloss.NewStyle().MarginLeft(2).Render(mainBox))
	fmt.Println()

	// Pergunta se deseja prosseguir
	prompt := promptui.Prompt{
		Label:     "  Deseja baixar e instalar agora?",
		IsConfirm: true,
	}

	if _, err := prompt.Run(); err != nil {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Render("  Update ignorado. Você pode atualizar mais tarde."))
		return
	}

	// Realiza o Update com Feedback de Etapas
	fmt.Println()
	printStep("active", "Baixando novo binário...")

	exe, err := os.Executable()
	if err != nil {
		HandleError(err, "Localização do executável")
		return
	}

	if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
		HandleError(err, "Processo de substituição")
		return
	}

	printStep("done", "Download e instalação finalizados")

	showSuccessBox(latest.Version.String(), "Atualização Concluída")

	fmt.Printf("\n  %s\n", lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("240")).
		Render("Por favor, reinicie seu terminal para carregar a nova versão."))
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
