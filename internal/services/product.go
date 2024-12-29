package services

import (
	"context"
	"go-graphql-test/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateProduct(client *mongo.Client, name string, price, stock float64) (*models.Product, error) {
	collection := client.Database("go_graphql").Collection("products")

	product := models.Product{Name: name, Price: price, Stock: stock}
	_, err := collection.InsertOne(context.Background(), product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func GetAllProducts(client *mongo.Client) ([]models.Product, error) {
	collection := client.Database("go_graphql").Collection("products")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var products []models.Product
	for cursor.Next(context.Background()) {
		var product models.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
