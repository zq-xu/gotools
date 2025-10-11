package gotools

import (
	"zq-xu/gotools/apperror"
	"zq-xu/gotools/bricks/cryptokit"
	"zq-xu/gotools/config"
	"zq-xu/gotools/logs"
	"zq-xu/gotools/router"
	"zq-xu/gotools/store"
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
	AuthMiddleware     = router.AuthMiddleware
	InitAuthMiddleware = router.InitAuthMiddleware

	NewRouter     = router.NewRouter
	RegisterGroup = router.RegisterGroup

	StartRouter = router.StartRouter
)

var (
	DB = store.DB
)

type ErrorInfo = apperror.ErrorInfo
