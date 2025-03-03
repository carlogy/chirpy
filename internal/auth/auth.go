package auth

import (
	"errors"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {

	if len(password) > 72 {
		return "", errors.New("Password length exceeds maximum number of bytes")
	}

	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}
	return string(hashedpassword), nil
}
func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		return err
	}

	return nil
}

func GetAPIKey(headers http.Header) (string, error) {

	if len(headers) < 1 {
		return "", errors.New("Error parsing headers, expecting non-empty headers")
	}

	key := headers.Get("Authorization")

	if key == "" {
		return "", errors.New("Authorization Key-Value header not in headers")
	}

	key = strings.TrimPrefix(key, "ApiKey")
	key = strings.TrimSpace(key)

	return key, nil
}
