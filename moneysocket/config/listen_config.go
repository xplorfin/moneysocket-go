package config

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	ozzo_validators "github.com/xplorfin/ozzo-validators"
	"github.com/xplorfin/tlsutils"
)

// ListenConfig specifies the configuration for terminus to listen on.
type ListenConfig struct {
	// BindHost defines listening bind setting. Defaults to 127.0.0.1 for localhost connections, 0.0.0.0
	// for allowing connections from other hosts
	BindHost string
	// BindPort defines the default port to listen for websocket connections port not specified.
	BindPort int
	// ExternalHost defines the host for other devices to connect via the beacon
	ExternalHost string
	// ExternalPort for other devices to connect via the beacon
	ExternalPort int
	// useTLS for websocket connections
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
	// defaultBind defines the default listening bind setting. 127.0.0.1 for localhost connections, 0.0.0.0
	// for allowing connections from other hosts
	defaultBind string
	// defaultPort to listen for websocket connections port not specified.
	defaultPort int
}

// certPath is the path to the certificate file.
func (l ListenConfig) certPath() string {
	return l.certFile
}

// certKeyPath is the path to the certificate key file.
func (l ListenConfig) certKeyPath() string {
	return l.certKey
}

// Validate the configuration.
func (l ListenConfig) Validate() error {
	err := validation.ValidateStruct(&l,
		// bind host cannot be null
		validation.Field(&l.BindHost, validation.Required, is.Host),
		// bind port cannot be null
		validation.Field(&l.BindPort, validation.Required, ozzo_validators.IsPortAvailable),
		// external host cannot be null and must be valid
		validation.Field(&l.ExternalHost, validation.Required, is.Host),
		// external port must be available
		validation.Field(&l.ExternalPort, validation.Required, ozzo_validators.IsPortAvailable),
		// use tls must be set (should always be set)
		validation.Field(&l.useTLS),
		// certFile is required when use tls is true
		validation.Field(&l.certFile, validation.When(l.useTLS, validation.Required, ozzo_validators.IsValidPath)),
		// certKey is required when use tls is true
		validation.Field(&l.certKey, validation.When(l.useTLS, validation.Required, ozzo_validators.IsValidPath)),
		// cert chain file must be a valid path
		validation.Field(&l.certChainFile, ozzo_validators.IsValidPath),
	)
	if err != nil {
		return err
	}
	// just validate the ssl certs
	if l.useTLS && validation.IsEmpty(l.certChainFile) {
		isValid, err := tlsutils.VerifyCertificate(getCertificate(l))
		if !isValid {
			return err
		}
	}
	// validate root certificate
	if l.useTLS && !validation.IsEmpty(l.certChainFile) {
		rawCert := getCertificate(l)
		// we already validated this exists
		chainFile, _ := ioutil.ReadFile(l.certChainFile)
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

var _ certConfig = &ListenConfig{}
