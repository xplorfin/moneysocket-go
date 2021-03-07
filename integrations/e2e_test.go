package integrations

import (
	"context"
	"testing"
	"time"

	"github.com/Flaque/filet"
	"github.com/xplorfin/moneysocket-go/moneysocket/config"
	"github.com/xplorfin/moneysocket-go/relay"
	"github.com/xplorfin/moneysocket-go/terminus"
	nettest "github.com/xplorfin/netutils/testutils"
)

func makeConfig(t *testing.T) *config.Config {
	testConfig := config.NewConfig()
	testConfig.AccountPersistDir = filet.TmpDir(t, "")
	testConfig.ListenConfig.BindPort = nettest.GetFreePort(t)
	testConfig.ListenConfig.BindHost = "localhost"
	testConfig.RpcConfig.BindHost = "localhost"
	testConfig.RpcConfig.BindPort = nettest.GetFreePort(t)
	return testConfig
}

func TestE2E(t *testing.T) {
	cfg := makeConfig(t)
	ctx := context.Background()

	// setup test relay
	testRelay := relay.NewRelay(cfg)
	go testRelay.RunApp()

	// setup test rpc server
	testRpcServer := terminus.NewTerminus(cfg)
	go testRpcServer.Start(ctx)

	// test rpc server hostname
	nettest.AssertConnected(cfg.GetRpcHostname(), t)

	terminusClient := terminus.NewClient(cfg)
	// create two accounts
	// -- account 1
	res, err := terminusClient.CreateAccount(1000000)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
	// -- acount 2
	res, err = terminusClient.CreateAccount(1000000)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)

	// grab info, TODO validate

	res, err = terminusClient.GetInfo()
	if err != nil {
		t.Error(err)
	}
	t.Log(res)

	// listen on account-1

	res, err = terminusClient.Listen("0")
	if err != nil {
		t.Error(err)
	}
	walletCon := NewWalletConsumer(cfg.GetBindHost(), cfg.GetUseTls(), cfg.GetBindPort())
	err = walletCon.DoConnect(walletCon.ConsumerBeacon)
	if true {
		walletCon := NewWalletConsumer(cfg.GetBindHost(), cfg.GetUseTls(), cfg.GetBindPort())
		err = walletCon.DoConnect(walletCon.ConsumerBeacon)
		if err != nil {
			t.Error(err)
		}
	}
	app := NewSellerApp(cfg.GetBindHost(), cfg.GetUseTls(), cfg.GetBindPort())
	err = app.ConsumerStack.DoConnect(makeConsumerBeacon(cfg.GetBindHost(), cfg.GetUseTls(), cfg.GetBindPort()))
	if err != nil {
		t.Error(err)
	}
	NewSellerApp(cfg.GetBindHost(), cfg.GetUseTls(), cfg.GetBindPort())
	time.Sleep(time.Second * 10)
}
