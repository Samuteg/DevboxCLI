package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

// Estrutura para organizar os checks
type Dependency struct {
	Name    string
	Command string
	Args    []string
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Verifica se o sistema possui as dependências necessárias",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("🩺 Iniciando check-up do sistema (%s/%s)...\n\n", runtime.GOOS, runtime.GOARCH)

		dependencies := []Dependency{
			{Name: "Git", Command: "git", Args: []string{"--version"}},
			{Name: "Docker", Command: "docker", Args: []string{"--version"}},
			{Name: "Go Compiler", Command: "go", Args: []string{"version"}},
			{Name: "Node.js", Command: "node", Args: []string{"-v"}},
		}

		allOk := true

		for _, dep := range dependencies {
			// Tenta encontrar o executável no PATH
			path, err := exec.LookPath(dep.Command)

			if err != nil {
				fmt.Printf("❌ %-12s | Não encontrado no PATH\n", dep.Name)
				allOk = false
				continue
			}

			// Tenta executar para garantir que está funcionando
			out, err := exec.Command(path, dep.Args...).Output()
			if err != nil {
				fmt.Printf("⚠️  %-12s | Instalado, mas retornou erro ao executar\n", dep.Name)
				allOk = false
			} else {
				fmt.Printf("✅ %-12s | %s", dep.Name, string(out))
			}
		}

		fmt.Println("\n---")
		if allOk {
			fmt.Println("🎉 Tudo pronto! Seu ambiente de desenvolvimento está saudável.")
		} else {
			fmt.Println("💡 Algumas ferramentas estão faltando ou mal configuradas. Instale-as para garantir o funcionamento total da DevBox.")
		}
		home, _ := os.UserHomeDir()
		configPath := home + "/.devbox.yaml"
		if _, err := os.Stat(configPath); err == nil {
			fmt.Println("✅ Config      | Arquivo .devbox.yaml encontrado")
		} else {
			fmt.Println("ℹ️  Config      | Arquivo .devbox.yaml não encontrado (usando padrões)")
		}
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
