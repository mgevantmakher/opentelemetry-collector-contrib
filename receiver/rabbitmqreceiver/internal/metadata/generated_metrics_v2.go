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

// MetricsSettings provides settings for rabbitmqreceiver metrics.
type MetricsSettings struct {
	RabbitmqConsumerCount       MetricSettings `mapstructure:"rabbitmq.consumer.count"`
	RabbitmqMessageAcknowledged MetricSettings `mapstructure:"rabbitmq.message.acknowledged"`
	RabbitmqMessageCurrent      MetricSettings `mapstructure:"rabbitmq.message.current"`
	RabbitmqMessageDelivered    MetricSettings `mapstructure:"rabbitmq.message.delivered"`
	RabbitmqMessageDropped      MetricSettings `mapstructure:"rabbitmq.message.dropped"`
	RabbitmqMessagePublished    MetricSettings `mapstructure:"rabbitmq.message.published"`
}

func DefaultMetricsSettings() MetricsSettings {
	return MetricsSettings{
		RabbitmqConsumerCount: MetricSettings{
			Enabled: true,
		},
		RabbitmqMessageAcknowledged: MetricSettings{
			Enabled: true,
		},
		RabbitmqMessageCurrent: MetricSettings{
			Enabled: true,
		},
		RabbitmqMessageDelivered: MetricSettings{
			Enabled: true,
		},
		RabbitmqMessageDropped: MetricSettings{
			Enabled: true,
		},
		RabbitmqMessagePublished: MetricSettings{
			Enabled: true,
		},
	}
}

type metricRabbitmqConsumerCount struct {
	data     pdata.Metric   // data buffer for generated metric.
	settings MetricSettings // metric settings provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills rabbitmq.consumer.count metric with initial data.
func (m *metricRabbitmqConsumerCount) init() {
	m.data.SetName("rabbitmq.consumer.count")
	m.data.SetDescription("The number of consumers currently reading from the queue.")
	m.data.SetUnit("{consumers}")
	m.data.SetDataType(pdata.MetricDataTypeSum)
	m.data.Sum().SetIsMonotonic(false)
	m.data.Sum().SetAggregationTemporality(pdata.MetricAggregationTemporalityCumulative)
}

func (m *metricRabbitmqConsumerCount) recordDataPoint(start pdata.Timestamp, ts pdata.Timestamp, val int64) {
	if !m.settings.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntVal(val)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricRabbitmqConsumerCount) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricRabbitmqConsumerCount) emit(metrics pdata.MetricSlice) {
	if m.settings.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricRabbitmqConsumerCount(settings MetricSettings) metricRabbitmqConsumerCount {
	m := metricRabbitmqConsumerCount{settings: settings}
	if settings.Enabled {
		m.data = pdata.NewMetric()
		m.init()
	}
	return m
}

type metricRabbitmqMessageAcknowledged struct {
	data     pdata.Metric   // data buffer for generated metric.
	settings MetricSettings // metric settings provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills rabbitmq.message.acknowledged metric with initial data.
func (m *metricRabbitmqMessageAcknowledged) init() {
	m.data.SetName("rabbitmq.message.acknowledged")
	m.data.SetDescription("The number of messages acknowledged by consumers.")
	m.data.SetUnit("{messages}")
	m.data.SetDataType(pdata.MetricDataTypeSum)
	m.data.Sum().SetIsMonotonic(true)
	m.data.Sum().SetAggregationTemporality(pdata.MetricAggregationTemporalityCumulative)
}

func (m *metricRabbitmqMessageAcknowledged) recordDataPoint(start pdata.Timestamp, ts pdata.Timestamp, val int64) {
	if !m.settings.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntVal(val)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricRabbitmqMessageAcknowledged) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricRabbitmqMessageAcknowledged) emit(metrics pdata.MetricSlice) {
	if m.settings.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricRabbitmqMessageAcknowledged(settings MetricSettings) metricRabbitmqMessageAcknowledged {
	m := metricRabbitmqMessageAcknowledged{settings: settings}
	if settings.Enabled {
		m.data = pdata.NewMetric()
		m.init()
	}
	return m
}

type metricRabbitmqMessageCurrent struct {
	data     pdata.Metric   // data buffer for generated metric.
	settings MetricSettings // metric settings provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills rabbitmq.message.current metric with initial data.
func (m *metricRabbitmqMessageCurrent) init() {
	m.data.SetName("rabbitmq.message.current")
	m.data.SetDescription("The total number of messages currently in the queue.")
	m.data.SetUnit("{messages}")
	m.data.SetDataType(pdata.MetricDataTypeSum)
	m.data.Sum().SetIsMonotonic(false)
	m.data.Sum().SetAggregationTemporality(pdata.MetricAggregationTemporalityCumulative)
	m.data.Sum().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricRabbitmqMessageCurrent) recordDataPoint(start pdata.Timestamp, ts pdata.Timestamp, val int64, messageStateAttributeValue string) {
	if !m.settings.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntVal(val)
	dp.Attributes().Insert(A.MessageState, pdata.NewAttributeValueString(messageStateAttributeValue))
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricRabbitmqMessageCurrent) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricRabbitmqMessageCurrent) emit(metrics pdata.MetricSlice) {
	if m.settings.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricRabbitmqMessageCurrent(settings MetricSettings) metricRabbitmqMessageCurrent {
	m := metricRabbitmqMessageCurrent{settings: settings}
	if settings.Enabled {
		m.data = pdata.NewMetric()
		m.init()
	}
	return m
}

type metricRabbitmqMessageDelivered struct {
	data     pdata.Metric   // data buffer for generated metric.
	settings MetricSettings // metric settings provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills rabbitmq.message.delivered metric with initial data.
func (m *metricRabbitmqMessageDelivered) init() {
	m.data.SetName("rabbitmq.message.delivered")
	m.data.SetDescription("The number of messages delivered to consumers.")
	m.data.SetUnit("{messages}")
	m.data.SetDataType(pdata.MetricDataTypeSum)
	m.data.Sum().SetIsMonotonic(true)
	m.data.Sum().SetAggregationTemporality(pdata.MetricAggregationTemporalityCumulative)
}

func (m *metricRabbitmqMessageDelivered) recordDataPoint(start pdata.Timestamp, ts pdata.Timestamp, val int64) {
	if !m.settings.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntVal(val)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricRabbitmqMessageDelivered) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricRabbitmqMessageDelivered) emit(metrics pdata.MetricSlice) {
	if m.settings.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricRabbitmqMessageDelivered(settings MetricSettings) metricRabbitmqMessageDelivered {
	m := metricRabbitmqMessageDelivered{settings: settings}
	if settings.Enabled {
		m.data = pdata.NewMetric()
		m.init()
	}
	return m
}

type metricRabbitmqMessageDropped struct {
	data     pdata.Metric   // data buffer for generated metric.
	settings MetricSettings // metric settings provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills rabbitmq.message.dropped metric with initial data.
func (m *metricRabbitmqMessageDropped) init() {
	m.data.SetName("rabbitmq.message.dropped")
	m.data.SetDescription("The number of messages dropped as unroutable.")
	m.data.SetUnit("{messages}")
	m.data.SetDataType(pdata.MetricDataTypeSum)
	m.data.Sum().SetIsMonotonic(true)
	m.data.Sum().SetAggregationTemporality(pdata.MetricAggregationTemporalityCumulative)
}

func (m *metricRabbitmqMessageDropped) recordDataPoint(start pdata.Timestamp, ts pdata.Timestamp, val int64) {
	if !m.settings.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntVal(val)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricRabbitmqMessageDropped) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricRabbitmqMessageDropped) emit(metrics pdata.MetricSlice) {
	if m.settings.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricRabbitmqMessageDropped(settings MetricSettings) metricRabbitmqMessageDropped {
	m := metricRabbitmqMessageDropped{settings: settings}
	if settings.Enabled {
		m.data = pdata.NewMetric()
		m.init()
	}
	return m
}

type metricRabbitmqMessagePublished struct {
	data     pdata.Metric   // data buffer for generated metric.
	settings MetricSettings // metric settings provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills rabbitmq.message.published metric with initial data.
func (m *metricRabbitmqMessagePublished) init() {
	m.data.SetName("rabbitmq.message.published")
	m.data.SetDescription("The number of messages published to a queue.")
	m.data.SetUnit("{messages}")
	m.data.SetDataType(pdata.MetricDataTypeSum)
	m.data.Sum().SetIsMonotonic(true)
	m.data.Sum().SetAggregationTemporality(pdata.MetricAggregationTemporalityCumulative)
}

func (m *metricRabbitmqMessagePublished) recordDataPoint(start pdata.Timestamp, ts pdata.Timestamp, val int64) {
	if !m.settings.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntVal(val)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricRabbitmqMessagePublished) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricRabbitmqMessagePublished) emit(metrics pdata.MetricSlice) {
	if m.settings.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricRabbitmqMessagePublished(settings MetricSettings) metricRabbitmqMessagePublished {
	m := metricRabbitmqMessagePublished{settings: settings}
	if settings.Enabled {
		m.data = pdata.NewMetric()
		m.init()
	}
	return m
}

// MetricsBuilder provides an interface for scrapers to report metrics while taking care of all the transformations
// required to produce metric representation defined in metadata and user settings.
type MetricsBuilder struct {
	startTime                         pdata.Timestamp
	metricRabbitmqConsumerCount       metricRabbitmqConsumerCount
	metricRabbitmqMessageAcknowledged metricRabbitmqMessageAcknowledged
	metricRabbitmqMessageCurrent      metricRabbitmqMessageCurrent
	metricRabbitmqMessageDelivered    metricRabbitmqMessageDelivered
	metricRabbitmqMessageDropped      metricRabbitmqMessageDropped
	metricRabbitmqMessagePublished    metricRabbitmqMessagePublished
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
		startTime:                         pdata.NewTimestampFromTime(time.Now()),
		metricRabbitmqConsumerCount:       newMetricRabbitmqConsumerCount(settings.RabbitmqConsumerCount),
		metricRabbitmqMessageAcknowledged: newMetricRabbitmqMessageAcknowledged(settings.RabbitmqMessageAcknowledged),
		metricRabbitmqMessageCurrent:      newMetricRabbitmqMessageCurrent(settings.RabbitmqMessageCurrent),
		metricRabbitmqMessageDelivered:    newMetricRabbitmqMessageDelivered(settings.RabbitmqMessageDelivered),
		metricRabbitmqMessageDropped:      newMetricRabbitmqMessageDropped(settings.RabbitmqMessageDropped),
		metricRabbitmqMessagePublished:    newMetricRabbitmqMessagePublished(settings.RabbitmqMessagePublished),
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
	mb.metricRabbitmqConsumerCount.emit(metrics)
	mb.metricRabbitmqMessageAcknowledged.emit(metrics)
	mb.metricRabbitmqMessageCurrent.emit(metrics)
	mb.metricRabbitmqMessageDelivered.emit(metrics)
	mb.metricRabbitmqMessageDropped.emit(metrics)
	mb.metricRabbitmqMessagePublished.emit(metrics)
}

// RecordRabbitmqConsumerCountDataPoint adds a data point to rabbitmq.consumer.count metric.
func (mb *MetricsBuilder) RecordRabbitmqConsumerCountDataPoint(ts pdata.Timestamp, val int64) {
	mb.metricRabbitmqConsumerCount.recordDataPoint(mb.startTime, ts, val)
}

// RecordRabbitmqMessageAcknowledgedDataPoint adds a data point to rabbitmq.message.acknowledged metric.
func (mb *MetricsBuilder) RecordRabbitmqMessageAcknowledgedDataPoint(ts pdata.Timestamp, val int64) {
	mb.metricRabbitmqMessageAcknowledged.recordDataPoint(mb.startTime, ts, val)
}

// RecordRabbitmqMessageCurrentDataPoint adds a data point to rabbitmq.message.current metric.
func (mb *MetricsBuilder) RecordRabbitmqMessageCurrentDataPoint(ts pdata.Timestamp, val int64, messageStateAttributeValue string) {
	mb.metricRabbitmqMessageCurrent.recordDataPoint(mb.startTime, ts, val, messageStateAttributeValue)
}

// RecordRabbitmqMessageDeliveredDataPoint adds a data point to rabbitmq.message.delivered metric.
func (mb *MetricsBuilder) RecordRabbitmqMessageDeliveredDataPoint(ts pdata.Timestamp, val int64) {
	mb.metricRabbitmqMessageDelivered.recordDataPoint(mb.startTime, ts, val)
}

// RecordRabbitmqMessageDroppedDataPoint adds a data point to rabbitmq.message.dropped metric.
func (mb *MetricsBuilder) RecordRabbitmqMessageDroppedDataPoint(ts pdata.Timestamp, val int64) {
	mb.metricRabbitmqMessageDropped.recordDataPoint(mb.startTime, ts, val)
}

// RecordRabbitmqMessagePublishedDataPoint adds a data point to rabbitmq.message.published metric.
func (mb *MetricsBuilder) RecordRabbitmqMessagePublishedDataPoint(ts pdata.Timestamp, val int64) {
	mb.metricRabbitmqMessagePublished.recordDataPoint(mb.startTime, ts, val)
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
	ilm.InstrumentationLibrary().SetName("otelcol/rabbitmqreceiver")
	return md
}

// Attributes contains the possible metric attributes that can be used.
var Attributes = struct {
	// MessageState (The state of messages in a queue.)
	MessageState string
	// RabbitmqNodeName (The name of the RabbitMQ node.)
	RabbitmqNodeName string
	// RabbitmqQueueName (The name of the RabbitMQ queue.)
	RabbitmqQueueName string
	// RabbitmqVhostName (The name of the RabbitMQ vHost.)
	RabbitmqVhostName string
}{
	"state",
	"rabbitmq.node.name",
	"rabbitmq.queue.name",
	"rabbitmq.vhost.name",
}

// A is an alias for Attributes.
var A = Attributes

// AttributeMessageState are the possible values that the attribute "message.state" can have.
var AttributeMessageState = struct {
	Ready          string
	Unacknowledged string
}{
	"ready",
	"unacknowledged",
}
