package config

import (
	"io/ioutil"

	"github.com/xplorfin/tlsutils"
)

// certConfig contains a config which can produce a tls cert object
type certConfig interface {
	// certFile is the path to the cert file
	certPath() string
	// certKeyPath is the path to the certKey
	certKeyPath() string
}

// getCertificate gets the certificate
// Note: this function makes no guarantees about files being present
// this should be verified separately. Reads certificates from filesystem
func getCertificate(c certConfig) tlsutils.TlsCert {
	pub, _ := ioutil.ReadFile(c.certPath())

	priv, _ := ioutil.ReadFile(c.certKeyPath())

	return tlsutils.TlsCert{
		PublicKey:  string(pub),
		PrivateKey: string(priv),
	}
}
