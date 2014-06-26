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

	"bitbucket.org/tebeka/strftime"
	"bytes"
	"github.com/vladimirvivien/gowfs"
	"time"

	"github.com/mikeszltd/chainsd/channels"
	"github.com/mikeszltd/chainsd/generic"
	"github.com/mikeszltd/chainsd/messages"
)

type Hdfsout struct {
	host    string
	user    string
	pattern string

	Terminal chan messages.Message

	stop		chan bool
	quit		chan bool	
}

func NewHdfsout() *Hdfsout {
	return &Hdfsout{
		stop: make(chan bool),
		quit: make(chan bool),
	}
}

func (hdfsout *Hdfsout) GetName() string {
	return "hdfsout"
}

func (hdfsout *Hdfsout) GetType() generic.LinkType {
	return generic.LinkTypeOutput
}

func (hdfsout *Hdfsout) Config(config *generic.ConfigurationLink) bool {
	hdfsout.Terminal = channels.GetGlobalChannel(channels.TERMINAL_CHANNEL)

	var err error

	hdfsout.host, err = generic.GetConfigStringValue(config, "host", "")
	if hdfsout.host == "" || err != nil {
		glog.Error("[hdfsout] Filename not provided.")
		return false
	}

	hdfsout.user, err = generic.GetConfigStringValue(config, "user", "")
	if hdfsout.user == "" || err != nil {
		glog.Error("[hdfsout] User not provided.")
		return false
	}

	hdfsout.pattern, err = generic.GetConfigStringValue(config, "pattern", "")
	if hdfsout.pattern == "" || err != nil {
		glog.Error("[hdfsout] Pattern not provided.")
		return false
	}

	return true
}

func (hdfsout *Hdfsout) Do(input chan messages.Message) chan messages.Message {
	go func() {
		conf := *gowfs.NewConfiguration()
		conf.Addr = hdfsout.host
		conf.User = hdfsout.user
		conf.ConnectionTimeout = time.Second * 60
		conf.DisableKeepAlives = false

		//TODO: connect on config time
		fs, err := gowfs.NewFileSystem(conf)

		if err != nil {
			glog.Errorf("[hdfsout] Error acquiring filesystem: %s", err)
		}

		shell := gowfs.FsShell{FileSystem: fs, WorkingPath: "/"}

		for {
			select {
			case m := <-input:
				pattern := m.Touch(hdfsout.pattern)
				pattern, err = strftime.Format(pattern, time.Now())
				path := pattern

				exists, err := shell.Exists(path)

				if !exists {
					_, err = fs.Create(
						bytes.NewBufferString(m.Content),
						gowfs.Path{Name: path},
						true,
						0,
						0,
						0700,
						0,
					)

					if err != nil {
						glog.Errorf("[hdfsout] Error creating file: %s", err)
					}

				} else {
					_, err = fs.Append(bytes.NewBufferString(m.Content+"\n"), gowfs.Path{Name: path}, 4096)

					if err != nil {
						glog.Errorf("[hdfsout] Error appending to file: %s", err)
					}
				}				
				hdfsout.Terminal <- *messages.NewMessage("ok", "ok")
			case <- hdfsout.stop:
				hdfsout.quit <- true
				return
			}
		}
	}()

	return nil
}

func (hdfsout *Hdfsout) Close() {
	hdfsout.stop <- true
	<- hdfsout.quit
	glog.Info("[hdfsout] Closed")
}
