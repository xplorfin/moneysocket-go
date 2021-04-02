package terminus

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"

	"github.com/xplorfin/moneysocket-go/moneysocket/lightning"

	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/location"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/util"
	"github.com/xplorfin/moneysocket-go/moneysocket/config"
	"github.com/xplorfin/moneysocket-go/moneysocket/wad"
	"github.com/xplorfin/moneysocket-go/terminus/account"
	"github.com/xplorfin/netutils/testutils"
	"golang.org/x/sync/errgroup"
)

// Terminus is the app
type Terminus struct {
	// config is the terminus app config
	config *config.Config
	// directory is the terminus directory
	directory *Directory
	// stack is the terminus stack
	stack *Stack
	// lightning is the lightning driver used to interact with the lnd node
	lightning *lightning.Lightning
}

// NewTerminus creates a Terminus node from a config
func NewTerminus(config *config.Config) (terminus Terminus, err error) {
	var lightningClient lightning.Lightning
	if config.LndConfig.HasLndConfig() {
		lightningClient, err = lightning.NewLnd(config)
		if err != nil {
			return terminus, err
		}
	}
	terminus = Terminus{
		config:    config,
		directory: NewTerminusDirectory(config),
		stack:     NewTerminusStack(config),
		lightning: &lightningClient,
	}

	terminus.stack.onAnnounce = terminus.OnAnnounce
	terminus.stack.onRevoke = terminus.OnRevoke

	return terminus, err
}

// OnAnnounce handles nexus announcements
func (t *Terminus) OnAnnounce(nexus nexus.Nexus) {
	// TODO register for messages and log errors if we get any not handled
	// by stack
}

// OnRevoke handles a nexus revoke attempt. Does nothing in terminus
func (t *Terminus) OnRevoke(nexus nexus.Nexus) {
	// do nothing
}

// ServeHTTP handles requests through the rpc server
// todo break this out into separate methods
func (t *Terminus) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var rpcReq RPCMessage
	err := decoder.Decode(&rpcReq)
	if err != nil {
		w.WriteHeader(500)
		_, _ = w.Write([]byte(err.Error()))
	}
	switch rpcReq.Method {
	// TODO actually decode these args
	case GetInfoMethod:
		_, _ = w.Write([]byte(t.GetInfo()))
	case CreateAccountMethod:
		adb := t.Create(1000)
		_, _ = w.Write([]byte(t.makeJSONResponse(fmt.Sprintf("created account %s, wad: %s", adb.Details.AccountName, adb.Details.Wad.FmtShort()))))
	case ConnectMethod:
		if len(rpcReq.Params) != 1 || len(rpcReq.Params[0]) != 2 {
			_, _ = w.Write([]byte(t.makeJSONResponse("error, account not passed")))
			return
		}
		params := rpcReq.Params[0]
		decodedBeacon, err := beacon.DecodeFromBech32Str(params[1])
		if err != nil {
			_, _ = w.Write([]byte(t.makeJSONResponse("error, beacon invalid")))
			return
		}

		acct := t.directory.LookupByName(params[0])
		if acct == nil {
			_, _ = w.Write([]byte(t.makeJSONResponse(fmt.Sprintf("account %s not found", params[0]))))
			return
		}
		ss := decodedBeacon.GetSharedSeed()
		_, err = t.stack.Connect(decodedBeacon.Locations()[0].(location.WebsocketLocation), &ss)
		acct.AddConnectionAttempt(decodedBeacon, err)
		acct.Details.AddBeacon(decodedBeacon)
		t.directory.ReindexAccount(*acct)
		_, _ = w.Write([]byte(t.makeJSONResponse(fmt.Sprintf("connected: %s to %s", params[0], params[1]))))
	case ListenMethod:
		if len(rpcReq.Params) != 1 || len(rpcReq.Params[0]) != 1 {
			_, _ = w.Write([]byte(t.makeJSONResponse("error, account not passed")))
		}
		acct := rpcReq.Params[0][0]
		beaconServer, err := t.Listen(acct, "")
		if err != nil {
			w.WriteHeader(500)
			_, _ = w.Write([]byte(err.Error()))
		}
		_, _ = w.Write([]byte(t.makeJSONResponse(fmt.Sprintf("listening to %s on %s", beaconServer, acct))))
	default:
		_, _ = w.Write([]byte("method not yet implemented"))
	}
}

// StartServer starts the terminus server
func (t *Terminus) StartServer() error {
	server := http.NewServeMux()
	server.Handle("/", t)
	return http.ListenAndServe(t.config.GetRPCHostname(), server)
}

// makeJSONResponse makes a raw json response from a string
func (t Terminus) makeJSONResponse(res string) string {
	jsonRes := []string{res}
	jsonResponse, err := json.Marshal(jsonRes)
	if err != nil {
		log.Println(err)
	}
	return string(jsonResponse)
}

// GetInfo gets formatted info (for python parity)
func (t *Terminus) GetInfo() (res string) {
	locations := t.stack.GetListenLocation()
	accounts := t.directory.GetAccountList()
	res += "ACCOUNTS:"
	if len(accounts) == 0 {
		res += "\nnone"
	}
	for _, account := range accounts {
		res += fmt.Sprintf("\n%s", account.GetSummaryString(locations))
	}
	return t.makeJSONResponse(res)
}

// Create creates an account with a given number of msats and add it to the directory
func (t *Terminus) Create(msats int) account.DB {
	name := t.directory.GenerateAccountName()
	acct := account.NewAccountDb(name, t.config)
	acct.Details.Wad = wad.BitcoinWad(float64(msats))
	_ = acct.Persist()
	t.directory.AddAccount(acct)
	return acct
}

// RetryConnectionLoop retries a connection on a loop
func (t *Terminus) RetryConnectionLoop() {
	for {
		for _, acct := range t.directory.GetAccountList() {
			disconnectedBeacons := acct.GetDisconnectedBeacons()
			for _, disconnectedBeacon := range disconnectedBeacons {
				loc := disconnectedBeacon.Locations()[0]
				if loc.Type() != util.WebsocketLocationTLVType {
					panic(fmt.Sprintf("location type %d is not supported", util.WebsocketLocationTLVType))
				}
				ss := disconnectedBeacon.GetSharedSeed()
				_, err := t.stack.Connect(loc.(location.WebsocketLocation), &ss)
				acct.AddConnectionAttempt(disconnectedBeacon, err)
			}
		}
		time.Sleep(time.Second * 3)
	}
}

// Listen listens on a port. Raw shared seed is optional
func (t *Terminus) Listen(rawAcct string, rawSharedSeed string) (encodedBeacon string, err error) {
	acct := t.directory.LookupByName(rawAcct)
	if acct == nil {
		return encodedBeacon, fmt.Errorf("could not find account of name: %s", rawAcct)
	}

	var sharedSeed beacon.SharedSeed
	if rawSharedSeed != "" {
		sharedSeed, err = beacon.HexToSharedSeed(rawSharedSeed)
		if err != nil {
			return encodedBeacon, err
		}
	} else {
		sharedSeed = beacon.NewSharedSeed()
	}

	bcn := beacon.NewBeaconFromSharedSeed(sharedSeed)
	for _, loc := range t.stack.GetListenLocation() {
		bcn.AddLocation(loc)
	}
	acct.Details.AddSharedSeed(sharedSeed)

	t.stack.LocalConnect(sharedSeed)
	t.directory.ReindexAccount(*acct)
	return bcn.ToBech32Str(), nil
}

// LoadPersisted loads persisted accounts from the disk
func (t *Terminus) LoadPersisted() {
	for _, adb := range account.GetPersistedAccounts(t.config) {
		t.directory.AddAccount(adb)
		for _, bcn := range adb.Details.Beacons {
			loc := bcn.Locations()[0]
			if loc.Type() != util.WebsocketLocationTLVType {
				panic("non-websocket loc types not yet supported")
			}
			sharedSeed := bcn.GetSharedSeed()
			// TODO
			conn, err := t.stack.Connect(loc.(location.WebsocketLocation), &sharedSeed)
			_ = conn
			if err != nil {
				adb.AddConnectionAttempt(bcn, err)
			}
		}
		for _, sharedSeed := range adb.Details.SharedSeeds {
			t.stack.localLayer.Connect(sharedSeed)
		}
	}
}

// Start starts the server
func (t *Terminus) Start(ctx context.Context) error {
	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		return t.StartServer()
	})
	g.Go(func() error {
		rpcStarted := testutils.WaitForConnectTimeout(t.config.GetRPCHostname(), t.config.RPCServerTimeout())
		if !rpcStarted {
			return fmt.Errorf("failed to detect rpc server at %s in %s", t.config.GetRPCHostname(), t.config.RPCServerTimeout().String())
		}
		t.LoadPersisted()
		return t.stack.Listen()
	})

	g.Go(func() error {
		t.RetryConnectionLoop()
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	// TODO start prune_loop (and maybe connect_loop?)
	return nil
}
