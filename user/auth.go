package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Hashes a raw password using bcrypt
func hashPassword(password string) (string, error) {
	// arguments: raw password string, computational complexity represented as an integer
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Compares raw password with a hashed password
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type claims struct {
	UserId string
	jwt.RegisteredClaims
}
// Generates a new JWT with the user id and secret key
func generateAccessToken(userId string) (string, error) {
	claims := &claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.access_token_expiry_mins) * time.Minute)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.jwt_secret_key))
}