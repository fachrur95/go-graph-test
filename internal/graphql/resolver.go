package graphql

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleGraphQLRequest(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			log.Printf("Error parsing body: %v", err)
			return
		}

		query, ok := requestBody["query"].(string)
		if !ok {
			http.Error(w, "Query missing", http.StatusBadRequest)
			log.Println("Query missing from request body")
			return
		}

		// Set up the GraphQL schema and context
		ctx := context.WithValue(r.Context(), "mongoClient", client)
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: query,
			Context:       ctx,
		})
		if len(result.Errors) > 0 {
			log.Printf("Failed to execute GraphQL operation: %v", result.Errors)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}
