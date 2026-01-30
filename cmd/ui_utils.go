package cmd

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

// Cores pré-definidas
var (
	info    = color.New(color.FgCyan).SprintFunc()
	success = color.New(color.FgGreen).SprintFunc()
	warning = color.New(color.FgYellow).SprintFunc()
	err     = color.New(color.FgRed).SprintFunc()
	bold    = color.New(color.Bold).SprintFunc()
)

// NewSpinner cria um carregamento padronizado
func NewSpinner(message string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond) // O 14 é um set bem moderno
	s.Suffix = " " + message
	s.Color("cyan")
	return s
}

func showSuccessBox(projectName, stack string) {
	fmt.Println(bold("\n  Próximos passos:"))
	fmt.Printf("  1. cd %s\n", projectName)
	if stack == "Go" {
		fmt.Println("  2. go run main.go")
	} else {
		fmt.Println("  2. npm run dev")
	}
	fmt.Println()
}

// executeCommandSilent roda o comando sem poluir o terminal, ideal para usar com Spinner
func executeCommandSilent(name string, args []string, dir string) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	// Redirecionamos a saída para nulo para não quebrar o Spinner
	cmd.Run()
}

func PrintBanner() {
	banner := `
    ____  _______    ______  ____ _  __
   / __ \/ ____/ |  / / __ )/ __ \ |/ /
  / / / / __/  | | / / __  / / / /   / 
 / /_/ / /___  | |/ / /_/ / /_/ /   |  
/_____/_____/  |___/_____/\____/_/|_|  
                                       `
	color.Cyan(banner)
	fmt.Printf("%s %s\n\n", info("v1.0.0"), bold("| Sua caixa de ferramentas de desenvolvimento"))
}
