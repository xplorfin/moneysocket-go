package terminus

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/xplorfin/moneysocket-go/moneysocket/config"
)

const (
	// get info rpc method
	GetInfoMethod = "getinfo"
	// create account rpc method
	CreateAccountMethod = "create"
	// connect rpc method
	ConnectMethod = "connect"
	// listen rpc method
	ListenMethod = "listen"
)

type Client struct {
	config *config.Config // terminus client config
}

type RPCMessage struct {
	Method string     `json:"method"`
	Params [][]string `json:"params"`
}

func NewClient(config *config.Config) Client {
	return Client{
		config: config,
	}
}

// get list of beacons, accounts, etc
func (t Client) GetInfo() (res string, err error) {
	resp, err := t.ExecCmd(GetInfoMethod, []string{})
	return t.decodeJSONResponse(resp), err
}

// create account with a given number of sats
func (t Client) CreateAccount(msats int) (res string, err error) {
	resp, err := t.ExecCmd(CreateAccountMethod, []string{strconv.Itoa(msats)})
	return t.decodeJSONResponse(resp), err
}

func (t Client) Listen(accountName string) (res string, err error) {
	resp, err := t.ExecCmd(ListenMethod, []string{accountName})
	return t.decodeJSONResponse(resp), err
}

func (t Client) Connect(accountName, beacon string) (res string, err error) {
	resp, err := t.ExecCmd(ConnectMethod, []string{accountName, beacon})
	return t.decodeJSONResponse(resp), err
}

func (t Client) ExecCmd(method string, argv []string) (res []byte, err error) {
	message := RPCMessage{
		Method: method,
		Params: [][]string{argv},
	}

	msg, err := json.Marshal(&message)
	if err != nil {
		return res, err
	}

	req, err := http.NewRequest(http.MethodPost, t.config.GetRPCAddress(), bytes.NewBuffer(msg))
	if err != nil {
		return res, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	responseReader := new(bytes.Buffer)
	_, _ = responseReader.ReadFrom(resp.Body)
	return responseReader.Bytes(), nil
}

// decode terminus response
func (t Client) decodeJSONResponse(res []byte) string {
	jsonRes := []string{}
	err := json.Unmarshal(res, &jsonRes)
	if err != nil {
		log.Println(err)
	}
	if len(jsonRes) == 0 {
		return ""
	}
	return jsonRes[0]
}
