/*
Copyright 2022 The deepauto-io LLC.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tls

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/deepauto-io/nsqcc/filepath/ifs"
	"github.com/kelseyhightower/envconfig"
	"log"

	"github.com/youmark/pkcs8"
)

// ClientCertConfig contains config fields for a client certificate.
type ClientCertConfig struct {
	CertFile string `json:"cert_file" envconfig:"TLS_CERT_FILE" json:"cert_file" yaml:"cert_file"`
	KeyFile  string `json:"key_file" envconfig:"TLS_KEY_FILE" json:"key_file" yaml:"key_file"`
	Cert     string `json:"cert" envconfig:"TLS_CERT" json:"cert" yaml:"cert"`
	Key      string `json:"key" envconfig:"TLS_KEY" json:"key" yaml:"key"`
	Password string `json:"password" envconfig:"TLS_PASSWORD" json:"password" yaml:"password"`
}

// Config contains configuration params for TLS.
type Config struct {
	Enabled             bool               `json:"enabled" envconfig:"TLS_ENABLE"   json:"enabled" yaml:"enabled"`
	RootCAs             string             `json:"root_c_as" envconfig:"TLS_ROOTCAS"  json:"root_cas" yaml:"root_cas"`
	RootCAsFile         string             `json:"root_c_as_file" envconfig:"TLS_ROOTCAS_FILE"             json:"root_cas_file" yaml:"root_cas_file"`
	InsecureSkipVerify  bool               `json:"insecure_skip_verify" envconfig:"TLS_INSECURE_SKIP_VERIFY" json:"skip_cert_verify" yaml:"skip_cert_verify"`
	ClientCertificates  []ClientCertConfig `json:"client_certs" yaml:"client_certs"`
	EnableRenegotiation bool               `json:"enable_renegotiation" envconfig:"TLS_ENABLE_RENEGOTIATION" json:"enable_renegotiation" yaml:"enable_renegotiation"`
}

// NewConfig creates a new Config with default values.
func NewConfig() Config {
	return Config{
		Enabled:             false,
		RootCAs:             "",
		RootCAsFile:         "",
		InsecureSkipVerify:  false,
		ClientCertificates:  []ClientCertConfig{},
		EnableRenegotiation: false,
	}
}

func MustLoadConfig(cfgPath string) Config {
	cfg := Config{}
	if cfgPath == "" {
		if err := envconfig.Process("", &cfg); err != nil {
			log.Fatal(err)
		}
		return cfg
	}
	return cfg
}

func defaultTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
}

// GetNonToggled returns a valid *tls.Config based on the configuration values
// of Config. If none of the config fields are set then a nil config is
// returned.
func (c *Config) GetNonToggled(f ifs.FS) (*tls.Config, error) {
	var tlsConf *tls.Config
	initConf := func() {
		if tlsConf != nil {
			return
		}
		tlsConf = defaultTLSConfig()
	}

	if len(c.RootCAs) > 0 && len(c.RootCAsFile) > 0 {
		return nil, errors.New("only one field between root_cas and root_cas_file can be specified")
	}

	if len(c.RootCAsFile) > 0 {
		caCert, err := ifs.ReadFile(f, c.RootCAsFile)
		if err != nil {
			return nil, err
		}
		initConf()
		tlsConf.RootCAs = x509.NewCertPool()
		tlsConf.RootCAs.AppendCertsFromPEM(caCert)
	}

	if len(c.RootCAs) > 0 {
		initConf()
		tlsConf.RootCAs = x509.NewCertPool()
		tlsConf.RootCAs.AppendCertsFromPEM([]byte(c.RootCAs))
	}

	for _, conf := range c.ClientCertificates {
		cert, err := conf.Load(f)
		if err != nil {
			return nil, err
		}
		initConf()
		tlsConf.Certificates = append(tlsConf.Certificates, cert)
	}

	if c.EnableRenegotiation {
		initConf()
		tlsConf.Renegotiation = tls.RenegotiateFreelyAsClient
	}

	if c.InsecureSkipVerify {
		initConf()
		tlsConf.InsecureSkipVerify = true
	}

	return tlsConf, nil
}

// Get returns a valid *tls.Config based on the configuration values of Config,
// or nil if tls is not enabled.
func (c *Config) Get(f ifs.FS) (*tls.Config, error) {
	if !c.Enabled {
		return nil, nil
	}
	tConf, err := c.GetNonToggled(f)
	if err != nil {
		return nil, err
	}
	if tConf == nil {
		tConf = defaultTLSConfig()
	}
	return tConf, nil
}

func getKeyPair(cert []byte, keyType string, keyBytes []byte) (tls.Certificate, error) {
	return tls.X509KeyPair(cert, pem.EncodeToMemory(&pem.Block{Type: keyType, Bytes: keyBytes}))
}

func loadKeyPair(cert, key []byte, password string) (tls.Certificate, error) {
	keyPem, _ := pem.Decode(key)
	if keyPem == nil {
		return tls.Certificate{}, errors.New("failed to decode private key")
	}

	var err error
	//nolint:staticcheck // SA1019 Disable linting for deprecated x509.IsEncryptedPEMBlock call
	if x509.IsEncryptedPEMBlock(keyPem) {
		if password == "" {
			return tls.Certificate{}, errors.New("missing password for PKCS#1 encrypted private key")
		}

		var decryptedKey []byte
		//nolint:staticcheck // SA1019 Disable linting for deprecated x509.DecryptPEMBlock call
		if decryptedKey, err = x509.DecryptPEMBlock(keyPem, []byte(password)); err != nil {
			return tls.Certificate{}, fmt.Errorf("failed to parse encrypted PKCS#1 private key: %s", err)
		}

		// x509.DecryptPEMBlock() can sometimes fail to detect invalid passwords and will return a nil error, so we
		// should validate the decrypted key. Otherwise, tls.X509KeyPair() will return an error anyway, but it
		// wouldn't be clear why. Details here: https://github.com/golang/go/issues/10171
		// and here https://cs.opensource.google/go/go/+/refs/tags/go1.19.3:src/crypto/tls/tls.go;l=339
		validKey := false
		if _, err = x509.ParsePKCS1PrivateKey(decryptedKey); err == nil {
			validKey = true
		}
		if _, err = x509.ParsePKCS8PrivateKey(decryptedKey); err == nil {
			validKey = true
		}
		if _, err := x509.ParseECPrivateKey(decryptedKey); err == nil {
			validKey = true
		}
		if !validKey {
			return tls.Certificate{}, fmt.Errorf("failed to decrypt PKCS#1 key: %s", x509.IncorrectPasswordError)
		}

		return getKeyPair(cert, keyPem.Type, decryptedKey)
	} else if keyPem.Type == "ENCRYPTED PRIVATE KEY" {
		if password == "" {
			return tls.Certificate{}, errors.New("missing password for PKCS#8 encrypted private key")
		}

		var decryptedKey *rsa.PrivateKey
		if decryptedKey, err = pkcs8.ParsePKCS8PrivateKeyRSA(keyPem.Bytes, []byte(password)); err != nil {
			return tls.Certificate{}, fmt.Errorf("failed to parse encrypted PKCS#8 private key: %s", err)
		}
		return getKeyPair(cert, keyPem.Type, x509.MarshalPKCS1PrivateKey(decryptedKey))
	}

	return tls.X509KeyPair(cert, key)
}

// Load returns a TLS certificate, based on either file paths in the
// config or the raw certs as strings.
func (c *ClientCertConfig) Load(f ifs.FS) (tls.Certificate, error) {
	if c.CertFile != "" || c.KeyFile != "" {
		if c.CertFile == "" {
			return tls.Certificate{}, errors.New("missing cert_file field in client certificate config")
		}
		if c.KeyFile == "" {
			return tls.Certificate{}, errors.New("missing key_file field in client certificate config")
		}

		cert, err := ifs.ReadFile(f, c.CertFile)
		if err != nil {
			return tls.Certificate{}, err
		}

		key, err := ifs.ReadFile(f, c.KeyFile)
		if err != nil {
			return tls.Certificate{}, err
		}

		return loadKeyPair(cert, key, c.Password)
	}

	if c.Cert == "" {
		return tls.Certificate{}, errors.New("missing cert field in client certificate config")
	}
	if c.Key == "" {
		return tls.Certificate{}, errors.New("missing key field in client certificate config")
	}

	return loadKeyPair([]byte(c.Cert), []byte(c.Key), c.Password)
}
