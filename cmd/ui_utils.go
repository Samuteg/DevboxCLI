package cmd

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

// --- COMPATIBILIDADE (Mantém seu código antigo funcionando) ---
var (
	// Estas variáveis existiam no seu código antigo.
	// Eu as mantive aqui para que arquivos como root.go ou init.go não quebrem.
	info     = color.New(color.FgCyan).SprintFunc()
	success  = color.New(color.FgGreen).SprintFunc()
	warning  = color.New(color.FgYellow).SprintFunc()
	errColor = color.New(color.FgRed).SprintFunc() // Renomeei levemente para evitar conflito com tipo error
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

// LogSuccess: Imprime uma mensagem de sucesso padronizada com ícone
func LogSuccess(message string) {
	fmt.Printf("%s %s\n", success(IconSuccess), message)
}

// LogError: Imprime erro padronizado
func LogError(message string) {
	fmt.Printf("%s %s\n", errColor(IconError), errColor(message))
}

// LogInfo: Imprime informação padronizada
func LogInfo(message string) {
	fmt.Printf("%s %s\n", info(IconInfo), message)
}

// LogWarning: Imprime aviso
func LogWarning(message string) {
	fmt.Printf("%s %s\n", warning("!"), warning(message))
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
// Mudei para Exportado (Letra Maiúscula) e retornando erro.
// Se seu código antigo chama "executeCommandSilent" (minúscula), você precisará mudar lá.
func ExecuteCommandSilent(name string, args []string, dir string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd.Run()
}

// showSuccessBox mantido com letra minúscula se seu init.go o chama assim,
// mas recomendo usar a versão nova ShowSuccessBox abaixo para projetos futuros.
func showSuccessBox(projectName, stack string) {
	ShowSuccessBox(projectName, stack) // Apenas repassa para a nova
}

// ShowSuccessBox (Versão Nova e Melhorada)
func ShowSuccessBox(projectName, stack string) {
	fmt.Println(bold("\n✨ Projeto criado com sucesso!"))
	fmt.Println(color.MagentaString("---------------------------------"))

	fmt.Printf("  %s %s %s\n", success(IconStep), bold("cd"), projectName)

	if stack == "Go" {
		fmt.Printf("  %s %s\n", success(IconStep), bold("go run main.go"))
	} else {
		fmt.Printf("  %s %s\n", success(IconStep), bold("npm run dev"))
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
	fmt.Printf("\n%s %s\n", success("●"), bold("DevBox CLI v0.0.1"))
	fmt.Printf("%s %s\n\n", info("ℹ"), "Pronto para otimizar sua rotina.\n")
}
