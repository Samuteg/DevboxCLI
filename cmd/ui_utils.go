package cmd

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
)

var (
	info     = color.New(color.FgCyan).SprintFunc()
	success  = color.New(color.FgGreen).SprintFunc()
	warning  = color.New(color.FgYellow).SprintFunc()
	errColor = color.New(color.FgRed).SprintFunc()
	bold     = color.New(color.Bold).SprintFunc()

	// Ícones de Estado
	iconStepTodo   = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).SetString("○") // Círculo cinza
	iconStepActive = lipgloss.NewStyle().Foreground(primaryColor).SetString("●")          // Círculo Roxo
	iconStepDone   = lipgloss.NewStyle().Foreground(successColor).SetString("✔")          // Check Verde

	// Texto dos Passos
	textStepTodo   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	textStepActive = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF")).Bold(true)
	textStepDone   = lipgloss.NewStyle().Foreground(lipgloss.Color("#CCC")).Strikethrough(false)

	// Paleta de Cores
	primaryColor   = lipgloss.Color("#7D56F4")
	grayColor      = lipgloss.Color("#626262")
	secondaryColor = lipgloss.Color("#00ADD8")
	successColor   = lipgloss.Color("#27AE60")
	errorColor     = lipgloss.Color("#E74C3C")

	// Estilo para Mensagens de Sucesso
	successBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(successColor).
			Padding(1, 2).
			Bold(true).
			MarginTop(1)

	// Estilo para Texto em Destaque
	highlight = lipgloss.NewStyle().Foreground(secondaryColor).Bold(true)

	// Cores para o Cleanup
	deleteColor  = lipgloss.Color("#E74C3C") // Vermelho
	neutralColor = lipgloss.Color("242")     // Cinza

	delStyle  = lipgloss.NewStyle().Foreground(deleteColor).Bold(true)
	pathStyle = lipgloss.NewStyle().Foreground(neutralColor).Italic(true)

	summaryBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(0, 2).
			MarginTop(1)

	// Cores para o Add
	addComponentColor = lipgloss.Color("#00ADD8") // Ciano para novos arquivos
	addDirColor       = lipgloss.Color("#F1C40F") // Amarelo para diretórios

	// Estilo da árvore
	treeBranch = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("├──")
	treeLast   = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("└──")

	// Cores por tipo de commit
	featColor     = lipgloss.Color("#A3BE8C") // Verde suave
	fixColor      = lipgloss.Color("#BF616A") // Vermelho/Vinho
	docsColor     = lipgloss.Color("#81A1C1") // Azul gelo
	refactorColor = lipgloss.Color("#B48EAD") // Roxo/Lilás

	// Estilo para a mensagem final de commit no log
	commitScopeStyle = lipgloss.NewStyle().Foreground(secondaryColor).Bold(true)
	commitTypeStyle  = lipgloss.NewStyle().Bold(true).Padding(0, 1).Foreground(lipgloss.Color("#FFF"))
)

// Ícones para feedback visual rápido
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

// --- SPINNER & EXECUÇÃO ---

// NewSpinner cria um carregamento padronizado
func NewSpinner(message string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " " + info(message)
	s.Color("cyan")
	return s
}

// ExecuteCommandSilent roda comando "escondido".
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

	if stack == "Go" {
		fmt.Printf(stringHandler, success(IconStep), bold("go run main.go"))
		if stack == "Python" {
			fmt.Printf(stringHandler, success(IconStep), bold("python main.py"))
		} else {
			fmt.Printf(stringHandler, success(IconStep), bold("npm run dev"))
		}
	}
	fmt.Println(color.MagentaString("---------------------------------"))
	fmt.Println(info("  Dúvidas? Acesse nosso GitHub! 🚀\n"))
}

func PrintBanner() {
	// Gradiente simulado (Texto Roxo + Setinha Ciano)
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
