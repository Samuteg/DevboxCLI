package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Gere as configurações da Devbox",
}

var configSetCmd = &cobra.Command{
	Use:   "set [chave] [valor]",
	Short: "Define um valor de configuração (ex: config set author \"João\")",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := strings.ToLower(args[0])
		value := args[1]

		// Guarda em memória e depois escreve no ficheiro
		viper.Set(key, value)
		if err := viper.WriteConfig(); err != nil {
			HandleError(err, "Guardar Configuração")
			return
		}

		printStep("done", "Configuração atualizada!")
		fmt.Printf("  %s %s\n\n",
			lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(fmt.Sprintf("%s =", key)),
			lipgloss.NewStyle().Foreground(secondaryColor).Bold(true).Render(value),
		)
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get [chave]",
	Short: "Recupera um valor de configuração (ex: config get author)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := strings.ToLower(args[0])

		value := viper.GetString(key)
		if value == "" {
			value = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("240")).Render("<vazio>")
		} else {
			value = lipgloss.NewStyle().Foreground(secondaryColor).Bold(true).Render(value)
		}

		fmt.Printf("  %s %s\n",
			lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(fmt.Sprintf("%s:", key)),
			value,
		)
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista todas as configurações definidas",
	Run: func(cmd *cobra.Command, args []string) {
		printStep("active", "Lendo ~/.devbox.yaml")

		settings := viper.AllSettings()
		if len(settings) == 0 {
			printStep("todo", "Nenhuma configuração definida.")
			return
		}

		fmt.Println()
		// Desenhamos uma tabela simples ou uma lista alinhada
		for k, v := range settings {
			kStyle := lipgloss.NewStyle().Width(15).Foreground(lipgloss.Color("242")).Render(k)
			vStyle := lipgloss.NewStyle().Foreground(secondaryColor).Render(fmt.Sprintf("%v", v))
			fmt.Printf("  %s %s\n", kStyle, vStyle)
		}
		fmt.Println()
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configListCmd)

	rootCmd.AddCommand(configCmd)
}
