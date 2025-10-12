package gotools

import (
	"github.com/zq-xu/gotools/apperror"
	"github.com/zq-xu/gotools/bricks/cryptokit"
	"github.com/zq-xu/gotools/config"
	"github.com/zq-xu/gotools/logs"
	"github.com/zq-xu/gotools/router"
	"github.com/zq-xu/gotools/store"
	"github.com/zq-xu/gotools/utils"
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

type UnixTime = utils.UnixTime
