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

	iconStepTodo   = lipgloss.NewStyle().Foreground(ColorSubtle).SetString("○")
	iconStepActive = lipgloss.NewStyle().Foreground(ColorPrimary).SetString("●")
	iconStepDone   = lipgloss.NewStyle().Foreground(ColorSuccess).SetString("✔")

	textStepTodo   = lipgloss.NewStyle().Foreground(ColorSubtle)
	textStepActive = lipgloss.NewStyle().Foreground(ColorWhite).Bold(true)
	textStepDone   = lipgloss.NewStyle().Foreground(lipgloss.Color("#CCC")).Strikethrough(false)

	successBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSuccess).
			Padding(1, 2).
			Bold(true).
			MarginTop(1)

	highlight = lipgloss.NewStyle().Foreground(ColorSecondary).Bold(true)

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

// Aliases de compatibilidade temporários: outros comandos ainda usam os nomes
// antigos e serão migrados nas tasks seguintes. Não usar em código novo —
// prefira os tokens Color*.
var (
	primaryColor      = ColorPrimary
	secondaryColor    = ColorSecondary
	successColor      = ColorSuccess
	featColor         = ColorFeat
	fixColor          = ColorFix
	docsColor         = ColorDocs
	refactorColor     = ColorRefactor
	addComponentColor = ColorSecondary
	addDirColor       = ColorWarning
	targetColor       = ColorStyle
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
	style := lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)

	asciiArt := `
    ____  _______    ______  ____  _  __
   / __ \/ ____/ |  / / __ )/ __ \| |/ /
  / / / / __/  | | / / __  / / / /   /
 / /_/ / /___  | |/ / /_/ / /_/ /   |
/_____/_____/  |___/_____/\____/_/|_|
`
	fmt.Println(style.Render(asciiArt))
	fmt.Println(lipgloss.NewStyle().Foreground(ColorSubtle).PaddingLeft(2).Render("v1.0.0 • Automation Tool"))
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

	tip := lipgloss.NewStyle().Foreground(ColorSubtle).Italic(true).Render("  💡 Dica: Verifique as permissões ou use 'devbox --help'")
	fmt.Println(tip)
	fmt.Println()
}
