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

// Compares raw password with a hashed password
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}