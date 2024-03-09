package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(input string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(input), 14)
	if err != nil {
		log.Fatalf("unable to hash a password: %s", err)
	}

	return string(hash)
}

func VerifyPassword(hash, input string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(input)) == nil
}
