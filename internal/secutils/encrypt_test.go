package secutils

import (
	"bytes"
	"crypto/rand"
	"math/big"
	"testing"
)

func TestEncrypt(t *testing.T) {

	var key []byte
	for i := 0; i < 32; i++ {
		randomNumber, err := rand.Int(rand.Reader, big.NewInt(256))
		if err != nil {
			t.Fatal("could not generate random numbers")
		}
		key = append(key, byte(randomNumber.Int64()))
	}

	var plainBytes []byte

	plainBytesLen, err := rand.Int(rand.Reader, big.NewInt(256))
	if err != nil {
		t.Fatal("could not generate random numbers")
	}

	for i := 0; i < int(plainBytesLen.Int64()); i++ {
		randomNumber, err := rand.Int(rand.Reader, big.NewInt(256))
		if err != nil {
			t.Fatal("could not generate random numbers")
		}

		plainBytes = append(plainBytes, byte(randomNumber.Int64()))
	}

	enceyptedBytes, err := EncryptAES(plainBytes, key)
	if err != nil {
		t.Fatal(err)
	}

	newPlainBytes, err := DecryptAES(enceyptedBytes, key)

	if !bytes.Equal(plainBytes, newPlainBytes) {
		t.Fatal("plaintexts are not same")

	}

	newPlainBytes, err = DecryptAES(enceyptedBytes[2:], key)

	if bytes.Equal(plainBytes, newPlainBytes) {
		t.Fatal("plaintexts match but they should not")

	}

}

func TestInvalidKey(t *testing.T) {

	key := []byte{1, 2, 3, 4, 5, 6}
	plainBytes := []byte{1, 2, 3, 4, 5, 5}
	_, err := EncryptAES(plainBytes, key)
	if err != ErrInvalidKey {
		t.Fatal(err)
	}

}
