package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestTokenCreation(t *testing.T) {
	testId, _ := uuid.NewRandom()
	testSecret := "Security"
	testDuration, _ := time.ParseDuration("2h30m")
	testString, err := MakeJWT(testId, testSecret, testDuration)
	if err != nil {
		t.Errorf("Failed to generate token: %s", err)
	}

	//Check token creation
	testClaims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(testString, &testClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(testSecret), nil
	})
	if err != nil {
		t.Errorf("Failed to retrieve token, so that's broken: %s", err)
	}

	tokenId, err := testClaims.GetSubject()
	if err != nil {
		t.Errorf("Failed to retrieve subject from token: %s", err)
	}

	tokenUuid, err := uuid.Parse(tokenId)
	if err != nil {
		t.Errorf("Failed to parse subject so it's probably not a UUID: %s", err)
	}

	//Check testid is in token
	if testId != tokenUuid {
		t.Errorf("Failed to get matched token ID: %s", err)
	}

	//Check duration
	tokenCreated, _ := testClaims.GetIssuedAt()
	tokenExpires, _ := testClaims.GetExpirationTime()
	tokenDuration := tokenExpires.Time.Sub(tokenCreated.Time)

	if testDuration != tokenDuration {
		t.Errorf("Somehow, token expires in %s but it actually expires in %s", testDuration, tokenDuration)
	}

}
