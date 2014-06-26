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
	"encoding/json"
	"github.com/golang/glog"
	"labix.org/v2/mgo"

	"github.com/mikeszltd/chainsd/channels"
	"github.com/mikeszltd/chainsd/generic"
	"github.com/mikeszltd/chainsd/messages"
)

type Mongoout struct {
	lines      int
	hosts      string
	database   string
	collection string

	session *mgo.Session

	Terminal chan messages.Message

	stop		chan bool
	quit		chan bool	
}

func NewMongoout() *Mongoout {
	return &Mongoout{
		stop: make(chan bool),
		quit: make(chan bool),
	}		
}

func (mongoout *Mongoout) GetName() string {
	return "mongoout"
}

func (mongoout *Mongoout) GetType() generic.LinkType {
	return generic.LinkTypeOutput
}

func (mongoout *Mongoout) Config(config *generic.ConfigurationLink) bool {
	mongoout.Terminal = channels.GetGlobalChannel(channels.TERMINAL_CHANNEL)

	var err error

	mongoout.hosts, err = generic.GetConfigStringValue(config, "hosts", "")
	if mongoout.hosts == "" || err != nil {
		glog.Error("[mongoout] \"hosts\" not provided")
		return false
	}

	mongoout.database, err = generic.GetConfigStringValue(config, "database", "")
	if mongoout.database == "" || err != nil {
		glog.Error("[mongoout] \"database\" not provided")
		return false
	}

	mongoout.collection, err = generic.GetConfigStringValue(config, "collection", "")
	if mongoout.database == "" || err != nil {
		glog.Error("[mongoout] \"collection\" not provided")
		return false
	}

	mongoout.session, err = mgo.Dial(mongoout.hosts)

	if err != nil {
		glog.Errorf("[mongoout] Unable to establish connection to: %s", mongoout.hosts)
		return false
	}

	return true
}

func (mongoout *Mongoout) Do(input chan messages.Message) chan messages.Message {
	go func() {
		for {
			select {				
				case m := <-input:
					c := mongoout.session.DB(mongoout.database).C(mongoout.collection)

					var data map[string]interface{}
					json.Unmarshal([]byte(m.Content), &data)

					c.Insert(data)

					mongoout.Terminal <- *messages.NewMessage("ok", "ok")
				case <- mongoout.stop:
					mongoout.quit <- true
					return
			}
		}
	}()

	return nil
}

func (mongoout *Mongoout) Close() {
	mongoout.stop <- true
	<- mongoout.quit

	if mongoout.session != nil {
		mongoout.session.Close()
	}
	glog.Info("[mongoout] Closed")
}
