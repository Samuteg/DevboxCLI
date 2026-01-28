# 📦 DevBox CLI

<p align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go" />
  <img src="https://img.shields.io/badge/Cobra-v1.8.0-blue?style=for-the-badge" alt="Cobra" />
  <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" alt="License" />
</p>

**DevBox** é uma CLI de alta performance desenvolvida em Go, projetada para automatizar fluxos de trabalho repetitivos, garantir padrões de commit e diagnosticar a saúde do seu ambiente de desenvolvimento.

---

## 🎨 O Projeto

```text
  _____  ______      __ ____   ____ __   __
 |  __ \|  ____\ \    / /  _ \ / __ \\ \ / /
 | |  | | |__   \ \  / /| |_) | |  | |\ V / 
 | |  | |  __|   \ \/ / |  _ <| |  | | > <  
 | |__| | |____   \  /  | |_) | |__| |/ . \ 
 |_____/|______|   \/   |____/ \____//_/ \_\

      >>> Sua Toolbox de Automação Pessoal <<<
🚀 Funcionalidades
🩺 System Doctor
Verifica instantaneamente se as dependências essenciais (Git, Docker, Go, Node) estão instaladas e configuradas no seu PATH.

Bash
devbox doctor
🛡️ Smart Save (Git Flow)
Commita e envia suas alterações com segurança.

Proteção de Branch: Impede commits diretos na main ou master.

Validação: Garante que as mensagens de commit sigam padrões.

Bash
devbox save "feat: nova funcionalidade incrível"
🧹 Workspace Cleanup
Remove branches locais que já foram mergeadas ou que não existem mais no repositório remoto, mantendo seu ambiente limpo.

Bash
devbox cleanup --dry-run
🔄 Self-Update
Mantenha sua ferramenta sempre atualizada com um único comando, baixando a versão mais recente diretamente do repositório.

Bash
devbox update
🛠️ Instalação
Certifique-se de que o diretório $GOPATH/bin está no seu PATH.

Bash
# Clone o repositório
git clone [https://github.com/seu-usuario/devbox.git](https://github.com/seu-usuario/devbox.git)

# Entre na pasta
cd devbox

# Instale globalmente
go install .
⚙️ Configuração
A DevBox utiliza um arquivo de configuração yaml para personalizar o comportamento:

Crie o arquivo em ~/.devbox.yaml:

YAML
repo: "[github.com/seu-usuario/devbox](https://github.com/seu-usuario/devbox)"
protected_branches:
  - "main"
  - "master"
  - "production"
workspace: "~/projects"
🏗️ Tecnologias Utilizadas
Go - Linguagem base.

Cobra - Framework para interfaces CLI.

Viper - Gerenciamento de configuração.

Go-Git - Manipulação nativa de repositórios Git.

Desenvolvido por Seu Nome
