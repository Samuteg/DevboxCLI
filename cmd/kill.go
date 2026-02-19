package cmd

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:   "kill [porta]",
	Short: "Termina o processo que está a ocupar uma porta específica",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		port := args[0]

		// Validação de segurança para a porta
		validPort := regexp.MustCompile(`^[0-9]+$`)
		if !validPort.MatchString(port) {
			HandleError(fmt.Errorf("porta '%s' é inválida", port), "Validação de Entrada")
			return
		}

		fmt.Printf("  %s %s %s\n\n",
			lipgloss.NewStyle().Foreground(targetColor).Render("🎯"),
			"Rastreando alvo na porta:",
			portStyle.Render(port),
		)

		if runtime.GOOS == "windows" {
			killWindows(port)
		} else {
			killUnix(port)
		}
	},
}

func killUnix(port string) {
	printStep("active", "Buscando PID via lsof...")

	cmdFind := exec.Command("lsof", "-t", "-i:"+port)
	out, err := cmdFind.Output()

	if err != nil || len(out) == 0 {
		printStep("todo", "Nenhum processo ativo encontrado")
		return
	}

	pid := strings.TrimSpace(string(out))

	// Segurança contra injeção no comando kill
	validPid := regexp.MustCompile(`^[0-9\s]+$`)
	if !validPid.MatchString(pid) {
		HandleError(fmt.Errorf("PID retornado é suspeito"), "Segurança")
		return
	}

	printStep("active", fmt.Sprintf("Encerrando processo %s", pidStyle.Render("("+pid+")")))

	pids := strings.Fields(pid)
	for _, p := range pids {
		cmdKill := exec.Command("kill", "-9", p)
		if err := cmdKill.Run(); err != nil {
			HandleError(err, "Falha ao matar processo "+p)
			return
		}
	}

	printStep("done", "Processo(s) terminado(s)")
	showKillFinal(port)
}

func killWindows(port string) {
	printStep("active", "Executando PowerShell Stop-Process...")

	// Comando robusto para pegar todos os processos na porta e forçar parada
	command := fmt.Sprintf("(Get-NetTCPConnection -LocalPort %s -ErrorAction SilentlyContinue).OwningProcess | ForEach-Object { Stop-Process -Id $_ -Force }", port)
	cmd := exec.Command("powershell", "-Command", command)

	if err := cmd.Run(); err != nil {
		printStep("todo", "Porta parece já estar livre ou acesso negado")
	} else {
		printStep("done", "Porta libertada")
		showKillFinal(port)
	}
}

// Helper para o feedback final estilizado
func showKillFinal(port string) {
	fmt.Println()
	msg := lipgloss.NewStyle().
		Bold(true).
		Foreground(successColor).
		Render(fmt.Sprintf("✨ Porta %s limpa e pronta para uso!", port))

	fmt.Printf("  %s %s\n", skullIcon, msg)
}

func init() {
	rootCmd.AddCommand(killCmd)
}
