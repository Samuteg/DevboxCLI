package cmd

import (
	"fmt"
	"os"

	"github.com/blang/semver"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Atualiza a devbox para a versão mais recente",
	Example: "  devbox update",
	Run: func(cmd *cobra.Command, args []string) {
		updateCLI()
	},
}

func updateCLI() {
	if !jsonOutput {
		printStep("active", "Buscando atualizações no GitHub...")
	}

	latest, found, err := selfupdate.DetectLatest("Samuteg/DevboxCLI")
	if err != nil {
		if jsonOutput {
			printJSON(JSONResult{Success: false, Command: "update", Message: err.Error()})
			return
		}
		HandleError(err, "Falha na conexão com GitHub")
		return
	}

	if !found {
		if jsonOutput {
			printJSON(JSONResult{Success: true, Command: "update", Message: "no release found"})
			return
		}
		fmt.Printf("\n  %s \n", lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Render("Nenhuma release encontrada no repositório."))
		return
	}

	// Compara as versões (SemVer)
	vCurrent, err := semver.Make(Version)
	if err != nil {
		if jsonOutput {
			printJSON(JSONResult{Success: false, Command: "update", Message: fmt.Sprintf("invalid local version %q: %s", Version, err)})
			return
		}
		HandleError(fmt.Errorf("versão local %q inválida: %w", Version, err), "Versão")
		return
	}
	if latest.Version.LTE(vCurrent) {
		if jsonOutput {
			printJSON(JSONResult{Success: true, Command: "update", Message: "already up to date", Data: map[string]string{"current": Version}})
			return
		}
		printStep("done", "Você já está na última versão!")

		fmt.Println(lipgloss.NewStyle().
			MarginLeft(4).
			Foreground(lipgloss.Color("240")).
			Render(fmt.Sprintf("Versão atual: v%s", Version)))
		return
	}

	if !jsonOutput {
		printStep("done", "Nova versão disponível!")
		fmt.Println()

		compareBox := lipgloss.JoinHorizontal(lipgloss.Center,
			lipgloss.NewStyle().Foreground(lipgloss.Color("246")).Render("v"+Version),
			lipgloss.NewStyle().Padding(0, 2).Foreground(lipgloss.Color("240")).Render("→"),
			lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Bold(true).Render("v"+latest.Version.String()),
		)

		banner := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF")).
			Background(ColorPrimary).
			Padding(0, 1).
			Bold(true).
			Render(" ATUALIZAÇÃO DISPONÍVEL ")

		mainBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 3).
			Align(lipgloss.Center).
			Render(fmt.Sprintf("%s\n\n%s\n\nNovas melhorias e correções esperam por você.", banner, compareBox))

		fmt.Println(lipgloss.NewStyle().MarginLeft(2).Render(mainBox))
		fmt.Println()
	}

	var confirmed bool
	if err := huh.NewConfirm().
		Title("Deseja baixar e instalar agora?").
		Affirmative("Sim").
		Negative("Não").
		Value(&confirmed).
		Run(); err != nil {
		if promptAborted(err) {
			os.Exit(0)
		}
		HandleErrorAndExit(err, "Confirmação de atualização")
	}
	if !confirmed {
		if jsonOutput {
			printJSON(JSONResult{Success: false, Command: "update", Message: "update skipped by user", Data: map[string]string{"current": Version, "available": latest.Version.String()}})
			return
		}
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Render("  Update ignorado. Você pode atualizar mais tarde."))
		return
	}

	if !jsonOutput {
		fmt.Println()
		printStep("active", "Baixando novo binário...")
	}

	exe, err := os.Executable()
	if err != nil {
		if jsonOutput {
			printJSON(JSONResult{Success: false, Command: "update", Message: err.Error()})
			return
		}
		HandleError(err, "Localização do executável")
		return
	}

	if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
		if jsonOutput {
			printJSON(JSONResult{Success: false, Command: "update", Message: err.Error()})
			return
		}
		HandleError(err, "Processo de substituição")
		return
	}

	if jsonOutput {
		printJSON(JSONResult{Success: true, Command: "update", Message: "update installed", Data: map[string]string{"from": Version, "to": latest.Version.String()}})
		return
	}

	printStep("done", "Download e instalação finalizados")

	ShowSuccessBox(latest.Version.String(), "Atualização Concluída")

	fmt.Printf("\n  %s\n", lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("240")).
		Render("Por favor, reinicie seu terminal para carregar a nova versão."))
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
