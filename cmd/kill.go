package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:   "kill [porta]",
	Short: "Termina o processo que está a ocupar uma porta específica",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		port := args[0]
		fmt.Printf("🔍 Procurando processo na porta %s...\n", port)

		if runtime.GOOS == "windows" {
			killWindows(port)
		} else {
			killUnix(port)
		}
	},
}

// Lógica para Linux e macOS
func killUnix(port string) {
	// lsof -t -i:PORTA retorna apenas o PID
	cmdFind := exec.Command("lsof", "-t", "-i:"+port)
	out, err := cmdFind.Output()

	if err != nil || len(out) == 0 {
		fmt.Printf("⚠️  Nenhum processo encontrado na porta %s.\n", port)
		return
	}

	pid := strings.TrimSpace(string(out))

	// kill -9 para forçar o encerramento
	cmdKill := exec.Command("kill", "-9", pid)
	if err := cmdKill.Run(); err != nil {
		fmt.Printf("❌ Erro ao matar processo %s: %v\n", pid, err)
	} else {
		fmt.Printf("✅ Processo %s na porta %s terminado com sucesso!\n", pid, port)
	}
}

// Lógica para Windows
func killWindows(port string) {
	// netstat para encontrar o PID e taskkill para encerrar
	command := fmt.Sprintf("(Get-NetTCPConnection -LocalPort %s).OwningProcess | ForEach-Object { Stop-Process -Id $_ -Force }", port)
	cmd := exec.Command("powershell", "-Command", command)

	if err := cmd.Run(); err != nil {
		fmt.Printf("⚠️  Não foi possível encontrar ou encerrar processos na porta %s.\n", port)
	} else {
		fmt.Printf("✅ Porta %s libertada com sucesso!\n", port)
	}
}

func init() {
	rootCmd.AddCommand(killCmd)
}
