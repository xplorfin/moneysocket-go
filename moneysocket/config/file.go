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
	config.RpcConfig.BindHost = rpcConfig.GetString("BindHost")
	config.RpcConfig.BindPort = rpcConfig.GetInt("BindPort")
	config.RpcConfig.ExternalHost = rpcConfig.GetString("ExternalHost")
	config.RpcConfig.ExternalPort = rpcConfig.GetInt("ExternalPort")

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
