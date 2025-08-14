package secutils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
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
		return nil, fmt.Errorf("%w: %v", ErrEncryptFail, err)
	}

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrEncryptFail, err)
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrEncryptFail, err)
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
		return nil, fmt.Errorf("%w: %v", ErrDecryptFail, err)
	}

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		return nil, ErrDecryptFail
	}
	nonce := encryptedBytes[:gcm.NonceSize()]
	ciphertext := encryptedBytes[gcm.NonceSize():]

	plainBytes, err := gcm.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDecryptFail, err)
	}

	return plainBytes, nil

}
