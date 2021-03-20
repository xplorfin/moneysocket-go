package config

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/macaroons"
	validators "github.com/xplorfin/ozzo-validators"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/macaroon.v2"
)

// LndConfig defines the lnd config (if lnd is in use)
type LndConfig struct {
	// LndDir is the lnd settings directory
	LndDir string
	// MacaroonPath is the path to the macaroon to use for the lnd connection (requires invoice permissions)
	MacaroonPath string
	// TLSCertPath for LND, which is different from the websocket listening TLS
	TLSCertPath string
	// Network the lnd node is running on (e.g. mainnet, testnet, regtest)
	// TODO this can be handled with a getinfo command
	Network string
	// GrpcHost is the grpc host of the lnd node
	GrpcHost string
	//GrpcPort is the grpc port of the lnd node
	GrpcPort int
}

// Validate validates the LndConfig is valid
func (l LndConfig) Validate() (err error) {
	err = validation.ValidateStruct(&l,
		// validate lnd dir exists
		validation.Field(&l.LndDir, validators.IsDir),
		// validate the macaroon path points to a file
		validation.Field(&l.MacaroonPath, validators.IsFile),
		// validate the tls cert path points to a file
		validation.Field(&l.TLSCertPath, validators.IsFile),
		// validate the grpc host name
		validation.Field(&l.GrpcHost, validators.IsSpaceless),
		// validate the grpc port
		validation.Field(&l.GrpcPort, validators.IsValidPort),
	)

	if err != nil {
		return err
	}

	if l.Network != "bitcoin" {
		return fmt.Errorf("network %s invalid, expected bitcoin", l.Network)
	}

	// we validate these seperately from Grpc for easier tracebacks
	if l.TLSCertPath != "" {
		_, err = l.GetTLSCert()
		if err != nil {
			return err
		}
	}

	if l.MacaroonPath != "" {
		_, err = l.GetMacaroon()
		if err != nil {
			return err
		}
	}

	// TODO validate grpc connection
	return err
}

// GetTLSCert fetches the tls cert from an LndConfig
func (l LndConfig) GetTLSCert() (cert *tls.Config, err error) {
	tlsCert, err := ioutil.ReadFile(l.TLSCertPath)
	if err != nil {
		return cert, err
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(tlsCert) {
		return cert, err
	}
	cert = &tls.Config{
		InsecureSkipVerify: true,
		RootCAs:            cp,
	}
	return cert, nil
}

// GetMacaroon fetches the encoded macaroon object from LndConfig
func (l LndConfig) GetMacaroon() (mac *macaroon.Macaroon, err error) {
	rawMac, err := ioutil.ReadFile(l.MacaroonPath)
	if err != nil {
		return nil, err
	}

	mac = &macaroon.Macaroon{}
	err = mac.UnmarshalBinary(rawMac)
	if err != nil {
		return nil, err
	}

	return mac, err
}

// LndHost gets the lnd hostname/port
func (l LndConfig) LndHost() string {
	return fmt.Sprintf("%s:%d", l.GrpcHost, l.GrpcPort)
}

// GRPCConnection returns the grpc connection object
func (l LndConfig) GRPCConnection(ctx context.Context) (conn *grpc.ClientConn, err error) {
	cert, err := l.GetTLSCert()
	if err != nil {
		return conn, err
	}
	mac, err := l.GetMacaroon()
	if err != nil {
		return conn, err
	}

	return grpc.DialContext(
		ctx,
		l.LndHost(),
		grpc.WithTransportCredentials(credentials.NewTLS(cert)),
		grpc.WithPerRPCCredentials(macaroons.NewMacaroonCredential(mac)))
}

// RPCClient generates an lnd rpc client from an LndConfig
func (l LndConfig) RPCClient(ctx context.Context) (rpcClient lnrpc.LightningClient, err error) {
	grpcConn, err := l.GRPCConnection(ctx)
	if err != nil {
		return nil, err
	}
	return lnrpc.NewLightningClient(grpcConn), nil
}
