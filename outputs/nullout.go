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
)

type Nullout struct {
	Terminal chan messages.Message

	stop		chan bool
	quit		chan bool	
}

func NewNullout() *Nullout {
	return &Nullout{
		stop: make(chan bool),
		quit: make(chan bool),
	}			
}

func (nullout *Nullout) GetName() string {
	return "nullout"
}

func (nullout *Nullout) GetType() generic.LinkType {
	return generic.LinkTypeOutput
}

func (nullout *Nullout) Config(config *generic.ConfigurationLink) bool {
	nullout.Terminal = channels.GetGlobalChannel(channels.TERMINAL_CHANNEL)
	return true
}

func (nullout *Nullout) Do(input chan messages.Message) chan messages.Message {
	go func() {
		for {
			select {
				case _ = <-input:
					nullout.Terminal <- *messages.NewMessage("ok", "ok")
				case <- nullout.stop:
					nullout.quit <- true
					return
			}
		}
	}()

	return nil
}

func (nullout *Nullout) Close() {
	nullout.stop <- true
	<- nullout.quit

	glog.Info("[nullout] Closed")
}
