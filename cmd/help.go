package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	sectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginTop(1).
			MarginBottom(0).
			PaddingLeft(1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderLeft(true).
			BorderForeground(ColorPrimary)

	commandStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Bold(true).
			Width(24)

	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("246"))

	flagStyle = lipgloss.NewStyle().
			Foreground(ColorWarning)

	usageStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Bold(true).
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
	fmt.Printf("  %s\n", usageStyle.Render(useLine))

	if len(cmd.Aliases) > 0 {
		fmt.Println(sectionStyle.Render("ALIASES"))
		for _, alias := range cmd.Aliases {
			fmt.Printf("  %s%s\n",
				commandStyle.Render(alias),
				descStyle.Render("(alias de "+cmd.Name()+")"),
			)
		}
	}

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

	if cmd.Example != "" {
		fmt.Println(sectionStyle.Render("EXAMPLES"))
		fmt.Println(lipgloss.NewStyle().MarginLeft(2).Foreground(ColorSubtle).Render(cmd.Example))
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
				flagStyle.Width(24).Render(name),
				descStyle.Render(f.Usage),
			)
		})
	}

	fmt.Println()
	footer := lipgloss.NewStyle().Italic(true).Foreground(ColorSubtle).Render("  Use '" + cmd.CommandPath() + " --help' para detalhes (ex: 'devbox init --help').")
	fmt.Println(footer)
	fmt.Println()
}
