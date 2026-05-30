package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestWebSocketUpgraderAllowsLoopbackBrowserOrigin(t *testing.T) {
	upgrader := websocket.Upgrader{CheckOrigin: WebSocketOriginAllowed}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade: %v (origin=%q)", err, r.Header.Get("Origin"))
			return
		}
		_ = conn.Close()
	}))
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	header := http.Header{}
	header.Set("Origin", "http://localhost:5173")

	conn, resp, err := websocket.DefaultDialer.Dial(wsURL, header)
	if err != nil {
		if resp != nil {
			t.Fatalf("dial failed: %v (HTTP %d)", err, resp.StatusCode)
		}
		t.Fatalf("dial failed: %v", err)
	}
	_ = conn.Close()
}
