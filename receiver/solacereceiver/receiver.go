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
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/Azure/go-amqp"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenterror"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/obsreport"
	"go.uber.org/zap"
)

// compiler check on solaceReceiver if it is a Receiver
var _ component.Receiver = (*solaceTraceConsumer)(nil)

// solaceTraceConsumer uses azure AMQP to consume and handle AMQP messages from Solace.
// solaceTraceConsumer implements component.Receiver interface
type solaceTraceConsumer struct {
	instanceID        config.ComponentID
	cfg               *Config
	traceReceiver     *config.Receiver
	nextConsumer      consumer.Traces
	settings          component.ReceiverCreateSettings
	obsrecv           *obsreport.Receiver
	cancel            context.CancelFunc
	shutdownWaitGroup *sync.WaitGroup
	unmarchaller      *TracesUnmarshaler
}

//newReceiver create a new TracesReceiver (trace receiver) with the given parameters
func newTracesReceiver(rcfg *Config, receiverCreateSettings component.ReceiverCreateSettings, nextConsumer consumer.Traces) (component.TracesReceiver, error) {

	if nextConsumer == nil {
		receiverCreateSettings.Logger.Error("Next consumer in pipeline is null, stopping receiver")
		return nil, componenterror.ErrNilNextConsumer
	}

	if valErr := rcfg.Validate(); valErr != nil {
		receiverCreateSettings.Logger.Error("Invalid configuration", zap.Any("error", valErr))
		return nil, valErr
	}
	receiverCreateSettings.Logger.Info("Starting to listen on " + rcfg.Queue)

	var u TracesUnmarshaler = SolaceMessageUnmarshaller{Logger: receiverCreateSettings.Logger}
	return &solaceTraceConsumer{
		instanceID:   rcfg.ID(),
		cfg:          rcfg,
		nextConsumer: nextConsumer,
		settings:     receiverCreateSettings,
		obsrecv: obsreport.NewReceiver(obsreport.ReceiverSettings{
			ReceiverID:             rcfg.ID(),
			Transport:              "solace",
			ReceiverCreateSettings: receiverCreateSettings,
		}),
		unmarchaller: &u,
	}, nil

}

func (s *solaceTraceConsumer) Start(ctx context.Context, host component.Host) error {

	s.settings.Logger.Info(" ---------- hello ---------- ")
	var cancelableContext context.Context
	cancelableContext, s.cancel = context.WithCancel(context.Background())
	s.settings.Logger.Info("solaceTraceConsumer Start " + s.instanceID.String() + " " + s.cfg.ID().String() + s.cfg.Broker[0])
	conConf, err := s.createConnectionConfig()

	if err != (nil) {
		// exit on config errors
		return err
	}

	s.shutdownWaitGroup = new(sync.WaitGroup)
	s.shutdownWaitGroup.Add(1)
	go s.connectAndReceive(cancelableContext, conConf)

	s.settings.Logger.Info(" solace receiver started")
	//RecordNeedRestart()
	return nil
}

func (s *solaceTraceConsumer) Shutdown(ctx context.Context) error {
	s.settings.Logger.Info("shutdown signal received ")

	s.cancel()

	s.settings.Logger.Info("solaceTraceConsumer shutdown, waiting for waiting group" + s.instanceID.String() + " " + s.cfg.ID().String() + s.cfg.Broker[0] + " transport:")

	s.shutdownWaitGroup.Wait()
	s.settings.Logger.Info("solaceTraceConsumer shutdown  wait group closed ")

	return nil
}

type connectConfig struct {
	addr       *string
	saslConfig *amqp.ConnOption
	tlsConfig  *tls.Config
}

func (s *solaceTraceConsumer) createConnectionConfig() (*connectConfig, error) {
	s.settings.Logger.Info("createConnectionConfig")
	broker := s.cfg.Broker[0]
	amqpHostAddress := fmt.Sprintf("amqps"+"://%s", broker)
	s.settings.Logger.Info("configured to connect to " + amqpHostAddress)

	tlsConfig, err := s.cfg.Transport.TLS.LoadTLSConfig()
	if err != nil {
		return (nil), err
	}
	s.settings.Logger.Info("tls config loaded")
	saslConnOption, authErr := ConfigureAuthentication(s.cfg)
	if authErr != (nil) {
		s.settings.Logger.Info("problem with auth config")
		return nil, authErr
	}
	s.settings.Logger.Info("auth config loaded")

	return &connectConfig{
		addr:       &amqpHostAddress,
		tlsConfig:  tlsConfig,
		saslConfig: &saslConnOption}, (nil)

}

func (s *solaceTraceConsumer) connectAndReceive(ctx context.Context, cConf *connectConfig) error {

	s.settings.Logger.Warn("enter connection loop")

	// Don't create done channel on every iteration.
	doneChan := (ctx).Done()
reconnectionLoop:
	for {
		// TODO from config take timeout between attempts
		s.settings.Logger.Warn("waiting for 1 sec before connection attempt")

		s.sleep(ctx, 1*time.Second)

		select {
		case <-doneChan: //(*ctx).Done():
			s.settings.Logger.Warn("shutdown a reconnection loop")
			break reconnectionLoop
		default:
			s.settings.Logger.Info(" ----default -----" + *cConf.addr + "    InsecureSkipVerify=" + strconv.FormatBool(cConf.tlsConfig.InsecureSkipVerify))
			client, err := amqp.Dial(*cConf.addr, *cConf.saslConfig,
				amqp.ConnTLSConfig(cConf.tlsConfig)) //, amqp.ConnTLSConfig(cConf.tlsConfig))

			if err != nil {
				s.settings.Logger.Info(" ----connection failed ", zap.Error(err))
				continue
			}
			s.settings.Logger.Info(" ----connected successful_ ")
			// TODO timeout + return on error

			s.settings.Logger.Info(" ----creating session ")
			session, err := client.NewSession()
			if err != nil {
				s.settings.Logger.Warn(" ----session creation failed", zap.Error(err))
				RecordFailedReconnection()
				continue
			}
			s.settings.Logger.Info(" ----session created ")
			receiver, err := session.NewReceiver(
				amqp.LinkSourceAddress(s.cfg.Queue),
				amqp.LinkCredit(s.cfg.MaxUnaked),
			)
			if err != nil {
				s.settings.Logger.Info(" ---- amqp receiver creation failed", zap.Error(err))
				continue
			}
			s.settings.Logger.Info(" ---- amqp receiver  created ")

			// when connected perform blocking receive untel eeror
			s.settings.Logger.Info(" ---- starting message receipt")
			err = s.receiveMessages(&ctx, receiver)
			if err != nil {
				// TODO add code for checking of version incompatibility error and stop reconnection loop
				// TODO report need upgrade stats
				s.settings.Logger.Info(" ---- message receipt failed")
				continue
			}

			defer func(client *amqp.Client) {
				if client != nil {
					s.settings.Logger.Info(" ----defer--- CLOSING AMQP client")
					_ = client.Close()
					s.settings.Logger.Info(" AMQP client is closed")
					// TODO deal with an error
				} else {
					s.settings.Logger.Info(" ----defer--- skip closing AMQP client")
				}
			}(client)
			// close receiver
			defer func(r *amqp.Receiver) {
				s.settings.Logger.Info(" ----defer---- closing receiver ")
				if r != nil {
					ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)

					err := r.Close(ctx)
					if err != nil {
						return
					}
					cancel()
				} //else {
				//	s.settings.Logger.Info(" ----defer---- SKIP closing receiver ")
				//}
			}(receiver)

		}
		//	s.settings.Logger.Info("  reconnection loop select is exited")
	}

	s.settings.Logger.Info("  reconnection loop 'for' is exited")

	defer func() {
		s.settings.Logger.Info(" --defer-- finishing wait group")
		s.shutdownWaitGroup.Done()
	}()

	return nil
}

func (s *solaceTraceConsumer) receiveMessages(ctx *context.Context, receiver *amqp.Receiver) error {

	for {

		msg, err := receiver.Receive(*ctx)
		if err != nil {
			return err
		}

		var td, e = (*s.unmarchaller).Unmarshal(msg)
		if e != nil {

			// TODO handle unmarshalling error  right way

			//
			//
			//NewVersionIncompatibleUnmarshallingError
			s.settings.Logger.Error("unmarshalling error", zap.Error(e))

			RecordFatalUnmarshallingError()
			// take it from a queue anyway

			if IsUnknownMessageError(e) {
				er := amqp.Error{
					Condition:   amqp.ErrorDecodeError,
					Description: e.Error(),
				}
				s.settings.Logger.Error("incompatible message error, Rejecting Message")
				err := receiver.RejectMessage(*ctx, msg, &er)
				if err != nil {
					return err
				}
				continue
			}

			if IsVersionIncompatibleUnmarshallingError(e) {
				er := amqp.Error{
					Condition:   amqp.ErrorDecodeError,
					Description: e.Error(),
				}
				s.settings.Logger.Error("incompatible version error, Rejecting Message")
				err := receiver.RejectMessage(*ctx, msg, &er)
				if err != nil {
					return err
				}
				continue
			}

			s.settings.Logger.Error("unmarshalling error, Accepting Message")
			err := receiver.AcceptMessage(*ctx, msg)
			if err != nil {
				return err
			}
			continue
		}

		// forward to next consumer
		forwardErr := s.nextConsumer.ConsumeTraces(*ctx, *td)
		if forwardErr != nil {
			s.settings.Logger.Error("can't forward traces to the next receiver, rejecting the message", zap.Error(forwardErr))
			err := receiver.RejectMessage(*ctx, msg, nil)
			if err != nil {
				return err
			}
		}

		s.settings.Logger.Debug("message successful processed and accepted")
		// accept message, done
		e = receiver.AcceptMessage(*ctx, msg)
		if e != nil {
			return e
		}

		select {
		case <-(*ctx).Done():
			s.settings.Logger.Warn("shutdown a receiver loop")
			return nil
		default:
			continue

		}
	}

}

func (s *solaceTraceConsumer) sleep(ctx context.Context, d time.Duration) {
	s.settings.Logger.Warn("-----sleep-----" + d.String())
	timer := time.NewTimer(d)
	select {
	case <-(ctx).Done():
		s.settings.Logger.Warn("-----sleep interrupted-----" + d.String())
		if !timer.Stop() {
			<-timer.C
		}
		s.settings.Logger.Warn("-----sleep exit 1-----" + d.String())
	case <-timer.C:
		s.settings.Logger.Warn("-----sleep exit 2-----" + d.String())
	}
}
