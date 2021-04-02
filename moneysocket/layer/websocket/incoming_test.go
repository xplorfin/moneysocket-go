package websocket

import (
	"fmt"
	"testing"

	. "github.com/stretchr/testify/assert"
	"github.com/xplorfin/moneysocket-go/moneysocket/config"
	"github.com/xplorfin/netutils/testutils"
)

func TestNewIncomingWebsocketLayerUnsecure(t *testing.T) {
	configuration := config.NewConfig()
	incomingLayer := NewIncomingWebsocketLayer(configuration)
	host := fmt.Sprintf("ws://localhost:%d", testutils.GetFreePort(t))
	go func() {
		err := incomingLayer.Listen(host, nil)
		Nil(t, err)
	}()
	testutils.AssertConnected(host, t)
}
