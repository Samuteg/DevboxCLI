package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Verifica a saúde do ambiente de desenvolvimento",
	Run:   runDoctor,
}

var (
	headerColor = lipgloss.Color("#7D56F4")
	subtleColor = lipgloss.Color("#5C5C5C")
	failColor   = lipgloss.Color("#FF4C4C")

	// Dimensões das Colunas

	colNameWidth   = 15
	colStatusWidth = 10
	colMsgWidth    = 40

	// Estilos das Células
	headerStyle = lipgloss.NewStyle().
			Foreground(headerColor).
			Bold(true).
			Padding(0, 1)

	checkStyle = lipgloss.NewStyle().Padding(0, 1)

	iconSuccess = lipgloss.NewStyle().Foreground(failColor).SetString("✔ PASSED")
	iconFail    = lipgloss.NewStyle().Foreground(failColor).SetString("✖ FAILED")
)

type CheckResult struct {
	Name    string
	Status  string
	Message string
}

func runDoctor(cmd *cobra.Command, args []string) {
	fmt.Println(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFF")).Render("  🩺  DEVBOX DOCTOR"))
	fmt.Println(lipgloss.NewStyle().Foreground(subtleColor).Render("  Verificando dependências do sistema..."))
	fmt.Println()

	checks := []struct {
		cmd  string
		name string
		url  string
	}{
		{"go", "Go Lang", "https://go.dev/dl/"},
		{"node", "Node.js", "https://nodejs.org/"},
		{"npm", "NPM", "install node"},
		{"pnpm", "PNPM", "npm install -g pnpm"},
		{"python", "Python", "https://python.org"},
		{"docker", "Docker", "https://docs.docker.com/get-docker/"},
		{"git", "Git", "https://git-scm.com/"},
	}

	headers := lipgloss.JoinHorizontal(lipgloss.Top,
		headerStyle.Width(colNameWidth).Render("FERRAMENTA"),
		headerStyle.Width(colStatusWidth).Render("STATUS"),
		headerStyle.Width(colMsgWidth).Render("DETALHES"),
	)

	border := lipgloss.NewStyle().Foreground(subtleColor).Render(strings.Repeat("─", colNameWidth+colStatusWidth+colMsgWidth+6))

	fmt.Println("  " + headers)
	fmt.Println("  " + border)

	hasError := false

	for _, c := range checks {
		var status, msg string

		path, err := exec.LookPath(c.cmd)
		if err != nil {
			status = iconFail.String()
			msg = lipgloss.NewStyle().Foreground(failColor).Render("Instale via: " + c.url)
			hasError = true
		} else {
			status = iconSuccess.String()
			version := getVersion(c.cmd)
			msg = lipgloss.NewStyle().Foreground(subtleColor).Render(fmt.Sprintf("%s (%s)", path, version))
		}

		row := lipgloss.JoinHorizontal(lipgloss.Top,
			checkStyle.Width(colNameWidth).Foreground(lipgloss.Color("#FFF")).Render(c.name),
			checkStyle.Width(colStatusWidth).Render(status),
			checkStyle.Width(colMsgWidth).Render(msg),
		)

		fmt.Println("  " + row)
	}

	fmt.Println()

	// Diagnóstico Final
	if hasError {
		box := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(failColor).
			Padding(0, 1).
			Render("⚠️  Algumas ferramentas essenciais estão faltando.\nPor favor, instale-as para garantir o funcionamento total.")
		fmt.Println(box)
	} else {
		box := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#04B575")).
			Padding(0, 1).
			Render("✨  Seu ambiente está perfeito! Tudo pronto para codar.")
		fmt.Println(box)
	}

	fmt.Println(lipgloss.NewStyle().Foreground(subtleColor).Render(fmt.Sprintf("\n  OS: %s | Arch: %s", runtime.GOOS, runtime.GOARCH)))
	fmt.Println()
}

// Função auxiliar simples para tentar pegar versão curta
func getVersion(cmd string) string {
	out, err := exec.Command(cmd, "--version").Output()
	if err != nil {
		return "detectado"
	}
	// Pega só a primeira linha e limita tamanho para não quebrar a tabela
	v := strings.Split(string(out), "\n")[0]
	if len(v) > 15 {
		return "v" + strings.TrimSpace(v[:15]) + "..."
	}
	return strings.TrimSpace(v)
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
