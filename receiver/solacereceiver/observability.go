// Copyright 2020 OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package solacereceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/solacereceiver"

import (
	"context"
	"strings"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

const (
	// ReceiverKey used to identify receivers in metrics and traces.
	receiverKey    = "receiver"
	nameSep        = "/"
	receiverPrefix = receiverKey + nameSep
)

// receiver will register internal telemetry views
func init() {
	view.Register(
		viewFailedReconnections,
		viewRecoverableUnmarshallingErrors,
		viewFatalUnmarshallingErrors,
		viewNeedRestart,
	)
}

var (
	failedReconnections            = stats.Int64("solacereceiver/failed_reconnections", "Number of failed broker reconnections", "1")
	recoverableUnmarshallingErrors = stats.Int64("solacereceiver/recoverable_unmarshalling_errors", "Number of recoverable errors by messages unmarshalling", "1")
	fatalUnmarshallingErrors       = stats.Int64("solacereceiver/fatal_unmarshalling_errors", "Number of fatal errors by messages unmarshalling", "1")
	needRestart                    = stats.Int64("solacereceiver/need_restart", "Indicates with value 1 that receiver require a restart", "1")
)

var viewFailedReconnections = &view.View{
	Name:        buildReceiverCustomMetricName(string(TypeStr), failedReconnections.Name()),
	Description: failedReconnections.Description(),
	Measure:     failedReconnections,
	Aggregation: view.Sum(),
}

var viewRecoverableUnmarshallingErrors = &view.View{
	Name:        buildReceiverCustomMetricName(string(TypeStr), recoverableUnmarshallingErrors.Name()),
	Description: recoverableUnmarshallingErrors.Description(),
	Measure:     recoverableUnmarshallingErrors,
	Aggregation: view.Sum(),
}

var viewFatalUnmarshallingErrors = &view.View{
	Name:        buildReceiverCustomMetricName(string(TypeStr), fatalUnmarshallingErrors.Name()),
	Description: fatalUnmarshallingErrors.Description(),
	Measure:     fatalUnmarshallingErrors,
	Aggregation: view.Sum(),
}

var viewNeedRestart = &view.View{
	Name:        buildReceiverCustomMetricName(string(TypeStr), needRestart.Name()),
	Description: needRestart.Description(),
	Measure:     needRestart,
	Aggregation: view.LastValue(),
}

func buildReceiverCustomMetricName(configType, metric string) string {
	componentPrefix := receiverPrefix
	if !strings.HasSuffix(componentPrefix, nameSep) {
		componentPrefix += nameSep
	}
	if configType == "" {
		return componentPrefix
	}
	return componentPrefix + configType + nameSep + metric
}

// RecordFailedReconnection increments the metric that records failed reconnection event.
func RecordFailedReconnection() {
	stats.Record(context.Background(), failedReconnections.M(int64(1)))
}

// RecordRecoverableUnmarshallingError increments the metric that records a recoverable error by trace message unmarshalling.
func RecordRecoverableUnmarshallingError() {
	stats.Record(context.Background(), recoverableUnmarshallingErrors.M(int64(1)))
}

// RecordFatalUnmarshallingError increments the metric that records a fatal arrow by trace message unmarshalling.
func RecordFatalUnmarshallingError() {
	stats.Record(context.Background(), fatalUnmarshallingErrors.M(int64(1)))
}

// RecordNeedRestart turns a need restart flag on
func RecordNeedRestart() {
	stats.Record(context.Background(), needRestart.M(1))
}
