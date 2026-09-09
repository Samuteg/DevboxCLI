package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Desfaz o último devbox init",
	Long:  "Remove arquivos criados por um devbox init usando snapshots armazenados em .devbox/undo/.",
	Run:   runUndo,
}

func init() {
	rootCmd.AddCommand(undoCmd)
}

func runUndo(cmd *cobra.Command, args []string) {
	undoDir := filepath.Join(".devbox", "undo")
	entries, err := os.ReadDir(undoDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println(lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Render("Nenhum snapshot de undo encontrado."))
			return
		}
		HandleError(err, "Leitura de snapshots")
		return
	}

	type snapshotInfo struct {
		name     string
		path     string
		modTime  os.FileInfo
	}

	var snapshots []snapshotInfo
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		snapshots = append(snapshots, snapshotInfo{
			name:    e.Name(),
			path:    filepath.Join(undoDir, e.Name()),
			modTime: info,
		})
	}

	if len(snapshots) == 0 {
		fmt.Println(lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Render("Nenhum snapshot de undo encontrado."))
		return
	}

	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].name > snapshots[j].name
	})

	latest := snapshots[0]
	manifestPath := filepath.Join(latest.path, "manifest.json")
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		HandleError(fmt.Errorf("leitura do manifest %s: %w", manifestPath, err), "Undo")
		return
	}

	var paths []string
	if err := json.Unmarshal(manifestData, &paths); err != nil {
		HandleError(fmt.Errorf("parse do manifest: %w", err), "Undo")
		return
	}

	fmt.Println(lipgloss.NewStyle().Bold(true).MarginLeft(2).Render("📋 Snapshot encontrado:"))
	fmt.Printf("  %s\n\n", lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(latest.name))

	// List what was created
	fmt.Println(lipgloss.NewStyle().Bold(true).MarginLeft(2).Render("Arquivos que serão removidos:"))
	for _, p := range paths {
		absPath := filepath.Join(".", p)
		icon := "📄"
		if info, err := os.Stat(absPath); err == nil && info.IsDir() {
			icon = "📁"
		}
		fmt.Printf("  %s %s\n", icon, p)
	}
	fmt.Println()

	var confirmed bool
	if err := huh.NewConfirm().
		Title("Remover estes arquivos?").
		Affirmative("Sim, desfazer").
		Negative("Cancelar").
		Value(&confirmed).
		Run(); err != nil {
		if promptAborted(err) {
			os.Exit(0)
		}
		HandleErrorAndExit(err, "Confirmação de undo")
	}
	if !confirmed {
		fmt.Println("Operação cancelada.")
		return
	}

	// Delete paths in reverse order (deepest first)
	for i := len(paths) - 1; i >= 0; i-- {
		absPath := filepath.Join(".", paths[i])
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			continue
		}
		if err := os.RemoveAll(absPath); err != nil {
			LogWarning(fmt.Sprintf("Falha ao remover %s: %v", paths[i], err))
			continue
		}
		fmt.Printf("  %s %s\n",
			lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render("✓"),
			paths[i])
	}

	// Remove snapshot directory
	_ = os.RemoveAll(latest.path)

	// Clean up .devbox/undo if empty
	if remaining, _ := os.ReadDir(undoDir); len(remaining) == 0 {
		_ = os.RemoveAll(undoDir)
		if _, err := os.Stat(".devbox"); err == nil {
			empty, _ := os.ReadDir(".devbox")
			if len(empty) == 0 {
				_ = os.RemoveAll(".devbox")
			}
		}
	}

	fmt.Println()
	fmt.Println(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2")).MarginLeft(2).Render("✓ Desfeito com sucesso!"))
}

// createSnapshot saves a manifest of paths created by devbox init.
// It scans the project directory and records all created entries.
func createSnapshot(projectDir string, createdDirs []string) {
	undoBase := filepath.Join(projectDir, ".devbox", "undo")
	if err := os.MkdirAll(undoBase, 0755); err != nil {
		LogWarning(fmt.Sprintf("Não foi possível criar diretório de undo: %v", err))
		return
	}

	ts := timeNow()
	snapshotDir := filepath.Join(undoBase, ts)
	if err := os.MkdirAll(snapshotDir, 0755); err != nil {
		LogWarning(fmt.Sprintf("Não foi possível criar snapshot: %v", err))
		return
	}

	var paths []string

	// Add explicitly created directories first (in creation order)
	for _, d := range createdDirs {
		paths = append(paths, d)
	}

	// Scan all files created by MaterializeTemplates
	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || path == projectDir {
			return nil
		}
		rel, err := filepath.Rel(".", path)
		if err != nil {
			return nil
		}
		// Skip .devbox itself — it's our undo metadata
		if strings.HasPrefix(rel, filepath.Join(projectDir, ".devbox")) {
			return filepath.SkipDir
		}
		paths = append(paths, rel)
		return nil
	})
	if err != nil {
		LogWarning(fmt.Sprintf("Escaneamento para undo: %v", err))
		return
	}

	// Sort: files first (by length, deeper paths first), then dirs
	sort.Slice(paths, func(i, j int) bool {
		iDepth := strings.Count(paths[i], string(filepath.Separator))
		jDepth := strings.Count(paths[j], string(filepath.Separator))
		if iDepth != jDepth {
			return iDepth > jDepth
		}
		return paths[i] > paths[j]
	})

	data, err := json.MarshalIndent(paths, "", "  ")
	if err != nil {
		LogWarning(fmt.Sprintf("Serialização do manifest: %v", err))
		return
	}

	manifestPath := filepath.Join(snapshotDir, "manifest.json")
	if err := os.WriteFile(manifestPath, data, 0644); err != nil {
		LogWarning(fmt.Sprintf("Escrita do manifest: %v", err))
	}
}

func timeNow() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
}
