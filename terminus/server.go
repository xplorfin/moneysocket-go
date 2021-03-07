package terminus

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/location"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/util"
	"github.com/xplorfin/moneysocket-go/moneysocket/config"
	"github.com/xplorfin/moneysocket-go/moneysocket/wad"
	"github.com/xplorfin/moneysocket-go/terminus/account"
	"github.com/xplorfin/netutils/testutils"
	"golang.org/x/sync/errgroup"
)

type Terminus struct {
	config    *config.Config
	directory *TerminusDirectory
	stack     *TerminusStack
}

func NewTerminus(config *config.Config) Terminus {
	return Terminus{
		config:    config,
		directory: NewTerminusDirectory(config),
		stack:     NewTerminusStack(config),
	}
}

// todo break this out into seperate methods
func (t *Terminus) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var rpcReq RpcMessage
	err := decoder.Decode(&rpcReq)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
	switch rpcReq.Method {
	// TODO actually decode these args
	case GetInfoMethod:
		w.Write([]byte(t.GetInfo()))
	case CreateAccountMethod:
		adb := t.Create(1000)
		w.Write([]byte(t.makeJsonResponse(fmt.Sprintf("created account %s, wad: %s", adb.Details.AccountName, adb.Details.Wad.FmtShort()))))
	case ConnectMethod:
		if len(rpcReq.Params) != 1 || len(rpcReq.Params[0]) != 2 {
			w.Write([]byte(t.makeJsonResponse("error, account not passed")))
			return
		}
		params := rpcReq.Params[0]
		decodedBeacon, err := beacon.DecodeFromBech32Str(params[1])
		if err != nil {
			w.Write([]byte(t.makeJsonResponse("error, beacon invalid")))
			return
		}

		acct := t.directory.LookupByName(params[0])
		if acct == nil {
			w.Write([]byte(t.makeJsonResponse(fmt.Sprintf("account %s not found", params[0]))))
			return
		}
		ss := decodedBeacon.GetSharedSeed()
		_, err = t.stack.Connect(decodedBeacon.Locations()[0].(location.WebsocketLocation), &ss)
		acct.AddConnectionAttempt(decodedBeacon, err)
		acct.Details.AddBeacon(decodedBeacon)
		t.directory.ReindexAccount(*acct)
		w.Write([]byte(t.makeJsonResponse(fmt.Sprintf("connected: %s to %s", params[0], params[1]))))
	case ListenMethod:
		if len(rpcReq.Params) != 1 || len(rpcReq.Params[0]) != 1 {
			w.Write([]byte(t.makeJsonResponse("error, account not passed")))
		}
		acct := rpcReq.Params[0][0]
		beaconServer, err := t.Listen(acct, "")
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
		w.Write([]byte(t.makeJsonResponse(fmt.Sprintf("listening to %s on %s", beaconServer, acct))))
	default:
		w.Write([]byte("method not yet implemented"))
	}
}

// start the server
func (t *Terminus) StartServer() {
	server := http.NewServeMux()
	server.Handle("/", t)
	http.ListenAndServe(t.config.GetRpcHostname(), server)
}

// make a raw json response from a string
func (t Terminus) makeJsonResponse(res string) string {
	jsonRes := []string{res}
	jsonResponse, err := json.Marshal(jsonRes)
	if err != nil {
		log.Println(err)
	}
	return string(jsonResponse)
}

// get formatted info (for python paritY)
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
	return t.makeJsonResponse(res)
}

// create an account with a given number of msats and add it to the directory
func (t *Terminus) Create(msats int) account.AccountDb {
	name := t.directory.GenerateAccountName()
	acct := account.NewAccountDb(name, t.config)
	acct.Details.Wad = wad.BitcoinWad(float64(msats))
	acct.Persist()
	t.directory.AddAccount(acct)
	return acct
}

// raw shared seed is optional
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

	bcn := beacon.NewBeaconFromSeed(sharedSeed)
	for _, loc := range t.stack.GetListenLocation() {
		bcn.AddLocation(loc)
	}
	acct.Details.AddSharedSeed(sharedSeed)

	t.stack.LocalConnect(sharedSeed)
	t.directory.ReindexAccount(*acct)
	return bcn.ToBech32Str(), nil
}

// load persisted accounts from disk
func (t *Terminus) LoadPersisted() {
	for _, adb := range account.GetPersistedAccounts(t.config) {
		t.directory.AddAccount(adb)
		for _, bcn := range adb.Details.Beacons {
			loc := bcn.Locations()[0]
			if loc.Type() != util.WebsocketLocationTlvType {
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

func (t *Terminus) Start(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		t.StartServer()
		return nil
	})
	g.Go(func() error {
		rpcStarted := testutils.WaitForConnectTimeout(t.config.GetRpcHostname(), t.config.RpcServerTimeout())
		if !rpcStarted {
			return fmt.Errorf("failed to detect rpc server at %s in %s", t.config.GetRpcHostname(), t.config.RpcServerTimeout().String())
		}
		t.LoadPersisted()
		t.stack.Listen()
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
	// TODO start prune_loop (and maybe connect_loop?)
}
