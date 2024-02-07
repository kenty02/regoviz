# regoviz

This is a very experimental tool to visualize policies written in Rego and their evaluation process.

Currently available at [vizrego.pages.dev](https://vizrego.pages.dev).

This is in very early stages of development and is not ready for production use.

## Getting Started

### Prerequisites

- Node.js 20
- Go 1.21
  - air
- GoLand

### Development

1. Run `./setup.sh` (`./setup.bat` on Windows)
   1. Run this again if you have problems after `git pull`
2. Run Configurations

If you edit openapi.yml, you need to run `go generate ./...` and it will generate both frontend and backend code.
