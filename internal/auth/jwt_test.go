package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {

	tests := []struct {
		UserID     uuid.UUID
		expiryTime time.Duration
		secret     string
	}{
		{
			UserID:     uuid.New(),
			expiryTime: time.Duration(time.Minute),
			secret:     "thisIsASecret",
		},
	}

	for _, test := range tests {

		_, err := MakeJWT(test.UserID, test.secret, test.expiryTime)

		if err != nil {
			t.Fatal("Error creating token:\t", err)
		}

	}

}

func TestValidateJWT(t *testing.T) {

	tests := []struct {
		ExpectedUserID uuid.UUID
		TokenSecret    string
	}{
		{
			ExpectedUserID: uuid.New(),
			TokenSecret:    "thisIsASecret",
		},
		{
			ExpectedUserID: uuid.New(),
			TokenSecret:    "",
		},
	}

	for _, test := range tests {

		token, err := MakeJWT(test.ExpectedUserID, test.TokenSecret, time.Duration(time.Minute))

		if err != nil {
			t.Fatal("Error creating token for test:\t", err)
		}

		gotUserID, err := ValidateJWT(token, test.TokenSecret)

		if gotUserID != test.ExpectedUserID {
			t.Fatalf("UserID's don't match:\nGot:\t%v\tExpected:\t%v\n", gotUserID, test.ExpectedUserID)
		}

	}

}

func TestCorrectErrValidateToken(t *testing.T) {

	tests := []struct {
		UserId      uuid.UUID
		TokenSecret string
	}{
		{
			UserId:      uuid.New(),
			TokenSecret: "FailedSecret",
		},
	}

	for _, test := range tests {

		token, err := MakeJWT(test.UserId, test.TokenSecret, time.Duration(time.Minute))

		if err != nil {
			t.Errorf("Error creating token:\t%v", err)
		}

		_, err = ValidateJWT(token, "IncorrectSecret")

		if err != nil && err.Error() != "token signature is invalid: signature is invalid" {

			t.Errorf("Error mismatch error:\t%v", err.Error())

		}

	}

}

func TestGetBearerToken(t *testing.T) {

	tests := []struct {
		h         http.Header
		expecting string
	}{
		{
			h:         http.Header{"Authorization": {"Bearer TokenStringHere123"}},
			expecting: "TokenStringHere123",
		}, {
			h:         http.Header{"Authorization": {"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflR86k0izqG_x7zdvI96rC23zhQ9m5wrj1Xv4sHhkQ"}},
			expecting: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflR86k0izqG_x7zdvI96rC23zhQ9m5wrj1Xv4sHhkQ",
		},
	}

	for _, test := range tests {

		got, err := GetBearerToken(test.h)
		if err != nil {
			t.Errorf("Expecting token:\t%s\tGot:\t%v", test.expecting, err)
		}

		if got != test.expecting {
			t.Errorf("Expecting:\t%s\tGot:%s", test.expecting, got)
		}
	}
}

func TestErrorGetBearerToken(t *testing.T) {
	tests := []struct {
		h        http.Header
		expected string
	}{
		{
			h:        http.Header{},
			expected: "Error parsing headers, expecting non-empty headers",
		}, {
			h:        http.Header{"Content-Type": {"Invalid"}},
			expected: "",
		},
	}

	for _, test := range tests {

		got, err := GetBearerToken(test.h)

		if err == nil {
			t.Error("Expecting error got nil")
		}

		if got != "" {
			t.Errorf("Expecting empty token string got: %s", got)
		}

	}
}
