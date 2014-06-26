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

	"github.com/mikeszltd/chainsd/channels"
	"github.com/mikeszltd/chainsd/generic"
	"github.com/mikeszltd/chainsd/messages"	
)

type Fanin struct {
	name    string
	Channel chan messages.Message

	stop		chan bool
	quit		chan bool	
}

func NewFanin() *Fanin {
	return &Fanin{
		stop: make(chan bool),
		quit: make(chan bool),
	}
}

func (fanin *Fanin) GetName() string {
	return "fanout"
}

func (fanin *Fanin) GetType() generic.LinkType {
	return generic.LinkTypeGenerator
}

func (fanin *Fanin) Config(config *generic.ConfigurationLink) bool {
	fanin.name = (*config)["name"].(string)
	fanin.Channel = channels.GetGlobalChannel(fanin.name)
	return true
}

func (fanin *Fanin) Do(input chan messages.Message) chan messages.Message {
	o := make(chan messages.Message)

	go func() {
		for {
			select {
				case m := <-fanin.Channel:
					o <- m
				case <- fanin.stop:
					fanin.quit <- true
					return
			}
		}
	}()

	return o
}

func (fanin *Fanin) Close() {
	<- fanin.quit
	glog.Info("Closed fanin")
}
