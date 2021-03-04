package config

import (
	"io/ioutil"

	"github.com/mvo5/goconfigparser"
)

// parse a config from a file. Note: this is meant for python parity only
// and may be deprecated at some point. Since we validate in our setters
// and not at runtime, we have to parse into an intermediary struct managed
// by goconf

type Section struct {
	sectionName string
	config      *goconfigparser.ConfigParser
}

// create a new section object that reads from a a section of the config
func NewSection(name string, cfg *goconfigparser.ConfigParser) Section {
	return Section{
		sectionName: name,
		config:      cfg,
	}
}

// get a string from the config section
// handles python None types
func (s Section) GetString(option string) string {
	res, _ := s.config.Get(s.sectionName, option)
	if res == "None" {
		return ""
	}
	return res
}

func (s Section) GetInt(option string) int {
	res, _ := s.config.Getint(s.sectionName, option)
	return res
}

func (s Section) GetBool(option string) bool {
	res, _ := s.config.Getbool(s.sectionName, option)
	return res
}

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
	config.ListenConfig.externalHost = listenConfig.GetString("ExternalHost")
	config.ListenConfig.externalPort = listenConfig.GetInt("ExternalPort")
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

func ParseConfigFromFile(filePath string) (conf Config, err error) {
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	return ParseConfig(string(contents))
}
