package lightning

import (
	"context"
	"log"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/xplorfin/moneysocket-go/moneysocket/config"
)

// Lnd is the Lightning interface for lnd
type Lnd struct {
	// config contains the terminus config
	config *config.Config
	// paidRecvCb is the callback for a paid invoice receipt
	paidRecvCb PaidCallback
	// client is the lnd client
	client lnrpc.LightningClient
	// pendingPaymentHashes contains a list of pending payments
	pendingPaymentHashes [][]byte
}

// RecvPaid receives an invoice payment
func (l *Lnd) RecvPaid(preimage string, msats int) {
	l.paidRecvCb(preimage, msats)
}

// NewLnd generates a new lnd client from a given config
// TODO use streaming for invoices
func NewLnd(config *config.Config) (lnd *Lnd, err error) {
	client, err := config.LndConfig.RPCClient(context.Background())
	if err != nil {
		return lnd, err
	}

	return &Lnd{
		config:     config,
		paidRecvCb: nil,
		client:     client,
	}, nil
}

// RegisterPaidRecvCb registers the callback
func (l *Lnd) RegisterPaidRecvCb(callback PaidCallback) {
	l.paidRecvCb = callback
}

// GetInvoice gets an invoice (paymentRequest) fora given amount msatAmount
func (l *Lnd) GetInvoice(msatAmount int) (paymentRequest string, err error) {
	log.Printf("getting invoice %d msats", msatAmount)
	satAmount := int64(msatAmount / 1000.0)

	invoice := &lnrpc.Invoice{
		Memo:      "",
		Value:     satAmount,
		ValueMsat: satAmount,
		Expiry:    3600,
		Private:   false,
	}

	resp, err := l.client.AddInvoice(context.Background(), invoice)
	if err != nil {
		return paymentRequest, err
	}
	l.pendingPaymentHashes = append(l.pendingPaymentHashes, resp.RHash)
	return resp.PaymentRequest, err
}

// PayInvoice sends a payment for a bolt-11 invoice
func (l *Lnd) PayInvoice(bolt11 string) (preimage []byte, msatAmount int, err error) {
	log.Printf("paying invoice %s", bolt11)
	// TODO figure out deprecation status hre
	payReq := lnrpc.SendRequest{
		PaymentRequest: bolt11,
	}
	resp, err := l.client.SendPaymentSync(context.Background(), &payReq)
	if err != nil {
		return preimage, msatAmount, err
	}

	log.Println(resp)
	log.Printf("paid %s", payReq.PaymentRequest)
	log.Printf("route %s", resp.PaymentRoute)
	return resp.PaymentPreimage, int(resp.PaymentRoute.TotalAmtMsat), nil

}

var _ Lightning = &Lnd{}
