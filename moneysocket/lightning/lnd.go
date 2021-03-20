package lightning

import (
	"github.com/xplorfin/moneysocket-go/moneysocket/config"
	"google.golang.org/grpc"
)

// Lnd is the Lightning interface for lnd
type Lnd struct {
	// config contains the terminus config
	config *config.Config
	// paidRecvCb is the callback for a paid invoice receipt
	paidRecvCb PaidCallback
	// dialContext is the grpc connection for
	dialContext *grpc.ClientConn
}

// NewLnd generates a new lnd client from a given config
func NewLnd(config *config.Config) *Lnd {
	return &Lnd{
		config:     config,
		paidRecvCb: nil,
	}
}

// TODO
func getGrpcConnection(config *config.Config) (conn *grpc.ClientConn, err error) {
	//cert, err := config.RpcConfig.BindPort
	return nil, err
}

func (l Lnd) RegisterPaidRecvCb(callback PaidCallback) {
	panic("implement me")
}

func (l Lnd) GetInvoice(msatAmount int) {
	panic("implement me")
}

func (l Lnd) PayInvoice(bolt11 string) {
	panic("implement me")
}

var _ Lightning = &Lnd{}
