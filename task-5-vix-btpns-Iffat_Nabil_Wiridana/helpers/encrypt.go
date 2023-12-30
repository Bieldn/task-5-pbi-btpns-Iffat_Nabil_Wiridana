package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
)

func NewRandomKey() []byte {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}

	return key
}

func Encrypt(key []byte) ([]byte, error) {
	secretKey := getSecret()

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aesgcm.NonceSize())
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(iv, iv, key, nil)

	return ciphertext, nil
}

func Decrypt(encryptedKey []byte) ([]byte, error) {
	secretKey := getSecret()

	block, err := aes.NewCipher(secretKey)

	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(encryptedKey) < aesgcm.NonceSize() {
		// worth panicking when encrypted key is bad
		panic("Malformed encrypted key")
	}

	return aesgcm.Open(
		nil,
		encryptedKey[:aesgcm.NonceSize()],
		encryptedKey[aesgcm.NonceSize():],
		nil,
	)
}

func getSecret() []byte {
	secret := GetEnv("ENCRYPT_SECRET_KEY")
	if secret == "" {
		panic("Error: Must provide a secret key under env variable SECRET")
	}

	secretbite, err := hex.DecodeString(secret)

	if err != nil {
		// probably malform secret, panic out
		panic(err)
	}

	return secretbite
}
