package cryptox

import (
	"encoding/base64"
	"testing"
)

func TestInitAndEncryptDecrypt(t *testing.T) {
	key := base64.StdEncoding.EncodeToString([]byte("0123456789abcdef")) // 16 字节

	if err := Init(Config{AesKey: key}); err != nil {
		t.Fatalf("Init error: %v", err)
	}

	ciphertext, err := Encrypt([]byte("hello"))
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
	}

	plaintext, err := Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decrypt error: %v", err)
	}

	if plaintext != "hello" {
		t.Fatalf("Decrypt() = %q, want %q", plaintext, "hello")
	}
}

func TestInitRejectsInvalidKeyLength(t *testing.T) {
	// 16/24/32 字节之外的密钥长度都不是合法 AES 密钥
	key := base64.StdEncoding.EncodeToString([]byte("too-short"))

	if err := Init(Config{AesKey: key}); err == nil {
		t.Fatalf("expect an error for a 9-byte key")
	}
}

func TestInitWithEmptyKeySkipsInitialization(t *testing.T) {
	prev := Crypto
	Crypto = nil

	defer func() { Crypto = prev }()

	if err := Init(Config{}); err != nil {
		t.Fatalf("Init with empty key should not fail, got %v", err)
	}

	if Crypto != nil {
		t.Fatalf("Crypto should stay nil when aesKey is empty")
	}
}

func TestEncryptWithoutInit(t *testing.T) {
	prev := Crypto
	Crypto = nil

	defer func() { Crypto = prev }()

	if _, err := Encrypt([]byte("x")); err == nil {
		t.Fatalf("expect an error when Crypto is not initialized")
	}
}
