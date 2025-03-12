package helpers

import (
	"encoding/base64"
	"testing"
)

func TestGenerateHash(t *testing.T) {
	plainText := "password"
	hash, err := generateHash(plainText)
	if err != nil {
		t.Error(err)
	}
	if hash == "" {
		t.Error("Hash is empty")
	}
}
func TestGenerateHashBase64URL(t *testing.T) {
	plainText := "password"
	hash, err := GenerateHashUsingBase64URL(plainText)
	if err != nil {
		t.Error(err)
	}
	if hash == "" {
		t.Error("Hash is empty")
	}
}

func TestVerifyHash(t *testing.T) {
	plainText := "password"
	hash, err := generateHash(plainText)
	if err != nil {
		t.Error(err)
	}
	if !verifyHash(plainText, hash) {
		t.Error("Hash verification failed")
	}
}
func TestVerifyHashBase64URL(t *testing.T) {
	plainText := "password"
	hash, err := GenerateHashUsingBase64URL(plainText)
	if err != nil {
		t.Error(err)
	}
	if !VerifyHashWithBase64URL(plainText, hash) {
		t.Error("Hash verification failed")
	}
}
func TestHashSaltRecovery(t *testing.T) {
	plainText := "password"
	salt := []byte("salt-special")[:SaltLen]
	hash := generateHashWithSaltUsingBase64URL(plainText, []byte(salt))

	s, err := base64.URLEncoding.DecodeString(hash)
	if err != nil {
		t.Error(err)
	}
	actual := string(s[:SaltLen])
	expected := string(salt)

	if actual != expected {
		t.Errorf("Salt recovery failed: expected: %s, got: %s", expected, actual)
	}
}
