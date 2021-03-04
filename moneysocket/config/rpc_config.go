package config

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	ozzo_validators "github.com/xplorfin/ozzo-validators"
)

type RpcConfig struct {
	// host for client to connect
	BindHost string
	// port for client to connect
	BindPort int
	// host for a client to connect
	ExternalHost string
	// port for client to connect
	ExternalPort int
}

func (l RpcConfig) Validate() error {
	return validation.ValidateStruct(&l,
		// bind host cannot be null
		validation.Field(&l.BindHost, validation.Required, is.Host),
		// bind port cannot be null
		validation.Field(&l.BindPort, validation.Required, ozzo_validators.IsValidPort),
		// required when port is not null
		validation.Field(&l.ExternalHost, validation.When(l.ExternalPort != 0, validation.Required), is.Host),
		// required when host is not null
		validation.Field(&l.ExternalPort, validation.When(l.ExternalHost != "", validation.Required), ozzo_validators.IsValidPort),
	)
}
