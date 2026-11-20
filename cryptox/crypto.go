package cryptox

import (
	"encoding/base64"

	"github.com/rotisserie/eris"
)

var (
	// Crypto 是包级默认实例，需要由 Init 或 SetCrypto 显式装配。
	Crypto PasswordCrypto
	// CryptoConfig 是 "crypto" 配置段的结构。本包不再自己读取配置，
	// 需要配置驱动的应用自行注册（见 _examples/cryptokit，或直接引用 storex/account）。
	CryptoConfig Config
)

type Config struct {
	AesKey string
}

type PasswordCrypto interface {
	Encrypt(plaintext []byte) (string, error)
	Decrypt(ciphertext string) (string, error)
}

// SetCrypto replaces the package-level default instance.
func SetCrypto(c PasswordCrypto) {
	Crypto = c
}

// Init builds the package-level Crypto from config.
// AesKey is a base64-encoded 16/24/32-byte key (AES-128/192/256).
// 未配置 AesKey 时不做初始化（也不报错），真正调用 Encrypt/Decrypt 时会返回 not initialized。
func Init(cfg Config) error {
	if cfg.AesKey == "" {
		return nil
	}

	keyBytes, err := base64.StdEncoding.DecodeString(cfg.AesKey)
	if err != nil {
		return eris.Wrap(err, "failed to decode aes key")
	}

	switch len(keyBytes) {
	case 16, 24, 32:
	default:
		return eris.Errorf("invalid aes key length %d bytes, want 16/24/32 bytes before base64 encoding", len(keyBytes))
	}

	Crypto = NewAesPasswordCrypto(keyBytes)
	return nil
}

func Encrypt(plaintext []byte) (string, error) {
	if Crypto == nil {
		return "", eris.New("cryptox is not initialized, call cryptox.Init or cryptox.SetCrypto first")
	}

	return Crypto.Encrypt(plaintext)
}

func Decrypt(ciphertext string) (string, error) {
	if Crypto == nil {
		return "", eris.New("cryptox is not initialized, call cryptox.Init or cryptox.SetCrypto first")
	}

	return Crypto.Decrypt(ciphertext)
}
