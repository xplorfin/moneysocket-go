package config

import (
	"testing"

	nettest "github.com/xplorfin/netutils/testutils"
	tlsmock "github.com/xplorfin/tlsutils/mock"
)

type RelayConfigTest struct {
	// wether or not config is valid
	isValid bool
	// what we're testing
	testCase string
	// config to test
	config RelayConfig
}

// test variations where tls is false.
func makeTLSFalseVariations(t *testing.T) []RelayConfigTest {
	return []RelayConfigTest{
		{
			isValid:  true,
			testCase: "test use false",
			config: RelayConfig{
				BindHost:       "127.0.0.1",
				BindPort:       nettest.GetFreePort(t),
				useTLS:         false,
				certFile:       "",
				certKey:        "",
				selfSignedCert: false,
				certChainFile:  "",
			},
		},
		{
			isValid:  false,
			testCase: "test taken bind port",
			config: RelayConfig{
				BindHost:       "127.0.0.1",
				BindPort:       nettest.GetUnfreePort(t),
				useTLS:         false,
				certFile:       "",
				certKey:        "",
				selfSignedCert: false,
				certChainFile:  "",
			},
		},
		{
			isValid:  false,
			testCase: "test invalid bind host",
			config: RelayConfig{
				BindHost:       "%$127.0.0.1",
				BindPort:       nettest.GetFreePort(t),
				useTLS:         false,
				certFile:       "",
				certKey:        "",
				selfSignedCert: false,
				certChainFile:  "",
			},
		},
	}
}

// test variations where tls is false.
func makeCertVariations(t *testing.T) []RelayConfigTest {
	validCertFile, validKeyFile := tlsmock.TemporaryCert(t)
	chainFile, serverCertFile, serverKeyFile := tlsmock.TemporaryCertInChain(t)
	return []RelayConfigTest{
		{
			isValid:  true,
			testCase: "test valid cert without cert chain",
			config: RelayConfig{
				BindHost:       "127.0.0.1",
				BindPort:       nettest.GetFreePort(t),
				useTLS:         true,
				certFile:       validCertFile,
				certKey:        validKeyFile,
				selfSignedCert: true,
				certChainFile:  "",
			},
		},
		{
			isValid:  true,
			testCase: "test valid cert with cert chain",
			config: RelayConfig{
				BindHost:       "127.0.0.1",
				BindPort:       nettest.GetFreePort(t),
				useTLS:         true,
				certFile:       serverCertFile,
				certKey:        serverKeyFile,
				selfSignedCert: true,
				certChainFile:  chainFile,
			},
		},
		{
			isValid:  false,
			testCase: "test valid cert with invalid cert chain",
			config: RelayConfig{
				BindHost:       "127.0.0.1",
				BindPort:       nettest.GetFreePort(t),
				useTLS:         true,
				certFile:       serverCertFile,
				certKey:        serverKeyFile,
				selfSignedCert: true,
				certChainFile:  nettest.MockFile(t),
			},
		},
		{
			isValid:  false,
			testCase: "test invalid (but existent) cert file",
			config: RelayConfig{
				BindHost:       "127.0.0.1",
				BindPort:       nettest.GetFreePort(t),
				useTLS:         true,
				certFile:       nettest.MockFile(t),
				certKey:        validKeyFile,
				selfSignedCert: true,
				certChainFile:  "",
			},
		},
		{
			isValid:  false,
			testCase: "test invalid (but existent) key file",
			config: RelayConfig{
				BindHost:       "127.0.0.1",
				BindPort:       nettest.GetFreePort(t),
				useTLS:         true,
				certFile:       validCertFile,
				certKey:        nettest.MockFile(t),
				selfSignedCert: true,
				certChainFile:  "",
			},
		},
		{
			isValid:  false,
			testCase: "swap ssl params",
			config: RelayConfig{
				BindHost:       "127.0.0.1",
				BindPort:       nettest.GetFreePort(t),
				useTLS:         true,
				certFile:       validKeyFile,
				certKey:        validCertFile,
				selfSignedCert: true,
				certChainFile:  "",
			},
		},
	}
}

func makeRelayConfigTests(t *testing.T) (configs []RelayConfigTest) {
	configs = append(configs, makeTLSFalseVariations(t)...)
	configs = append(configs, makeCertVariations(t)...)
	return configs
}

func TestConfig(t *testing.T) {
	testConfigs := makeRelayConfigTests(t)
	for _, tc := range testConfigs {
		err := tc.config.Validate()
		if tc.isValid && err != nil {
			t.Errorf("test case '%s' failed with error %s. expected error to be nil", tc.testCase, err.Error())
		} else if !tc.isValid && err == nil {
			t.Errorf("expected error for test case %s, got none", tc.testCase)
		}
	}
}
