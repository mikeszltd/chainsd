/*
   Copyright 2014 Mikesz Ltd (https://www.mikesz.com/)

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

package outputs

import (
	"github.com/golang/glog"

	"github.com/mikeszltd/chainsd/channels"
	"github.com/mikeszltd/chainsd/generic"
	"github.com/mikeszltd/chainsd/messages"
	"github.com/streadway/amqp"
)

type Amqpout struct {
	route string
	uri   string

	AmqpConnection *amqp.Connection
	AmqpChannel    *amqp.Channel
	AmqpQueue      amqp.Queue

	Terminal chan messages.Message

	stop		chan bool
	quit		chan bool	
}

func NewAmqpout() *Amqpout {
	return &Amqpout{
		stop: make(chan bool),
		quit: make(chan bool),
	}		
}

func (amqpout *Amqpout) GetName() string {
	return "Amqpout"
}

func (amqpout *Amqpout) GetType() generic.LinkType {
	return generic.LinkTypeOutput
}

func (amqpout *Amqpout) Config(config *generic.ConfigurationLink) bool {
	amqpout.Terminal = channels.GetGlobalChannel(channels.TERMINAL_CHANNEL)

	var err error

	amqpout.route = (*config)["route"].(string)
	amqpout.uri = (*config)["uri"].(string)

	amqpout.AmqpConnection, err = amqp.Dial(amqpout.route)
	if err != nil {
		glog.Fatalf("Can't connect to amqp queue [%s]: %v", amqpout.route, err)
		panic(err)
	}
	amqpout.AmqpChannel, err = amqpout.AmqpConnection.Channel()
	if err != nil {
		glog.Fatalf("Can't open channel: %v", err)
		panic(err)
	}
	amqpout.AmqpChannel.Confirm(false)

	amqpout.AmqpQueue, err = amqpout.AmqpChannel.QueueDeclare(
		amqpout.route,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		glog.Fatalf("Can't declare queue [%v]: %v", amqpout.route, err)
		panic(err)
	}

	return true
}

func (amqpout *Amqpout) Do(input chan messages.Message) chan messages.Message {
	go func() {
		for {
			select {
				case m := <-input:
					amqpout.AmqpChannel.Publish(
						"",
						amqpout.route,
						false,
						false,
						amqp.Publishing{
							DeliveryMode: amqp.Persistent,
							ContentType:  "application/json",
							Body:         []byte(m.Content),
						},
					)
					amqpout.Terminal <- *messages.NewMessage("ok", "ok")
				case <- amqpout.stop:
					amqpout.quit <- true
					return
			}
		}
	}()

	return nil
}

func (amqpout *Amqpout) Close() {
	amqpout.stop <- true
	<- amqpout.quit
	
	if err := amqpout.AmqpConnection.Close(); err != nil {
		glog.Fatalf("AMQP connection close error: %s", err)		
		panic(err)
	}
	glog.Info("[amqpout] Closed")
}
