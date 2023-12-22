package main

import (
	"github.com/rs/cors"
	"log"
	"net/http"
	vizrego "vizrego-poc/vizrego"
)

func main() {
	// Create service instance.
	service := NewService()
	// Create generated server.
	srv, err := vizrego.NewServer(service)
	if err != nil {
		log.Fatal(err)
	}

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation below for more options.
	handler := cors.Default().Handler(srv)
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}

// corsMiddleware applies CORS headers to the response
func corsMiddleware(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		h.ServeHTTP(w, r)
	}
}
