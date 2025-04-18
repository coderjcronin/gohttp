package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", fmt.Errorf("FAILED TO FIND AUTHORIZATION HEADER")
	}

	splitString := strings.Fields(authHeader)

	if splitString[1] == "" {
		return "", fmt.Errorf("FAILED TO FIND VALID TOKEN STRING")
	}

	if !strings.EqualFold("apikey", splitString[0]) {
		return "", fmt.Errorf("FAILED TO FIND APIKEY IN AUTHORIZATION HEADER")
	}

	return splitString[1], nil
}
