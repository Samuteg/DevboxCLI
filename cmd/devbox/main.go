package main

import (
	"os"

	"github.com/Samuteg/DevboxCLI/cmd"
	"github.com/Samuteg/DevboxCLI/internal/term"
)

func main() {
	// Banner só em terminal interativo: não poluir pipes/CI.
	if !term.Plain() && !hasHelpFlag() {
		cmd.PrintBanner()
	}
	cmd.Execute()
}

func hasHelpFlag() bool {
	for _, a := range os.Args[1:] {
		if a == "--help" || a == "-h" || a == "help" || a == "--version" || a == "-v" {
			return true
		}
	}
	return false
}
