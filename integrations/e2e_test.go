package integrations

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	. "github.com/stretchr/testify/assert"
	"github.com/xplorfin/filet"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/config"
	"github.com/xplorfin/moneysocket-go/relay"
	"github.com/xplorfin/moneysocket-go/terminus"
	nettest "github.com/xplorfin/netutils/testutils"
)

// makeConfig creates a mock config for e2e tests.
func makeConfig(t *testing.T) *config.Config {
	testConfig := config.NewConfig()
	testConfig.AccountPersistDir = filet.TmpDir(t, "")
	testConfig.ListenConfig.BindPort = nettest.GetFreePort(t)
	testConfig.ListenConfig.BindHost = "0.0.0.0"
	testConfig.ListenConfig.ExternalHost = "127.0.0.1"
	testConfig.ListenConfig.ExternalPort = testConfig.GetBindPort()

	testConfig.RPCConfig.BindHost = "localhost"
	testConfig.RPCConfig.BindPort = nettest.GetFreePort(t)

	testConfig.RelayConfig.BindHost = "localhost"
	testConfig.RelayConfig.BindPort = nettest.GetFreePort(t)

	return testConfig
}

// TestE2E attempts to run an end-to-end test of the moneysockte opinion app.
func TestE2E(t *testing.T) {
	cfg := makeConfig(t)
	ctx := context.Background()

	// setup test relay
	testRelay := relay.NewRelay(cfg)
	go func() {
		err := testRelay.RunApp()
		Nil(t, err)
	}()

	// setup test rpc server
	testRPCServer, err := terminus.NewTerminus(cfg)
	Nil(t, err)
	go func() {
		err = testRPCServer.Start(ctx)
		Nil(t, err)
	}()

	// test rpc server hostname
	nettest.AssertConnected(cfg.GetRPCHostname(), t)

	terminusClient := terminus.NewClient(cfg)
	// create two accounts
	// -- account 1
	account1Beacon, err := terminusClient.CreateAccount(1000000)
	Nil(t, err)

	t.Log(account1Beacon)
	// -- acount 2
	account1Beacon, err = terminusClient.CreateAccount(1000000)
	Nil(t, err)

	t.Log(account1Beacon)

	// start the seller wallet consumer on account-1
	account1Listen := getBeacon(t, terminusClient, "1")
	app := NewSellerApp(account1Listen)
	err = app.ConsumerStack.DoConnect(account1Listen)
	Nil(t, err)

	// start the wallet consumer on account-0
	account0Listen := getBeacon(t, terminusClient, "0")
	walletCon := NewWalletConsumer(account0Listen)
	err = walletCon.DoConnect(account0Listen)
	Nil(t, err)

	// generate a new beacon and call connect
	providerBeacon := generateNewBeacon(cfg.GetExternalHost(), cfg.GetUseTLS(), cfg.GetExternalPort())
	_, err = terminusClient.Connect("0", providerBeacon.ToBech32Str())
	Nil(t, err)

	// check if incoming socket is there
	fmt.Print(terminusClient.GetInfo())
	time.Sleep(time.Second * 20)
}

// getBeacon mocks a new beacon for a given account.
func getBeacon(t *testing.T, terminusClient terminus.Client, account string) beacon.Beacon {
	accountBeacon, err := terminusClient.Listen(account)
	Nil(t, err)

	// get the beacon
	acc, err := beacon.DecodeFromBech32Str(extractBeacon(accountBeacon))
	Nil(t, err)

	return acc
}

// extractBeacon extracts a beacon from a terminus response.
func extractBeacon(response string) string {
	return strings.Split(response[strings.Index(response, beacon.MoneysocketHrp):], " ")[0]
}
