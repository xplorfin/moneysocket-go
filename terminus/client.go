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
	// GetInfoMethod is the rpc method for getting info
	GetInfoMethod = "getinfo"
	// CreateAccountMethod is the create account rpc method
	CreateAccountMethod = "create"
	// ConnectMethod is the connect rpc method
	ConnectMethod = "connect"
	// ListenMethod is the listen rpc method
	ListenMethod = "listen"
)

// RPCMessage is used for getting messages from the terminus client
type RPCMessage struct {
	// Method is the rpc method (e.g. connect, listen, etc)
	Method string `json:"method"`
	// Params are any parameters passed in RPCMessage
	Params [][]string `json:"params"`
}

// Client is the terminus client
type Client struct {
	// config is used for the terminus client
	config *config.Config
}

// NewClient creates a new Client for querying the rpc clint
func NewClient(config *config.Config) Client {
	return Client{
		config: config,
	}
}

// GetInfo gets a list of beacons, accounts, etc
func (t Client) GetInfo() (res string, err error) {
	resp, err := t.ExecCmd(GetInfoMethod, []string{})
	return t.decodeJSONResponse(resp), err
}

// CreateAccount creates account with a given number of sats
func (t Client) CreateAccount(msats int) (res string, err error) {
	resp, err := t.ExecCmd(CreateAccountMethod, []string{strconv.Itoa(msats)})
	return t.decodeJSONResponse(resp), err
}

// Listen listens on a given account
func (t Client) Listen(accountName string) (res string, err error) {
	resp, err := t.ExecCmd(ListenMethod, []string{accountName})
	return t.decodeJSONResponse(resp), err
}

// Connect connects an account to a beacon
func (t Client) Connect(accountName, beacon string) (res string, err error) {
	resp, err := t.ExecCmd(ConnectMethod, []string{accountName, beacon})
	return t.decodeJSONResponse(resp), err
}

// ExecCmd executes a command from the Client
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

// decodeJSONResponse decodes a terminus response for the client
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
