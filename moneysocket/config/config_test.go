package config

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	. "github.com/stretchr/testify/assert"
	"github.com/xplorfin/filet"
	nettest "github.com/xplorfin/netutils/testutils"
	tlsmock "github.com/xplorfin/tlsutils/mock"
)

// TestGetters tests the getter methods for the config files
func TestGetters(t *testing.T) {
	chainFile, serverCertFile, serverKeyFile := tlsmock.TemporaryCertInChain(t)
	accountDir := filet.TmpDir(t, "")
	config := Config{
		AccountPersistDir: accountDir,
		ListenConfig: ListenConfig{
			BindHost:       "127.0.0.1",
			BindPort:       nettest.GetFreePort(t),
			ExternalHost:   gofakeit.Word(),
			ExternalPort:   nettest.GetFreePort(t),
			useTLS:         true,
			certFile:       serverCertFile,
			certKey:        serverKeyFile,
			selfSignedCert: true,
			certChainFile:  chainFile,
		},
		RpcConfig: RpcConfig{
			BindHost:     "127.0.0.1",
			BindPort:     nettest.GetFreePort(t),
			ExternalHost: "127.0.0.1",
			ExternalPort: nettest.GetFreePort(t),
		},
	}

	// just make sure this is valid
	if config.Validate() != nil {
		t.Errorf("expected config to be valid for testing getters, found errors %s", config.Validate())
	}

	Equal(t, config.GetAccountPersistDir(), config.AccountPersistDir)
	Equal(t, config.GetBindHost(), config.ListenConfig.BindHost)
	Equal(t, config.GetBindPort(), config.ListenConfig.BindPort)
	Equal(t, config.GetExternalHost(), config.ListenConfig.ExternalHost)
	Equal(t, config.GetExternalPort(), config.ListenConfig.ExternalPort)
	Equal(t, config.GetUseTls(), config.ListenConfig.useTLS)
	Equal(t, config.GetCertFile(), config.ListenConfig.certFile)
	Equal(t, config.GetKeyFile(), config.ListenConfig.certKey)
	Equal(t, config.GetSelfSignedCert(), config.ListenConfig.selfSignedCert)
	Equal(t, config.GetCertChainFile(), config.ListenConfig.certChainFile)
}
