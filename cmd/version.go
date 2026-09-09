package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version é injetada via ldflags no release (-X ...Version={{.Version}}).
// Fallback local para builds sem ldflags.
var Version = "1.0.1-dev"

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Mostra a versão do Devbox CLI",
	Example: "  devbox version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("devbox v%s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
