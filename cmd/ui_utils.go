package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/Samuteg/DevboxCLI/internal/term"
	"github.com/briandowns/spinner"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func styleRender(st lipgloss.Style) func(...interface{}) string {
	return func(a ...interface{}) string {
		if term.Plain() {
			return fmt.Sprint(a...)
		}
		return st.Render(fmt.Sprint(a...))
	}
}

var (
	info     = styleRender(lipgloss.NewStyle().Foreground(ColorSecondary))
	success  = styleRender(lipgloss.NewStyle().Foreground(ColorSuccess))
	warning  = styleRender(lipgloss.NewStyle().Foreground(ColorWarning))
	errColor = styleRender(lipgloss.NewStyle().Foreground(ColorError))
	bold     = styleRender(lipgloss.NewStyle().Bold(true))

	iconStepTodo   = lipgloss.NewStyle().Foreground(ColorSubtle).SetString("○")
	iconStepActive = lipgloss.NewStyle().Foreground(ColorPrimary).SetString("●")
	iconStepDone   = lipgloss.NewStyle().Foreground(ColorSuccess).SetString("✔")

	textStepTodo   = lipgloss.NewStyle().Foreground(ColorSubtle)
	textStepActive = lipgloss.NewStyle().Foreground(ColorWhite).Bold(true)
	textStepDone   = lipgloss.NewStyle().Foreground(lipgloss.Color("#CCC")).Strikethrough(false)

	delStyle  = lipgloss.NewStyle().Foreground(ColorError).Bold(true)
	pathStyle = lipgloss.NewStyle().Foreground(ColorMuted).Italic(true)

	summaryBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSecondary).
			Padding(0, 2).
			MarginTop(1)

	treeBranch = lipgloss.NewStyle().Foreground(ColorSubtle).Render("├──")
	treeLast   = lipgloss.NewStyle().Foreground(ColorSubtle).Render("└──")

	commitScopeStyle = lipgloss.NewStyle().Foreground(ColorSecondary).Bold(true)
	commitTypeStyle  = lipgloss.NewStyle().Bold(true).Padding(0, 1).Foreground(ColorWhite)

	errorBanner = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorWhite).
			Background(ColorError).
			Padding(0, 1)

	errorContextStyle = lipgloss.NewStyle().
				Foreground(ColorWhite).
				Bold(true)

	errorMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("246"))

	errorIcon = lipgloss.NewStyle().
			Foreground(ColorError).
			SetString("✘")

	pidStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true)
	portStyle = lipgloss.NewStyle().Foreground(ColorSecondary).Bold(true)

	skullIcon = lipgloss.NewStyle().Foreground(ColorFix).SetString("☠")
)

const (
	IconSuccess = "✔"
	IconError   = "✖"
	IconInfo    = "ℹ"
	IconWait    = "⚡"
	IconStep    = "➜"
)

const stringHandler = "%s %s\n"

func LogSuccess(message string) {
	fmt.Printf(stringHandler, success(IconSuccess), message)
}

func LogError(message string) {
	fmt.Printf(stringHandler, errColor(IconError), errColor(message))
}

func LogInfo(message string) {
	fmt.Printf(stringHandler, info(IconInfo), message)
}

func LogWarning(message string) {
	fmt.Printf(stringHandler, warning("!"), warning(message))
}

func NewSpinner(message string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	if term.Plain() {
		s.Disable()
		fmt.Println("  " + message)
		return s
	}
	s.Suffix = " " + info(message)
	s.Color("cyan")
	return s
}

// withSpinner executa fn com spinner (ou log simples em modo plain/CI).
func withSpinner(message string, fn func() error) error {
	if term.Plain() {
		fmt.Println("  " + message)
		return fn()
	}
	s := NewSpinner(message)
	s.Start()
	defer s.Stop()
	return fn()
}

func ShowSuccessBox(projectName, stack string) {
	rule := lipgloss.NewStyle().Foreground(ColorSecondary).Render("---------------------------------")
	if term.Plain() {
		rule = "---------------------------------"
	}
	fmt.Println(bold("Projeto criado com sucesso!"))
	fmt.Println(rule)

	fmt.Printf("  %s %s %s\n", success(IconStep), bold("cd"), projectName)

	switch stack {
	case "Go":
		fmt.Printf(stringHandler, success(IconStep), bold("go run ./..."))
	case "Python":
		fmt.Printf(stringHandler, success(IconStep), bold("python main.py"))
	case "Vite", "Next", "Next.js":
		fmt.Printf(stringHandler, success(IconStep), bold("pnpm dev"))
	default:
		fmt.Printf(stringHandler, success(IconStep), bold("veja o README do projeto"))
	}
	fmt.Println(rule)
	fmt.Println(info("Duvidas? Acesse nosso GitHub!\n"))
}

func PrintBanner() {
	style := lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)

	asciiArt := `
    ____  _______    ______  ____  _  __
   / __ \/ ____/ |  / / __ )/ __ \| |/ /
  / / / / __/  | | / / __  / / / /   /
 / /_/ / /___  | |/ / /_/ / /_/ /   |
/_____/_____/  |___/_____/\____/_/|_|
`
	fmt.Println(style.Render(asciiArt))
	fmt.Println(lipgloss.NewStyle().Foreground(ColorSubtle).PaddingLeft(2).Render(fmt.Sprintf("v%s • Devbox CLI", Version)))
	fmt.Println()
}

func printStep(status string, text string) {
	var icon, msg string

	switch status {
	case "todo":
		icon = iconStepTodo.String()
		msg = textStepTodo.Render(text)
	case "active":
		icon = iconStepActive.String()
		msg = textStepActive.Render(text)
	case "done":
		icon = iconStepDone.String()
		msg = textStepDone.Render(text)
	case "warn":
		icon = lipgloss.NewStyle().Foreground(ColorWarning).SetString("⚠").String()
		msg = lipgloss.NewStyle().Foreground(ColorWarning).Render(text)
	}

	fmt.Printf("  %s  %s\n", icon, msg)
}

// promptAborted diz se o usuário cancelou (Ctrl+C) e já imprime aviso.
func promptAborted(err error) bool {
	if err == huh.ErrUserAborted {
		fmt.Println("  Operação cancelada.")
		return true
	}
	return false
}

func HandleError(err error, context string) {
	if err == nil {
		return
	}

	fmt.Println()
	fmt.Printf(stringHandler, errorBanner.Render(" ERROR "), errorContextStyle.Render(context))

	fmt.Printf(stringHandler, errorIcon, errorMessageStyle.Render(err.Error()))

	tip := lipgloss.NewStyle().Foreground(ColorSubtle).Italic(true).Render("  💡 Dica: Verifique as permissões ou use 'devbox --help'")
	fmt.Println(tip)
	fmt.Println()
}

func HandleErrorAndExit(err error, context string) {
	HandleError(err, context)
	os.Exit(1)
}
