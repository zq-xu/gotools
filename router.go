package gotools

// THIS FILE IS AUTO-GENERATED. DO NOT EDIT MANUALLY.

import "github.com/zq-xu/gotools/router"

var (
    // From config.go
    RouterConfig = router.RouterConfig

    // From auth.go
    Login = router.Login

    Logout = router.Logout

    AuthMiddleware = router.AuthMiddleware

    InitAuthMiddleware = router.InitAuthMiddleware

    GetUserInfoFromToken = router.GetUserInfoFromToken

    GetAccountInfoHandler = router.GetAccountInfoHandler

    // From health.go
    SetHealthHandler = router.SetHealthHandler

    Health = router.Health

    HealthPath = router.HealthPath

    // From group.go
    NewGroup = router.NewGroup

    // From log_filter.go
    FormatterMiddleWare = router.FormatterMiddleWare

    LogInterrupt = router.LogInterrupt

    // From router.go
    RegisterGroup = router.RegisterGroup

    StartRouter = router.StartRouter

    NewRouter = router.NewRouter

    MaxMultipartMemory = router.MaxMultipartMemory

)

type (
    // From config.go
    Config = router.Config

    // From group.go
    APIGroup = router.APIGroup

    API = router.API

)
