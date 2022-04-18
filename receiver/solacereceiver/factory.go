// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package solacereceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/solacereceiver"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver/receiverhelper"
)

const (
	TypeStr = "solace"
	// default value for max unaked messages
	maxUnaked = 10
)

// NewFactory creates a factory for Solace receiver.
func NewFactory() component.ReceiverFactory {
	return receiverhelper.NewFactory(
		TypeStr,
		createDefaultConfig,
		receiverhelper.WithTraces(createTracesReceiver))
}

// createDefaultConfig creates the default configuration for receiver.
func createDefaultConfig() config.Receiver {

	return &Config{
		ReceiverSettings: config.NewReceiverSettings(config.NewComponentID(TypeStr)),
		Broker:           createDefaultBrokerList(),
		MaxUnaked:        maxUnaked,
		Auth:             Authentication{},
		Transport: Transport{
			TLS: configtls.TLSClientSetting{
				InsecureSkipVerify: false,
				Insecure:           false,
			},
		},
	}
}

// CreateTracesReceiver creates a default list of brokerr with a single entry 'localhost:5671' .
func createDefaultBrokerList() []string {
	return []string{"localhost:5671"}
}

// CreateTracesReceiver creates a trace receiver based on provided config. Component is not shared
func createTracesReceiver(
	_ context.Context,
	params component.ReceiverCreateSettings,
	cfg config.Receiver,
	nextConsumer consumer.Traces) (component.TracesReceiver, error) {
	rcfg := cfg.(*Config)
	return newTracesReceiver(rcfg, params, nextConsumer)
}
