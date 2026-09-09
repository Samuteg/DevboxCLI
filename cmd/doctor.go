package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:     "doctor",
	Short:   "Verifica a saúde do ambiente de desenvolvimento",
	Example: "  devbox doctor",
	Run:     runDoctor,
}

var (
	colNameWidth   = 15
	colStatusWidth = 10
	colMsgWidth    = 40

	headerStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true).
			Padding(0, 1)

	checkStyle = lipgloss.NewStyle().Padding(0, 1)

	iconSuccess = lipgloss.NewStyle().Foreground(ColorSuccess).SetString("OK")
	iconFail    = lipgloss.NewStyle().Foreground(ColorError).SetString("FALHOU")
	iconMissing = lipgloss.NewStyle().Foreground(ColorWarning).SetString("AUSENTE (opc)")
)

type CheckResult struct {
	Name    string
	Status  string
	Message string
}

func runDoctor(cmd *cobra.Command, args []string) {
	if !jsonOutput {
		fmt.Println(lipgloss.NewStyle().Bold(true).Foreground(ColorWhite).Render("  🩺  DEVBOX DOCTOR"))
		fmt.Println(lipgloss.NewStyle().Foreground(ColorMuted).Render("  Verificando dependências do sistema..."))
		fmt.Println()
	}

	checks := []struct {
		cmd      string
		name     string
		url      string
		optional bool
	}{
		{"go", "Go Lang", "https://go.dev/dl/", false},
		{"node", "Node.js", "https://nodejs.org/", false},
		{"npm", "NPM", "instale o Node.js: https://nodejs.org/", false},
		{"pnpm", "PNPM", "npm install -g pnpm", true},
		{"python", "Python", "https://python.org", true},
		{"docker", "Docker", "https://docs.docker.com/get-docker/", true},
		{"git", "Git", "https://git-scm.com/", false},
	}

	if !jsonOutput {
		headers := lipgloss.JoinHorizontal(lipgloss.Top,
			headerStyle.Width(colNameWidth).Render("FERRAMENTA"),
			headerStyle.Width(colStatusWidth).Render("STATUS"),
			headerStyle.Width(colMsgWidth).Render("DETALHES"),
		)

		border := lipgloss.NewStyle().Foreground(ColorMuted).Render(strings.Repeat("─", colNameWidth+colStatusWidth+colMsgWidth+6))

		fmt.Println("  " + headers)
		fmt.Println("  " + border)
	}

	results := make([]CheckResult, len(checks))
	var wg sync.WaitGroup
	for i, c := range checks {
		wg.Add(1)
		go func(i int, tool, url string, optional bool) {
			defer wg.Done()
			results[i] = runOneCheck(tool, url, optional)
		}(i, c.cmd, c.url, c.optional)
	}
	wg.Wait()

	if jsonOutput {
		type DoctorCheck struct {
			Name    string `json:"name"`
			Status  string `json:"status"`
			Message string `json:"message"`
			OK      bool   `json:"ok"`
		}
		var data []DoctorCheck
		for i, r := range results {
			data = append(data, DoctorCheck{
				Name:    checks[i].name,
				Status:  r.Status,
				Message: r.Message,
				OK:      r.Status == "ok",
			})
		}
		printJSON(JSONResult{
			Success: true,
			Command: "doctor",
			Data:    data,
		})
		return
	}

	hasError := false

	for i, r := range results {
		var status, msg string
		if r.Status == "ok" {
			status = iconSuccess.String()
			msg = lipgloss.NewStyle().Foreground(ColorMuted).Render(r.Message)
		} else if checks[i].optional {
			status = iconMissing.String()
			msg = lipgloss.NewStyle().Foreground(ColorMuted).Render("Opcional. Instale via: " + r.Message)
		} else {
			hasError = true
			status = iconFail.String()
			msg = lipgloss.NewStyle().Foreground(ColorError).Render("Instale via: " + r.Message)
		}

		row := lipgloss.JoinHorizontal(lipgloss.Top,
			checkStyle.Width(colNameWidth).Foreground(ColorWhite).Render(checks[i].name),
			checkStyle.Width(colStatusWidth).Render(status),
			checkStyle.Width(colMsgWidth).Render(msg),
		)

		fmt.Println("  " + row)
	}

	fmt.Println()

	if hasError {
		box := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorError).
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

	fmt.Println(lipgloss.NewStyle().Foreground(ColorMuted).Render(fmt.Sprintf("\n  OS: %s | Arch: %s", runtime.GOOS, runtime.GOARCH)))
	fmt.Println()
}

func runOneCheck(tool, url string, _ bool) CheckResult {
	path, err := exec.LookPath(tool)
	if err != nil {
		return CheckResult{Name: tool, Status: "missing", Message: url}
	}
	version := getVersion(tool)
	return CheckResult{Name: tool, Status: "ok", Message: fmt.Sprintf("%s (%s)", path, version)}
}

func getVersion(cmd string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, cmd, "--version").Output()
	if err != nil {
		return "detectado"
	}
	v := strings.TrimSpace(strings.Split(string(out), "\n")[0])
	return v
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
