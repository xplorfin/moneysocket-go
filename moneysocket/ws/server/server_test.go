package server

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	nettest "github.com/xplorfin/netutils/testutils"
	"golang.org/x/sync/errgroup"
)

// TODO test stuff in here.
func TestListenUnsecure(t *testing.T) {
	g, _ := errgroup.WithContext(context.Background())

	protocol := NewTestWebsocketServiceProtocol()
	port := nettest.GetFreePort(t)
	testURL := fmt.Sprintf("%slocalhost:%d", "ws://", port)
	g.Go(func() error {
		return Listen(testURL, nil, func(writer http.ResponseWriter, request *http.Request) {
			ServeHTTP(&protocol, writer, request)
		})
	})

	nettest.AssertConnected(fmt.Sprintf("localhost:%d", port), t)
}

func TestListenSecure(t *testing.T) {
	t.Skip("TODO implement")
}
