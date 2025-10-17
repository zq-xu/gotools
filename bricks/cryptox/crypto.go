package cryptox

import (
	"encoding/base64"

	"github.com/rotisserie/eris"
	"github.com/zq-xu/gotools/configx"
)

var (
	Crypto       PasswordCrypto
	CryptoConfig Config
)

type Config struct {
	AesKey string
}

type PasswordCrypto interface {
	Encrypt(plaintext []byte) (string, error)
	Decrypt(ciphertext string) (string, error)
}

func init() {
	configx.Register("crypto", &CryptoConfig, initPasswordCrypto)
}

func initPasswordCrypto() error {
	keyBytes, err := base64.StdEncoding.DecodeString(CryptoConfig.AesKey)
	if err != nil {
		return eris.Wrap(err, "failed to decode aes key")
	}

	Crypto = NewAesPasswordCrypto(keyBytes)
	return nil
}

func Encrypt(plaintext []byte) (string, error) { return Crypto.Encrypt(plaintext) }

func Decrypt(ciphertext string) (string, error) { return Crypto.Decrypt(ciphertext) }
