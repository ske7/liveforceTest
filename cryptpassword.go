package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword return string with hash of the user input password
func HashPassword(password []byte) string {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

// CheckPassword compares user input password with saved in DB
func CheckPassword(hash string, password []byte) bool {
	byteHash := []byte(hash)
	err := bcrypt.CompareHashAndPassword(byteHash, password)

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
