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
	//	"fmt"
	"bytes"
	"net/http"

	"github.com/mikeszltd/chainsd/channels"
	"github.com/mikeszltd/chainsd/generic"
	"github.com/mikeszltd/chainsd/messages"
)

type Httpout struct {
	Terminal chan messages.Message

	host        string
	requestType string
	body        string

	stop		chan bool
	quit		chan bool
}

func NewHttpout() *Httpout {
	return &Httpout{
		stop: make(chan bool),
		quit: make(chan bool),
	}		
}

func (httpout *Httpout) GetName() string {
	return "httpout"
}

func (httpout *Httpout) GetType() generic.LinkType {
	return generic.LinkTypeOutput
}

func (httpout *Httpout) Config(config *generic.ConfigurationLink) bool {
	httpout.Terminal = channels.GetGlobalChannel(channels.TERMINAL_CHANNEL)

	var err error

	httpout.host, err = generic.GetConfigStringValue(config, "host", "")
	if httpout.host == "" || err != nil {
		glog.Error("[httpout] \"host\" not provided")
		return false
	}
	//	httpout.requestType = (*config)["requestType"].(string)
	//httpout.query = (*config)["query"].(string)

	return true
}

func (httpout *Httpout) Do(input chan messages.Message) chan messages.Message {
	go func() {
		for {
			select {
				case m := <-input:
					resp, _ := http.Post(httpout.host, "application/json", bytes.NewBufferString(m.Content))
					resp.Body.Close()
					httpout.Terminal <- *messages.NewMessage("ok", "ok")
				case <- httpout.stop:
					httpout.quit <- true
					return
			}
		}
	}()

	return nil
}

func (httpout *Httpout) Close() {
	httpout.stop <- true
	<- httpout.quit
	glog.Info("[httpout] Closed")
}
