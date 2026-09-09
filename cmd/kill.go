package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	killForce bool
	killYes   bool
)

var killCmd = &cobra.Command{
	Use:   "kill [porta]",
	Short: "Termina o processo que está ocupando uma porta específica",
	Example: "  devbox kill 8080\n  devbox kill 3000 --force\n  devbox kill 8080 --yes",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var port string
		if len(args) > 0 {
			port = args[0]
		} else {
			port = viper.GetString("default-port")
			fmt.Printf("  %s %s\n\n",
				lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true).Render("Nenhuma porta informada. Usando porta padrão:"),
				lipgloss.NewStyle().Foreground(ColorSecondary).Render(port),
			)
		}

		clean, err := parsePort(port)
		if err != nil {
			HandleError(err, "Validação de Entrada")
			fmt.Println("  Próximo passo: devbox kill 8080  (porta 1-65535)")
			return
		}
		port = clean

		fmt.Printf("  %s %s %s\n\n",
			lipgloss.NewStyle().Foreground(ColorStyle).Render("🎯"),
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

func parsePort(raw string) (string, error) {
	if !regexp.MustCompile(`^[0-9]{1,5}$`).MatchString(raw) {
		return "", fmt.Errorf("porta %q é inválida: use número de 1 a 65535 (ex: devbox kill 8080)", raw)
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 || n > 65535 {
		return "", fmt.Errorf("porta %q é inválida: use número de 1 a 65535 (ex: devbox kill 8080)", raw)
	}
	return strconv.Itoa(n), nil
}

func killUnix(port string) {
	printStep("active", "Buscando PID via lsof...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmdFind := exec.CommandContext(ctx, "lsof", "-t", "-i:"+port)
	out, err := cmdFind.Output()

	if err != nil || len(strings.TrimSpace(string(out))) == 0 {
		printStep("warn", "Nenhum processo ativo encontrado na porta")
		return
	}

	pid := strings.TrimSpace(string(out))

	pids := strings.Fields(pid)
	validPid := regexp.MustCompile(`^[0-9]+$`)
	for _, p := range pids {
		if !validPid.MatchString(p) {
			HandleError(fmt.Errorf("PID retornado é suspeito: %q", p), "Segurança")
			return
		}
	}

	printStep("active", fmt.Sprintf("Processo(s) na porta %s: %s", port, pidStyle.Render(strings.Join(pids, ", "))))

	if !killYes && !killForce {
		var confirmed bool
		if err := huh.NewConfirm().
			Title(fmt.Sprintf("Encerrar %d processo(s) na porta %s?", len(pids), port)).
			Description("Será enviado SIGTERM. Use --force para SIGKILL imediato.").
			Affirmative("Sim").
			Negative("Não").
			Value(&confirmed).
			Run(); err != nil {
			if promptAborted(err) {
				os.Exit(0)
			}
			HandleError(err, "Confirmação de kill")
			return
		}
		if !confirmed {
			printStep("warn", "Operação cancelada pelo usuário")
			return
		}
	}

	killed := 0
	for _, p := range pids {
		sig := "SIGTERM"
		if killForce {
			sig = "SIGKILL"
		}
		printStep("active", fmt.Sprintf("Enviando %s para PID %s...", sig, p))
		if killForce {
			if err := exec.Command("kill", "-9", p).Run(); err != nil {
				HandleError(err, "Falha ao matar processo "+p)
				return
			}
			killed++
			continue
		}
		if exec.Command("kill", p).Run() == nil {
			killed++
			continue
		}
		HandleError(fmt.Errorf("PID %s não respondeu a SIGTERM; repita com --force para SIGKILL", p), "Kill")
		return
	}

	if killed == 0 {
		printStep("warn", "Nenhum processo foi terminado")
		return
	}

	printStep("done", fmt.Sprintf("Processo(s) terminado(s): %d", killed))
	showKillFinal(port)
}

func killWindows(port string) {
	printStep("active", "Executando PowerShell Stop-Process...")

	command := fmt.Sprintf("(Get-NetTCPConnection -LocalPort %s -ErrorAction SilentlyContinue).OwningProcess | ForEach-Object { Stop-Process -Id $_ -Force }", port)
	cmd := exec.Command("powershell", "-Command", command)

	if cmd.Run() != nil {
		printStep("warn", "Porta parece já estar livre ou acesso negado")
	} else {
		printStep("done", "Porta liberada")
		showKillFinal(port)
	}
}

func showKillFinal(port string) {
	fmt.Println()
	msg := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorSuccess).
		Render(fmt.Sprintf("✨ Porta %s limpa e pronta para uso!", port))

	fmt.Printf("  %s %s\n", skullIcon, msg)
}

func init() {
	killCmd.Flags().BoolVarP(&killForce, "force", "f", false, "usa SIGKILL imediato (sem tentar SIGTERM)")
	killCmd.Flags().BoolVarP(&killYes, "yes", "y", false, "pula confirmação interativa")
	rootCmd.AddCommand(killCmd)
}
