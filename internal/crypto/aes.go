package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

// NewAesCbcPkcs7Padding creates a new AesCbcCrypt with PKCS7 padding.
func NewAesCbcPkcs7Padding(key []byte) *AesCbcCrypt {
	return &AesCbcCrypt{
		key:     key,
		padding: pkcs7{},
	}
}

// AesCbcCrypt is the AES-CBC encryption method with configurable padding.
type AesCbcCrypt struct {
	key     []byte
	padding padding
}

func (c *AesCbcCrypt) Encode(plain []byte) ([]byte, error) {
	return c.Encrypt(plain)
}

func (c *AesCbcCrypt) Decode(encoded []byte) ([]byte, error) {
	return c.Decrypt(encoded)
}

// Encrypt encrypts the given plain bytes and returns the encrypted bytes.
func (c *AesCbcCrypt) Encrypt(plainBytes []byte) ([]byte, error) {
	paddedBytes, err := c.padding.pad(plainBytes, 32)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	iv := c.key[:aes.BlockSize]

	mode := cipher.NewCBCEncrypter(block, iv)

	encryptedBytes := make([]byte, len(paddedBytes))
	mode.CryptBlocks(encryptedBytes, paddedBytes)

	return encryptedBytes, nil
}

// Decrypt decrypts the given encrypted bytes and return plain bytes.
func (c *AesCbcCrypt) Decrypt(encryptedBytes []byte) ([]byte, error) {
	switch {
	case len(encryptedBytes) < aes.BlockSize:
		return nil, errors.New("invalid encrypted bytes size: smaller than block size")
	case len(encryptedBytes)%aes.BlockSize != 0:
		return nil, errors.New("invalid encrypted bytes size: not multiple of block size")
	}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	iv := c.key[:aes.BlockSize]

	mode := cipher.NewCBCDecrypter(block, iv)

	decryptedBytes := make([]byte, len(encryptedBytes))
	mode.CryptBlocks(decryptedBytes, encryptedBytes)

	return c.padding.unPad(decryptedBytes, 32)
}

type padding interface {
	pad(src []byte, blockSize int) ([]byte, error)
	unPad(src []byte, blockSize int) ([]byte, error)
}

type pkcs7 struct{}

func (_ pkcs7) pad(src []byte, blockSize int) ([]byte, error) {
	padLen := blockSize - (len(src) % blockSize)
	padded := bytes.Repeat([]byte{byte(padLen)}, padLen)

	var buffer bytes.Buffer
	buffer.Write(src)
	buffer.Write(padded)

	return buffer.Bytes(), nil
}

func (_ pkcs7) unPad(src []byte, blockSize int) ([]byte, error) {
	switch {
	case len(src) == 0:
		return nil, errors.New("zero length padded bytes in pkcs7 padding")
	case len(src)%blockSize != 0:
		return nil, errors.New("padded bytes length not multiple of block size in pkcs7 padding")
	}

	return src[:len(src)-int(src[len(src)-1])], nil
}
