#!/usr/bin/env bash

set -e

node -v
go version
if ! [ -x "$(command -v air)" ]; then
    read -p "air is not installed. Do you want to install it via 'go install'? (y/N) " -n 1 -r
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        go install github.com/cosmtrek/air@latest
    fi
fi
if ! [ -x "$(command -v lefthook)" ]; then
    read -p "lefthook is not installed. Do you want to install it via 'go install'? (y/N) " -n 1 -r
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

cd frontend

if [ -f .env.local ]; then
    echo "frontend/.env.local already exists"
else
    cp .env.local.example .env.local
fi

if [ -f .dev.vars ]; then
    echo "frontend/.dev.vars already exists"
else
    cp .env.local.example .dev.vars
fi

npm install
