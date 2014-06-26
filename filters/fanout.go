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

package filters

import (
	"github.com/golang/glog"

	"github.com/mikeszltd/chainsd/channels"
	"github.com/mikeszltd/chainsd/generic"
	"github.com/mikeszltd/chainsd/messages"
)

type Fanout struct {
	name string

	Channel chan messages.Message
}

func NewFanout() *Fanout {
	return &Fanout{}
}

func (fanout *Fanout) GetName() string {
	return "fanout"
}

func (fanout *Fanout) GetType() generic.LinkType {
	return generic.LinkTypeFilter
}

func (fanout *Fanout) Config(config *generic.ConfigurationLink) bool {
	fanout.name = (*config)["name"].(string)

	var err error

	fanout.name, err = generic.GetConfigStringValue(config, "name", "")

	if fanout.name == "" || err != nil {
		glog.Error("[fanout] \"name\" not provided")
		return false
	}

	fanout.Channel = channels.GetGlobalChannel(fanout.name)
	return true
}

func (fanout *Fanout) Do(input chan messages.Message) chan messages.Message {
	o := make(chan messages.Message)
	go func() {
		for {
			m := <-input

			/*
				select {
					case fanout.Channel <- *m.Clone():
					default:
				}*/
			fanout.Channel <- m
			o <- m
		}
	}()

	return o
}

func (fanout *Fanout) Close() {
	//TODO: free up resources
}
