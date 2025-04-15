package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	})

	key := []byte(tokenSecret)
	signedString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return signedString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	stringID, err := claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	tokenID, err := uuid.Parse(stringID)
	if err != nil {
		return uuid.Nil, err
	}

	return tokenID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", fmt.Errorf("FAILED TO FIND AUTHORIZATION HEADER")
	}

	if !strings.Contains(authHeader, "Bearer") {
		return "", fmt.Errorf("FAILED TO FIND BEARER IN AUTHORIZATION HEADER")
	}

	splitString := strings.Fields(authHeader)

	if splitString[1] == "" {
		return "", fmt.Errorf("FAILED TO FIND VALID TOKEN STRING")
	}

	return splitString[1], nil
}
