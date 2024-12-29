package main

import (
	"context"
	"go-graphql-test/internal/auth"
	"go-graphql-test/internal/graphql"
	"go-graphql-test/internal/mongodb"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Connect to MongoDB
	client, err := mongodb.ConnectMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Initialize GraphQL server
	r := mux.NewRouter()
	r.Use(auth.InjectHTTPRequestMiddleware) // Tambahkan middleware di sini
	r.HandleFunc("/graphql", graphql.HandleGraphQLRequest(client)).Methods("POST")

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", r)
}
