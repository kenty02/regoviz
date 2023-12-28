# vizrego

## Getting Started

### Prerequisites

- Node.js 20
- Go 1.21
  - air
- GoLand

### Development

1. `cp .env.local.example .env.local`
2. `cd frontend && cp .env.local.example .env.local && cp .env.local.example .dev.vars`
3. `npm install`
4. Run Configurations

If you edit openapi.yml, you need to run `go generate` and it will generate both frontend and backend code.
