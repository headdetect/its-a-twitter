package utils

import (
	"crypto/rand"
	"encoding/base64"
	
	"golang.org/x/crypto/bcrypt"
)

func RandomString(length int) string {
	randomBytes := make([]byte, length)
	rand.Read(randomBytes)
	return base64.URLEncoding.EncodeToString(randomBytes)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}