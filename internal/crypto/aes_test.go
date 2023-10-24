package crypto

import (
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAesCbcPkcs7Crypt(t *testing.T) {
	plain := "abcdefghijklmnopqrstuvwxyz0123456789"

	key := make([]byte, 32)
	if _, err := rand.Read(key); !assert.NoError(t, err) {
		return
	}

	aes := NewAesCbcPkcs7Padding(key)

	encryptedBytes, err := aes.Encrypt([]byte(plain))
	if !assert.NoError(t, err) {
		return
	}

	decryptedBytes, err := aes.Decrypt(encryptedBytes)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, plain, string(decryptedBytes))
}
