package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"

	"golang.org/x/crypto/argon2"
)

func hashKey(key, salt []byte) []byte {
	return argon2.IDKey(key, salt, 1, 64*1024, 4, 32)
}

func Encrypt(data, key []byte) ([]byte, error) {
	salt := make([]byte, 32)
	if _, err := rand.Read(salt[:]); err != nil {
		return nil, err
	}

	key = hashKey(key, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	out := append(salt, ciphertext...)

	return out, nil
}

func Decrypt(data, key []byte) ([]byte, error) {
	if len(data) < 32 {
		return nil, errors.New("ciphertext too short")
	}

	key = hashKey(key, data[:32])

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[32:nonceSize+32], data[nonceSize+32:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
