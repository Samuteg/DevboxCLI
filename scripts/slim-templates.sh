#!/usr/bin/env bash
# Reduz cmd/templates/node/{ts,js}.zip removendo node_modules, .pnpm e *.db.
# Esperado: 153M -> <5M. Rode na raiz do repo e commite os zips resultantes.
set -euo pipefail

cd "$(dirname "$0")/.."

for name in ts js; do
  zip="cmd/templates/node/$name.zip"
  [ -f "$zip" ] || { echo "ausente: $zip"; continue; }
  tmp=$(mktemp -d)
  echo "== $zip =="
  unzip -q "$zip" -d "$tmp"
  rm -rf "$tmp/$name/node_modules" "$tmp/$name/data" "$tmp"/**/.pnpm 2>/dev/null || true
  find "$tmp/$name" -name "*.db" -delete
  find "$tmp/$name" -name ".pnpm" -prune -exec rm -rf {} + 2>/dev/null || true
  rm -f "$zip"
  (cd "$tmp" && zip -qr "$OLDPWD/$zip" "$name")
  rm -rf "$tmp"
  ls -lh "$zip"
done

echo "OK. Confira com: unzip -l cmd/templates/node/ts.zip | head"
