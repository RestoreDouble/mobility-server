package service

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type UserClaims struct {
	Phone int  `json:"Phone"`
	IsNew bool `json:"IsNew"`
	jwt.MapClaims
}

func NewAccessToken(claims jwt.MapClaims) (string, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return accessToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func GenerateAccessToken(userClaims *UserClaims) (string, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Create a new token with custom claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)

	// Sign the token with a secret
	tokenString, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseAccessToken(accessToken string) (*UserClaims, error) {
	parsedAccessToken, _ := jwt.Parse(
		accessToken,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("TOKEN_SECRET")), nil
		})

	res, ok := parsedAccessToken.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("token claims could not be parsed as UserClaims")
	}

	return res, nil
}
