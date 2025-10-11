package cryptokit

import (
	"encoding/base64"

	"github.com/rotisserie/eris"

	"zq-xu/gotools/config"
	"zq-xu/gotools/setup"
)

var (
	Crypto PasswordCrypto
)

type PasswordCrypto interface {
	Encrypt(plaintext []byte) (string, error)
	Decrypt(ciphertext string) (string, error)
}

func init() {
	setup.RegisterSetup("CryptoKit", initPasswordCrypto)
}

func initPasswordCrypto(cfg *config.Config) error {
	keyBytes, err := base64.StdEncoding.DecodeString(cfg.AesKey)
	if err != nil {
		return eris.Wrap(err, "failed to decode aes key")
	}

	Crypto = NewAesPasswordCrypto(keyBytes)
	return nil
}
