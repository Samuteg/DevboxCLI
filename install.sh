#!/usr/bin/env bash
set -e

# Cores para o output combinar com a identidade da Devbox
CYAN='\033[0;36m'
PURPLE='\033[0;35m'
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

REPO="Samuteg/DevboxCLI"
BINARY_NAME="devbox"

echo -e "${PURPLE}🚀 A preparar a instalação da Devbox CLI...${NC}"

# ──────────────────────────────────────────────
# 1. Detetar o Sistema Operativo e Arquitetura
# ──────────────────────────────────────────────
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

if [ "$ARCH" = "x86_64" ]; then
    ARCH="x86_64"
elif [ "$ARCH" = "arm64" ] || [ "$ARCH" = "aarch64" ]; then
    ARCH="arm64"
else
    echo -e "${RED}❌ Arquitetura não suportada: $ARCH${NC}"
    exit 1
fi

if [ "$OS" != "linux" ] && [ "$OS" != "darwin" ]; then
    echo -e "${RED}❌ Sistema operativo não suportado por este script: $OS${NC}"
    echo "Por favor, utilize a instalação via Go ou transfira o binário manualmente."
    exit 1
fi

# ──────────────────────────────────────────────
# 2. Obter a versão mais recente (Latest Release)
# ──────────────────────────────────────────────
echo -e "${CYAN}🔍 A procurar a versão mais recente...${NC}"
LATEST_VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo -e "${RED}❌ Não foi possível determinar a versão mais recente. Verifique a sua ligação à internet.${NC}"
    exit 1
fi

echo -e "📦 Versão encontrada: ${GREEN}${LATEST_VERSION}${NC}"

# ──────────────────────────────────────────────
# 3. Montar URL de download e transferir
# ──────────────────────────────────────────────
TAR_FILE="DevboxCLI_${OS^}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/$REPO/releases/download/${LATEST_VERSION}/${TAR_FILE}"

TMP_DIR=$(mktemp -d)
echo -e "${CYAN}⬇️ A transferir o binário de $DOWNLOAD_URL...${NC}"

if curl -sL "$DOWNLOAD_URL" -o "$TMP_DIR/$TAR_FILE"; then
    tar -xzf "$TMP_DIR/$TAR_FILE" -C "$TMP_DIR"
else
    echo -e "${RED}❌ Falha ao transferir o ficheiro. Verifique se a release existe no GitHub.${NC}"
    exit 1
fi

# ──────────────────────────────────────────────
# 4. Instalar e garantir que está no PATH
# ──────────────────────────────────────────────

detect_shell_config() {
    case "$SHELL" in
        */zsh)  echo "${ZDOTDIR:-$HOME}/.zshrc" ;;
        */bash) echo "$HOME/.bashrc" ;;
        */fish) echo "$HOME/.config/fish/config.fish" ;;
        *)      echo "$HOME/.profile" ;;
    esac
}

add_to_path() {
    local dir="$1"
    local rc_file="$2"

    # Já está no PATH? Nada a fazer.
    if echo ":$PATH:" | grep -q ":$dir:"; then
        return 0
    fi

    echo -e "${CYAN}📝 A adicionar $dir ao PATH em $rc_file...${NC}"

    case "$(basename "$rc_file")" in
        *.fish)
            echo "set -gx PATH $dir \$PATH" >> "$rc_file"
            ;;
        *)
            echo "export PATH=\"$dir:\$PATH\"" >> "$rc_file"
            ;;
    esac
}

# Tenta instalação global (sudo).
install_system_wide() {
    local dst="/usr/local/bin"
    if sudo mv "$TMP_DIR/$BINARY_NAME" "$dst/$BINARY_NAME" 2>/dev/null && sudo chmod +x "$dst/$BINARY_NAME"; then
        echo -e "${GREEN}✅ Instalado em $dst/$BINARY_NAME${NC}"
        add_to_path "$dst" "$(detect_shell_config)"
        return 0
    fi
    return 1
}

# Instalação local (sem sudo) e garante PATH.
install_user_local() {
    local dst="${XDG_DATA_HOME:-$HOME/.local}/bin"
    mkdir -p "$dst"

    if mv "$TMP_DIR/$BINARY_NAME" "$dst/$BINARY_NAME" 2>/dev/null; then
        chmod +x "$dst/$BINARY_NAME"
        echo -e "${GREEN}✅ Instalado em $dst/$BINARY_NAME${NC}"

        local rc_file
        rc_file="$(detect_shell_config)"
        add_to_path "$dst" "$rc_file"

        echo ""
        echo -e "${PURPLE}💡 O diretório foi adicionado ao ficheiro de configuração do teu shell.${NC}"
        echo -e "${PURPLE}   Reinicia o terminal ou executa: source $rc_file${NC}"
        return 0
    fi
    return 1
}

echo -e "${CYAN}⚙️  A instalar o binário...${NC}"

if command -v sudo &>/dev/null; then
    if install_system_wide; then
        INSTALL_OK=true
    else
        echo -e "${YELLOW}⚠️  Instalação global falhou. A tentar instalação local...${NC}"
        if install_user_local; then
            INSTALL_OK=true
        fi
    fi
else
    if install_user_local; then
        INSTALL_OK=true
    fi
fi

if [ "${INSTALL_OK:-false}" != true ]; then
    echo -e "${RED}❌ Não foi possível instalar o binário.${NC}"
    exit 1
fi

# ──────────────────────────────────────────────
# 5. Limpeza
# ──────────────────────────────────────────────
rm -rf "$TMP_DIR"

# ──────────────────────────────────────────────
# 6. Verificação final
# ──────────────────────────────────────────────
echo ""
if command -v "$BINARY_NAME" &>/dev/null; then
    echo -e "${GREEN}✨ Devbox CLI instalada com sucesso!${NC}"
    echo -e "Execute ${PURPLE}devbox --help${NC} para começar."
else
    echo -e "${GREEN}✨ Devbox CLI instalada!${NC}"
    echo -e "${PURPLE}⚠️  Reinicia o terminal ou executa 'source $(detect_shell_config)' para usar o comando.${NC}"
fi
