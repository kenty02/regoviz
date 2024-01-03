package gen

//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --target ../internal/api --clean ../api/openapi.yml
//go:generate go run generate_frontend.go
