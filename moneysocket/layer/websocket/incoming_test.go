package websocket

import (
	"testing"

	"github.com/xplorfin/moneysocket-go/moneysocket/config"
)

func TestNewIncomingWebsocketLayerUnsecure(t *testing.T) {
	configuration := config.NewConfig()
	incomingLayer := NewIncomingWebsocketLayer(configuration)
	incomingLayer.Listen("ws://localhost/", nil)
}
