package message

import (
	"encoding/hex"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	mock "github.com/xplorfin/lndmock"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/notification"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/request"
)

func TestPythonDecode(t *testing.T) {
	testSeed, err := beacon.BytesToSharedSeed([]byte("hello from light"))
	if err != nil {
		t.Error(err)
	}
	const ToDecode = "181dbd096191cd7825fe8cc97991d9374cd89840b8a4fad9f5f89f0eb32b9ee9fc69e0e29485c5f4d011e34645f2f07e30d879cf3d1ca9d7971f2768ef4f70d79ddd22352cd687852128796bc80a47bd8ec5b40424a8a32161f264ef2d39a6b17f2dbd2274abd254200ddb8c97c23e6da4b8527ab54e1306b8c8d3d21886ae1f9136ecf25d2fb6998d3d3e9045868888eb148232eb4cebcbfa6c882608d19054d20859b3505565011e897f86f1a5b6dca25fdd317afc84c301e1de2d440be77de958801c29971e1206e285f359edee38a51863efdd678509130fb017a4fb3032"
	decoded, err := hex.DecodeString(ToDecode)
	if err != nil {
		t.Error(err)
	}
	_, _, err = WireDecode(decoded, &testSeed)
	if err != nil {
		t.Error(err)
	}
}

func EncodeDecode(t *testing.T, msg base.MoneysocketMessage, seed beacon.SharedSeed) {
	res, err := WireEncode(msg, &seed)
	if err != nil {
		t.Error(err)
	}
	// TODO test types individually
	_, _, err = WireDecode(res, &seed)
	if err != nil {
		t.Error(err)
	}
}

// test parity with https://git.io/JtT1f
func TestPythonParity(t *testing.T) {
	testSeed, err := beacon.BytesToSharedSeed([]byte("hello from light"))
	if err != nil {
		t.Error(err)
	}

	// requests
	pingRequest := request.NewPingRequest()
	EncodeDecode(t, pingRequest, testSeed)
	providerRequest := request.NewRequestProvider()
	EncodeDecode(t, providerRequest, testSeed)
	msatsRequest := request.NewRequestInvoice(1000)
	EncodeDecode(t, msatsRequest, testSeed)
	bolt, _ := mock.MockLndInvoiceMainnet(t)
	boltTest := request.NewRequestPay(bolt)
	EncodeDecode(t, boltTest, testSeed)

	rendezvousRequest := request.NewRendezvousRequest(string(testSeed.DeriveRendezvousID()))
	EncodeDecode(t, rendezvousRequest, testSeed)

	// notifications
	rendezvousEndNotify := notification.NewRendezvousEnd(string(testSeed.DeriveRendezvousID()), uuid.NewV4().String())
	EncodeDecode(t, rendezvousEndNotify, testSeed)

	bolt, _ = mock.MockLndInvoiceMainnet(t)
	notifyInvoice := notification.NewNotifyInvoice(bolt, uuid.NewV4().String())
	EncodeDecode(t, notifyInvoice, testSeed)

	notifyOpinionSeller := notification.NewNotifyOpinionSeller(uuid.NewV4().String(), []notification.Item{{
		ItemID: uuid.NewV4().String(),
		Name:   gofakeit.Word(),
		Msats:  50,
	}}, uuid.NewV4().String())
	EncodeDecode(t, notifyOpinionSeller, testSeed)

	notifyOpinionSellerNotReady := notification.NewNotifyOpinionSellerNotReady(uuid.NewV4().String())
	EncodeDecode(t, notifyOpinionSellerNotReady, testSeed)

	requestOpinionInvoice := request.NewRequestOpinionInvoice(uuid.NewV4().String(), uuid.NewV4().String())
	EncodeDecode(t, requestOpinionInvoice, testSeed)

	notifyOpinionInvoice := notification.NewNotifyOpinionInvoice(uuid.NewV4().String(), bolt)
	EncodeDecode(t, notifyOpinionInvoice, testSeed)

	notifyProviderNotReady := notification.NewNotifyProviderNotReady(uuid.NewV4().String())
	EncodeDecode(t, &notifyProviderNotReady, testSeed)

	notifyPong := notification.NewNotifyPong(uuid.NewV4().String())
	EncodeDecode(t, notifyPong, testSeed)

	notifyRendezvous := notification.NewNotifyRendezvous(uuid.NewV4().String(), uuid.NewV4().String())
	EncodeDecode(t, notifyRendezvous, testSeed)

	notifyRendezvousNotReady := notification.NewRendezvousNotReady(uuid.NewV4().String(), uuid.NewV4().String())
	EncodeDecode(t, notifyRendezvousNotReady, testSeed)
}
