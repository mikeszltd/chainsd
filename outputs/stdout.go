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
	"fmt"

	"github.com/mikeszltd/chainsd/channels"
	"github.com/mikeszltd/chainsd/generic"
	"github.com/mikeszltd/chainsd/messages"
)

type Stdout struct {
	Terminal chan messages.Message

	stop		chan bool
	quit		chan bool	
}

func NewStdout() *Stdout {
	return &Stdout{
		stop: make(chan bool),
		quit: make(chan bool),
	}
}

func (stdout *Stdout) GetName() string {
	return "stdout"
}

func (stdout *Stdout) GetType() generic.LinkType {
	return generic.LinkTypeOutput
}

func (stdout *Stdout) Config(config *generic.ConfigurationLink) bool {
	stdout.Terminal = channels.GetGlobalChannel(channels.TERMINAL_CHANNEL)
	return true
}

func (stdout *Stdout) Do(input chan messages.Message) chan messages.Message {
	go func() {
		for {
			select {
				case m := <-input:
					fmt.Println(m.Content)					
					stdout.Terminal <- *messages.NewMessage("ok", "ok")
				case <- stdout.stop:
					stdout.quit <- true
					return
			}
		}
	}()

	return nil
}

func (stdout *Stdout) Close() {
	stdout.stop <- true
	<- stdout.quit
	glog.Info("[stdout] Closed")
}
