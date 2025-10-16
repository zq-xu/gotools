package account

import (
	"context"
	"errors"
	"net/http"

	"gorm.io/gorm"

	"github.com/zq-xu/gotools/apperror"
	"github.com/zq-xu/gotools/bricks/cryptokit"
	"github.com/zq-xu/gotools/router/auth"
	"github.com/zq-xu/gotools/store"
)

const (
	AccountTableName = "accounts"
)

type Account struct {
	store.Model `json:",inline"`

	Username string `description:"the username for login, like: alan_123"`
	Password string

	Name  string `description:"the real or nick name, like: Alan,Bob,Tommy"`
	Roles string
}

func init() {
	store.RegisterTable(&Account{})
}

func (u *Account) TableName() string {
	return AccountTableName
}

func (u *Account) GetID() string       { return u.Model.GetID() }
func (u *Account) GetUsername() string { return u.Username }
func (u *Account) GetName() string     { return u.Name }
func (u *Account) GetRoles() string    { return u.Roles }
func (u *Account) GetStatus() string   { return u.Model.GetStatus() }

func LoadLoginAccount(ctx context.Context, username, password string) (auth.AuthAccount, apperror.ErrorInfo) {
	obj := &Account{}

	err := store.DB(ctx).GetByField(obj, "username", username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewError(http.StatusUnauthorized, "invailid accountname or passowrd", nil)
		}

		return nil, apperror.NewError(http.StatusBadRequest, "unexpected error", err)
	}

	str, err := cryptokit.Crypto.Decrypt(obj.Password)
	if err != nil {
		return nil, apperror.NewError(http.StatusInternalServerError, "failed to decrypt passowrd", err)
	}

	if password != str {
		return nil, apperror.NewError(http.StatusUnauthorized, "invailid accountname or passowrd", nil)
	}

	return obj, nil
}
