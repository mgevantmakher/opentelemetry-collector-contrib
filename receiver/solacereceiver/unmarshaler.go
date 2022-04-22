// Copyright  The OpenTelemetry Authors
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

package solacereceiver

import (
	"reflect"
	"strings"

	"github.com/Azure/go-amqp"
	"go.opentelemetry.io/collector/model/pdata"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	model_v1 "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/solacereceiver/model/pdata/v1"
)

// TracesUnmarshaler deserializes the message body.
type TracesUnmarshaler interface {
	// Unmarshal the amqp-message into traces.
	// Only valid traces are produced or error is returned
	Unmarshal(message *amqp.Message) (*pdata.Traces, error)

	unmarshalV1(message *amqp.Message) (*pdata.Traces, error)
	unmarshalV2(message *amqp.Message) (*pdata.Traces, error)
}

type SolaceMessageUnmarshaller struct {
	Logger *zap.Logger
}

func (u SolaceMessageUnmarshaller) Unmarshal(message *amqp.Message) (*pdata.Traces, error) {
	// for now just v1 exists, if strings.Contains(*to, "receive/v1") { will be moved here
	return u.unmarshalV1(message)
}

func (u SolaceMessageUnmarshaller) unmarshalV1(message *amqp.Message) (*pdata.Traces, error) {
	messageProperties := message.Properties
	var to = messageProperties.To
	// # substituted with _
	if strings.Contains(*to, "_telemetry/trace/broker/v") {
		if strings.Contains(*to, "v1/receive") {
			return u.unmarshalReceiveSpanV1(message)
		} else {
			return nil, NewVersionIncompatibleUnmarshallingError("unknown trace message type and version", (nil))
		}

	}
	u.Logger.Info("unknown trace message type and version " + *to)
	return nil, NewUnknownMessageError("unknown trace message type and version", (nil))
}

func (SolaceMessageUnmarshaller) unmarshalV2(message *amqp.Message) (*pdata.Traces, error) {
	// not implemented yet, just a placeholder for future use
	return &pdata.Traces{}, nil
}

const clientSpanName = "(topic) receive"

//https://solacedotcom.slack.com/archives/C025E138JDA/p1631654565079700
const serviceName = "receive"

func (u SolaceMessageUnmarshaller) unmarshalReceiveSpanV1(message *amqp.Message) (*pdata.Traces, error) {
	rSpan := &model_v1.SpanData{}

	if err := proto.Unmarshal(message.GetData(), rSpan); err != nil {
		// TODO proper error handling

		return nil, err
	}
	td := pdata.NewTraces()
	rs := td.ResourceSpans().AppendEmpty()
	rs.Resource().Attributes().InsertString("service.name", serviceName)
	instrLibrarySpans := rs.InstrumentationLibrarySpans().AppendEmpty()

	instrLibrarySpans.InstrumentationLibrary().SetName("com.solace.pubsubplus")
	instrLibrarySpans.InstrumentationLibrary().SetVersion("1")
	clientSpan := instrLibrarySpans.Spans().AppendEmpty()
	clientSpan.SetName(clientSpanName)

	var spanId [8]byte
	copy(spanId[:8], rSpan.SpanId)
	clientSpan.SetSpanID(pdata.NewSpanID(spanId))

	var traceId [16]byte
	copy(traceId[:16], rSpan.TraceId)
	clientSpan.SetTraceID(pdata.NewTraceID(traceId))

	// conditional parent-span-id
	if len(rSpan.ParentSpanId) == 8 {
		var parentSpanId [8]byte
		copy(parentSpanId[:8], rSpan.ParentSpanId)
		clientSpan.SetParentSpanID(pdata.NewSpanID(parentSpanId))
	} else {
		u.Logger.Error("parent span id is not length 8: "+string(rSpan.ParentSpanId[:])+" ", zap.Int("length", len(rSpan.ParentSpanId)))
	}
	//SPAN_KIND_CONSUMER == 5
	clientSpan.SetKind(5)
	clientSpan.SetStartTimestamp(pdata.Timestamp(rSpan.GetStartTimeUnixNano()))
	clientSpan.SetEndTimestamp(pdata.Timestamp(rSpan.GetEndTimeUnixNano()))

	//spanEvent := clientSpan.Events().AppendEmpty()
	//spanEvent.SetName("eventName")
	//spanEvent.SetTimestamp(pdata.Timestamp(rSpan.StartTimeUnixNano))

	u.unmarshalAttributesV1(&clientSpan, rSpan)
	u.unmarshalEventsV1(&clientSpan, rSpan)

	u.Logger.Info("Unmarshalling receive span with span " + rSpan.String())

	//switch solaceTelemetryObject := solaceTelemetryData.Object.(type) {
	//case *model_v1.TelemetryData_TelemetryDataObject_ReceiveSpan:
	//	rSpan:=solaceTelemetryObject.ReceiveSpan
	//
	//	rSpan.
	//
	//	fmt.Printf("Copy Operation start: %d, end : %d\n", solaceTelemetryObject.ReceiveSpan)
	//
	//default:
	//	fmt.Println("No matching operations")
	//}

	return &td, nil
}

const adEventName = "AD Receive"
const messagingDestinationEvtKey = "messaging.destination"
const messagingSolaceDestinationTypeEvtKey = "messaging.solace.destination_type"
const statusCodeEvtKey = "Status.Code"
const statusMessageEvtKey = "Status.Message"

// converts and injects all events into pdata.SpanEventSlice
func (u SolaceMessageUnmarshaller) unmarshalEventsV1(s *pdata.Span, spanData *model_v1.SpanData) {

	evts := s.Events()
	enqEvents := spanData.EnqueueEvents
	for _, enqEvent := range enqEvents {
		evt := evts.AppendEmpty()
		evt.SetName(enqEvent.GetQueueName() + " enqueue")
		evt.SetTimestamp(pdata.Timestamp(enqEvent.TimeUnixNano))
		// attributes
		u.Logger.Info("-DESTINATION TYPE: " + reflect.TypeOf(enqEvent.Dest).String())
		d := getDestination(enqEvent.Dest)
		if d != nil {
			u.Logger.Info("-HAS-getDestination: " + *d.name)
			evt.Attributes().InsertString(messagingDestinationEvtKey, *d.name)
			evt.Attributes().InsertString(messagingSolaceDestinationTypeEvtKey, *d.name)
		}
		if enqEvent.GetErrorDescription() != "" {
			evt.Attributes().InsertString(statusMessageEvtKey, enqEvent.GetErrorDescription())
			evt.Attributes().InsertString(statusCodeEvtKey, "Err")
		}

	}

	adReceiveEvent := spanData.AdReceiveEvent
	if adReceiveEvent != nil {
		evt := evts.AppendEmpty()
		evt.SetName(adEventName)
		evt.SetTimestamp(pdata.Timestamp(adReceiveEvent.TimeUnixNano))
	}

	// TODO: add transaction event decoding

}

// attributes
const systemAttrKey = "messaging.system"
const systemAttrValue = "SolacePubSub+"
const operationAttrKey = "messaging.operation"
const operationAttrValue = "receive"

const protocolAttrKey = "messaging.protocol"

const protocolVersionAttrKey = "messaging.protocol_version"
const messageIDAttrKey = "messaging.message_id"
const conversationIDAttrKey = "messaging.conversation_id"
const payloadSizeBytesAttrKey = "messaging.message_payload_size_bytes"
const routerNameAttrKey = "messaging.solace.router_name"
const messageVpnNameAttrKey = "messaging.solace.message_vpn_name"
const clientUsernameAttrKey = "messaging.solace.client_username"
const clientNameAttrKey = "messaging.solace.client_name"
const replicationGroupMessageIDAttrKey = "messaging.solace.replication_group_message_id"
const priorityAttrKey = "messaging.solace.priority"
const ttlAttrKey = "messaging.solace.ttl"
const dmqEligibleAttrKey = "messaging.solace.dmq_eligible"
const droppedEnqueueEventsSuccessAttrKey = "messaging.solace.dropped_enqueue_events_success"
const droppedEnqueueFailedEventsAttrKey = "messaging.solace.dropped_enqueue_failed_events"
const hostIPAttrKey = "net.host.ip"
const hostPortAttrKey = "net.host.port"
const peerIPAttrKey = "net.peer.ip"
const peerPortAttrKey = "net.peer.port"
const userPropertiesPrefixAttrKey = "messaging.solace.user_properteis."

// omitted for now
//const AttrKey="messaging.solace.user_data"

// converts and injects all standard solace and user properties key value pairs into otel pdata.AttributeMap
func (u SolaceMessageUnmarshaller) unmarshalAttributesV1(s *pdata.Span, spanData *model_v1.SpanData) {
	atrMap := s.Attributes()
	// first add default attributes with constant values
	// solaceUserProperties map[string]*model_v1.SpanData_UserPropertyValue
	solaceUserProperties := spanData.GetUserProperties()
	atrMap.InsertString(systemAttrKey, systemAttrValue)
	atrMap.InsertString(operationAttrKey, operationAttrValue)
	atrMap.InsertString(protocolAttrKey, spanData.Protocol)
	if spanData.ProtocolVersion != nil {
		atrMap.InsertString(protocolVersionAttrKey, *spanData.ProtocolVersion)
	}
	if spanData.ApplicationMessageId != nil {
		atrMap.InsertString(messageIDAttrKey, *spanData.ApplicationMessageId)
	}
	if spanData.CorrelationId != nil {
		atrMap.InsertString(conversationIDAttrKey, *spanData.CorrelationId)
	}
	atrMap.InsertInt(payloadSizeBytesAttrKey, int64(spanData.BinaryAttachmentSize+spanData.XmlAttachmentSize+spanData.MetadataSize))
	if spanData.RouterName != nil {
		atrMap.InsertString(routerNameAttrKey, *spanData.RouterName)
	}
	if spanData.MessageVpnName != nil {
		atrMap.InsertString(messageVpnNameAttrKey, *spanData.MessageVpnName)
	}

	atrMap.InsertString(clientUsernameAttrKey, spanData.ClientUsername)

	atrMap.InsertString(clientNameAttrKey, spanData.ClientName)

	atrMap.InsertBytes(replicationGroupMessageIDAttrKey, spanData.ReplicationGroupMessageId)

	if spanData.Priority != nil {
		atrMap.InsertInt(priorityAttrKey, int64(*spanData.Priority))
	}
	if spanData.Ttl != nil {
		atrMap.InsertInt(ttlAttrKey, *spanData.Ttl)
	}
	atrMap.InsertBool(dmqEligibleAttrKey, spanData.DmqEligible)
	atrMap.InsertInt(droppedEnqueueEventsSuccessAttrKey, int64(spanData.DroppedEnqeueueEventsSuccess))
	atrMap.InsertInt(droppedEnqueueFailedEventsAttrKey, int64(spanData.DroppedEnqueueEventsFailed))

	hostIpLen := len(spanData.HostIp)
	if hostIpLen == 4 || hostIpLen == 16 {
		atrMap.InsertString(hostIPAttrKey, string(spanData.HostIp[:]))
	} else {
		u.Logger.Warn("Host ip attribute has an illegallength", zap.Int("length", hostIpLen))
	}
	atrMap.InsertInt(hostPortAttrKey, int64(spanData.HostPort))

	peerIpLen := len(spanData.HostIp)
	if peerIpLen == 4 || peerIpLen == 16 {
		atrMap.InsertString(peerIPAttrKey, string(spanData.PeerIp[:]))
	} else {
		u.Logger.Warn("Peer ip attribute has an illegal length", zap.Int("length", peerIpLen))
	}
	atrMap.InsertInt(peerPortAttrKey, int64(spanData.PeerPort))
	// messaging.solace.user_data omitted

	// user properties
	for key, value := range solaceUserProperties {
		//u.Logger.Info("reading attribute: " + key)
		if value != nil {
			u.insertUserPropertyV1(&atrMap, key, value.Value)
		}
	}
}

type destination struct {
	name *string
	kind string
}

const queueKind = "queue"
const topicEndpoindKind = "topic-endpoint"

func getDestination(value interface{}) *destination {
	switch v := value.(type) {
	case model_v1.SpanData_EnqueueEvent_TopicEndpointName:
		return &destination{kind: topicEndpoindKind, name: &v.TopicEndpointName}

	case model_v1.SpanData_EnqueueEvent_QueueName:
		return &destination{kind: queueKind, name: &v.QueueName}
	}
	return nil
}

func (u SolaceMessageUnmarshaller) insertUserPropertyV1(toMap *pdata.AttributeMap, key string, value interface{}) {
	k := userPropertiesPrefixAttrKey + key
	switch v := value.(type) {
	case *model_v1.SpanData_UserPropertyValue_NullValue:
		toMap.InsertNull(k)

	case *model_v1.SpanData_UserPropertyValue_BoolValue:
		toMap.InsertBool(k, v.BoolValue)

	case *model_v1.SpanData_UserPropertyValue_DoubleValue:
		toMap.InsertDouble(k, v.DoubleValue)

	case *model_v1.SpanData_UserPropertyValue_ByteArrayValue:
		toMap.InsertBytes(k, v.ByteArrayValue)

	case *model_v1.SpanData_UserPropertyValue_FloatValue:
		toMap.InsertDouble(k, float64(v.FloatValue))

	case *model_v1.SpanData_UserPropertyValue_Int8Value:
		toMap.InsertInt(k, int64(v.Int8Value))

	case *model_v1.SpanData_UserPropertyValue_Int16Value:
		toMap.InsertInt(k, int64(v.Int16Value))

	case *model_v1.SpanData_UserPropertyValue_Int32Value:
		toMap.InsertInt(k, int64(v.Int32Value))

	case *model_v1.SpanData_UserPropertyValue_Int64Value:
		toMap.InsertInt(k, v.Int64Value)

	case *model_v1.SpanData_UserPropertyValue_Uint8Value:
		toMap.InsertInt(k, int64(v.Uint8Value))

	case *model_v1.SpanData_UserPropertyValue_Uint16Value:
		toMap.InsertInt(k, int64(v.Uint16Value))

	case *model_v1.SpanData_UserPropertyValue_Uint32Value:
		toMap.InsertInt(k, int64(v.Uint32Value))

	case *model_v1.SpanData_UserPropertyValue_Uint64Value:
		toMap.InsertInt(k, int64(v.Uint64Value))

	case *model_v1.SpanData_UserPropertyValue_StringValue:
		toMap.InsertString(k, v.StringValue)

	case *model_v1.SpanData_UserPropertyValue_DestinationValue:
		toMap.InsertString(k, v.DestinationValue)

	default:
		u.Logger.Warn("unknown type of user property" + reflect.TypeOf(v).String())
	}

}
