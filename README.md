# regoviz

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

If you edit openapi.yml, you need to run `go generate ./gen` and it will generate both frontend and backend code.
