package lightning

import (
	"crypto/tls"
	"testing"

	. "github.com/stretchr/testify/assert"
	"github.com/xplorfin/filet"
	mock "github.com/xplorfin/lndmock"
	terminusConfig "github.com/xplorfin/moneysocket-go/moneysocket/config"
	"gopkg.in/macaroon.v2"
)

// LndTestNode is a test node in lnd.
type LndTestNode struct {
	mock.LndContainer
	// t is a pointer to the tessting object
	t *testing.T
	// Mocker is a pointer to the lightning Mocker (this should be initializeD)
	Mocker *mock.LightningMocker
	// name is the container name
	name string
	// address is an address generated by the lnd node
	address string
	// pubkey of the node
	pubkey string
	// admin macaroon
	mac *macaroon.Macaroon
	// macFile: path to the macaroon file for this node
	macFile string
	// tls details
	tls *tls.Config
	// tlsFile path
	tlsFile string
}

// NewLndTestNode generates a new lnd test node
// note: the btcd container should be created and the Mocker should be initialized before this step
func NewLndTestNode(t *testing.T, mocker *mock.LightningMocker, name string) LndTestNode {
	var err error
	var tlsCert string
	node := LndTestNode{
		// LndContainer this will be replaced later in the function and is not usable
		LndContainer: mock.LndContainer{},
		t:            t,
		name:         name,
	}
	// start alice's lnd instance
	node.LndContainer, err = mocker.CreateLndContainer(name)
	Nil(t, err)

	// get address
	node.address, err = node.Address()
	Nil(t, err)

	node.pubkey, err = node.GetPubKey()
	Nil(t, err)

	// get macaroon
	node.mac, err = node.GetAdminMacaroon()
	Nil(t, err)

	rawMac, err := node.mac.MarshalBinary()
	Nil(t, err)
	node.macFile = filet.TmpBinFile(t, "", rawMac).Name()

	// get alices tls cert
	node.tls, tlsCert, err = node.GetTLSCert()
	Nil(t, err)
	node.tlsFile = filet.TmpFile(t, "", tlsCert).Name()
	return node
}

func (l LndTestNode) LndConfig() terminusConfig.LndConfig {
	return terminusConfig.LndConfig{
		LndDir:       filet.TmpDir(l.t, ""),
		MacaroonPath: l.macFile,
		TLSCertPath:  l.tlsFile,
		Network:      "bitcoin",
		GrpcHost:     "localhost",
		GrpcPort:     l.PortMap.GetHostPort(10009),
	}
}

func TestLnd(t *testing.T) {
	mocker := mock.NewLightningMocker()
	defer func() {
		Nil(t, mocker.Teardown())
	}()

	err := mocker.Initialize()
	Nil(t, err)

	// start btcd as a prereq to lnd
	btcdContainer, err := mocker.CreateBtcdContainer()
	Nil(t, err)

	// create an alice node
	alice := NewLndTestNode(t, &mocker, "alice")

	err = btcdContainer.MineToAddress(alice.address, 600)
	Nil(t, err)

	// create an bob node
	bob := NewLndTestNode(t, &mocker, "bob")

	err = btcdContainer.MineToAddress(bob.address, 600)
	Nil(t, err)

	// segwit activates at block 400
	// wait for sync
	err = alice.WaitForSync(true, false)
	Nil(t, err)

	err = bob.WaitForSync(true, false)
	Nil(t, err)

	// open bob->alice channel
	err = bob.OpenChannel(alice.pubkey, alice.name, 1000000)
	Nil(t, err)

	// mine blocks to process channels
	err = btcdContainer.Mine(4)
	Nil(t, err)

	// wait for sync
	err = alice.WaitForCondition(func(res *lnrpc.GetInfoResponse) bool {
		return res.NumActiveChannels == 1
	})
	Nil(t, err)

	// open alice->bob channel
	err = alice.OpenChannel(bob.pubkey, bob.name, 1000000)
	Nil(t, err)

	// mine blocks to process channels
	err = btcdContainer.Mine(4)
	Nil(t, err)

	// wait for sync
	err = alice.WaitForCondition(func(res *lnrpc.GetInfoResponse) bool {
		return res.NumActiveChannels == 2
	})
	Nil(t, err)

	// create alice lnd node
	aliceLnd, err := NewLnd(&terminusConfig.Config{LndConfig: alice.LndConfig()})
	Nil(t, err)

	// get an invoice from alice
	payRequest, err := aliceLnd.GetInvoice(1000)
	Nil(t, err)

	// create bob lnd node
	bobLnd, err := NewLnd(&terminusConfig.Config{LndConfig: bob.LndConfig()})
	Nil(t, err)

	_, _, err = bobLnd.PayInvoice(payRequest)
	Nil(t, err)
}
