package terminus

import (
	"context"
	"testing"

	"github.com/Flaque/filet"
	"github.com/xplorfin/moneysocket-go/moneysocket/config"
	nettest "github.com/xplorfin/netutils/testutils"
)

// create a new test server on port specified in config.port on localhost for testing
func GetTestServer(t *testing.T) (server Terminus, configuration *config.Config) {
	configuration = config.NewConfig()
	configuration.AccountPersistDir = filet.TmpDir(t, "")
	configuration.ListenConfig.BindHost = "localhost"
	configuration.ListenConfig.BindPort = nettest.GetFreePort(t)
	configuration.RpcConfig.BindHost = "localhost"
	configuration.RpcConfig.BindPort = nettest.GetFreePort(t)

	server = NewTerminus(configuration)
	return server, configuration
}

// start a test server on port specified in config on server
func GetStartedTestServer(t *testing.T) (server Terminus, config *config.Config) {
	server, config = GetTestServer(t)
	go server.Start(context.Background())
	nettest.AssertConnected(config.GetRpcHostname(), t)
	return server, config
}

func TestGetInfo(t *testing.T) {
	server, _ := GetStartedTestServer(t)
	client := NewClient(server.config)
	res, err := client.GetInfo()
	if err != nil {
		t.Error(err)
	}
	_ = res
	// TODO validate
}

func TestCreateAccountListen(t *testing.T) {
	server, _ := GetStartedTestServer(t)
	client := NewClient(server.config)
	res, err := client.CreateAccount(1000)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)

	res, err = client.Listen("0")
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
	_ = res
}
