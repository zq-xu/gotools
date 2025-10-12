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

	Name string `description:"the real or nick name, like: Alan,Bob,Tommy"`
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
func (u *Account) GetStatus() string   { return u.Model.GetStatus() }

func GetAccount(ctx context.Context, id string) (*Account, apperror.ErrorInfo) {
	obj := &Account{}
	ei := apperror.NewErrorInfoForDBGetError(store.DB().Get(ctx, obj, id))
	return obj, ei
}

func LoadLoginAccount(ctx context.Context, username, password string) (auth.AuthAccount, apperror.ErrorInfo) {
	obj := &Account{}

	err := store.DB().GetByField(ctx, obj, "username", username)
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
