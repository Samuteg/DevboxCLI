package cmd

import (
	"fmt"
	"os"
	"sync" // Essencial para concorrência
	"time"

	"github.com/spf13/cobra"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Limpeza profunda de arquivos desnecessários",
	Run: func(cmd *cobra.Command, args []string) {
		s := NewSpinner(info("Limpando a casa..."))
		s.Start()

		start := time.Now()
		folders := []string{"node_modules", "dist", ".next", "bin"}

		var wg sync.WaitGroup
		for _, folder := range folders {
			if _, e := os.Stat(folder); e == nil {
				wg.Add(1)
				go func(f string) {
					defer wg.Done()
					os.RemoveAll(f)
				}(folder)
			}
		}

		wg.Wait()
		s.Stop() // Para o spinner exatamente quando o WaitGroup libera

		duration := time.Since(start)
		fmt.Printf("%s Limpeza concluída em %v\n", success("✔"), bold(duration))
	},
}

func init() {
	rootCmd.AddCommand(cleanupCmd)
}
