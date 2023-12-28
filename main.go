package main

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	vizrego "vizrego-poc/vizrego"
)

type SecurityHandler struct{}

func (s SecurityHandler) HandleBearerAuth(ctx context.Context, operationName string, t vizrego.BearerAuth) (context.Context, error) {
	// get token from env
	token := os.Getenv("TOKEN")
	if token == "" {
		return ctx, fmt.Errorf("token is empty")
	}

	if t.GetToken() == token {
		return ctx, nil
	} else {
		return ctx, fmt.Errorf("invalid token")
	}
}

func main() {
	err := godotenv.Load(".env.local")

	// To initialize Sentry's handler, you need to initialize Sentry itself beforehand
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:           "https://b2d9c7ea4a09baf4e3ad530d14d2ab1e@o4504839999848448.ingest.sentry.io/4506472246345728",
		EnableTracing: true,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v", err)
	} // Create an instance of sentryhttp
	sentryHandler := sentryhttp.New(sentryhttp.Options{})

	// Create service instance.
	service := NewService()

	var securityHandler vizrego.SecurityHandler = SecurityHandler{}
	srv, err := vizrego.NewServer(service, securityHandler)
	if err != nil {
		log.Fatal(err)
	}
	var handler http.Handler = srv

	var allowedOrigins []string
	allowedOrigins = append(allowedOrigins, os.Getenv("FRONTEND_URL"))
	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: true,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"*",
		},
	})
	handler = c.Handler(srv)

	handler = sentryHandler.Handle(handler)
	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
