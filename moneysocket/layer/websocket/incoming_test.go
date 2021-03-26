package websocket

import (
	"testing"

	. "github.com/stretchr/testify/assert"
	"github.com/xplorfin/moneysocket-go/moneysocket/config"
)

func TestNewIncomingWebsocketLayerUnsecure(t *testing.T) {
	configuration := config.NewConfig()
	incomingLayer := NewIncomingWebsocketLayer(configuration)
	err := incomingLayer.Listen("ws://localhost/", nil)
	Nil(t, err)
}
