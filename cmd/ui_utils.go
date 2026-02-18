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

	// Estilo para o nome do arquivo na árvore
	fileStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ADD8"))            // Ciano
	dirStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#F4D03F")).Bold(true) // Amarelo

	// Paleta de Cores
	primaryColor   = lipgloss.Color("#7D56F4")
	grayColor      = lipgloss.Color("#626262")
	secondaryColor = lipgloss.Color("#00ADD8")
	successColor   = lipgloss.Color("#27AE60")
	errorColor     = lipgloss.Color("#E74C3C")

	// Estilo do Título (Banner)
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).         // Cor do texto
			Border(lipgloss.RoundedBorder()). // Borda arredondada
			BorderForeground(primaryColor).   // Cor da borda
			Padding(0, 1).                    // Espaço interno (respiro)
			MarginLeft(1)                     // Espaço da esquerda da tela1)

	// Estilo para Mensagens de Sucesso
	successBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(successColor).
			Padding(1, 2).
			Bold(true).
			MarginTop(1)

	// Estilo para Texto em Destaque
	highlight = lipgloss.NewStyle().Foreground(secondaryColor).Bold(true)
)

// --- NOVAS FUNCIONALIDADES (Use estas daqui para frente) ---

// Ícones para feedback visual rápido
const (
	IconSuccess = "✔"
	IconError   = "✖"
	IconInfo    = "ℹ"
	IconWait    = "⚡"
	IconStep    = "➜"
)

const stringHandler = "%s %s\n"

// LogSuccess: Imprime uma mensagem de sucesso padronizada com ícone
func LogSuccess(message string) {
	fmt.Printf(stringHandler, success(IconSuccess), message)
}

// LogError: Imprime erro padronizado
func LogError(message string) {
	fmt.Printf(stringHandler, errColor(IconError), errColor(message))
}

// LogInfo: Imprime informação padronizada
func LogInfo(message string) {
	fmt.Printf(stringHandler, info(IconInfo), message)
}

// LogWarning: Imprime aviso
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

// ShowSuccessBox (Versão Nova e Melhorada)
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

// Função auxiliar para imprimir um passo
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
