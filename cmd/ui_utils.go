package cmd

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

var (
	info     = color.New(color.FgCyan).SprintFunc()
	success  = color.New(color.FgGreen).SprintFunc()
	warning  = color.New(color.FgYellow).SprintFunc()
	errColor = color.New(color.FgRed).SprintFunc()
	bold     = color.New(color.Bold).SprintFunc()

	iconStepTodo   = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).SetString("○")
	iconStepActive = lipgloss.NewStyle().Foreground(primaryColor).SetString("●")
	iconStepDone   = lipgloss.NewStyle().Foreground(successColor).SetString("✔")

	textStepTodo   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	textStepActive = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF")).Bold(true)
	textStepDone   = lipgloss.NewStyle().Foreground(lipgloss.Color("#CCC")).Strikethrough(false)

	errorColor     = lipgloss.Color("#E74C3C")
	primaryColor   = lipgloss.Color("#7D56F4")
	secondaryColor = lipgloss.Color("#00ADD8")
	successColor   = lipgloss.Color("#27AE60")

	successBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(successColor).
			Padding(1, 2).
			Bold(true).
			MarginTop(1)

	highlight = lipgloss.NewStyle().Foreground(secondaryColor).Bold(true)

	deleteColor  = errorColor
	neutralColor = lipgloss.Color("242")

	delStyle  = lipgloss.NewStyle().Foreground(deleteColor).Bold(true)
	pathStyle = lipgloss.NewStyle().Foreground(neutralColor).Italic(true)

	summaryBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(0, 2).
			MarginTop(1)

	addComponentColor = lipgloss.Color("#00ADD8")
	addDirColor       = lipgloss.Color("#F1C40F")

	treeBranch = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("├──")
	treeLast   = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("└──")

	featColor     = lipgloss.Color("#A3BE8C")
	fixColor      = lipgloss.Color("#BF616A")
	docsColor     = lipgloss.Color("#81A1C1")
	refactorColor = lipgloss.Color("#B48EAD")

	commitScopeStyle = lipgloss.NewStyle().Foreground(secondaryColor).Bold(true)
	commitTypeStyle  = lipgloss.NewStyle().Bold(true).Padding(0, 1).Foreground(lipgloss.Color("#FFF"))

	errorBanner = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFF")).
			Background(errorColor).
			Padding(0, 1)

	errorContextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFF")).
				Bold(true)

	errorMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("246"))

	errorIcon = lipgloss.NewStyle().
			Foreground(errorColor).
			SetString("✘")

	targetColor = lipgloss.Color("#EBCB8B")
	killColor   = lipgloss.Color("#BF616A")

	pidStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true)
	portStyle = lipgloss.NewStyle().Foreground(secondaryColor).Bold(true)

	skullIcon = lipgloss.NewStyle().Foreground(killColor).SetString("☠")
)

const (
	IconSuccess = "✔"
	IconError   = "✖"
	IconInfo    = "ℹ"
	IconWait    = "⚡"
	IconStep    = "➜"
)

const stringHandler = "%s %s\n"

// Templates do promptui sem o estilo faint.
//
// Causa raiz: os templates padrão do promptui v0.9.0 renderizam a linha
// confirmada/selecionada com faint (SGR 2) — `{{ . | faint }}` — o que
// deixava emojis e texto "apagados" após a escolha (ex: tipos de commit).
// Estas funções devolvem instâncias novas a cada chamada; os demais campos
// continuam com os padrões do promptui. Mantêm o ✔ verde, sem esmaecer.
func newSelectTemplates() *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Selected: promptui.IconGood + ` {{ . }}`,
	}
}

func newPromptTemplates() *promptui.PromptTemplates {
	return &promptui.PromptTemplates{
		Success: `{{ . }}: `,
	}
}

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
	s.Suffix = " " + info(message)
	s.Color("cyan")
	return s
}

func ExecuteCommandSilent(name string, args []string, dir string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd.Run()
}

func showSuccessBox(projectName, stackName string) {
	content := fmt.Sprintf(
		"🚀 Projeto %s criado com sucesso!\n\nStack: %s\nPróximo passo: %s",
		highlight.Render(projectName),
		highlight.Render(stackName),
		highlight.Render("cd "+projectName+" && code ."),
	)

	fmt.Println(successBox.Render(content))
}

func ShowSuccessBox(projectName, stack string) {
	fmt.Println(bold("\n✨ Projeto criado com sucesso!"))
	fmt.Println(color.MagentaString("---------------------------------"))

	fmt.Printf("  %s %s %s\n", success(IconStep), bold("cd"), projectName)

	switch stack {
	case "Go":
		fmt.Printf(stringHandler, success(IconStep), bold("go run main.go"))
	case "Python":
		fmt.Printf(stringHandler, success(IconStep), bold("python main.py"))
	default:
		fmt.Printf(stringHandler, success(IconStep), bold("npm run dev"))
	}
	fmt.Println(color.MagentaString("---------------------------------"))
	fmt.Println(info("  Dúvidas? Acesse nosso GitHub! 🚀\n"))
}

func PrintBanner() {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true)

	asciiArt := `
    ____  _______    ______  ____  _  __
   / __ \/ ____/ |  / / __ )/ __ \| |/ /
  / / / / __/  | | / / __  / / / /   /
 / /_/ / /___  | |/ / /_/ / /_/ /   |
/_____/_____/  |___/_____/\____/_/|_|
`
	fmt.Println(style.Render(asciiArt))
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).PaddingLeft(2).Render("v1.0.0 • Automation Tool"))
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
	}

	fmt.Printf("  %s  %s\n", icon, msg)
}

func HandleError(err error, context string) {
	if err == nil {
		return
	}

	fmt.Println()
	fmt.Printf(stringHandler, errorBanner.Render(" ERROR "), errorContextStyle.Render(context))

	fmt.Printf(stringHandler, errorIcon, errorMessageStyle.Render(err.Error()))

	tip := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true).Render("  💡 Dica: Verifique as permissões ou use 'devbox --help'")
	fmt.Println(tip)
	fmt.Println()
}
