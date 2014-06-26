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
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/mikeszltd/chainsd/generic"
	"github.com/mikeszltd/chainsd/messages"	
)

type Httpin struct {
	name string
	host string
	port string

	stop		chan bool
	quit		chan bool
}

func NewHttpin() *Httpin {
	return &Httpin{
		stop: make(chan bool),
		quit: make(chan bool),		
	}
}

func (httpin *Httpin) GetName() string {
	return "httpin"
}

func (httpin *Httpin) GetType() generic.LinkType {
	return generic.LinkTypeGenerator
}

func (httpin *Httpin) Config(config *generic.ConfigurationLink) bool {
	httpin.name = (*config)["name"].(string)
	httpin.host = (*config)["host"].(string)
	httpin.port = (*config)["port"].(string)

	return true
}

func (httpin *Httpin) Serve(handler http.Handler) {
	addr := httpin.host + ":" + httpin.port
	s := http.Server{
		Addr: addr,
		Handler: nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	l, err := net.Listen("tcp", addr)

	if err != nil {
		glog.Fatalf("Can't listern on [%s]: %v", addr, err)
		panic(err)
	}

	go s.Serve(l)

	select {
	case <- httpin.stop:
		l.Close()
		httpin.quit <- true
	}
}

func (httpin *Httpin) Do(input chan messages.Message) chan messages.Message {
	o := make(chan messages.Message)

	go func() {
		mux := http.DefaultServeMux

		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "OPTIONS" {
				// handle CORS stuff
				if r.ContentLength != 0 {
					w.Header().Set("Content-Length", "0")
					mb := http.MaxBytesReader(w, r.Body, 4<<10)
  					io.Copy(ioutil.Discard, mb)
  				}
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "POST")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				w.Header().Set("Access-Control-Max-Age", "31536000")
  				return
			}
			w.Header().Set("Access-Control-Allow-Origin", "*")

			fmt.Println(r.Method, ":", r.URL.Query)

			message, _ := ioutil.ReadAll(r.Body)

			o <- *messages.NewMessage(r.URL.Path, string(message))
		})
		httpin.Serve(mux)
	}()

	return o
}

func (httpin *Httpin) Close() {
	//TODO: free up resources
	httpin.stop <- true
	<- httpin.quit 
}
