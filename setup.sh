#!/usr/bin/env bash

set -e

node -v
go version
if ! [ -x "$(command -v air)" ]; then
    read -p "air is not installed. Do you want to install it via 'go install'? (y/N) " -n 1 -r
    read -r
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        go install github.com/cosmtrek/air@latest
    fi
fi
if ! [ -x "$(command -v golangci-lint)" ]; then
    read -p "golangci-lint is not installed. Do you want to install it via 'curl'? (y/N) " -n 1 -r
    read -r
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        # https://golangci-lint.run/usage/install/#binaries
        # shellcheck disable=SC2046
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2
    else
        echo "Please install golangci-lint first."
        echo "https://golangci-lint.run/usage/install/"
        exit 1
    fi
fi
if ! [ -x "$(command -v lefthook)" ]; then
    read -p "lefthook is not installed. Do you want to install it via 'go install'? (y/N) " -n 1 -r
    read -r
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        go install github.com/evilmartians/lefthook@latest
    else
        echo "Please install lefthook first."
        exit 1
    fi
fi
lefthook install

if [ -f .env.local ]; then
    echo ".env.local already exists"
else
    cp .env.local.example .env.local
fi

cd web

if [ -f .env.local ]; then
    echo "web/.env.local already exists"
else
    cp .env.local.example .env.local
fi

if [ -f .dev.vars ]; then
    echo "web/.dev.vars already exists"
else
    cp .env.local.example .dev.vars
fi

npm install
