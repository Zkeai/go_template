package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

func GenerateSalt(size int) (string, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

func HashPassword(password, salt string) string {
	saltedPassword := salt + password
	hash := sha256.New()
	hash.Write([]byte(saltedPassword))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func CheckPassword(storedHash, password, salt string) bool {
	hashedPassword := HashPassword(password, salt)
	return hashedPassword == storedHash
}
