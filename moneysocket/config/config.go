// we store config seperately from terminus to prevent circular dependency errors
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
	// RpcConfig defines the configuration for the rpc server to listen on
	RpcConfig RpcConfig
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
		RpcConfig: RpcConfig{
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

	err = c.RpcConfig.Validate()
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

// GetUseTls determines wether or not to use tls
func (c *Config) GetUseTls() bool {
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

// GetRpcHostname fetches the hostname of the rpc server
func (c *Config) GetRpcHostname() string {
	return fmt.Sprintf("%s:%d", c.RpcConfig.BindHost, c.RpcConfig.BindPort)

}

// GetRpcAddress fetches the address (w/ connection schema) of the rpc server
func (c *Config) GetRpcAddress() string {
	return fmt.Sprintf("http://%s", c.GetRpcHostname())
}

// GetAddress fetches the address of the rpc server
func (c *Config) GetAddress() string {
	return fmt.Sprintf("http://%s/", c.GetHostName())
}

// GetRelayUrl gets the relay url
func (c *Config) GetRelayUrl() string {
	u := url.URL{
		Scheme: "ws",
		Host:   fmt.Sprintf("%s:%d", c.RelayConfig.BindHost, c.RelayConfig.BindPort),
	}

	if c.GetUseTls() {
		u.Scheme = "wss"
	}

	return u.String()
}

// RpcServerTimeout fetches the server timeout
func (c *Config) RpcServerTimeout() time.Duration {
	return time.Second * 10
}
