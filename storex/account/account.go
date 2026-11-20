package account

import (
	"context"
	"errors"
	"net/http"

	"gorm.io/gorm"

	"github.com/zq-xu/gotools/configx"
	"github.com/zq-xu/gotools/cryptox"
	"github.com/zq-xu/gotools/errorx"
	"github.com/zq-xu/gotools/routerx/auth"
	"github.com/zq-xu/gotools/storex"
)

const (
	AccountTableName = "accounts"
)

type Account struct {
	storex.Model `json:",inline"`

	Username string `description:"the username for login, like: alan_123"`
	Password string

	Name  string `description:"the real or nick name, like: Alan,Bob,Tommy"`
	Roles string
}

func init() {
	storex.RegisterTable(&Account{})

	// 账号密码以 AES 密文存储，因此 crypto 配置段由本包负责注册。
	// 只用 cryptox、不引用本包的应用，自行注册同样的一行即可。
	configx.Register("crypto", &cryptox.CryptoConfig, func() error {
		return cryptox.Init(cryptox.CryptoConfig)
	})
}

func (u *Account) TableName() string {
	return AccountTableName
}

func (u *Account) GetID() string       { return u.Model.GetID() }
func (u *Account) GetUsername() string { return u.Username }
func (u *Account) GetName() string     { return u.Name }
func (u *Account) GetRoles() string    { return u.Roles }
func (u *Account) GetStatus() string   { return u.Model.GetStatus() }

func LoadLoginAccount(ctx context.Context, username, password string) (auth.AuthAccount, errorx.ErrorInfo) {
	obj := &Account{}

	err := storex.DB(ctx).GetByField(obj, "username", username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NewError(http.StatusUnauthorized, "invailid accountname or passowrd", nil)
		}

		return nil, errorx.NewError(http.StatusBadRequest, "unexpected error", err)
	}

	str, err := cryptox.Decrypt(obj.Password)
	if err != nil {
		return nil, errorx.NewError(http.StatusInternalServerError, "failed to decrypt passowrd", err)
	}

	if password != str {
		return nil, errorx.NewError(http.StatusUnauthorized, "invailid accountname or passowrd", nil)
	}

	return obj, nil
}
