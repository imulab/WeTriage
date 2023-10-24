package route

import (
	"absurdlab.io/WeTriage/internal/crypto"
	"encoding/binary"
)

func decodeEncryptedMessage(aes *crypto.AesCbcCrypt, encrypted string) (decrypted []byte, receiveId string, err error) {
	const (
		discardBytes = 16
		lengthBytes  = 4
	)

	decrypted, err = crypto.Decode([]byte(encrypted),
		crypto.Base64Std,
		aes,
	)
	if err != nil {
		return
	}

	messageLen := binary.BigEndian.Uint32(decrypted[discardBytes : discardBytes+lengthBytes])

	receiveId = string(decrypted[discardBytes+lengthBytes+messageLen:])
	decrypted = decrypted[discardBytes+lengthBytes : discardBytes+lengthBytes+messageLen]

	return
}
