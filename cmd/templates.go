package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/Samuteg/DevboxCLI/internal/scaffold"
	"github.com/spf13/cobra"
)

var templatesCmd = &cobra.Command{
	Use:     "templates",
	Short:   "Lista templates disponíveis (built-in e overrides)",
	Example: "  devbox templates",
	Run:     runTemplates,
}

func runTemplates(cmd *cobra.Command, args []string) {
	stacks := scaffold.DefaultStacks()

	type tpl struct {
		name     string
		source   string
		override bool
	}

	var built []tpl
	seen := map[string]bool{}

	for stackName, s := range stacks {
		if s.IsBackend {
			sources := []string{s.Source}
			for _, v := range s.Variants {
				sources = append(sources, v.Source)
			}
			for _, src := range sources {
				key := stackName + ":" + filepath.Base(src)
				if seen[key] {
					continue
				}
				seen[key] = true
				built = append(built, tpl{
					name:     fmt.Sprintf("%s/%s", stackName, filepath.Base(src)),
					source:   src,
					override: scaffold.HasOverride(src),
				})
			}
		}
	}

	sort.Slice(built, func(i, j int) bool { return built[i].name < built[j].name })

	if jsonOutput {
		type jsonTpl struct {
			Name     string `json:"name"`
			Source   string `json:"source"`
			Override bool   `json:"override"`
		}
		var items []jsonTpl
		for _, t := range built {
			items = append(items, jsonTpl{Name: t.name, Source: t.source, Override: t.override})
		}
		printJSON(JSONResult{Success: true, Command: "templates", Data: items})
		return
	}

	fmt.Println(bold("Templates disponíveis:"))
	fmt.Println()
	for _, t := range built {
		tag := ""
		if t.override {
			tag = " " + success("[override]")
		}
		fmt.Printf("  %s  %s%s\n", info("•"), t.name, tag)
	}

	fmt.Println()
	home, err := os.UserHomeDir()
	if err == nil {
		dir := filepath.Join(home, ".devbox", "templates")
		if scaffold.LoadOverrideFS() != nil {
			fmt.Printf("  %s %s\n", info("Override dir:"), dir)
		} else {
			fmt.Printf("  %s %s\n", info("Override dir:"), pathStyle.Render("(não configurado)"))
			fmt.Printf("  %s\n", pathStyle.Render("Crie ~/.devbox/templates/ para adicionar overrides"))
		}
	}
}

func init() {
	registerDual(templatesCmd)
}
