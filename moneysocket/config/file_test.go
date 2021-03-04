package config

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/Flaque/filet"
	"github.com/brianvoe/gofakeit/v5"
	. "github.com/stretchr/testify/assert"
	nettest "github.com/xplorfin/netutils/testutils"
	tlsmock "github.com/xplorfin/tlsutils/mock"
)

var TestConfigs = []string{TerminusLndConf, TerminusClConf}

func TestConfigFromFile(t *testing.T) {
	for i, config := range TestConfigs {
		mockConfigFile, testConfig := MakeConfig(t, config)
		cfg, err := ParseConfigFromFile(mockConfigFile)
		if err != nil {
			panic(err)
		}
		// app config
		Equal(t, testConfig.AppAccountPersistDir, cfg.GetAccountPersistDir())
		// listen config
		Equal(t, testConfig.ListenBindHost, cfg.GetBindHost())
		Equal(t, testConfig.ListenBindPort, cfg.GetBindPort())
		Equal(t, testConfig.ListenUseTLS, cfg.GetUseTls())
		Equal(t, testConfig.ListenExternalHost, cfg.GetExternalHost())
		Equal(t, testConfig.ListenExternalPort, cfg.GetExternalPort())
		Equal(t, testConfig.ListenCertFile, cfg.GetCertFile())
		Equal(t, testConfig.ListenCertChainFile, cfg.GetCertChainFile())
		Equal(t, testConfig.ListenCertKey, cfg.GetKeyFile())
		Equal(t, testConfig.ListenSelfSignedCert, cfg.GetSelfSignedCert())
		if i != 1 { // second file does not have default bind or port
			Equal(t, testConfig.ListenDefaultBind, cfg.ListenConfig.defaultBind)
			Equal(t, testConfig.ListenDefaultPort, cfg.ListenConfig.defaultPort)
		}
		// rpc config
		Equal(t, testConfig.RpcBindHost, cfg.RpcConfig.BindHost)
		Equal(t, testConfig.RpcBindPort, cfg.RpcConfig.BindPort)
		Equal(t, testConfig.RpcExternalHost, cfg.RpcConfig.ExternalHost)
		Equal(t, testConfig.RpcExternalPort, cfg.RpcConfig.ExternalPort)

		err = cfg.Validate()
		if err != nil {
			t.Error(err)
		}
	}
}
func MakeConfig(t *testing.T, config string) (newConfigPath string, fileConfig TestFileConfig) {
	buf := new(bytes.Buffer)
	templ := template.Must(template.New("config").Parse(config))
	fileConfig = MakeMockConfig(t)
	err := templ.Execute(buf, fileConfig)
	if err != nil {
		panic(err)
	}
	tmpPath := filet.TmpFile(t, "", buf.String())
	return tmpPath.Name(), fileConfig
}

func MakeMockConfig(t *testing.T) TestFileConfig {
	lndCert, _ := tlsmock.TemporaryCert(t)
	chainFile, serverCertFile, serverKeyFile := tlsmock.TemporaryCertInChain(t)
	lndDir := filet.TmpDir(t, "")
	// TODO mock macaroon: https://git.io/JLj5Q
	macaroonPath := filet.TmpFile(t, lndDir, "")
	return TestFileConfig{
		AppAccountPersistDir: filet.TmpDir(t, ""),
		ListenBindHost:       gofakeit.DomainName(),
		ListenBindPort:       nettest.GetFreePort(t),
		ListenExternalHost:   gofakeit.DomainName(),
		ListenExternalPort:   nettest.GetFreePort(t),
		ListenUseTLS:         true,
		ListenCertFile:       serverCertFile,
		ListenCertKey:        serverKeyFile,
		ListenSelfSignedCert: true,
		ListenCertChainFile:  chainFile,
		ListenDefaultBind:    "127.0.0.1",
		ListenDefaultPort:    nettest.GetFreePort(t),
		RpcBindHost:          "127.0.0.1",
		RpcBindPort:          nettest.GetFreePort(t),
		RpcExternalHost:      gofakeit.DomainName(),
		RpcExternalPort:      nettest.GetFreePort(t),
		LndDir:               lndDir,
		LndMacaroonPath:      macaroonPath.Name(),
		LndTlsCertPath:       lndCert,
		LndNetwork:           "bitcoin",
		GrpcHost:             "127.0.0.1",
		GrpcPort:             nettest.GetFreePort(t),
	}
}

type TestFileConfig struct {
	AppAccountPersistDir string `goconf:"App:AccountPersistDir"`
	ListenBindHost       string `goconf:"Listen:BindHost"`
	ListenBindPort       int    `goconf:"Listen:BindPort"`
	ListenExternalHost   string `goconf:"Listen:ExternalHost"`
	ListenExternalPort   int    `goconf:"Listen:ExternalPort"`
	ListenUseTLS         bool   `goconf:"Listen:UseTLS"`
	ListenCertFile       string `goconf:"Listen:CertFile"`
	ListenCertKey        string `goconf:"Listen:CertKey"`
	ListenSelfSignedCert bool   `goconf:"Listen:SelfSignedCert"`
	ListenCertChainFile  string `goconf:"Listen:CertChainFile"`
	ListenDefaultBind    string `goconf:"Listen:DefaultBind"`
	ListenDefaultPort    int    `goconf:"Listen:DefaultPort"`

	RpcBindHost     string `goconf:"Rpc:BindHost"`
	RpcBindPort     int    `goconf:"Rpc:BindPort"`
	RpcExternalHost string `goconf:"Rpc:BindHost"`
	RpcExternalPort int    `goconf:"Rpc:ExternalPort"`

	LndDir          string `goconf:"LND:LndDir"`
	LndMacaroonPath string `goconf:"LND:LndMacaroonPath"`
	LndTlsCertPath  string `goconf:"LND:LndTlsCertPath"`
	LndNetwork      string `goconf:"LND:LndNetwork"`

	GrpcHost string `goconf:"LND:GrpcHost"`
	GrpcPort int    `goconf:"LND:GrpcPort"`
}

var TerminusClConf = `[App]

# account state is persisted in json format here
AccountPersistDir = {{.AppAccountPersistDir}}

[Listen]

# Default listening bind setting. 127.0.0.1 for localhost connections, 0.0.0.0
# for allowing connections from other hosts
BindHost = {{.ListenBindHost}}

# default port to listen for websocket connections port not specified.
BindPort = {{.ListenBindPort}}

# host for other devices to connect via the beacon
ExternalHost = {{.ListenExternalHost}}

# port for other devices to connect via the beacon
ExternalPort = {{.ListenExternalPort}}


# Use TLS for websocket connections
UseTLS = {{.ListenUseTLS}}

# if UseTLS is True, use this cert file
CertFile = {{.ListenCertFile}}

# if UseTLS is True, use this key file
CertKey = {{.ListenCertKey}}

# if UseTLS is True and we have a self-made cert for testing use this key file
# we don't need to provide a cert chain
SelfSignedCert = {{.ListenSelfSignedCert}}

# If we have a 'real' cert, we typically need to provide the cert chain file to
# make the browser clients happy.
CertChainFile = {{.ListenCertChainFile}}


[Rpc]

# host for client to connect
BindHost = {{.RpcBindHost}}

# port for client to connect
BindPort = {{.RpcBindPort}}

# host for client to connect
ExternalHost = {{.RpcExternalHost}}

# port for client to connect
ExternalPort = {{.RpcExternalPort}}
`

// https://github.com/moneysocket/terminus/blob/main/config/terminus-lnd.conf
var TerminusLndConf = `[App]

# account state is persisted in json format here
AccountPersistDir = {{.AppAccountPersistDir}}

[Listen]

# Default listening bind setting. 127.0.0.1 for localhost connections, 0.0.0.0
# for allowing connections from other hosts
BindHost = {{.ListenBindHost}}

# default port to listen for websocket connections port not specified.
BindPort = {{.ListenBindPort}}

# host for other devices to connect via the beacon
ExternalHost = {{.ListenExternalHost}}

# port for other devices to connect via the beacon
ExternalPort = {{.ListenExternalPort}}

# Default listening bind setting. 127.0.0.1 for localhost connections, 0.0.0.0
# for allowing connections from other hosts
DefaultBind = {{.ListenDefaultBind}}

# default port to listen for websocket connections port not specified.
DefaultPort = {{.ListenDefaultPort}}

# Use TLS for websocket connections
UseTLS = {{.ListenUseTLS}}

# if UseTLS is True, use this cert file
CertFile = {{.ListenCertFile}}

# if UseTLS is True, use this key file
CertKey = {{.ListenCertKey}}

# if UseTLS is True and we have a self-made cert for testing use this key file
# we don't need to provide a cert chain
SelfSignedCert = {{.ListenSelfSignedCert}}

# If we have a 'real' cert, we typically need to provide the cert chain file to
# make the browser clients happy.
CertChainFile = {{.ListenCertChainFile}}

[LND]

# LND settings directory
LndDir = {{.LndDir}}

# path to macaroon for grpc permissions
MacaroonPath = {{.LndMacaroonPath}}

# TLS cert for LND, which is different from the websocket listening TLS
TlsCertPath = {{.LndTlsCertPath}}

# LND network
Network = {{.LndNetwork}}

# GRPC connection
GrpcHost = {{.GrpcHost}}
GrpcPort = {{.GrpcPort}}

[Rpc]

# host for client to connect
BindHost = {{.RpcBindHost}}

# port for client to connect
BindPort = {{.RpcBindPort}}

# host for client to connect
ExternalHost = {{.RpcExternalHost}}

# port for client to connect
ExternalPort = {{.RpcExternalPort}}
`
