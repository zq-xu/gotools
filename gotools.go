package gotools

import (
	"github.com/zq-xu/gotools/apperror"
	"github.com/zq-xu/gotools/bricks/cryptokit"
	"github.com/zq-xu/gotools/config"
	"github.com/zq-xu/gotools/logs"
	"github.com/zq-xu/gotools/store"
)

var (
	Register = config.Register
	Setup    = config.Setup
)

var (
	Logger = logs.Logger
)

var (
	Encrypt = cryptokit.Encrypt
	Decrypt = cryptokit.Decrypt
)

var (
	Errorf   = apperror.Errorf
	NewError = apperror.NewError

	NewErrorInfoForCreateError = apperror.NewErrorInfoForCreateError
	NewErrorInfoForUpdateError = apperror.NewErrorInfoForUpdateError
	NewErrorInfoForDeleteError = apperror.NewErrorInfoForDeleteError
	NewErrorInfoForDBGetError  = apperror.NewErrorInfoForDBGetError
)

var (
	DB = store.DB
)

type ErrorInfo = apperror.ErrorInfo
