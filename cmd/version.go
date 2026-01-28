package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Mostra a versão atual da Devbox",
	Run: func(cmd *cobra.Command, args []string) {
		// Aqui você pode imprimir o banner de novo se quiser, ou apenas a versão
		fmt.Printf("🚀 DevBox CLI\n")
		fmt.Printf("Versão: %s\n", AppVersion)
		fmt.Printf("Ambiente: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
