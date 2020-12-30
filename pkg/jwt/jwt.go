package jwt

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// secret key being used to sign tokens
var (
	// SecretKey see above
	SecretKey = []byte("secret")
)

// GenerateToken generates a jwt token and assigns a username to its claims
func GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)
	/* Set token claims */
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		log.Fatal("Error in Generating key")
		return "", err
	}
	return tokenString, nil
}

// ParseToken parses a jwt token and returns the username in its claims
func ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		username := claims["username"].(string)
		return username, nil
	}
	return "", err
}
