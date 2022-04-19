// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package solacereceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/solacereceiver"

import (
	"errors"
	"strings"

	"github.com/Azure/go-amqp"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/configtls"
)

var (
	// ErrNilNextConsumer indicates an error on nil next consumer.
	ErrInvalidConfig = errors.New("nil nextConsumer")

	// ErrDataTypeIsNotSupported can be returned by receiver, exporter or processor
	// factory methods that create the entity if the particular telemetry
	// data type is not supported by the receiver, exporter or processor.
	ErrDataTypeIsNotSupported = errors.New("telemetry type is not supported")
)

const (
	// 8Kb
	saslMaxInitFrameSizeOverride = 8000
	// default queue
)

// Config defines configuration for Solace receiver.
type Config struct {
	config.ReceiverSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct
	// The list of solace brokers (default localhost:5671)
	Broker []string `mapstructure:"broker"`

	// The name of the solace queue to consume from, it is required parameter
	Queue string `mapstructure:"queue"`

	// The maximum number of unacknowledged messages the Solace broker can transmit, to configure AMQP Link
	MaxUnaked uint32 `mapstructure:"max_unacknowledged"`

	Transport Transport `mapstructure:"transport"`

	Auth Authentication `mapstructure:"auth"`
}

var _ config.Receiver = (*Config)(nil)

// Validate checks the receiver configuration is valid
func (cfg *Config) Validate() error {
	if cfg.Auth.PlainText == nil && cfg.Auth.External == nil && cfg.Auth.XAuth2 == nil {
		missingFileds := []string{}
		return NewMissingConfigurationError("Authentication details are required, either for plain user name password or XOAUTH2 or client certificate", nil, &missingFileds)
	}

	if len(strings.TrimSpace(cfg.Queue)) == 0 {
		return NewMissingConfigurationError("Queue definition is required, queue definition has format queue://<queuename>", nil, nil)
	}

	// TODO add check for missing required parameter when auth strategie is defined but parameter are mossing
	// no error
	return nil
}

type Transport struct {
	TLS configtls.TLSClientSetting `mapstructure:"tls,omitempty"`
}

// Authentication defines authentication strategies.
type Authentication struct {
	PlainText *SaslPlainTextConfig `mapstructure:"sasl_plain_text"`
	XAuth2    *SaslXAuth2Config    `mapstructure:"sasl_xauth2"`
	External  *SaslExternalConfig  `mapstructure:"sasl_external"`
}

// SaslPlainTextConfig defines SASL PLAIN authentication.
type SaslPlainTextConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// SaslXAuth2Config defines the configuration for the SASL XAUTH2 authentication.
type SaslXAuth2Config struct {
	Username string `mapstructure:"username"`
	Bearer   string `mapstructure:"bearer"`
}

// SaslExternalConfig defines the configuration for the SASL External used in conjunction with TLS client authentication.
type SaslExternalConfig struct {
}

// ConfigureAuthentication configures authentication in amqp.ConnOption slice
func ConfigureAuthentication(config *Config) (amqp.ConnOption, error) {

	if config.Auth.PlainText != nil {
		conOption, err := configurePlaintext(config.Auth.PlainText)
		if err != nil {
			return nil, err
		}
		return conOption, nil
	}

	if config.Auth.XAuth2 != nil {
		conOption, err := configureXAuth2(config.Auth.XAuth2)
		if err != nil {
			return nil, err
		}
		return conOption, nil
	}

	if config.Auth.External != nil {
		err, conOption := configureExternal()
		if err != nil {
			return nil, err
		}
		return conOption, nil
	}

	return nil, errors.New("auth params not correctly configured")
}

func configurePlaintext(config *SaslPlainTextConfig) (amqp.ConnOption, error) {
	if config.Password == "" || config.Username == "" {
		return nil, errors.New("missing plain text auth params: Username, Password")
	}
	return amqp.ConnSASLPlain(config.Username, config.Password), nil
}

func configureXAuth2(config *SaslXAuth2Config) (amqp.ConnOption, error) {
	if config.Bearer == "" || config.Username == "" {
		return nil, errors.New("missing xauth2 text auth params: Username, Bearer")
	}
	return amqp.ConnSASLXOAUTH2(config.Username, config.Bearer, saslMaxInitFrameSizeOverride), nil
}

func ConfigureTLS(config *configtls.TLSClientSetting) (error, amqp.ConnOption) {
	tlsConfig, err := config.LoadTLSConfig()
	if err != nil {
		return err, nil
	}
	return nil, amqp.ConnTLSConfig(tlsConfig)
}

func configureExternal() (error, amqp.ConnOption) {
	// empty string is common for TLS
	return nil, amqp.ConnSASLExternal("")
}
