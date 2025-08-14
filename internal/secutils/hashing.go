package secutils

import (
	"crypto/rand"
	"crypto/subtle"

	"golang.org/x/crypto/argon2"
)

type KeyGenParams struct {
	Iters       uint32
	MemoryBytes uint32
	Threads     uint8
	KeyLen      uint32
}

func GetDefaultKeyGenParams() *KeyGenParams {
	params := &KeyGenParams{
		Iters:       4,
		MemoryBytes: 1024 * 64,
		Threads:     4,
		KeyLen:      32,
	}

	return params
}

func GenerateKey(password, salt []byte, params *KeyGenParams) []byte {

	// derive the key
	key := argon2.IDKey(
		password,
		salt,
		params.Iters,
		params.MemoryBytes,
		params.Threads,
		params.KeyLen,
	)

	return key

}

func HashPassword(password, salt []byte, params *KeyGenParams) []byte {

	// derive the key
	key := argon2.IDKey(
		password,
		salt,
		params.Iters,
		params.MemoryBytes,
		params.Threads,
		params.KeyLen,
	)

	return key

}

func GenerateSalt() ([]byte, error) {
	b, err := GenerateRandomBytes(16)
	if err != nil {
		return nil, err

	}

	return b, nil
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func ComparePasswordAndHash(password, passwordSalt, passwordHash []byte, params *KeyGenParams) bool {

	computedHash := HashPassword(password, passwordSalt, params)

	return subtle.ConstantTimeCompare(passwordHash, computedHash) == 1

}
