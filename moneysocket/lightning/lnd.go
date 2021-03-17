package lightning

type Lnd struct {
	paidRecvCb PaidCallback
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