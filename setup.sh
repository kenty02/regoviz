#!/usr/bin/env bash

set -e

deps_ok=true
function check_installed() {
    if ! [ -x "$(command -v "$1")" ]; then
        echo "🔴 $1 is not installed. See $2"
        deps_ok=false
    fi
}

node -v
go version
check_installed "air" "https://github.com/cosmtrek/air?tab=readme-ov-file#installation"
check_installed "golangci-lint" "https://golangci-lint.run/usage/install/#local-installation"
check_installed "lefthook" "https://github.com/evilmartians/lefthook?tab=readme-ov-file#install"
check_installed "actionlint" "https://github.com/rhysd/actionlint?tab=readme-ov-file#quick-start"

if ! $deps_ok; then
    echo "Please install missing dependencies"
    exit 1
fi

lefthook install

cd web

npm install

echo "Setup complete"
