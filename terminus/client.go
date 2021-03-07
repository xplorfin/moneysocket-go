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

type TerminusClient struct {
	config *config.Config // terminus client config
}

type RpcMessage struct {
	Method string     `json:"method"`
	Params [][]string `json:"params"`
}

func NewClient(config *config.Config) TerminusClient {
	return TerminusClient{
		config: config,
	}
}

// get list of beacons, accounts, etc
func (t TerminusClient) GetInfo() (res string, err error) {
	resp, err := t.ExecCmd(GetInfoMethod, []string{})
	return t.decodeJsonResponse(resp), err
}

// create account with a given number of sats
func (t TerminusClient) CreateAccount(msats int) (res string, err error) {
	resp, err := t.ExecCmd(CreateAccountMethod, []string{strconv.Itoa(msats)})
	return t.decodeJsonResponse(resp), err
}

func (t TerminusClient) Listen(accountName string) (res string, err error) {
	resp, err := t.ExecCmd(ListenMethod, []string{accountName})
	return t.decodeJsonResponse(resp), err
}

func (t TerminusClient) ExecCmd(method string, argv []string) (res []byte, err error) {
	message := RpcMessage{
		Method: method,
		Params: [][]string{argv},
	}

	msg, err := json.Marshal(&message)
	if err != nil {
		return res, err
	}

	req, err := http.NewRequest(http.MethodPost, t.config.GetRpcAddress(), bytes.NewBuffer(msg))
	if err != nil {
		return res, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	responseReader := new(bytes.Buffer)
	responseReader.ReadFrom(resp.Body)
	return responseReader.Bytes(), nil
}

// decode terminus response
func (t TerminusClient) decodeJsonResponse(res []byte) string {
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
