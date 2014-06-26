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

package generators

import (
	"github.com/golang/glog"

	"github.com/mikeszltd/chainsd/generic"
	"github.com/mikeszltd/chainsd/messages"
	"github.com/streadway/amqp"
)

type Amqpin struct {
	route string
	uri   string

	AmqpConnection *amqp.Connection
	AmqpChannel    *amqp.Channel
	AmqpQueue      amqp.Queue

	stop		chan bool
	quit		chan bool
}

func NewAmqpin() *Amqpin {
	return &Amqpin{
		stop: make(chan bool),
		quit: make(chan bool),
	}
}

func (amqpin *Amqpin) GetName() string {
	return "amqpin"
}

func (amqpin *Amqpin) GetType() generic.LinkType {
	return generic.LinkTypeGenerator
}

func (amqpin *Amqpin) Config(config *generic.ConfigurationLink) bool {
	var err error

	amqpin.route = (*config)["route"].(string)
	amqpin.uri = (*config)["uri"].(string)

	amqpin.AmqpConnection, err = amqp.Dial(amqpin.route)
	if err != nil {
		glog.Fatalf("Can't connect to amqp queue [%s]: %v", amqpin.route, err)
		panic(err)
	}
	amqpin.AmqpChannel, err = amqpin.AmqpConnection.Channel()
	if err != nil {
		glog.Fatalf("Can't open channel: %v", err)
		panic(err)
	}
	amqpin.AmqpChannel.Confirm(false)

	amqpin.AmqpQueue, err = amqpin.AmqpChannel.QueueDeclare(
		amqpin.route,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		glog.Fatalf("Can't declare queue [%v]: %v", amqpin.route, err)
		panic(err)
	}

	return true
}

func (amqpin *Amqpin) Do(input chan messages.Message) chan messages.Message {
	o := make(chan messages.Message)

	items, err := amqpin.AmqpChannel.Consume(amqpin.route, "", false, false, false, false, nil)
	if err != nil {
		glog.Fatalf("Can't consume queue: %v", err)
		panic(err)
	}

	go func() {
		for item := range items {
			select {
				case o <- *messages.NewMessage(amqpin.route, string(item.Body)):
					//TODO: make this configurable
					item.Ack(true)
				case <- amqpin.stop:
					amqpin.quit <- true
					return
			}			
		}
	}()

	return o
}

func (amqpin *Amqpin) Close() {
	amqpin.stop <- true
	<- amqpin.quit
	//TODO: free up resources
	if err := amqpin.AmqpConnection.Close(); err != nil {
		glog.Fatalf("AMQP connection close error: %s", err)
		panic(err)
	}
}
