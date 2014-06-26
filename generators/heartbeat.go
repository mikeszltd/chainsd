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
   "encoding/json"
	"github.com/golang/glog"

	"github.com/mikeszltd/chainsd/generic"
	"github.com/mikeszltd/chainsd/messages"
	"time"
)

const (
	HEARTBEAT_TAG = "heartbeat"
)

type Heartbeat struct {
	durationRaw string
	duration    time.Duration
	message     string
	tag         string
	ticker      *time.Ticker

	stop		chan bool
	quit		chan bool
}

type Beat struct {
	Timestamp int64 `json:"timestamp"`
	Message string `json:"message"`
}

func NewHeartbeat() *Heartbeat {
	return &Heartbeat{
		stop: make(chan bool),
		quit: make(chan bool),
	}
}

func (heartbeat *Heartbeat) GetName() string {
	return "heartbeat"
}

func (heartbeat *Heartbeat) GetType() generic.LinkType {
	return generic.LinkTypeGenerator
}

func (heartbeat *Heartbeat) Config(config *generic.ConfigurationLink) bool {
	var err interface{}

	heartbeat.durationRaw = (*config)["interval"].(string)
	heartbeat.duration, _ = time.ParseDuration(heartbeat.durationRaw)

	heartbeat.tag, err = (*config)["tag"].(string)

	if err != nil {
		glog.Warningf("[heartbeat] \"tag\" value missing. Using default \"%s\".", HEARTBEAT_TAG)
		heartbeat.tag = HEARTBEAT_TAG
	}

	heartbeat.ticker = time.NewTicker(heartbeat.duration)

	m, ok := (*config)["message"].(string)
	if ok {
		heartbeat.message = m
	} else {
		heartbeat.message = "heartbeat"
	}

	return true
}

func (heartbeat *Heartbeat) Do(input chan messages.Message) chan messages.Message {
	message := make(chan messages.Message)

	go func() {				
		for {
			select {
			case <-heartbeat.ticker.C:
				m := &Beat{
					Timestamp: time.Now().Unix(),
					Message: heartbeat.message,
				}
				b, err := json.Marshal(m)
				if err != nil {
					glog.Errorf("Heart problem: %v", err)					
				} else {
					message <- *messages.NewMessage(heartbeat.tag, string(b))
				}
			case <- heartbeat.stop:
				glog.Info("Stopping heartbeat")
				heartbeat.quit <- true
				return
			}
		}
	}()

	return message
}

func (heartbeat *Heartbeat) Close() {	
	heartbeat.stop <- true
	<- heartbeat.quit
}

