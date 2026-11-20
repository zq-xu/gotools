package cryptox

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// implement interface PasswordCrypto
type aesPasswordCrypto struct {
	key []byte
}

func NewAesPasswordCrypto(key []byte) PasswordCrypto {
	return &aesPasswordCrypto{key: key}
}

// Encrypt
func (a *aesPasswordCrypto) Encrypt(plaintext []byte) (string, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", err
	}

	// PKCS#7
	plaintext = pkcs7Pad(plaintext, block.BlockSize())

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt
func (a *aesPasswordCrypto) Decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(data) < aes.BlockSize || len(data)%aes.BlockSize != 0 {
		return "", fmt.Errorf("invalid ciphertext length")
	}

	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", err
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(data, data)

	// remove PKCS#7 padding
	data, err = pkcs7Unpad(data)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// ---------------- PKCS#7 padding ----------------

func pkcs7Pad(data []byte, blockSize int) []byte {
	pad := blockSize - len(data)%blockSize
	padding := bytes.Repeat([]byte{byte(pad)}, pad)
	return append(data, padding...)
}

func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("invalid padding size")
	}
	pad := int(data[len(data)-1])
	if pad <= 0 || pad > len(data) {
		return nil, fmt.Errorf("invalid padding size")
	}
	for i := len(data) - pad; i < len(data); i++ {
		if int(data[i]) != pad {
			return nil, fmt.Errorf("invalid padding bytes")
		}
	}
	return data[:len(data)-pad], nil
}
