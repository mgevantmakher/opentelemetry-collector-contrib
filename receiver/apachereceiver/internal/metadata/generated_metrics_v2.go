// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"time"

	"go.opentelemetry.io/collector/model/pdata"
)

// MetricSettings provides common settings for a particular metric.
type MetricSettings struct {
	Enabled bool `mapstructure:"enabled"`
}

// MetricsSettings provides settings for apachereceiver metrics.
type MetricsSettings struct {
	ApacheCurrentConnections MetricSettings `mapstructure:"apache.current_connections"`
	ApacheRequests           MetricSettings `mapstructure:"apache.requests"`
	ApacheScoreboard         MetricSettings `mapstructure:"apache.scoreboard"`
	ApacheTraffic            MetricSettings `mapstructure:"apache.traffic"`
	ApacheUptime             MetricSettings `mapstructure:"apache.uptime"`
	ApacheWorkers            MetricSettings `mapstructure:"apache.workers"`
}

func DefaultMetricsSettings() MetricsSettings {
	return MetricsSettings{
		ApacheCurrentConnections: MetricSettings{
			Enabled: true,
		},
		ApacheRequests: MetricSettings{
			Enabled: true,
		},
		ApacheScoreboard: MetricSettings{
			Enabled: true,
		},
		ApacheTraffic: MetricSettings{
			Enabled: true,
		},
		ApacheUptime: MetricSettings{
			Enabled: true,
		},
		ApacheWorkers: MetricSettings{
			Enabled: true,
		},
	}
}

type metricApacheCurrentConnections struct {
	data     pdata.Metric   // data buffer for generated metric.
	settings MetricSettings // metric settings provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills apache.current_connections metric with initial data.
func (m *metricApacheCurrentConnections) init() {
	m.data.SetName("apache.current_connections")
	m.data.SetDescription("The number of active connections currently attached to the HTTP server.")
	m.data.SetUnit("connections")
	m.data.SetDataType(pdata.MetricDataTypeSum)
	m.data.Sum().SetIsMonotonic(false)
	m.data.Sum().SetAggregationTemporality(pdata.MetricAggregationTemporalityCumulative)
	m.data.Sum().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricApacheCurrentConnections) recordDataPoint(start pdata.Timestamp, ts pdata.Timestamp, val int64, serverNameAttributeValue string) {
	if !m.settings.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntVal(val)
	dp.Attributes().Insert(A.ServerName, pdata.NewAttributeValueString(serverNameAttributeValue))
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricApacheCurrentConnections) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricApacheCurrentConnections) emit(metrics pdata.MetricSlice) {
	if m.settings.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricApacheCurrentConnections(settings MetricSettings) metricApacheCurrentConnections {
	m := metricApacheCurrentConnections{settings: settings}
	if settings.Enabled {
		m.data = pdata.NewMetric()
		m.init()
	}
	return m
}

type metricApacheRequests struct {
	data     pdata.Metric   // data buffer for generated metric.
	settings MetricSettings // metric settings provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills apache.requests metric with initial data.
func (m *metricApacheRequests) init() {
	m.data.SetName("apache.requests")
	m.data.SetDescription("The number of requests serviced by the HTTP server per second.")
	m.data.SetUnit("1")
	m.data.SetDataType(pdata.MetricDataTypeSum)
	m.data.Sum().SetIsMonotonic(true)
	m.data.Sum().SetAggregationTemporality(pdata.MetricAggregationTemporalityCumulative)
	m.data.Sum().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricApacheRequests) recordDataPoint(start pdata.Timestamp, ts pdata.Timestamp, val int64, serverNameAttributeValue string) {
	if !m.settings.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntVal(val)
	dp.Attributes().Insert(A.ServerName, pdata.NewAttributeValueString(serverNameAttributeValue))
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricApacheRequests) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricApacheRequests) emit(metrics pdata.MetricSlice) {
	if m.settings.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricApacheRequests(settings MetricSettings) metricApacheRequests {
	m := metricApacheRequests{settings: settings}
	if settings.Enabled {
		m.data = pdata.NewMetric()
		m.init()
	}
	return m
}

type metricApacheScoreboard struct {
	data     pdata.Metric   // data buffer for generated metric.
	settings MetricSettings // metric settings provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills apache.scoreboard metric with initial data.
func (m *metricApacheScoreboard) init() {
	m.data.SetName("apache.scoreboard")
	m.data.SetDescription("The number of connections in each state.")
	m.data.SetUnit("scoreboard")
	m.data.SetDataType(pdata.MetricDataTypeSum)
	m.data.Sum().SetIsMonotonic(false)
	m.data.Sum().SetAggregationTemporality(pdata.MetricAggregationTemporalityCumulative)
	m.data.Sum().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricApacheScoreboard) recordDataPoint(start pdata.Timestamp, ts pdata.Timestamp, val int64, serverNameAttributeValue string, scoreboardStateAttributeValue string) {
	if !m.settings.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntVal(val)
	dp.Attributes().Insert(A.ServerName, pdata.NewAttributeValueString(serverNameAttributeValue))
	dp.Attributes().Insert(A.ScoreboardState, pdata.NewAttributeValueString(scoreboardStateAttributeValue))
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricApacheScoreboard) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricApacheScoreboard) emit(metrics pdata.MetricSlice) {
	if m.settings.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricApacheScoreboard(settings MetricSettings) metricApacheScoreboard {
	m := metricApacheScoreboard{settings: settings}
	if settings.Enabled {
		m.data = pdata.NewMetric()
		m.init()
	}
	return m
}

type metricApacheTraffic struct {
	data     pdata.Metric   // data buffer for generated metric.
	settings MetricSettings // metric settings provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills apache.traffic metric with initial data.
func (m *metricApacheTraffic) init() {
	m.data.SetName("apache.traffic")
	m.data.SetDescription("Total HTTP server traffic.")
	m.data.SetUnit("By")
	m.data.SetDataType(pdata.MetricDataTypeSum)
	m.data.Sum().SetIsMonotonic(true)
	m.data.Sum().SetAggregationTemporality(pdata.MetricAggregationTemporalityCumulative)
	m.data.Sum().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricApacheTraffic) recordDataPoint(start pdata.Timestamp, ts pdata.Timestamp, val int64, serverNameAttributeValue string) {
	if !m.settings.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntVal(val)
	dp.Attributes().Insert(A.ServerName, pdata.NewAttributeValueString(serverNameAttributeValue))
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricApacheTraffic) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricApacheTraffic) emit(metrics pdata.MetricSlice) {
	if m.settings.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricApacheTraffic(settings MetricSettings) metricApacheTraffic {
	m := metricApacheTraffic{settings: settings}
	if settings.Enabled {
		m.data = pdata.NewMetric()
		m.init()
	}
	return m
}

type metricApacheUptime struct {
	data     pdata.Metric   // data buffer for generated metric.
	settings MetricSettings // metric settings provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills apache.uptime metric with initial data.
func (m *metricApacheUptime) init() {
	m.data.SetName("apache.uptime")
	m.data.SetDescription("The amount of time that the server has been running in seconds.")
	m.data.SetUnit("s")
	m.data.SetDataType(pdata.MetricDataTypeSum)
	m.data.Sum().SetIsMonotonic(true)
	m.data.Sum().SetAggregationTemporality(pdata.MetricAggregationTemporalityCumulative)
	m.data.Sum().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricApacheUptime) recordDataPoint(start pdata.Timestamp, ts pdata.Timestamp, val int64, serverNameAttributeValue string) {
	if !m.settings.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntVal(val)
	dp.Attributes().Insert(A.ServerName, pdata.NewAttributeValueString(serverNameAttributeValue))
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricApacheUptime) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricApacheUptime) emit(metrics pdata.MetricSlice) {
	if m.settings.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricApacheUptime(settings MetricSettings) metricApacheUptime {
	m := metricApacheUptime{settings: settings}
	if settings.Enabled {
		m.data = pdata.NewMetric()
		m.init()
	}
	return m
}

type metricApacheWorkers struct {
	data     pdata.Metric   // data buffer for generated metric.
	settings MetricSettings // metric settings provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills apache.workers metric with initial data.
func (m *metricApacheWorkers) init() {
	m.data.SetName("apache.workers")
	m.data.SetDescription("The number of workers currently attached to the HTTP server.")
	m.data.SetUnit("connections")
	m.data.SetDataType(pdata.MetricDataTypeSum)
	m.data.Sum().SetIsMonotonic(false)
	m.data.Sum().SetAggregationTemporality(pdata.MetricAggregationTemporalityCumulative)
	m.data.Sum().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricApacheWorkers) recordDataPoint(start pdata.Timestamp, ts pdata.Timestamp, val int64, serverNameAttributeValue string, workersStateAttributeValue string) {
	if !m.settings.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntVal(val)
	dp.Attributes().Insert(A.ServerName, pdata.NewAttributeValueString(serverNameAttributeValue))
	dp.Attributes().Insert(A.WorkersState, pdata.NewAttributeValueString(workersStateAttributeValue))
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricApacheWorkers) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricApacheWorkers) emit(metrics pdata.MetricSlice) {
	if m.settings.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricApacheWorkers(settings MetricSettings) metricApacheWorkers {
	m := metricApacheWorkers{settings: settings}
	if settings.Enabled {
		m.data = pdata.NewMetric()
		m.init()
	}
	return m
}

// MetricsBuilder provides an interface for scrapers to report metrics while taking care of all the transformations
// required to produce metric representation defined in metadata and user settings.
type MetricsBuilder struct {
	startTime                      pdata.Timestamp
	metricApacheCurrentConnections metricApacheCurrentConnections
	metricApacheRequests           metricApacheRequests
	metricApacheScoreboard         metricApacheScoreboard
	metricApacheTraffic            metricApacheTraffic
	metricApacheUptime             metricApacheUptime
	metricApacheWorkers            metricApacheWorkers
}

// metricBuilderOption applies changes to default metrics builder.
type metricBuilderOption func(*MetricsBuilder)

// WithStartTime sets startTime on the metrics builder.
func WithStartTime(startTime pdata.Timestamp) metricBuilderOption {
	return func(mb *MetricsBuilder) {
		mb.startTime = startTime
	}
}

func NewMetricsBuilder(settings MetricsSettings, options ...metricBuilderOption) *MetricsBuilder {
	mb := &MetricsBuilder{
		startTime:                      pdata.NewTimestampFromTime(time.Now()),
		metricApacheCurrentConnections: newMetricApacheCurrentConnections(settings.ApacheCurrentConnections),
		metricApacheRequests:           newMetricApacheRequests(settings.ApacheRequests),
		metricApacheScoreboard:         newMetricApacheScoreboard(settings.ApacheScoreboard),
		metricApacheTraffic:            newMetricApacheTraffic(settings.ApacheTraffic),
		metricApacheUptime:             newMetricApacheUptime(settings.ApacheUptime),
		metricApacheWorkers:            newMetricApacheWorkers(settings.ApacheWorkers),
	}
	for _, op := range options {
		op(mb)
	}
	return mb
}

// Emit appends generated metrics to a pdata.MetricsSlice and updates the internal state to be ready for recording
// another set of data points. This function will be doing all transformations required to produce metric representation
// defined in metadata and user settings, e.g. delta/cumulative translation.
func (mb *MetricsBuilder) Emit(metrics pdata.MetricSlice) {
	mb.metricApacheCurrentConnections.emit(metrics)
	mb.metricApacheRequests.emit(metrics)
	mb.metricApacheScoreboard.emit(metrics)
	mb.metricApacheTraffic.emit(metrics)
	mb.metricApacheUptime.emit(metrics)
	mb.metricApacheWorkers.emit(metrics)
}

// RecordApacheCurrentConnectionsDataPoint adds a data point to apache.current_connections metric.
func (mb *MetricsBuilder) RecordApacheCurrentConnectionsDataPoint(ts pdata.Timestamp, val int64, serverNameAttributeValue string) {
	mb.metricApacheCurrentConnections.recordDataPoint(mb.startTime, ts, val, serverNameAttributeValue)
}

// RecordApacheRequestsDataPoint adds a data point to apache.requests metric.
func (mb *MetricsBuilder) RecordApacheRequestsDataPoint(ts pdata.Timestamp, val int64, serverNameAttributeValue string) {
	mb.metricApacheRequests.recordDataPoint(mb.startTime, ts, val, serverNameAttributeValue)
}

// RecordApacheScoreboardDataPoint adds a data point to apache.scoreboard metric.
func (mb *MetricsBuilder) RecordApacheScoreboardDataPoint(ts pdata.Timestamp, val int64, serverNameAttributeValue string, scoreboardStateAttributeValue string) {
	mb.metricApacheScoreboard.recordDataPoint(mb.startTime, ts, val, serverNameAttributeValue, scoreboardStateAttributeValue)
}

// RecordApacheTrafficDataPoint adds a data point to apache.traffic metric.
func (mb *MetricsBuilder) RecordApacheTrafficDataPoint(ts pdata.Timestamp, val int64, serverNameAttributeValue string) {
	mb.metricApacheTraffic.recordDataPoint(mb.startTime, ts, val, serverNameAttributeValue)
}

// RecordApacheUptimeDataPoint adds a data point to apache.uptime metric.
func (mb *MetricsBuilder) RecordApacheUptimeDataPoint(ts pdata.Timestamp, val int64, serverNameAttributeValue string) {
	mb.metricApacheUptime.recordDataPoint(mb.startTime, ts, val, serverNameAttributeValue)
}

// RecordApacheWorkersDataPoint adds a data point to apache.workers metric.
func (mb *MetricsBuilder) RecordApacheWorkersDataPoint(ts pdata.Timestamp, val int64, serverNameAttributeValue string, workersStateAttributeValue string) {
	mb.metricApacheWorkers.recordDataPoint(mb.startTime, ts, val, serverNameAttributeValue, workersStateAttributeValue)
}

// Reset resets metrics builder to its initial state. It should be used when external metrics source is restarted,
// and metrics builder should update its startTime and reset it's internal state accordingly.
func (mb *MetricsBuilder) Reset(options ...metricBuilderOption) {
	mb.startTime = pdata.NewTimestampFromTime(time.Now())
	for _, op := range options {
		op(mb)
	}
}

// NewMetricData creates new pdata.Metrics and sets the InstrumentationLibrary
// name on the ResourceMetrics.
func (mb *MetricsBuilder) NewMetricData() pdata.Metrics {
	md := pdata.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	ilm := rm.InstrumentationLibraryMetrics().AppendEmpty()
	ilm.InstrumentationLibrary().SetName("otelcol/apachereceiver")
	return md
}

// Attributes contains the possible metric attributes that can be used.
var Attributes = struct {
	// ScoreboardState (The state of a connection.)
	ScoreboardState string
	// ServerName (The name of the Apache HTTP server.)
	ServerName string
	// WorkersState (The state of workers.)
	WorkersState string
}{
	"state",
	"server_name",
	"state",
}

// A is an alias for Attributes.
var A = Attributes

// AttributeScoreboardState are the possible values that the attribute "scoreboard_state" can have.
var AttributeScoreboardState = struct {
	Open        string
	Waiting     string
	Starting    string
	Reading     string
	Sending     string
	Keepalive   string
	Dnslookup   string
	Closing     string
	Logging     string
	Finishing   string
	IdleCleanup string
}{
	"open",
	"waiting",
	"starting",
	"reading",
	"sending",
	"keepalive",
	"dnslookup",
	"closing",
	"logging",
	"finishing",
	"idle_cleanup",
}

// AttributeWorkersState are the possible values that the attribute "workers_state" can have.
var AttributeWorkersState = struct {
	Busy string
	Idle string
}{
	"busy",
	"idle",
}
