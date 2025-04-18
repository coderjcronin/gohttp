package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	rand.Read(key)

	returnString := hex.EncodeToString(key)
	if returnString == "" {
		return "", fmt.Errorf("FAILED TO GENERATE TOKEN STRING")
	}
	return returnString, nil
}
