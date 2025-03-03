package auth

import (
	"errors"
	"testing"
)

func TestHashingPassword(t *testing.T) {

	tests := []struct {
		password string
		gotError error
	}{
		{
			password: "123456",
			gotError: nil,
		}, {
			password: "LongerPassword123",
			gotError: nil,
		}, {
			password: "ThisPasswordHasSpecial!@12$",
			gotError: nil,
		},
	}

	for _, test := range tests {
		_, err := HashPassword(test.password)

		if test.gotError != nil {
			t.Error("Error hashing password: %w", err)
		}
	}

}

func TestCheckPasswords(t *testing.T) {
	tests := []struct {
		password      string
		expectedError error
	}{
		{
			password:      "123456",
			expectedError: nil,
		}, {
			password:      "LongerPassword123",
			expectedError: nil,
		}, {
			password:      "ThisPasswordHasSpecial!@12$",
			expectedError: nil,
		},
	}

	for _, test := range tests {
		hashPW, _ := HashPassword(test.password)

		gotErr := CheckPasswordHash(test.password, hashPW)

		if gotErr != nil {
			t.Fatal(gotErr)
		}
	}
}

func TestHashPasswordLength(t *testing.T) {
	tests := []struct {
		password string
		wantErr  error
	}{
		{
			password: "ThisisalongpasswordstringtotestthelimitsofthehashfunctionIneedTomakeitpast72chars",
			wantErr:  errors.New("Password length exceeds maximum number of bytes"),
		},
	}

	for _, test := range tests {

		_, err := HashPassword(test.password)

		if test.wantErr.Error() != err.Error() {
			t.Fatal("Error does not meet expected:\t", err)
		}
	}
}
