package config

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	ozzo_validators "github.com/xplorfin/ozzo-validators"
)

// RPCConfig defines the config for the rpc server.
type RPCConfig struct {
	// BindHost for client to connect
	BindHost string
	// BindPort for client to connect
	BindPort int
	// ExternalHost for a client to connect
	ExternalHost string
	// ExternalPort for client to connect
	ExternalPort int
}

// Validate  validates the RPCConfig.
func (r RPCConfig) Validate() error {
	return validation.ValidateStruct(&r,
		// bind host cannot be null
		validation.Field(&r.BindHost, validation.Required, is.Host),
		// bind port cannot be null
		validation.Field(&r.BindPort, validation.Required, ozzo_validators.IsValidPort),
		// required when port is not null
		validation.Field(&r.ExternalHost, validation.When(r.ExternalPort != 0, validation.Required), is.Host),
		// required when host is not null
		validation.Field(&r.ExternalPort, validation.When(r.ExternalHost != "", validation.Required), ozzo_validators.IsValidPort),
	)
}
