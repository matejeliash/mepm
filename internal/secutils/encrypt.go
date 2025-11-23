package secutils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"

	"github.com/matejeliash/mepm/internal/other"
)

var (
	ErrInvalidKey  = errors.New("invalid AES key, must be 32 bytes")
	ErrEncryptFail = errors.New("AES encryption failed")
	ErrDecryptFail = errors.New("AES decription failed")
)

// this function encrypts bytes and return encrypted bytes !!! (no strings)
func EncryptAES(plainBytes, key []byte) ([]byte, error) {

	if len(key) != 32 {
		return nil, ErrInvalidKey

	}

	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, other.WrapErr("EncryptAES", err)
	}

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		return nil, other.WrapErr("EncryptAES", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)

	if err != nil {
		return nil, other.WrapErr("EncryptAES", err)
	}

	encryptedBytes := gcm.Seal(nonce, nonce, plainBytes, nil)

	return encryptedBytes, nil

}

func DecryptAES(encryptedBytes, key []byte) ([]byte, error) {

	if len(key) != 32 {
		return nil, ErrInvalidKey

	}

	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, other.WrapErr("DecryptAES", err)
	}

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		return nil, other.WrapErr("DecryptAES", err)
	}
	nonce := encryptedBytes[:gcm.NonceSize()]
	ciphertext := encryptedBytes[gcm.NonceSize():]

	plainBytes, err := gcm.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		return nil, other.WrapErr("DecryptAES", err)
	}

	return plainBytes, nil

}
