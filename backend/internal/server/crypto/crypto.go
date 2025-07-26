package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
)

func EncryptAESGCM(key []byte, plaintext []byte) (ivBase64 string, ciphertextBase64 string, err error) {
	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", err
	}

	// Use GCM mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}

	// Generate a 12-byte IV (required size for AES-GCM)
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", "", err
	}

	// Encrypt
	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)

	// Return both nonce (IV) and ciphertext as base64
	return base64.StdEncoding.EncodeToString(nonce), base64.StdEncoding.EncodeToString(ciphertext), nil
}

func EncryptJSON(key []byte, data interface{}) (ivBase64 string, ciphertextBase64 string, err error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", "", err
	}
	return EncryptAESGCM(key, jsonData)
}

func DecryptAESGCM(key []byte, ivBase64 string, ciphertextBase64 string) (plaintext []byte, err error) {
	// Decode IV and ciphertext from base64
	iv, err := base64.StdEncoding.DecodeString(ivBase64)
	if err != nil {
		return nil, err
	}
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return nil, err
	}

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Use GCM mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Decrypt
	plaintext, err = aesGCM.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func DecryptJSON(key []byte, ivBase64 string, ciphertextBase64 string, target interface{}) error {
	plaintext, err := DecryptAESGCM(key, ivBase64, ciphertextBase64)
	if err != nil {
		return err
	}
	err = json.Unmarshal(plaintext, target)
	if err != nil {
		return err
	}
	return nil
}