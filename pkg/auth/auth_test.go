package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "secret_password"
	hash, err := HashPassword(password)

	if err != nil {
		t.Fatalf("HashPassword failed with error: %v", err)
	}

	if hash == "" {
		t.Fatal("HashPassword returned an empty hash")
	}
}

func TestValidatePassword(t *testing.T) {
	password := "secret_password"
	hash, err := HashPassword(password)

	if err != nil {
		t.Fatalf("HashPassword failed with error: %v", err)
	}

	if !ValidatePassword(password, hash) {
		t.Fatal("ValidatePassword failed, password and hash do not match")
	}
}
