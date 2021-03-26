package config

import (
	"io/ioutil"

	"github.com/mvo5/goconfigparser"
)

// Section defines a section of the config
type Section struct {
	sectionName string
	config      *goconfigparser.ConfigParser
}

// NewSection create a new section object that reads from a a section of the config
func NewSection(name string, cfg *goconfigparser.ConfigParser) Section {
	return Section{
		sectionName: name,
		config:      cfg,
	}
}

// GetString gets a string from the config section
// handles python None types
func (s Section) GetString(option string) string {
	res, _ := s.config.Get(s.sectionName, option)
	if res == "None" {
		return ""
	}
	return res
}

// GetInt gets an int from the config section
// handles python None types
func (s Section) GetInt(option string) int {
	res, _ := s.config.Getint(s.sectionName, option)
	return res
}

// GetBool gets an bool from the config section
// handles python None types (returning false)
func (s Section) GetBool(option string) bool {
	res, _ := s.config.Getbool(s.sectionName, option)
	return res
}

// ParseConfig parses a Config from the file contents
// returns an error if parsing fails. Validation can be done on the Config object
func ParseConfig(fileContents string) (config Config, err error) {
	cfg := goconfigparser.New()
	err = cfg.ReadString(fileContents)
	if err != nil {
		return Config{}, err
	}

	// app config
	appConfig := NewSection("App", cfg)
	config.AccountPersistDir = appConfig.GetString("AccountPersistDir")

	// listen config
	listenConfig := NewSection("Listen", cfg)
	config.ListenConfig.BindHost = listenConfig.GetString("BindHost")
	config.ListenConfig.BindPort = listenConfig.GetInt("BindPort")
	config.ListenConfig.useTLS = listenConfig.GetBool("UseTLS")
	config.ListenConfig.ExternalHost = listenConfig.GetString("ExternalHost")
	config.ListenConfig.ExternalPort = listenConfig.GetInt("ExternalPort")
	config.ListenConfig.certFile = listenConfig.GetString("CertFile")
	config.ListenConfig.certChainFile = listenConfig.GetString("CertChainFile")
	config.ListenConfig.certKey = listenConfig.GetString("CertKey")
	config.ListenConfig.selfSignedCert = listenConfig.GetBool("SelfSignedCert")
	config.ListenConfig.defaultBind = listenConfig.GetString("DefaultBind")
	config.ListenConfig.defaultPort = listenConfig.GetInt("DefaultPort")

	rpcConfig := NewSection("Rpc", cfg)
	config.RPCConfig.BindHost = rpcConfig.GetString("BindHost")
	config.RPCConfig.BindPort = rpcConfig.GetInt("BindPort")
	config.RPCConfig.ExternalHost = rpcConfig.GetString("ExternalHost")
	config.RPCConfig.ExternalPort = rpcConfig.GetInt("ExternalPort")

	relayConfig := NewSection("Relay", cfg)
	config.RelayConfig.BindHost = relayConfig.GetString("ListenBind")
	config.RelayConfig.BindPort = relayConfig.GetInt("ListenPort")
	config.RelayConfig.useTLS = relayConfig.GetBool("UseTLS")
	config.RelayConfig.certFile = relayConfig.GetString("CertFile")
	config.RelayConfig.certKey = relayConfig.GetString("CertKey")
	config.RelayConfig.selfSignedCert = relayConfig.GetBool("SelfSignedCert")
	config.RelayConfig.certChainFile = relayConfig.GetString("CertChainFile")

	lndConfig := NewSection("LND", cfg)
	config.LndConfig.LndDir = lndConfig.GetString("LNDDir")
	config.LndConfig.MacaroonPath = lndConfig.GetString("MacaroonPath")
	config.LndConfig.TLSCertPath = lndConfig.GetString("TlsCertPath")
	config.LndConfig.Network = lndConfig.GetString("Network")
	config.LndConfig.GrpcHost = lndConfig.GetString("GRPCHost")
	config.LndConfig.GrpcPort = lndConfig.GetInt("GRPCPort")

	return config, err
}

// ParseConfigFromFile is a wrapper around ParseConfig that also reads from the file
// returns an error if parsing fails. Validation can be done on the Config object
func ParseConfigFromFile(filePath string) (conf Config, err error) {
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	return ParseConfig(string(contents))
}
