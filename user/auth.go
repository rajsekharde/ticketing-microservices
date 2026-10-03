package main

import (
	"golang.org/x/crypto/bcrypt"
)

// Hashes a raw password using bcrypt
func hashPassword(password string) (string, error) {
	// arguments: raw password string, computational complexity represented as an integer
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}