package websocket

import (
	"net/http"
	"testing"
)

func newOriginReq(origin string) *http.Request {
	r, _ := http.NewRequest(http.MethodGet, "http://example.com/ws", nil)
	r.Host = "example.com"
	if origin != "" {
		r.Header.Set("Origin", origin)
	}
	return r
}

func TestSameOriginCheck(t *testing.T) {
	if !sameOriginCheck(newOriginReq("")) {
		t.Fatal("request without Origin should pass (non-browser client)")
	}
	if !sameOriginCheck(newOriginReq("http://example.com")) {
		t.Fatal("same origin should pass")
	}
	if sameOriginCheck(newOriginReq("http://evil.com")) {
		t.Fatal("cross origin should be rejected")
	}
}

func TestSetCheckOrigin(t *testing.T) {
	defer SetCheckOrigin(nil)

	AllowAllOrigins()
	checkOriginLock.RLock()
	fn := checkOrigin
	checkOriginLock.RUnlock()
	if !fn(newOriginReq("http://evil.com")) {
		t.Fatal("AllowAllOrigins should accept any origin")
	}

	SetCheckOrigin(nil)
	checkOriginLock.RLock()
	fn = checkOrigin
	checkOriginLock.RUnlock()
	if fn(newOriginReq("http://evil.com")) {
		t.Fatal("SetCheckOrigin(nil) should restore same-origin check")
	}
}
