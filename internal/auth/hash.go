package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	passByte := []byte(password)

	hashByte, err := bcrypt.GenerateFromPassword(passByte, bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hashByte), nil

}

func CheckPassword(hash, password string) error {
	if hash == "" || password == "" {
		return fmt.Errorf("EXPECTED TWO NON-NIL VALUES")
	}

	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return fmt.Errorf("INVALID PASSWORD")
	}

	return nil
}
