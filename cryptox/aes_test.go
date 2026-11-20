package cryptox

import (
	"bytes"
	"testing"
)

func TestAesPasswordCrypto_EncryptDecrypt(t *testing.T) {
	key := []byte("1234567890123456") // 16-byte AES key
	crypto := NewAesPasswordCrypto(key)

	tests := []struct {
		name      string
		plaintext string
	}{
		{"NormalText", "HelloWorld123"},
		{"EmptyText", ""},
		{"ShortText", "abc"},
		{"LongText", "This is a very long plaintext string to test AES CBC with PKCS7 padding."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ciphertext, err := crypto.Encrypt([]byte(tt.plaintext))
			if err != nil {
				t.Fatalf("Encrypt error: %v", err)
			}

			if len(ciphertext) == 0 {
				t.Fatalf("Ciphertext is empty")
			}

			decrypted, err := crypto.Decrypt(ciphertext)
			if err != nil {
				t.Fatalf("Decrypt error: %v", err)
			}

			if !bytes.Equal([]byte(tt.plaintext), []byte(decrypted)) {
				t.Fatalf("Decrypted text does not match original.\nGot: %s\nWant: %s", decrypted, tt.plaintext)
			}
		})
	}
}
