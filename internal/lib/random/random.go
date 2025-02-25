package random

import (
	"crypto/rand"
)

const Charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_"

func RandomString(size int) (string, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	result := make([]rune, size)
	for i := range b {
		result[i] = rune(Charset[int(b[i])%len(Charset)])
	}

	return string(result), nil
}
