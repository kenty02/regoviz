package main

import (
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	vizrego "vizrego-poc/vizrego"
)

func authHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		// get token from env
		token := os.Getenv("TOKEN")
		if token == "" {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("500 Internal Server Error\n"))
			return
		}

		if auth == "Bearer "+token {
			// 200
			h.ServeHTTP(w, r)
		} else {
			// 403
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte("403 Forbidden\n"))
		}
	})
}
func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal("Error loading .env.local file")
	}
	// Create service instance.
	service := NewService()
	// Create generated server.
	srv, err := vizrego.NewServer(service)
	if err != nil {
		log.Fatal(err)
	}
	handler := authHandler(srv)

	var allowedOrigins []string
	allowedOrigins = append(allowedOrigins, os.Getenv("FRONTEND_URL"))
	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: true,
		AllowedHeaders: []string{
			"*",
		},
	})
	handler = c.Handler(handler)

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
