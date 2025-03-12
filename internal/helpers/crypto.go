package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
)

const (
	SaltLen = 10
)

func Random(byteCount int) (b []byte, err error) {
	b = make([]byte, byteCount)
	_, err = io.ReadFull(rand.Reader, b)
	return
}

func generateHash(plainText string) (s string, err error) {
	slt, err := Random(SaltLen)
	if err != nil {
		return
	}

	s = generateHashWithSalt(plainText, slt)

	return
}

func generateHashWithSalt(plainText string, slt []byte) string {
	h := sha256.New()
	h.Write(slt)
	h.Write([]byte(plainText))

	return string(h.Sum(slt))
}

func verifyHash(plaintText, hash string) bool {
	s := []byte(hash)[:SaltLen]
	return hash == generateHashWithSalt(plaintText, s)
}

func GenerateHashUsingBase64URL(plainText string) (s string, err error) {
	slt, err := Random(SaltLen)
	if err != nil {
		return
	}

	s = generateHashWithSaltUsingBase64URL(plainText, slt)

	return
}

func generateHashWithSaltUsingBase64URL(plainText string, slt []byte) string {
	h := sha256.New()
	h.Write(slt)
	h.Write([]byte(plainText))

	return base64.URLEncoding.EncodeToString(h.Sum(slt))
}
func VerifyHashWithBase64URL(plaintText, hash string) bool {
	s, err := base64.URLEncoding.DecodeString(hash)
	if err != nil {
		return false
	}
	return hash == generateHashWithSaltUsingBase64URL(plaintText, s[:SaltLen])
}
