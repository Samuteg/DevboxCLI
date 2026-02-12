package cmd

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

var (
	info     = color.New(color.FgCyan).SprintFunc()
	success  = color.New(color.FgGreen).SprintFunc()
	warning  = color.New(color.FgYellow).SprintFunc()
	errColor = color.New(color.FgRed).SprintFunc()
	bold     = color.New(color.Bold).SprintFunc()
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

func showSuccessBox(projectName, stack string) {
	ShowSuccessBox(projectName, stack)
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
	banner := `
    ____  _______    ______  ____ _  __
   / __ \/ ____/ |  / / __ )/ __ \ |/ /
  / / / / __/  | | / / __  / / / /   / 
 / /_/ / /___  | |/ / /_/ / /_/ /   |  
/_____/_____/  |___/_____/\____/_/|_|  `

	fmt.Println(info(banner))
	fmt.Printf("\n%s %s\n", success("●"), bold("DevBox CLI v1.0.0"))
	fmt.Printf("%s %s\n\n", info("ℹ"), "Pronto para otimizar sua rotina.\n")
}
