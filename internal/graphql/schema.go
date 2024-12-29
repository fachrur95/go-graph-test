package graphql

import (
	"go-graphql-test/internal/auth"
	"go-graphql-test/internal/services"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/mongo"
)

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id":    &graphql.Field{Type: graphql.String},
		"name":  &graphql.Field{Type: graphql.String},
		"email": &graphql.Field{Type: graphql.String},
	},
})

var productType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Product",
	Fields: graphql.Fields{
		"id":    &graphql.Field{Type: graphql.String},
		"name":  &graphql.Field{Type: graphql.String},
		"price": &graphql.Field{Type: graphql.Float},
		"stock": &graphql.Field{Type: graphql.Float},
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"users": &graphql.Field{
				Type: graphql.NewList(userType),
				Resolve: auth.AuthMiddleware(func(p graphql.ResolveParams) (interface{}, error) {
					client := p.Context.Value("mongoClient").(*mongo.Client)
					return services.GetAllUsers(client)
				}),
			},
			"products": &graphql.Field{
				Type: graphql.NewList(productType),
				Resolve: auth.AuthMiddleware(func(p graphql.ResolveParams) (interface{}, error) {
					client := p.Context.Value("mongoClient").(*mongo.Client)
					return services.GetAllProducts(client)
				}),
			},
		},
	}),
	Mutation: graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"register": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"name":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"email":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"password": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					name := p.Args["name"].(string)
					email := p.Args["email"].(string)
					password := p.Args["password"].(string)
					client := p.Context.Value("mongoClient").(*mongo.Client)
					_, err := services.RegisterUser(client, name, email, password)
					if err != nil {
						return nil, err
					}
					return "User registered successfully", nil

				},
			},
			"login": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"email":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"password": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					email := p.Args["email"].(string)
					password := p.Args["password"].(string)
					client := p.Context.Value("mongoClient").(*mongo.Client)
					return services.LoginUser(client, email, password)
				},
			},
			"createProduct": &graphql.Field{
				Type: productType,
				Args: graphql.FieldConfigArgument{
					"name":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"price": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Float)},
					"stock": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Float)},
				},
				Resolve: auth.AuthMiddleware(func(p graphql.ResolveParams) (interface{}, error) {
					name := p.Args["name"].(string)
					price := p.Args["price"].(float64)
					stock := p.Args["stock"].(float64)
					client := p.Context.Value("mongoClient").(*mongo.Client)
					return services.CreateProduct(client, name, price, stock)
				}),
			},
		},
	}),
})
