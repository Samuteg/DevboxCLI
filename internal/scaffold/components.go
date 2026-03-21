package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var ValidComponentTypes = []string{"controller", "usecase", "repository", "handler"}

func IsValidComponentType(componentType string) bool {
	for _, candidate := range ValidComponentTypes {
		if candidate == componentType {
			return true
		}
	}
	return false
}

func CreateComponent(componentType, componentName, authorName string) (string, error) {
	modelName := toModelName(componentName)
	fileName := strings.ToLower(componentName)

	path, contentTemplate, err := resolveComponentTemplate(componentType, fileName)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", err
	}

	if _, err := os.Stat(path); err == nil {
		return "", fmt.Errorf("o arquivo %s já existe", path)
	}

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

func resolveComponentTemplate(componentType, fileName string) (string, string, error) {
	switch componentType {
	case "controller":
		return filepath.Join("internal", "controllers", fileName+"_controller.go"), controllerTemplate, nil
	case "usecase":
		return filepath.Join("internal", "usecase", fileName+"_usecase.go"), usecaseTemplate, nil
	case "repository":
		return filepath.Join("internal", "repository", fileName+"_repository.go"), repositoryTemplate, nil
	case "handler":
		return filepath.Join("internal", "handlers", fileName+"_handler.go"), handlerTemplate, nil
	default:
		return "", "", fmt.Errorf("tipo não implementado: %s", componentType)
	}
}

func toModelName(rawName string) string {
	parts := strings.FieldsFunc(strings.ToLower(strings.TrimSpace(rawName)), func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})

	var b strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		b.WriteString(strings.ToUpper(part[:1]))
		if len(part) > 1 {
			b.WriteString(part[1:])
		}
	}

	if b.Len() == 0 {
		return "Component"
	}

	return b.String()
}

const controllerTemplate = `

// Criado por: {{.Author}}

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

const usecaseTemplate = `
// Criado por: {{.Author}}

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

const repositoryTemplate = `
// Criado por: {{.Author}}

// Gerado via Devbox CLI

package repository

type {{.Name}}Repository interface {
	Save() error
}
`

const handlerTemplate = `
// Criado por: {{.Author}}

// Gerado via Devbox CLI

package handlers

import "net/http"

func Get{{.Name}}(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello {{.Name}}"))
}
`
