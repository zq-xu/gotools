package gotools

import (
	"github.com/zq-xu/gotools/apperror"
	"github.com/zq-xu/gotools/bricks/cryptokit"
	"github.com/zq-xu/gotools/config"
	"github.com/zq-xu/gotools/logs"
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
)

type ErrorInfo = apperror.ErrorInfo
