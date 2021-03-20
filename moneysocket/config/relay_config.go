package config

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	validators "github.com/xplorfin/ozzo-validators"
	"github.com/xplorfin/tlsutils"
)

// RelayConfig defines the configuration for the relay
type RelayConfig struct {
	// BindHost defines listening bind setting. Defaults to 127.0.0.1 for localhost connections, 0.0.0.0
	// for allowing connections from other hosts
	BindHost string
	// BindPort defines the default port to listen for websocket connections port not specified.
	BindPort int
	// useTLS for relay connections
	useTLS bool
	// certFile defines the file if useTLS is True, use this cert file
	certFile string
	// certKey if useTLS is True, use this key file
	certKey string
	// selfSignedCert if useTLS is True and we have a self-made cert for testing use this key file
	// we don't need to provide a cert chain
	selfSignedCert bool
	// certChainFile if we have a 'real' cert, we typically need to provide the cert chain file to
	// make the browser clients happy.
	certChainFile string
}

func (r RelayConfig) certPath() string {
	return r.certFile
}

func (r RelayConfig) certKeyPath() string {
	return r.certKey
}

// Validate the configuration
func (r RelayConfig) Validate() error {
	err := validation.ValidateStruct(&r,
		// bind host cannot be null
		validation.Field(&r.BindHost, validation.Required, is.Host),
		// bind port cannot be null
		validation.Field(&r.BindPort, validation.Required, validators.IsPortAvailable),
		// use tls must be set (should always be set)
		validation.Field(&r.useTLS),
		// certFile is required when use tls is true
		validation.Field(&r.certFile, validation.When(r.useTLS, validation.Required, validators.IsValidPath)),
		// certKey is required when use tls is true
		validation.Field(&r.certKey, validation.When(r.useTLS, validation.Required, validators.IsValidPath)),
		// cert chain file must be a valid path
		validation.Field(&r.certChainFile, validators.IsValidPath),
	)
	if err != nil {
		return err
	}
	// just validate the ssl certs
	if r.useTLS && validation.IsEmpty(r.certChainFile) {
		isValid, err := tlsutils.VerifyCertificate(getCertificate(r))
		if !isValid {
			return err
		}
	}
	// validate root certificate
	if r.useTLS && !validation.IsEmpty(r.certChainFile) {
		rawCert := getCertificate(r)
		// we already validated this exists
		chainFile, _ := ioutil.ReadFile(r.certChainFile)
		block, _ := pem.Decode(chainFile)
		if block == nil {
			return fmt.Errorf("expected chainfile block to not be null")
		}
		chain, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return err
		}
		block, _ = pem.Decode([]byte(rawCert.PublicKey))
		server, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return err
		}
		isValid, err := tlsutils.VerifyLowNoDca(chain, server)
		if !isValid {
			return err
		}
	}
	// TODO validate ssl certs
	return nil
}

var _ certConfig = &RelayConfig{}
