# regoviz

This is a very experimental tool to visualize policies written in Rego and their evaluation process.

Currently available at [vizrego.pages.dev](https://vizrego.pages.dev).

This is in very early stages of development and is not ready for production use.

本ソフトウェアの詳細については[Wiki](https://github.com/kenty02/regoviz/wiki/%E8%AA%AC%E6%98%8E)をご覧ください。

## Getting Started

### Prerequisites

- Node.js 20
- Go 1.21
- Git
- (More details are instructed by setup script)

### Development

1. Run `./setup.sh` (`./setup.bat` on Windows)
   1. Run this again if you have problems after `git pull`
2. Run `air` in the root directory to start the backend server, and `npm run start` in the `web` directory to start the frontend server.

If you edit openapi.yml, you need to run `go generate ./...` and it will generate both frontend and backend code.
