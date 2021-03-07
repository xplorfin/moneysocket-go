package config

import (
	"testing"

	"github.com/brianvoe/gofakeit/v5"
	nettest "github.com/xplorfin/netutils/testutils"
	tlsmock "github.com/xplorfin/tlsutils/mock"
)

type ListenConfigTest struct {
	// wether or not config is valid
	isValid bool
	// what we're testing
	testCase string
	// config to test
	config ListenConfig
}

// test variations where tls is false
func makeTlsFalseVariations(t *testing.T) []ListenConfigTest {
	return []ListenConfigTest{
		{
			isValid:  true,
			testCase: "test use false",
			config: ListenConfig{
				BindHost:       "127.0.0.1",
				BindPort:       nettest.GetFreePort(t),
				ExternalHost:   gofakeit.Word(),
				ExternalPort:   nettest.GetFreePort(t),
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
			config: ListenConfig{
				BindHost:       "127.0.0.1",
				BindPort:       nettest.GetUnfreePort(t),
				ExternalHost:   gofakeit.Word(),
				ExternalPort:   nettest.GetFreePort(t),
				useTLS:         false,
				certFile:       "",
				certKey:        "",
				selfSignedCert: false,
				certChainFile:  "",
			},
		},
		{
			isValid:  false,
			testCase: "test taken external port",
			config: ListenConfig{
				BindHost:       "127.0.0.1",
				BindPort:       nettest.GetFreePort(t),
				ExternalHost:   gofakeit.Word(),
				ExternalPort:   nettest.GetUnfreePort(t),
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
			config: ListenConfig{
				BindHost:       "%$127.0.0.1",
				BindPort:       nettest.GetFreePort(t),
				ExternalHost:   gofakeit.Word(),
				ExternalPort:   nettest.GetUnfreePort(t),
				useTLS:         false,
				certFile:       "",
				certKey:        "",
				selfSignedCert: false,
				certChainFile:  "",
			},
		},
		{
			isValid:  false,
			testCase: "test invalid external host",
			config: ListenConfig{
				BindHost:       "127.0.0.1",
				BindPort:       nettest.GetFreePort(t),
				ExternalHost:   "%$127.0.0.1",
				ExternalPort:   nettest.GetUnfreePort(t),
				useTLS:         false,
				certFile:       "",
				certKey:        "",
				selfSignedCert: false,
				certChainFile:  "",
			},
		},
	}
}

// test variations where tls is false
func makeCertVariations(t *testing.T) []ListenConfigTest {
	validCertFile, validKeyFile := tlsmock.TemporaryCert(t)
	chainFile, serverCertFile, serverKeyFile := tlsmock.TemporaryCertInChain(t)
	return []ListenConfigTest{
		{
			isValid:  true,
			testCase: "test valid cert without cert chain",
			config: ListenConfig{
				BindHost:       "127.0.0.1",
				BindPort:       nettest.GetFreePort(t),
				ExternalHost:   gofakeit.Word(),
				ExternalPort:   nettest.GetFreePort(t),
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
			config: ListenConfig{
				BindHost:       "127.0.0.1",
				BindPort:       nettest.GetFreePort(t),
				ExternalHost:   gofakeit.Word(),
				ExternalPort:   nettest.GetFreePort(t),
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
			config: ListenConfig{
				BindHost:       "127.0.0.1",
				BindPort:       nettest.GetFreePort(t),
				ExternalHost:   gofakeit.Word(),
				ExternalPort:   nettest.GetFreePort(t),
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
			config: ListenConfig{
				BindHost:       "127.0.0.1",
				BindPort:       nettest.GetFreePort(t),
				ExternalHost:   gofakeit.Word(),
				ExternalPort:   nettest.GetFreePort(t),
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
			config: ListenConfig{
				BindHost:       "127.0.0.1",
				BindPort:       nettest.GetFreePort(t),
				ExternalHost:   gofakeit.Word(),
				ExternalPort:   nettest.GetFreePort(t),
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
			config: ListenConfig{
				BindHost:       "127.0.0.1",
				BindPort:       nettest.GetFreePort(t),
				ExternalHost:   gofakeit.Word(),
				ExternalPort:   nettest.GetFreePort(t),
				useTLS:         true,
				certFile:       validKeyFile,
				certKey:        validCertFile,
				selfSignedCert: true,
				certChainFile:  "",
			},
		},
	}
}

func makeListenConfigTests(t *testing.T) (configs []ListenConfigTest) {
	configs = append(configs, makeTlsFalseVariations(t)...)
	configs = append(configs, makeCertVariations(t)...)
	return configs
}

func TestConfig(t *testing.T) {
	testConfigs := makeListenConfigTests(t)
	for _, tc := range testConfigs {
		err := tc.config.Validate()
		if tc.isValid && err != nil {
			t.Errorf("test case '%s' failed with error %s. expected error to be nil", tc.testCase, err.Error())
		} else if !tc.isValid && err == nil {
			t.Errorf("expected error for test case %s, got none", tc.testCase)
		}
	}
}
