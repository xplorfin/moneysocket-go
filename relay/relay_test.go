package relay

import (
	"testing"

	"github.com/xplorfin/moneysocket-go/moneysocket/config"
	nettest "github.com/xplorfin/netutils/testutils"
)

func TestRelay(t *testing.T) {
	testConfig := config.NewConfig()
	testConfig.ListenConfig.BindPort = nettest.GetFreePort(t)
	testConfig.ListenConfig.BindPort = 11060
	relay := NewRelay(testConfig)
	t.Skip("todo")
	err := relay.RunApp()
	if err != nil {
		t.Error(err)
	}
}
