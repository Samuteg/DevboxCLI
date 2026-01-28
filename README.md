📦 DevBox CLI
The Ultimate Developer Swiss-Army-Knife
DevBox is a high-performance Command Line Interface (CLI) built in Go designed to eliminate the friction of daily development workflows. From scaffolding clean architectures to enforcing git best practices, DevBox automates the "boring stuff" so you can focus on writing code.

🛠 Built With
Go (Golang): For lightning-fast execution and a single static binary.

Cobra: The industry-standard library for modern CLI applications.

Viper: Seamless configuration management (YAML/JSON/ENV).

Go-Git: Direct Git integration for repository manipulation.

PromptUI: Interactive terminal prompts for a better user experience.

🚀 Key Features
1. Project Scaffolding (devbox init)
Instantly bootstrap projects with production-ready structures:

--frontend: Generates a Vite project with optimized presets.

--backend: Sets up a Go Clean Architecture layout (cmd/, internal/, pkg/) including a docker-compose.yml for database services and a pre-configured .env.

2. Smart Sync (devbox save)
A unified command to add, commit, and push changes:

Conventional Commits Enforcement: Validates messages against standard types (feat, fix, chore, etc.).

Interactive Wizard: A devbox commit mode that guides you through building the perfect commit message.

3. Safety Guardrails
Branch Protection: Prevents accidental commits to main, master, or production branches unless the --force flag is explicitly used.

📥 Installation
Bash
# Clone the repository
git clone https://github.com/youruser/devbox

# Install globally
go install .
