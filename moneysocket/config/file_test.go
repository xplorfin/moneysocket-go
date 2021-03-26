package config

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/brianvoe/gofakeit/v6"
	. "github.com/stretchr/testify/assert"
	"github.com/xplorfin/filet"
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
		Equal(t, testConfig.ListenUseTLS, cfg.GetUseTLS())
		Equal(t, testConfig.ListenExternalHost, cfg.GetExternalHost())
		Equal(t, testConfig.ListenExternalPort, cfg.GetExternalPort())
		Equal(t, testConfig.ListenCertFile, cfg.GetCertFile())
		Equal(t, testConfig.ListenCertChainFile, cfg.GetCertChainFile())
		Equal(t, testConfig.ListenCertKey, cfg.GetKeyFile())
		Equal(t, testConfig.ListenSelfSignedCert, cfg.GetSelfSignedCert())
		if i != 1 { // second file does not have default bind or port/lnd config
			Equal(t, testConfig.ListenDefaultBind, cfg.ListenConfig.defaultBind)
			Equal(t, testConfig.ListenDefaultPort, cfg.ListenConfig.defaultPort)

			// lnd config
			Equal(t, testConfig.LNDDir, cfg.LndConfig.LndDir)
			Equal(t, testConfig.LNDMacaroonPath, cfg.LndConfig.MacaroonPath)
			Equal(t, testConfig.LNDTlsCertPath, cfg.LndConfig.TLSCertPath)
			Equal(t, testConfig.LNDNetwork, cfg.LndConfig.Network)
			Equal(t, testConfig.GRPCHost, cfg.LndConfig.GrpcHost)
			Equal(t, testConfig.GRPCPort, cfg.LndConfig.GrpcPort)
		}
		// rpc config
		Equal(t, testConfig.RPCBindHost, cfg.RPCConfig.BindHost)
		Equal(t, testConfig.RPCBindPort, cfg.RPCConfig.BindPort)
		Equal(t, testConfig.RPCExternalHost, cfg.RPCConfig.ExternalHost)
		Equal(t, testConfig.RPCExternalPort, cfg.RPCConfig.ExternalPort)
		// relay config
		Equal(t, testConfig.RelayUseTLS, cfg.RelayConfig.useTLS)
		Equal(t, testConfig.RelayCertFile, cfg.RelayConfig.certFile)
		Equal(t, testConfig.RelayCertKey, cfg.RelayConfig.certKey)
		Equal(t, testConfig.RelaySelfSignedCert, cfg.RelayConfig.selfSignedCert)
		Equal(t, testConfig.RelayCertChainFile, cfg.RelayConfig.certChainFile)

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
		RPCBindHost:          "127.0.0.1",
		RPCBindPort:          nettest.GetFreePort(t),
		RPCExternalHost:      gofakeit.DomainName(),
		RPCExternalPort:      nettest.GetFreePort(t),
		LNDDir:               lndDir,
		LNDMacaroonPath:      macaroonPath.Name(),
		LNDTlsCertPath:       lndCert,
		LNDNetwork:           "bitcoin",
		GRPCHost:             "127.0.0.1",
		GRPCPort:             nettest.GetFreePort(t),
		RelayBindHost:        "127.0.0.1",
		RelayBindPort:        nettest.GetFreePort(t),
		RelayUseTLS:          true,
		RelayCertFile:        serverCertFile,
		RelayCertKey:         serverKeyFile,
		RelaySelfSignedCert:  true,
		RelayCertChainFile:   chainFile,
	}
}

type TestFileConfig struct {
	// App
	AppAccountPersistDir string `goconf:"App:AccountPersistDir"`
	// Listen
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
	// Rpc
	RPCBindHost     string `goconf:"Rpc:BindHost"`
	RPCBindPort     int    `goconf:"Rpc:BindPort"`
	RPCExternalHost string `goconf:"Rpc:BindHost"`
	RPCExternalPort int    `goconf:"Rpc:ExternalPort"`
	// LND
	LNDDir          string `goconf:"LND:LndDir"`
	LNDMacaroonPath string `goconf:"LND:LndMacaroonPath"`
	LNDTlsCertPath  string `goconf:"LND:LndTlsCertPath"`
	LNDNetwork      string `goconf:"LND:LndNetwork"`
	GRPCHost        string `goconf:"LND:GrpcHost"`
	GRPCPort        int    `goconf:"LND:GrpcPort"`
	// Relay
	RelayBindHost string `goconf:"Relay:ListenBind"`
	RelayBindPort int    `goconf:"Relay:ListenPort"`

	RelayUseTLS         bool   `goconf:"Relay:UseTLS"`
	RelayCertFile       string `goconf:"Relay:CertFile"`
	RelayCertKey        string `goconf:"Relay:CertKey"`
	RelaySelfSignedCert bool   `goconf:"Relay:SelfSignedCert"`
	RelayCertChainFile  string `goconf:"Relay:CertChainFile"`
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
BindHost = {{.RPCBindHost}}

# port for client to connect
BindPort = {{.RPCBindPort}}

# host for client to connect
ExternalHost = {{.RPCExternalHost}}

# port for client to connect
ExternalPort = {{.RPCExternalPort}}

[Relay]
ListenBind = {{.RelayBindHost}}

ListenPort = {{.RelayBindPort}}

# Use TLS for websocket connections
UseTLS = {{.RelayUseTLS}}

# if UseTLS is True, use this cert file
CertFile = {{.RelayCertFile}}

# if UseTLS is True, use this key file
CertKey = {{.RelayCertKey}}

# if UseTLS is True and we have a self-made cert for testing use this key file
# we don't need to provide a cert chain
SelfSignedCert = {{.RelaySelfSignedCert}}

# If we have a 'real' cert, we typically need to provide the cert chain file to
# make the browser clients happy.
CertChainFile = {{.RelayCertChainFile}}
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
LndDir = {{.LNDDir}}

# path to macaroon for grpc permissions
MacaroonPath = {{.LNDMacaroonPath}}

# TLS cert for LND, which is different from the websocket listening TLS
TlsCertPath = {{.LNDTlsCertPath}}

# LND network
Network = {{.LNDNetwork}}

# GRPC connection
GrpcHost = {{.GRPCHost}}
GrpcPort = {{.GRPCPort}}

[Rpc]

# host for client to connect
BindHost = {{.RPCBindHost}}

# port for client to connect
BindPort = {{.RPCBindPort}}

# host for client to connect
ExternalHost = {{.RPCExternalHost}}

# port for client to connect
ExternalPort = {{.RPCExternalPort}}

[Relay]
ListenBind = {{.RelayBindHost}}

ListenPort = {{.RelayBindPort}}

# Use TLS for websocket connections
UseTLS = {{.RelayUseTLS}}

# if UseTLS is True, use this cert file
CertFile = {{.RelayCertFile}}

# if UseTLS is True, use this key file
CertKey = {{.RelayCertKey}}

# if UseTLS is True and we have a self-made cert for testing use this key file
# we don't need to provide a cert chain
SelfSignedCert = {{.RelaySelfSignedCert}}

# If we have a 'real' cert, we typically need to provide the cert chain file to
# make the browser clients happy.
CertChainFile = {{.RelayCertChainFile}}
`
