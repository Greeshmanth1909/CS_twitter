package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey []byte = []byte("top-secret-key-load-from-env") // This will usually be loaded from a .env file

// generateHash takes a given password and returns a hex encoded sha256 hash of the same.
func generateHash(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// generateJWT function returns a jwt token string signed with secretkey and an error if anything fails.
func generateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// verifyToken function takes a jwt token string and returns the corresponding claims. It returns a non nil error if the validation fails.
func verifyToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// Zip function takes two equal length string slices and returns a zipped up version of the same. It is similar to python's zip function.
func Zip(A []string, B []string) [][]string {
	out := [][]string{}

	for i := range A {
		new := []string{A[i], B[i]}
		out = append(out, new)
	}
	return out
}
