package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var validTypes = []string{"controller", "usecase", "repository", "handler"}

var addCmd = &cobra.Command{
	Use:   "add [tipo] [nome]",
	Short: "Adiciona um novo componente ao projeto",
	Args:  cobra.MaximumNArgs(2),
	Run:   runAdd,
}

func runAdd(cmd *cobra.Command, args []string) {
	var resourceType, resourceName string

	if len(args) > 0 {
		resourceType = strings.ToLower(args[0])
	}
	if !isValidType(resourceType) {
		if resourceType != "" {
			fmt.Printf("  %s \n", lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Render("Tipo '"+resourceType+"' desconhecido."))
		}
		resourceType = promptSelect("  O que você quer criar?", validTypes)
	}

	if len(args) > 1 {
		resourceName = args[1]
	} else {
		resourceName = promptInput("  Qual o nome do componente?", "O nome é obrigatório", 2)
	}

	printStep("active", fmt.Sprintf("Gerando %s: %s", resourceType, resourceName))

	path, err := createComponent(resourceType, resourceName)
	if err != nil {
		os.Exit(1)
	}

	printStep("done", "Componente criado com sucesso!")

	fmt.Println()
	fmt.Println(lipgloss.NewStyle().Bold(true).MarginLeft(2).Render("📂 Arquivo gerado:"))
	renderDynamicTree(path)

	showSuccessBox(resourceName, strings.Title(resourceType))
}

// --- Lógica de Criação de Arquivos (Modificada para retornar o path) ---

func createComponent(rType, rName string) (string, error) {
	modelName := strings.Title(strings.ToLower(rName))
	fileName := strings.ToLower(rName)

	var path, contentTemplate string

	switch rType {
	case "controller":
		path = filepath.Join("internal", "controllers", fileName+"_controller.go")
		contentTemplate = controllerTmpl
	case "usecase":
		path = filepath.Join("internal", "usecase", fileName+"_usecase.go")
		contentTemplate = usecaseTmpl
	case "repository":
		path = filepath.Join("internal", "repository", fileName+"_repository.go")
		contentTemplate = repositoryTmpl
	case "handler":
		path = filepath.Join("internal", "handlers", fileName+"_handler.go")
		contentTemplate = handlerTmpl
	default:
		return "", fmt.Errorf("tipo não implementado: %s", rType)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", err
	}

	if _, err := os.Stat(path); err == nil {
		return "", fmt.Errorf("o arquivo %s já existe", path)
	}

	authorName := viper.GetString("author")

	data := map[string]string{
		"Name":      modelName,
		"LowerName": fileName,
		"Author":    authorName,
	}

	tmpl, err := template.New("component").Parse(contentTemplate)
	if err != nil {
		return "", err
	}

	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return "", err
	}

	return path, nil
}

// --- UI: Árvore Dinâmica ---

func renderDynamicTree(path string) {
	parts := strings.Split(filepath.ToSlash(path), "/")

	folderStyle := lipgloss.NewStyle().Foreground(addDirColor).Bold(true)
	fileStyle := lipgloss.NewStyle().Foreground(addComponentColor)
	indent := "  "

	for i, part := range parts {
		isLast := i == len(parts)-1

		if i == 0 {
			fmt.Printf("%s%s\n", indent, folderStyle.Render(part+"/"))
		} else if !isLast {
			fmt.Printf("%s%s %s\n", indent, treeBranch, folderStyle.Render(part+"/"))
			indent += "│   "
		} else {
			fmt.Printf("%s%s %s\n", indent, treeLast, fileStyle.Render(part))
		}
	}
	fmt.Println()
}

// --- Templates Embutidos (Strings) ---

const controllerTmpl = `

/ Criado por: {{.Author}}

// Gerado via Devbox CLI

package controllers

import "github.com/gin-gonic/gin"

type {{.Name}}Controller struct{}

func New{{.Name}}Controller() *{{.Name}}Controller {
	return &{{.Name}}Controller{}
}

func (c *{{.Name}}Controller) Create(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "Create {{.Name}}"})
}
`

const usecaseTmpl = `
/ Criado por: {{.Author}}

// Gerado via Devbox CLI

package usecase

type {{.Name}}UseCase struct{}

func New{{.Name}}UseCase() *{{.Name}}UseCase {
	return &{{.Name}}UseCase{}
}

func (u *{{.Name}}UseCase) Execute() error {
	return nil
}
`

const repositoryTmpl = `
/ Criado por: {{.Author}}

// Gerado via Devbox CLI

package repository

type {{.Name}}Repository interface {
	Save() error
}
`

const handlerTmpl = `
/ Criado por: {{.Author}}

// Gerado via Devbox CLI

package handlers

import "net/http"

func Get{{.Name}}(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello {{.Name}}"))
}
`

func isValidType(t string) bool {
	for _, v := range validTypes {
		if v == t {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(addCmd)
}
