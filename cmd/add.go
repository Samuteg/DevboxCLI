package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
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

	// 1. Resolve o TIPO
	if len(args) > 0 {
		resourceType = strings.ToLower(args[0])
	}
	// Validação simples: se não for válido ou estiver vazio, pergunta
	if !isValidType(resourceType) {
		if resourceType != "" {
			fmt.Printf("⚠️  Tipo '%s' desconhecido.\n", resourceType)
		}
		resourceType = promptSelect("O que você quer criar?", validTypes)
	}

	// 2. Resolve o NOME
	if len(args) > 1 {
		resourceName = args[1]
	} else {
		resourceName = promptInput("Qual o nome do componente?", "O nome é obrigatório", 2)
	}

	// 3. Executa a criação real
	err := createComponent(resourceType, resourceName)
	if err != nil {
		fmt.Printf("❌ Erro ao criar componente: %v\n", err)
		os.Exit(1)
	}
}

// --- Lógica de Criação de Arquivos ---

func createComponent(rType, rName string) error {
	// Formata o nome (ex: user -> User)
	modelName := strings.Title(strings.ToLower(rName)) // Ex: User
	fileName := strings.ToLower(rName)                 // Ex: user

	var path, contentTemplate string

	// Define o caminho e o conteúdo baseado no tipo
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
		return fmt.Errorf("tipo não implementado: %s", rType)
	}

	// 1. Garante que a pasta existe
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("erro ao criar pasta: %w", err)
	}

	// 2. Verifica se o arquivo já existe para não sobrescrever
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("o arquivo %s já existe", path)
	}

	// 3. Prepara os dados para o template
	data := map[string]string{
		"Name":      modelName, // User
		"LowerName": fileName,  // user
	}

	// 4. Processa o template e escreve no arquivo
	tmpl, err := template.New("component").Parse(contentTemplate)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return err
	}

	fmt.Printf("✅ Arquivo criado com sucesso: %s\n", path)
	return nil
}

// --- Templates Embutidos (Strings) ---

const controllerTmpl = `package controllers

import "github.com/gin-gonic/gin"

type {{.Name}}Controller struct{}

func New{{.Name}}Controller() *{{.Name}}Controller {
	return &{{.Name}}Controller{}
}

func (c *{{.Name}}Controller) Create(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "Create {{.Name}}"})
}
`

const usecaseTmpl = `package usecase

type {{.Name}}UseCase struct{}

func New{{.Name}}UseCase() *{{.Name}}UseCase {
	return &{{.Name}}UseCase{}
}

func (u *{{.Name}}UseCase) Execute() error {
	return nil
}
`

const repositoryTmpl = `package repository

type {{.Name}}Repository interface {
	Save() error
}
`

const handlerTmpl = `package handlers

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
