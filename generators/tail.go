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
	"bufio"
	"fmt"
	"github.com/golang/glog"
	"os/exec"
	"encoding/json"
	//	"os"

	"github.com/mikeszltd/chainsd/generic"
	"github.com/mikeszltd/chainsd/messages"
	"github.com/mikeszltd/chainsd/parsers"	
)

type Tail struct {
	filename string
	tag      string
	parser   parsers.Parser

	cmd *exec.Cmd

	stop		chan bool
	quit		chan bool	
}

type Line struct {
	Line string `json:"line"`
}

func NewTail() *Tail {
	return &Tail{
		stop: make(chan bool),
		quit: make(chan bool),		
	}
}

func (tail *Tail) GetName() string {
	return "tail"
}

func (tail *Tail) GetType() generic.LinkType {
	return generic.LinkTypeGenerator
}

func (tail *Tail) Config(config *generic.ConfigurationLink) bool {
	var err error

	tail.filename, err = generic.GetConfigStringValue(config, "filename", "")

	if tail.filename == "" || err != nil {
		glog.Error("[tail] \"filename\" not provided")
		return false
	}

	tail.tag, _ = generic.GetConfigStringValue(config, "tag", tail.filename)

	parserName, _ := generic.GetConfigStringValue(config, "parser", "body")
	tail.parser = parsers.ParserFactory(parserName)

	return true
}

func (tail *Tail) ProcessLine(line string) (output string, err error) {
	if tail.parser != nil {
		output, err = tail.parser.Parse(line)		
	} else {
		data := &Line{
			Line: line,
		}
		var b []byte
		b, err = json.Marshal(data)
		if err != nil {
			glog.Errorf("[tail] Problem serializing line [%v]", data, err)						
		}
		output = string(b)		
	}
	return
}

func (tail *Tail) Do(input chan messages.Message) chan messages.Message {
	message := make(chan messages.Message)

	go func() {	
		tail.cmd = exec.Command("tail", "-f", tail.filename)
		out, err := tail.cmd.StdoutPipe()

		if err != nil {
			msg := fmt.Sprintf("[tail] Tail couldn't start: %v", err)
			glog.Fatalf(msg)
			panic(msg)
		}

		err = tail.cmd.Start()
		if err != nil {
			glog.Fatalf("Can't start tail: %v", err)
			panic(err)
		}

		reader := bufio.NewReader(out)

		glog.Infof("Started tail [%s]", tail.filename)
		for err == nil {
			select {
			case <- tail.stop:
				glog.Info("Stopping tail")
				tail.quit <- true
				return
			default:
				line, err := reader.ReadString('\n')

				if err != nil {
					glog.Errorf("[tail] Problem reading from \"%s\"", tail.filename)
				} else {
					processed, err := tail.ProcessLine(line)
					if err != nil {
						glog.Errorf("[tail] Problem processing line \"%s\"", line)
					} else {
						message <- *messages.NewMessage(tail.tag, processed)
					}
				}
			}
		}
	}()

	return message
}

func (tail *Tail) Close() {
	tail.stop <- true
	<- tail.quit
	glog.Infof("Closing tail [%s]", tail.filename)
	tail.cmd.Process.Kill()
}
