#!/usr/bin/env bash
set -e

# Cores para o output combinar com a identidade da Devbox
CYAN='\033[0;36m'
PURPLE='\033[0;35m'
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

REPO="Samuteg/DevboxCLI"
BINARY_NAME="devbox"

echo -e "${PURPLE}🚀 A preparar a instalação da Devbox CLI...${NC}"

# 1. Detetar o Sistema Operativo e Arquitetura
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

# 2. Obter a versão mais recente (Latest Release) do GitHub
echo -e "${CYAN}🔍 A procurar a versão mais recente...${NC}"
LATEST_VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo -e "${RED}❌ Não foi possível determinar a versão mais recente. Verifique a sua ligação à internet.${NC}"
    exit 1
fi

echo -e "📦 Versão encontrada: ${GREEN}${LATEST_VERSION}${NC}"

# 3. Construir o URL de transferência (Baseado no padrão do GoReleaser)
# Nota: O nome do ficheiro exato dependerá de como configurar o GoReleaser no futuro.
# Normalmente é algo como: DevboxCLI_Linux_x86_64.tar.gz
TAR_FILE="DevboxCLI_${OS^}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/$REPO/releases/download/${LATEST_VERSION}/${TAR_FILE}"

# 4. Transferir e extrair
TMP_DIR=$(mktemp -d)
echo -e "${CYAN}⬇️ A transferir o binário de $DOWNLOAD_URL...${NC}"

if curl -sL "$DOWNLOAD_URL" -o "$TMP_DIR/$TAR_FILE"; then
    tar -xzf "$TMP_DIR/$TAR_FILE" -C "$TMP_DIR"
else
    echo -e "${RED}❌ Falha ao transferir o ficheiro. Verifique se a release existe no GitHub.${NC}"
    exit 1
fi

# 5. Instalar no sistema (requer permissões de root normalmente)
INSTALL_DIR="/usr/local/bin"
echo -e "${CYAN}⚙️ A instalar em $INSTALL_DIR... (Pode pedir a sua password)${NC}"

sudo mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"

# Limpar ficheiros temporários
rm -rf "$TMP_DIR"

echo ""
echo -e "${GREEN}✨ Devbox instalada com sucesso!${NC}"
echo -e "Execute ${PURPLE}devbox --help${NC} para começar."