<div align="center">

# 📦 Devbox CLI

**Acelere seu desenvolvimento.**<br>
Crie projetos Backend e Frontend configurados com as melhores práticas<br>
(Clean Architecture, DDD, SOLID) em segundos.

  <p>
    <img src="https://img.shields.io/github/go-mod/go-version/Samuteg/DevboxCLI?style=flat-square&color=00ADD8" alt="Go Version" />
    <img src="https://img.shields.io/github/license/Samuteg/DevboxCLI?style=flat-square&color=5D6D7E" alt="License" />
    <img src="https://img.shields.io/github/v/release/Samuteg/DevboxCLI?style=flat-square&color=2ECC71" alt="Release" />
  </p>

  <br>

[**Instalar Agora**](#-instalação) • [**Como Usar**](#-comandos) • [**Stacks**](#-stacks-suportadas)

</div>

---

## ✨ Por que usar o Devbox?

A **Devbox CLI** elimina a fadiga de configuração inicial. Em vez de gastar horas configurando pastas, linters e Docker, inicie uma aplicação robusta com um comando.

<table>
  <tr>
    <td width="50%">
      <h3>🚀 Setup Instantâneo</h3>
      <p>Gere projetos completos com dependências instaladas e git inicializado automaticamente.</p>
    </td>
    <td width="50%">
      <h3>🏗️ Clean Architecture</h3>
      <p>Templates de Backend já nascem com estrutura de <i>Domain-Driven Design</i> e separação de camadas.</p>
    </td>
  </tr>
  <tr>
    <td width="50%">
      <h3>🎨 Frontend Moderno</h3>
      <p>Integração nativa com <b>Vite</b> (React/TS) e <b>Next.js</b>.</p>
    </td>
    <td width="50%">
      <h3>🩺 Devbox Doctor</h3>
      <p>O comando <code>doctor</code> verifica seu ambiente e avisa o que falta instalar.</p>
    </td>
  </tr>
</table>

---

## 📋 Pré-requisitos

| Ferramenta   | Versão Mínima | Obrigatório para              |
| :----------- | :------------ | :---------------------------- |
| **Go**       | 1.22+         | Instalação via `go install`   |
| **Node.js**  | 18+           | Stacks Node.js / Frontend     |
| **pnpm**     | 8+            | Frontend (Vite / Next.js)     |
| **Python**   | 3.11+         | Stack Python                  |
| **Ruby**     | 3.0+          | Stack Ruby                    |

> 💡 O comando `DevboxCLI doctor` detecta automaticamente o que está instalado e mostra o que precisa ser configurado.

---

## 🚀 Instalação

### Opção 1: Via Go Install (Recomendado)

```bash
go install github.com/Samuteg/DevboxCLI@latest
```

### Opção 2: Via Script (Linux / macOS)

```bash
curl -sSL https://raw.githubusercontent.com/Samuteg/DevboxCLI/main/install.sh | bash
```

> Após a instalação, o binário `DevboxCLI` estará disponível no seu `$PATH`.

---

## 🎮 Comandos

| Comando                         | Descrição                                                                 |
| :------------------------------ | :-----------------------------------------------------------------------  |
| `DevboxCLI init`                | Inicia um novo projeto com estrutura profissional                         |
| `DevboxCLI add [tipo] [nome]`   | Cria componentes como Controllers e Usecases (Clean Arch)                |
| `DevboxCLI commit`              | Wizard interativo para mensagens de commit padronizadas                   |
| `DevboxCLI kill [porta]`        | Encerra o processo ocupando uma porta (ex: `8080`)                       |
| `DevboxCLI config`              | Gerencia preferências no arquivo `~/.devbox.yaml`                        |
| `DevboxCLI doctor`              | Verifica o estado das dependências instaladas                            |
| `DevboxCLI cleanup`             | Limpa caches, `node_modules` e binários para liberar espaço              |

### Exemplo: Criando um projeto Go com Clean Architecture

```bash
DevboxCLI init
# → Nome do Projeto: meu-api
# → Tipo de Projeto: Backend
# → Escolha a Tech: Go
# → ⚡ Escolha uma variante: Gin
# ✔ Estrutura criada!
```

---

## 🛠️ Stacks Suportadas

### Backend

| Stack       | Variantes                  | Estrutura Gerada                                 |
| :---------- | :------------------------- | :----------------------------------------------- |
| **Go**      | `Gin`, `Simple`            | `cmd/`, `internal/entity`, `internal/usecase`    |
| **Python**  | `FastAPI`                  | `src/api`, `src/core`, `src/models`, `tests/`    |
| **Node.js** | `TypeScript`, `JavaScript` | `src/controllers`, `src/routes`, `src/models`    |
| **Ruby**    | `Simple`, `Rails`          | `app/controllers`, `app/routes`, `app/models`    |

### Frontend

| Stack       | Detalhes                        |
| :---------- | :------------------------------ |
| **React**   | Via Vite (TypeScript + SWC)     |
| **Next.js** | App Router, TailwindCSS, ESLint |

---

## ⚙️ Configuração

O Devbox CLI cria um arquivo de configuração em `~/.devbox.yaml`:

```yaml
author: "Seu Nome"
default-port: "8080"
update-channel: "stable"
template-style: "clean"
```

Use o comando `DevboxCLI config` para gerenciar essas preferências.

---

## 🧱 Como Funciona

### Sistema de Templates

A Devbox CLI utiliza **arquivos zip** embutidos no binário para gerar projetos:

```
cmd/templates/
├── golang/
│   ├── simple.zip      → Projeto Go simples
│   └── Gin.zip         → Projeto Go com Gin + Clean Arch
├── node/
│   ├── js.zip          → Node.js com JavaScript
│   └── ts.zip          → Node.js com TypeScript
├── python/
│   ├── simple.zip      → Python simples
│   └── FastAPI.zip     → API Python com FastAPI
└── ruby/
    ├── simple.zip      → Ruby simples
    └── On_rails.zip    → Ruby on Rails
```

Quando você executa `DevboxCLI init`:

1. O CLI identifica a stack e variante escolhidas
2. Localiza o zip correspondente dentro do binário (via `//go:embed`)
3. **Extrai o zip** diretamente no diretório do projeto
4. Cria diretórios extras e instala dependências (se aplicável)

Isso torna o binário **autocontido** e elimina a necessidade de baixar templates da internet.

### Arquitetura Interna

```
DevboxCLI/
├── cmd/                  # Comandos Cobra e fluxo de entrada/saída
│   ├── init.go           # Comando `init` — criação de projetos
│   ├── add.go            # Comando `add` — criação de componentes
│   ├── commit.go         # Wizard de commit
│   ├── kill.go           # Killer de porta
│   ├── config.go         # Gerenciamento de configuração
│   ├── doctor.go         # Diagnóstico do ambiente
│   ├── cleanup.go        # Limpeza de caches
│   └── templates/        # Templates zipados embutidos no binário
├── internal/
│   ├── scaffold/         # Lógica de geração (projetos, stacks, componentes)
│   └── system/           # Utilitários de sistema e execução de comandos
├── main.go               # Ponto de entrada
└── go.mod                # Dependências Go
```

> A separação entre `cmd/` e `internal/` reduz acoplamento e facilita manutenção e testes.

---

## ❓ Perguntas Frequentes

<details>
<summary><b>Preciso de internet para criar um projeto?</b></summary>
<br>
Para stacks <b>Backend</b>, não! Os templates já estão embutidos no binário via <code>//go:embed</code>.
Para stacks <b>Frontend</b> (Vite/Next.js), é necessário acesso à internet pois o CLI delega a criação aos scaffolds oficiais via <code>pnpm create</code>.
</details>

<details>
<summary><b>Como atualizar o Devbox CLI?</b></summary>
<br>
Reinstale com o mesmo comando de instalação:

```bash
go install github.com/Samuteg/DevboxCLI@latest
```
</details>

<details>
<summary><b>O comando <code>init</code> não encontra meus templates</b></summary>
<br>
Verifique se o binário foi compilado corretamente — os templates são embutidos em tempo de compilação via <code>//go:embed</code>. Se você fez <code>go build</code>, certifique-se de que os arquivos <code>.zip</code> estão em <code>cmd/templates/</code>.
</details>

<details>
<summary><b>Posso contribuir com novos templates?</b></summary>
<br>
Sim! Crie um diretório com os arquivos do template, execute o gerador de zip e abra um Pull Request. Veja o guia de contribuição no repositório.
</details>

---

<div align="center">
<p>Desenvolvido com 💜 por <a href="https://github.com/Samuteg">Samuteg</a></p>
</div>
