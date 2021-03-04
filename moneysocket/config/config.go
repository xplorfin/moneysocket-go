// we store config seperately from terminus to prevent circular dependency errors
package config

import (
	"fmt"
	"net/url"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	ozzo_validators "github.com/xplorfin/ozzo-validators"
)

// config for terminus
type Config struct {
	// directory to store accounts in
	AccountPersistDir string
	ListenConfig      ListenConfig
	RpcConfig         RpcConfig
}

// create a new terminus config
func NewConfig() *Config {
	return &Config{
		ListenConfig: ListenConfig{
			BindHost:     "localhost",
			BindPort:     5000,
			externalHost: "localhost",
			externalPort: 50001,
		},
		RpcConfig: RpcConfig{
			BindHost:     "localhost",
			BindPort:     5003,
			ExternalHost: "localhost",
			ExternalPort: 5004,
		},
	}
}

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

// port to listen for websocket connections
func (c *Config) GetBindPort() int {
	return c.ListenConfig.BindPort
}

// Default listening bind setting. 127.0.0.1 for localhost connections, 0.0.0.0
// for allowing connections from other hosts
func (c *Config) GetBindHost() string {
	return c.ListenConfig.BindHost
}

func (c *Config) GetAccountPersistDir() string {
	return c.AccountPersistDir
}

//  host for other devices to connect via the beacon
func (c *Config) GetExternalHost() string {
	return c.ListenConfig.externalHost
}

//  host for other devices to connect via the beacon
func (c *Config) GetExternalPort() int {
	return c.ListenConfig.externalPort
}

// wether or not to use tls
func (c *Config) GetUseTls() bool {
	return c.ListenConfig.useTLS
}

// ssl certificate file
func (c *Config) GetCertFile() string {
	return c.ListenConfig.certFile
}

// ssl certificate file
func (c *Config) GetKeyFile() string {
	return c.ListenConfig.certKey
}

// ssl certificate file
func (c *Config) GetSelfSignedCert() bool {
	return c.ListenConfig.selfSignedCert
}

// ssl certificate file
func (c *Config) GetCertChainFile() string {
	return c.ListenConfig.certChainFile
}

// get host:port
func (c *Config) GetHostName() string {
	return fmt.Sprintf("%s:%d", c.GetBindHost(), c.GetBindPort())
}

func (c *Config) GetRpcHostname() string {
	return fmt.Sprintf("%s:%d", c.RpcConfig.BindHost, c.RpcConfig.BindPort)

}

func (c *Config) GetRpcAddress() string {
	return fmt.Sprintf("http://%s", c.GetRpcHostname())
}

func (c *Config) GetAddress() string {
	return fmt.Sprintf("http://%s/", c.GetHostName())
}

func (c *Config) GetRelayUrl() string {
	u := url.URL{
		Scheme: "ws",
		Host:   fmt.Sprintf("%s:%d", c.GetBindHost(), c.GetBindPort()),
	}

	if c.GetUseTls() {
		u.Scheme = "wss"
	}

	return u.String()
}

func (c *Config) RpcServerTimeout() time.Duration {
	return time.Second * 10
}
