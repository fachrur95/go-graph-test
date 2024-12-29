package services

import (
	"context"
	"errors"
	"fmt"
	"go-graphql-test/internal/auth"
	"go-graphql-test/internal/models"
	"go-graphql-test/pkg/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterUser(client *mongo.Client, name, email, password string) (*models.User, error) {
	collection := client.Database("go_graphql").Collection("users")

	// Check user exists
	var existingUser models.User
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&existingUser)
	if err == nil {
		return nil, fmt.Errorf("user already exists")
	}
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// New Insert
	user := models.User{Name: name, Email: email, Password: hashedPassword}
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func LoginUser(client *mongo.Client, email, password string) (string, error) {
	collection := client.Database("go_graphql").Collection("users")

	var user *models.User
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("invalid credentials")
		}
		return "", err
	}
	if err := utils.VerifyPassword(user.Password, password); err != nil { // Verifikasi password dengan bcrypt
		return "", errors.New("invalid password")
	}
	// Create JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString(auth.JwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetAllUsers(client *mongo.Client) ([]models.User, error) {
	collection := client.Database("go_graphql").Collection("users")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []models.User
	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
