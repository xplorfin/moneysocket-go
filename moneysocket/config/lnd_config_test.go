package config

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jinzhu/copier"
	"github.com/lightningnetwork/lnd/lnrpc"
	. "github.com/stretchr/testify/assert"
	"github.com/xplorfin/filet"
	mock "github.com/xplorfin/lndmock"
)

// Copy copies the LndConfig for testing
func (l LndConfig) Copy(t *testing.T) LndConfig {
	var newConfig LndConfig
	err := copier.Copy(&newConfig, &l)
	Nil(t, err)
	return newConfig
}

func testConfigValidation(validConfig LndConfig, t *testing.T) {
	err := validConfig.Validate()
	Nil(t, err)

	// invalid lnd dir
	newConfig := validConfig.Copy(t)
	newConfig.LndDir = fmt.Sprintf("%s/%s", filet.TmpDir(t, ""), gofakeit.Word())
	NotNil(t, newConfig.Validate())

	// invalid tls path
	newConfig = validConfig.Copy(t)
	newConfig.TLSCertPath = fmt.Sprintf("%s/%s", filet.TmpDir(t, ""), gofakeit.Word())
	NotNil(t, newConfig.Validate())

	// invalid network
	newConfig = validConfig.Copy(t)
	newConfig.Network = fmt.Sprintf("not-%s", newConfig.Network)
	NotNil(t, newConfig.Validate())
}

func TestLndConfig(t *testing.T) {
	mocker := mock.NewLightningMocker()
	defer func() {
		Nil(t, mocker.Teardown())
	}()

	err := mocker.Initialize()
	Nil(t, err)

	// start btcd as a prereq to lnd
	_, err = mocker.CreateBtcdContainer()
	Nil(t, err)

	// start alice's lnd instance
	aliceContainer, err := mocker.CreateLndContainer("alice")
	Nil(t, err)

	lndDir := filet.TmpDir(t, "")

	// get alices macaroon
	aliceMac, err := aliceContainer.GetAdminMacaroon()
	Nil(t, err)
	rawAliceMac, err := aliceMac.MarshalBinary()
	Nil(t, err)

	macaroonFile := filet.TmpBinFile(t, lndDir, rawAliceMac)
	Nil(t, err)

	// get alices tls cert
	_, aliceTls, err := aliceContainer.GetTLSCert()
	Nil(t, err)
	tlsFile := filet.TmpFile(t, lndDir, aliceTls)

	config := LndConfig{
		LndDir:       lndDir,
		MacaroonPath: macaroonFile.Name(),
		TLSCertPath:  tlsFile.Name(),
		Network:      "bitcoin",
		GrpcHost:     "localhost",
		GrpcPort:     aliceContainer.PortMap.GetHostPort(10009),
	}

	testConfigValidation(config, t)
	Nil(t, err)

	lnclient, err := config.RPCClient(context.Background())
	Nil(t, err)

	// wait for lnd boot. TODO add a getter to make sure lnd is unlocked
	time.Sleep(time.Second * 10)

	req := lnrpc.GetInfoRequest{}
	_, err = lnclient.GetInfo(context.Background(), &req)
	Nil(t, err)
}