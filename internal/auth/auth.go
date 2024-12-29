package auth

import (
	"context"
	"errors"
	"go-graphql-test/internal/configs"
	"go-graphql-test/internal/models"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/graphql-go/graphql"
	"golang.org/x/crypto/bcrypt"
)

var JwtSecret = []byte(configs.EnvJwtSecret())

// HashPassword hashes the password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword verifies the hashed password
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// GenerateJWT generates JWT token for the user
func GenerateJWT(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString(JwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

type contextKey string

const UserIDKey contextKey = "email"

func AuthMiddleware(next func(p graphql.ResolveParams) (interface{}, error)) func(p graphql.ResolveParams) (interface{}, error) {
	return func(p graphql.ResolveParams) (interface{}, error) {
		req, ok := p.Context.Value(HTTPRequestKey).(*http.Request)
		if !ok {
			return nil, errors.New("invalid HTTP request context")
		}

		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			return nil, errors.New("missing Authorization header")
		}

		// Validate token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(JwtSecret), nil
		})
		if err != nil || !token.Valid {
			return nil, errors.New("invalid token")
		}

		userID, ok := claims["email"].(string)
		if !ok {
			return nil, errors.New("invalid token payload")
		}
		ctx := context.WithValue(p.Context, UserIDKey, userID)

		return next(graphql.ResolveParams{
			Args:    p.Args,
			Context: ctx,
			Info:    p.Info,
			Source:  p.Source,
		})
	}
}

type ContextKey string

const HTTPRequestKey ContextKey = "httpRequest"

func InjectHTTPRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), HTTPRequestKey, r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
