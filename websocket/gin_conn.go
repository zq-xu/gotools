package websocket

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	checkOriginLock sync.RWMutex
	checkOrigin     = sameOriginCheck
)

// SetCheckOrigin 覆盖握手时的 Origin 校验（可用于校验 token、跨域白名单等），
// 传 nil 恢复默认的同源校验。
func SetCheckOrigin(fn func(r *http.Request) bool) {
	checkOriginLock.Lock()
	defer checkOriginLock.Unlock()

	if fn == nil {
		checkOrigin = sameOriginCheck
		return
	}
	checkOrigin = fn
}

// AllowAllOrigins 关闭 Origin 校验，允许任意站点握手（公网服务请勿开启）。
func AllowAllOrigins() {
	SetCheckOrigin(func(r *http.Request) bool { return true })
}

// sameOriginCheck 同源校验：没有 Origin 头（非浏览器客户端）直接放行，
// 有则要求 Origin 的 host 与请求 host 一致。
func sameOriginCheck(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}

	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	return strings.EqualFold(u.Host, r.Host)
}

func NewConnFromGin(ctx *gin.Context) (*websocket.Conn, error) {
	checkOriginLock.RLock()
	fn := checkOrigin
	checkOriginLock.RUnlock()

	upGrader := websocket.Upgrader{
		CheckOrigin:  fn,
		Subprotocols: []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
	}

	conn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return nil, fmt.Errorf("websocket connect error: %v", err)
	}

	return conn, nil
}
