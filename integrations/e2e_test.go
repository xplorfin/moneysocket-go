package integrations

import (
	"context"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"strings"
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
	testConfig.ListenConfig.ExternalHost = testConfig.GetBindHost()
	testConfig.ListenConfig.ExternalPort = testConfig.GetBindPort()

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
	account1Beacon, err := terminusClient.CreateAccount(1000000)
	if err != nil {
		t.Error(err)
	}
	t.Log(account1Beacon)
	// -- acount 2
	account1Beacon, err = terminusClient.CreateAccount(1000000)
	if err != nil {
		t.Error(err)
	}
	t.Log(account1Beacon)

	// start the seller wallet consumer on account-1
	account1Listen := getBeacon(t, terminusClient, "1")
	app := NewSellerApp(account1Listen)
	err = app.ConsumerStack.DoConnect(account1Listen)
	if err != nil {
		t.Error(err)
	}

	// connect to wallet consumer
	account0Listen := getBeacon(t, terminusClient, "0")
	walletCon := NewWalletConsumer(account0Listen)
	err = walletCon.DoConnect(account0Listen)

	time.Sleep(time.Second * 10)
}

// get new beacon for account
func getBeacon(t *testing.T, terminusClient terminus.TerminusClient, account string) beacon.Beacon {
	accountBeacon, err := terminusClient.Listen(account)
	if err != nil {
		t.Error(err)

	}

	// get the beacon
	acc, err := beacon.DecodeFromBech32Str(extractBeacon(accountBeacon))
	if err != nil {
		t.Error(err)
	}
	return acc
}

// extract a beacon from a terminus response
func extractBeacon(response string) string {
	return strings.Split(response[strings.Index(response, beacon.MoneysocketHrp):], " ")[0]
}
