package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// --- Estilos Específicos do Help ---
var (
	// Títulos das seções (USAGE, COMMANDS, FLAGS)
	sectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			MarginTop(1).
			MarginBottom(0).
			PaddingLeft(1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderLeft(true).
			BorderForeground(primaryColor)

	// Estilo do comando (coluna da esquerda)
	commandStyle = lipgloss.NewStyle().
			Foreground(secondaryColor). // Ciano
			Bold(true).
			Width(20) // Largura fixa para alinhar

	// Estilo da descrição (coluna da direita)
	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("246")) // Cinza

	// Estilo das Flags
	flagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F1C40F")) // Amarelo para flags (--help)

	// Exemplo de uso
	usageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF")).
			Background(lipgloss.Color("#333")).
			Padding(0, 1)
)

func helpFunc(cmd *cobra.Command, args []string) {
	if cmd.Long != "" {
		fmt.Println(lipgloss.NewStyle().MarginLeft(2).Width(60).Render(cmd.Long))
	} else {
		fmt.Println(lipgloss.NewStyle().MarginLeft(2).Render(cmd.Short))
	}
	fmt.Println()

	fmt.Println(sectionStyle.Render("USAGE"))

	useLine := cmd.UseLine()
	if !strings.HasPrefix(useLine, "devbox") {
		useLine = "devbox " + useLine
	}
	fmt.Printf("  %s\n", usageStyle.Render(useLine))

	if len(cmd.Commands()) > 0 {
		fmt.Println(sectionStyle.Render("COMMANDS"))

		for _, c := range cmd.Commands() {
			if !c.IsAvailableCommand() || c.Hidden {
				continue
			}

			fmt.Printf("  %s%s\n",
				commandStyle.Render(c.Name()),
				descStyle.Render(c.Short),
			)
		}
	}

	if cmd.Flags().HasFlags() {
		fmt.Println(sectionStyle.Render("FLAGS"))

		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			var name string
			if f.Shorthand != "" {
				name = fmt.Sprintf("-%s, --%s", f.Shorthand, f.Name)
			} else {
				name = fmt.Sprintf("    --%s", f.Name)
			}

			fmt.Printf("  %s%s\n",
				flagStyle.Width(20).Render(name),
				descStyle.Render(f.Usage),
			)
		})
	}

	fmt.Println()
	footer := lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("240")).Render("  Use 'devbox [command] --help' para mais informações.")
	fmt.Println(footer)
	fmt.Println()
}
