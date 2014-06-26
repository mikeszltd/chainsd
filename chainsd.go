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

package main

import (
	"github.com/golang/glog"

	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"io/ioutil"
	"runtime"

	"github.com/mikeszltd/chainsd/channels"
	"github.com/mikeszltd/chainsd/generic"
	"github.com/mikeszltd/chainsd/system"
)

const (
	CHAINSD_RELEASE = "0.0.1"
)

type Configuration struct {
	properties map[string][]map[string][]map[string]interface{}

	stores []*system.ChainStore
}

func main() {
	fmt.Printf("Chainsd by Mikesz Ltd\n")
	glog.Infof("Starting Chainsd")

	filename := flag.String("config", "config.json", "Path to configuration file")

	flag.Parse()

	data, err := ioutil.ReadFile(*filename)

	if err != nil {
		glog.Fatalf("[Chainsd] Couldn't find configuration file \"%s\"", *filename)
		panic(err)
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	configuration := new(Configuration)

	err = json.Unmarshal(data, &configuration.properties)

	if err != nil {
		glog.Fatal("[Chainsd] Couldn't parse configuration file")
		panic(err)
	}

	channels.Init()
	
	linkFactory := system.NewLinkFactory()

	for _, chain := range configuration.properties["chains"] {
		//load chain
		chainStore := system.NewChainStore()

		configuration.stores = append(configuration.stores, chainStore)

		for name, elements := range chain {
			glog.Infof("[Chainsd] Loading chain: %s", name)
			for _, config := range elements {
				cl := generic.ConfigurationLink(config)

				l := linkFactory.Build(&cl)

				if l != nil {
					chainStore.Join(&l)
				} else {
					msg := fmt.Sprintf("Unknown link: %s", cl)
					glog.Fatalf(msg)
					panic(msg)
				}
			}
		}
		chainStore.Run()
	}

	sig := make(chan os.Signal, 1)

	signal.Notify(sig,
    	syscall.SIGHUP,
    	syscall.SIGINT,
    	syscall.SIGTERM,
    	syscall.SIGQUIT)

	//check for signals, so we can clean-up stuff if needed (close processes etc.)
	go func() {
		s := <- sig
		glog.Infof("Received signal [%s]", s)

		fmt.Println(s)

		for _, store := range configuration.stores {			
			(*store).Close()
		}
		os.Exit(0)
	}()

	Terminal := channels.GetGlobalChannel(channels.TERMINAL_CHANNEL)
	for {
		x := <-Terminal
		glog.Info(x)
	}
}
