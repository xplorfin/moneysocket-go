package config

import (
	"fmt"
	"net/url"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	ozzo_validators "github.com/xplorfin/ozzo-validators"
)

// Config for terminus/relayd
type Config struct {
	// AccountPersistDir defines the directory to store accounts in
	AccountPersistDir string
	// ListenConfig defines the config to have generated beacons
	ListenConfig ListenConfig
	// RPCConfig defines the configuration for the rpc server to listen on
	RPCConfig RPCConfig
	// RelayConfig defines the configuration for the relay server (see: https://git.io/JmrYJ )
	RelayConfig RelayConfig
	// LndConfig defines the configuration for hooking up with the lnd server
	LndConfig LndConfig
}

// NewConfig creates a new config
func NewConfig() *Config {
	return &Config{
		ListenConfig: ListenConfig{
			BindHost:     "localhost",
			BindPort:     5000,
			ExternalHost: "localhost",
			ExternalPort: 50001,
		},
		RPCConfig: RPCConfig{
			BindHost:     "localhost",
			BindPort:     5003,
			ExternalHost: "localhost",
			ExternalPort: 5004,
		},
		RelayConfig: RelayConfig{
			BindHost: "localhost",
			BindPort: 5004,
		},
	}
}

// Validate validates the configuration
func (c Config) Validate() error {
	err := validation.ValidateStruct(&c,
		validation.Field(&c.AccountPersistDir, validation.Required, ozzo_validators.IsValidPath),
	)
	if err != nil {
		return err
	}

	err = c.ListenConfig.Validate()
	if err != nil {
		return err
	}

	err = c.RPCConfig.Validate()
	if err != nil {
		return err
	}
	return nil
}

// GetBindPort gets the port to listen for websocket connections
func (c *Config) GetBindPort() int {
	return c.ListenConfig.BindPort
}

// GetBindHost gets the listening bind setting
// for allowing connections from other hosts
func (c *Config) GetBindHost() string {
	return c.ListenConfig.BindHost
}

// GetAccountPersistDir gets the account persist dir
func (c *Config) GetAccountPersistDir() string {
	return c.AccountPersistDir
}

// GetExternalHost gets other devices to connect via the beacon
func (c *Config) GetExternalHost() string {
	return c.ListenConfig.ExternalHost
}

// GetExternalPort gets the externally binded port
// for other devices to connect via the beacon
func (c *Config) GetExternalPort() int {
	return c.ListenConfig.ExternalPort
}

// GetUseTLS determines wether or not to use tls
func (c *Config) GetUseTLS() bool {
	return c.ListenConfig.useTLS
}

// GetCertFile fetches the ssl certificate file
func (c *Config) GetCertFile() string {
	return c.ListenConfig.certFile
}

// GetKeyFile gets the ssl certificate file
func (c *Config) GetKeyFile() string {
	return c.ListenConfig.certKey
}

// GetSelfSignedCert fetches the ssl certificate file
func (c *Config) GetSelfSignedCert() bool {
	return c.ListenConfig.selfSignedCert
}

// GetCertChainFile fetches the ssl certificate chain file
func (c *Config) GetCertChainFile() string {
	return c.ListenConfig.certChainFile
}

// GetHostName fetches the hostname
func (c *Config) GetHostName() string {
	return fmt.Sprintf("%s:%d", c.GetBindHost(), c.GetBindPort())
}

// GetRPCHostname fetches the hostname of the rpc server
func (c *Config) GetRPCHostname() string {
	return fmt.Sprintf("%s:%d", c.RPCConfig.BindHost, c.RPCConfig.BindPort)

}

// GetRPCAddress fetches the address (w/ connection schema) of the rpc server
func (c *Config) GetRPCAddress() string {
	return fmt.Sprintf("http://%s", c.GetRPCHostname())
}

// GetAddress fetches the address of the rpc server
func (c *Config) GetAddress() string {
	return fmt.Sprintf("http://%s/", c.GetHostName())
}

// GetRelayURL gets the relay url
func (c *Config) GetRelayURL() string {
	u := url.URL{
		Scheme: "ws",
		Host:   fmt.Sprintf("%s:%d", c.RelayConfig.BindHost, c.RelayConfig.BindPort),
	}

	if c.GetUseTLS() {
		u.Scheme = "wss"
	}

	return u.String()
}

// RPCServerTimeout fetches the server timeout
func (c *Config) RPCServerTimeout() time.Duration {
	return time.Second * 10
}
