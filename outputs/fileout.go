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
	//"fmt"
	"os"
	"time"

	"github.com/mikeszltd/chainsd/channels"
	"github.com/mikeszltd/chainsd/generic"
	"github.com/mikeszltd/chainsd/messages"
)

type Fileout struct {
	lines    int
	filename string
	Terminal chan messages.Message

	stop		chan bool
	quit		chan bool	
}

func NewFileout() *Fileout {
	return &Fileout{
		stop: make(chan bool),
		quit: make(chan bool),		
	}
}

func (fileout *Fileout) GetName() string {
	return "stdout"
}

func (fileout *Fileout) GetType() generic.LinkType {
	return generic.LinkTypeOutput
}

func (fileout *Fileout) Config(config *generic.ConfigurationLink) bool {
	fileout.Terminal = channels.GetGlobalChannel(channels.TERMINAL_CHANNEL)

	var err error

	fileout.filename, err = generic.GetConfigStringValue(config, "filename", "")
	if fileout.filename == "" || err != nil {
		glog.Error("[fileout] \"filename\" not provided.")
		return false
	}

	return true
}

func (fileout *Fileout) Do(input chan messages.Message) chan messages.Message {

	go func() {
		for {
			select {
			case m := <-input:
				filename := m.Touch(fileout.filename)

				file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
				if err != nil {
					glog.Errorf("[fileout] Can't open/create file \"%s\"", filename)
				} else {
					t := time.Unix(m.CreatedAt, 0)

					if _, err = file.WriteString("[" + t.Format("2006-01-02 15:04:05") + "] " + m.Content + "\n"); err != nil {
						glog.Errorf("[fileout] Can't write to file \"%s\"", filename)
					}

					file.Close()
				}

				fileout.Terminal <- *messages.NewMessage("ok", "ok")
			case <- fileout.stop:
				fileout.quit <- true
				return
			}
		}
	}()

	return nil
}

func (fileout *Fileout) Close() {
	fileout.stop <- true
	<- fileout.quit
	glog.Info("[fileout] Closed")
}
