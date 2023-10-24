package crypto

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
)

var (
	NoEncoding Encoding = noEnc{}
	Base64Std  Encoding = base64Std{}
	SHA1       Encoding = sha1Enc{}
	Hex        Encoding = hexEnc{}
)

// Encode chain encods the plain bytes with sequential encodings.
func Encode(plain []byte, encodings ...Encoding) (encoded []byte, err error) {
	encoded = plain

	for _, enc := range encodings {
		encoded, err = enc.Encode(encoded)
		if err != nil {
			return
		}
	}

	return
}

// Decode chain decodes the encoded bytes with sequential encodings.
func Decode(encoded []byte, encodings ...Encoding) (decoded []byte, err error) {
	decoded = encoded

	for _, enc := range encodings {
		decoded, err = enc.Decode(decoded)
		if err != nil {
			return
		}
	}

	return
}

// Encoding abstracts transformation between encoding of bytes
type Encoding interface {
	// Encode transforms plain bytes to encoded bytes.
	Encode(plain []byte) ([]byte, error)
	// Decode transforms encoded bytes to decoded bytes.
	Decode(encoded []byte) ([]byte, error)
}

type noEnc struct{}

func (_ noEnc) Encode(plain []byte) ([]byte, error) {
	return plain, nil
}

func (_ noEnc) Decode(encoded []byte) ([]byte, error) {
	return encoded, nil
}

type base64Std struct{}

func (_ base64Std) Encode(plain []byte) ([]byte, error) {
	return []byte(base64.StdEncoding.EncodeToString(plain)), nil
}

func (_ base64Std) Decode(encoded []byte) ([]byte, error) {
	if m := len(encoded) % 4; m != 0 {
		encoded = append(encoded, bytes.Repeat([]byte("="), 4-m)...)
	}

	return base64.StdEncoding.DecodeString(string(encoded))
}

type sha1Enc struct{}

func (_ sha1Enc) Encode(plain []byte) ([]byte, error) {
	h := sha1.New()
	h.Write(plain)
	return h.Sum(nil), nil
}

func (_ sha1Enc) Decode(_ []byte) ([]byte, error) {
	panic("SHA-1 hash does not support decoding")
}

type hexEnc struct{}

func (_ hexEnc) Encode(plain []byte) ([]byte, error) {
	return []byte(hex.EncodeToString(plain)), nil
}

func (_ hexEnc) Decode(encoded []byte) ([]byte, error) {
	return hex.DecodeString(string(encoded))
}
