package terminus

import (
	"fmt"
	"testing"

	"github.com/xplorfin/moneysocket-go/moneysocket/config"
)

func TestClientListen(t *testing.T) {
	t.Skip("for local testing only")
	configuration := config.NewConfig()
	configuration.RPCConfig.BindHost = "127.0.0.1"
	configuration.RPCConfig.BindPort = 11054
	client := NewClient(configuration)

	res, err := client.Listen("account-0")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)
}

func TestClientCreate(t *testing.T) {
	t.Skip("for local testing only")
	configuration := config.NewConfig()
	configuration.RPCConfig.BindHost = "127.0.0.1"
	configuration.RPCConfig.BindPort = 11054
	client := NewClient(configuration)

	res, err := client.CreateAccount(1000)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)
}

func TestClientInfo(t *testing.T) {
	t.Skip("for local testing only")
	configuration := config.NewConfig()
	configuration.RPCConfig.BindHost = "127.0.0.1"
	configuration.RPCConfig.BindPort = 11054
	client := NewClient(configuration)

	res, err := client.GetInfo()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)
}
