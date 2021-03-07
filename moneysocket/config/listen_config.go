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

type ListenConfig struct {
	// Default listening bind setting. 127.0.0.1 for localhost connections, 0.0.0.0
	// for allowing connections from other hosts
	BindHost string
	// default port to listen for websocket connections port not specified.
	BindPort int
	//  host for other devices to connect via the beacon
	ExternalHost string
	// port for other devices to connect via the beacon
	ExternalPort int
	// Use TLS for websocket connections
	useTLS bool
	// if UseTLS is True, use this cert file
	certFile string
	// if UseTLS is True, use this key file
	certKey string
	// if UseTLS is True and we have a self-made cert for testing use this key file
	// we don't need to provide a cert chain
	selfSignedCert bool
	// If we have a 'real' cert, we typically need to provide the cert chain file to
	// make the browser clients happy.
	certChainFile string
	// Default listening bind setting. 127.0.0.1 for localhost connections, 0.0.0.0
	// for allowing connections from other hosts
	defaultBind string
	// default port to listen for websocket connections port not specified.
	defaultPort int
}

// this function makes no gurantees about files being present
// this should be verified seperately. Reads certificates from filesystem
func (l ListenConfig) getCertificate() tlsutils.TlsCert {
	pub, _ := ioutil.ReadFile(l.certFile)

	priv, _ := ioutil.ReadFile(l.certKey)

	return tlsutils.TlsCert{
		PublicKey:  string(pub),
		PrivateKey: string(priv),
	}
}

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
		isValid, err := tlsutils.VerifyCertificate(l.getCertificate())
		if !isValid {
			return err
		}
	}
	// validate root certificate
	if l.useTLS && !validation.IsEmpty(l.certChainFile) {
		rawCert := l.getCertificate()
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
