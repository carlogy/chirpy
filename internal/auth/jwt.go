package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userid uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  &jwt.NumericDate{Time: time.Now()},
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(expiresIn)},
		Subject:   userid.String(),
	})

	token, err := tokenClaims.SignedString([]byte(tokenSecret))

	if err != nil {
		return "", err
	}
	return token, err

}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.MapClaims{}
	t, err := jwt.ParseWithClaims(tokenString, &claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		})

	if err != nil {
		log.Print("Invalid token used:\t", err)
		return uuid.Nil, err
	}

	subject, err := t.Claims.GetSubject()
	if err != nil {
		log.Fatal(err)
		return uuid.Nil, err
	}

	userId, err := uuid.Parse(subject)
	if err != nil {
		log.Fatal(err)
		return uuid.Nil, err
	}

	return userId, nil

}

func GetBearerToken(headers http.Header) (string, error) {

	if len(headers) < 1 {


		return "", errors.New("Error parsing headers, expecting non-empty headers")
	}

	token := headers.Get("Authorization")

	if token == "" {


		return "", errors.New("Authorization KeyValue header not in headers")
	}

	token = strings.TrimPrefix(token, "Bearer")
	token = strings.TrimSpace(token)

	return token, nil
}
