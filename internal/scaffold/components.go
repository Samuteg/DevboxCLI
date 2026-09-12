package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"github.com/Samuteg/DevboxCLI/internal/validation"
)

type Language string

const (
	LangGo       Language = "go"
	LangNodeJS   Language = "javascript"
	LangNodeTS   Language = "typescript"
	LangPython   Language = "python"
	LangRuby     Language = "ruby"
	LangJava     Language = "java"
)

var ValidComponentTypes = []string{"controller", "usecase", "repository", "handler"}

func IsValidComponentType(componentType string) bool {
	return slices.Contains(ValidComponentTypes, componentType)
}

func DetectLanguage(dir string) Language {
	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
		return LangGo
	}
	if _, err := os.Stat(filepath.Join(dir, "tsconfig.json")); err == nil {
		return LangNodeTS
	}
	if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
		return LangNodeJS
	}
	if _, err := os.Stat(filepath.Join(dir, "requirements.txt")); err == nil {
		return LangPython
	}
	if _, err := os.Stat(filepath.Join(dir, "pyproject.toml")); err == nil {
		return LangPython
	}
	if _, err := os.Stat(filepath.Join(dir, "Gemfile")); err == nil {
		return LangRuby
	}
	if _, err := os.Stat(filepath.Join(dir, "pom.xml")); err == nil {
		return LangJava
	}
	return LangGo
}

func CreateComponent(componentType, componentName, authorName string) (string, error) {
	if !validation.IsValidComponentName(componentName) {
		return "", fmt.Errorf("nome de componente inválido %q: use 2-64 caracteres alfanuméricos, '-', '_' (sem path ou espaços)", componentName)
	}
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
	lang := DetectLanguage(".")
	return resolveForLanguage(componentType, fileName, lang)
}

func resolveForLanguage(componentType, fileName string, lang Language) (string, string, error) {
	switch lang {
	case LangGo:
		return resolveGo(componentType, fileName)
	case LangNodeJS:
		return resolveNodeJS(componentType, fileName)
	case LangNodeTS:
		return resolveNodeTS(componentType, fileName)
	case LangPython:
		return resolvePython(componentType, fileName)
	case LangRuby:
		return resolveRuby(componentType, fileName)
	case LangJava:
		return resolveJava(componentType, fileName)
	default:
		return resolveGo(componentType, fileName)
	}
}

func resolveGo(componentType, fileName string) (string, string, error) {
	switch componentType {
	case "controller":
		return filepath.Join("internal", "controllers", fileName+"_controller.go"), controllerGoTemplate, nil
	case "usecase":
		return filepath.Join("internal", "usecase", fileName+"_usecase.go"), usecaseGoTemplate, nil
	case "repository":
		return filepath.Join("internal", "repository", fileName+"_repository.go"), repositoryGoTemplate, nil
	case "handler":
		return filepath.Join("internal", "handlers", fileName+"_handler.go"), handlerGoTemplate, nil
	default:
		return "", "", fmt.Errorf("tipo não implementado: %s", componentType)
	}
}

func resolveNodeJS(componentType, fileName string) (string, string, error) {
	switch componentType {
	case "controller":
		return filepath.Join("src", "controllers", fileName+".js"), controllerJSTemplate, nil
	case "usecase":
		return filepath.Join("src", "usecases", fileName+".js"), usecaseJSTemplate, nil
	case "repository":
		return filepath.Join("src", "repositories", fileName+".js"), repositoryJSTemplate, nil
	case "handler":
		return filepath.Join("src", "handlers", fileName+".js"), handlerJSTemplate, nil
	default:
		return "", "", fmt.Errorf("tipo não implementado: %s", componentType)
	}
}

func resolveNodeTS(componentType, fileName string) (string, string, error) {
	switch componentType {
	case "controller":
		return filepath.Join("src", "controllers", fileName+".ts"), controllerTSTemplate, nil
	case "usecase":
		return filepath.Join("src", "usecases", fileName+".ts"), usecaseTSTemplate, nil
	case "repository":
		return filepath.Join("src", "repositories", fileName+".ts"), repositoryTSTemplate, nil
	case "handler":
		return filepath.Join("src", "handlers", fileName+".ts"), handlerTSTemplate, nil
	default:
		return "", "", fmt.Errorf("tipo não implementado: %s", componentType)
	}
}

func resolvePython(componentType, fileName string) (string, string, error) {
	switch componentType {
	case "controller":
		return filepath.Join("src", "api", "routes", fileName+".py"), controllerPythonTemplate, nil
	case "usecase":
		return filepath.Join("src", "usecases", fileName+".py"), usecasePythonTemplate, nil
	case "repository":
		return filepath.Join("src", "repository", "crud", fileName+".py"), repositoryPythonTemplate, nil
	case "handler":
		return filepath.Join("src", "api", "handlers", fileName+".py"), handlerPythonTemplate, nil
	default:
		return "", "", fmt.Errorf("tipo não implementado: %s", componentType)
	}
}

func resolveRuby(componentType, fileName string) (string, string, error) {
	switch componentType {
	case "controller":
		return filepath.Join("app", "controllers", fileName+".rb"), controllerRubyTemplate, nil
	case "usecase":
		return filepath.Join("app", "services", fileName+".rb"), usecaseRubyTemplate, nil
	case "repository":
		return filepath.Join("app", "repositories", fileName+".rb"), repositoryRubyTemplate, nil
	case "handler":
		return filepath.Join("app", "handlers", fileName+".rb"), handlerRubyTemplate, nil
	default:
		return "", "", fmt.Errorf("tipo não implementado: %s", componentType)
	}
}

func resolveJava(componentType, fileName string) (string, string, error) {
	pkg := filepath.Join("src", "main", "java", "com")
	switch componentType {
	case "controller":
		return filepath.Join(pkg, "controllers", fileName+".java"), controllerJavaTemplate, nil
	case "usecase":
		return filepath.Join(pkg, "usecases", fileName+".java"), usecaseJavaTemplate, nil
	case "repository":
		return filepath.Join(pkg, "repositories", fileName+".java"), repositoryJavaTemplate, nil
	case "handler":
		return filepath.Join(pkg, "handlers", fileName+".java"), handlerJavaTemplate, nil
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

// ==================== GO ====================

const controllerGoTemplate = `// Criado por: {{.Author}}
// Gerado via Devbox CLI

package controllers

type {{.Name}}Controller struct{}

func New{{.Name}}Controller() *{{.Name}}Controller {
	return &{{.Name}}Controller{}
}

func (c *{{.Name}}Controller) Create(ctx interface{}) {
	// Ajuste o tipo de ctx e a chamada para seu framework
}
`

const usecaseGoTemplate = `// Criado por: {{.Author}}
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

const repositoryGoTemplate = `// Criado por: {{.Author}}
// Gerado via Devbox CLI

package repository

type {{.Name}}Repository interface {
	Save() error
}
`

const handlerGoTemplate = `// Criado por: {{.Author}}
// Gerado via Devbox CLI

package handlers

import "net/http"

func Get{{.Name}}(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello {{.Name}}"))
}
`

// ==================== JAVASCRIPT ====================

const controllerJSTemplate = `// Criado por: {{.Author}}
// Gerado via Devbox CLI

class {{.Name}}Controller {
  create(req, res) {
    // Implemente sua lógica aqui
    res.status(201).json({ message: "{{.Name}} created" });
  }
}

module.exports = new {{.Name}}Controller();
`

const usecaseJSTemplate = `// Criado por: {{.Author}}
// Gerado via Devbox CLI

class {{.Name}}UseCase {
  execute() {
    // Implemente sua lógica aqui
  }
}

module.exports = new {{.Name}}UseCase();
`

const repositoryJSTemplate = `// Criado por: {{.Author}}
// Gerado via Devbox CLI

class {{.Name}}Repository {
  save(data) {
    // Implemente sua lógica aqui
  }
}

module.exports = new {{.Name}}Repository();
`

const handlerJSTemplate = `// Criado por: {{.Author}}
// Gerado via Devbox CLI

function get{{.Name}}(req, res) {
  res.send("Hello {{.Name}}");
}

module.exports = { get{{.Name}} };
`

// ==================== TYPESCRIPT ====================

const controllerTSTemplate = `// Criado por: {{.Author}}
// Gerado via Devbox CLI

import { Request, Response } from "express";

export class {{.Name}}Controller {
  create(req: Request, res: Response): void {
    // Implemente sua lógica aqui
    res.status(201).json({ message: "{{.Name}} created" });
  }
}
`

const usecaseTSTemplate = `// Criado por: {{.Author}}
// Gerado via Devbox CLI

export class {{.Name}}UseCase {
  execute(): void {
    // Implemente sua lógica aqui
  }
}
`

const repositoryTSTemplate = `// Criado por: {{.Author}}
// Gerado via Devbox CLI

export class {{.Name}}Repository {
  save(data: unknown): void {
    // Implemente sua lógica aqui
  }
}
`

const handlerTSTemplate = `// Criado por: {{.Author}}
// Gerado via Devbox CLI

import { Request, Response } from "express";

export function get{{.Name}}(req: Request, res: Response): void {
  res.send("Hello {{.Name}}");
}
`

// ==================== PYTHON ====================

const controllerPythonTemplate = `# Criado por: {{.Author}}
# Gerado via Devbox CLI

class {{.Name}}Controller:
    def create(self, data: dict) -> dict:
        # Implemente sua lógica aqui
        return {"message": "{{.Name}} created"}
`

const usecasePythonTemplate = `# Criado por: {{.Author}}
# Gerado via Devbox CLI

class {{.Name}}UseCase:
    def execute(self):
        # Implemente sua lógica aqui
        pass
`

const repositoryPythonTemplate = `# Criado por: {{.Author}}
# Gerado via Devbox CLI

class {{.Name}}Repository:
    def save(self, data: dict):
        # Implemente sua lógica aqui
        pass
`

const handlerPythonTemplate = `# Criado por: {{.Author}}
# Gerado via Devbox CLI

from fastapi import APIRouter

router = APIRouter()

@router.get("/{{.LowerName}}")
def get_{{.LowerName}}():
    return {"message": "Hello {{.Name}}"}
`

// ==================== RUBY ====================

const controllerRubyTemplate = `# Criado por: {{.Author}}
# Gerado via Devbox CLI

class {{.Name}}Controller
  def create
    # Implemente sua lógica aqui
    render json: { message: "{{.Name}} created" }, status: :created
  end
end
`

const usecaseRubyTemplate = `# Criado por: {{.Author}}
# Gerado via Devbox CLI

class {{.Name}}UseCase
  def execute
    # Implemente sua lógica aqui
  end
end
`

const repositoryRubyTemplate = `# Criado por: {{.Author}}
# Gerado via Devbox CLI

class {{.Name}}Repository
  def save(data)
    # Implemente sua lógica aqui
  end
end
`

const handlerRubyTemplate = `# Criado por: {{.Author}}
# Gerado via Devbox CLI

require "sinatra"

get "/{{.LowerName}}" do
  "Hello {{.Name}}"
end
`

// ==================== JAVA ====================

const controllerJavaTemplate = `// Criado por: {{.Author}}
// Gerado via Devbox CLI

package com.controllers;

public class {{.Name}}Controller {

    public String create() {
        // Implemente sua lógica aqui
        return "{{.Name}} created";
    }
}
`

const usecaseJavaTemplate = `// Criado por: {{.Author}}
// Gerado via Devbox CLI

package com.usecases;

public class {{.Name}}UseCase {

    public void execute() {
        // Implemente sua lógica aqui
    }
}
`

const repositoryJavaTemplate = `// Criado por: {{.Author}}
// Gerado via Devbox CLI

package com.repositories;

public interface {{.Name}}Repository {

    void save();
}
`

const handlerJavaTemplate = `// Criado por: {{.Author}}
// Gerado via Devbox CLI

package com.handlers;

import jakarta.servlet.http.*;
import java.io.IOException;

public class {{.Name}}Handler extends HttpServlet {

    @Override
    protected void doGet(HttpServletRequest req, HttpServletResponse resp) throws IOException {
        resp.getWriter().write("Hello {{.Name}}");
    }
}
`
