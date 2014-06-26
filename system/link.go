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

package system

import (
	"github.com/mikeszltd/chainsd/filters"
	"github.com/mikeszltd/chainsd/generators"
	"github.com/mikeszltd/chainsd/generic"
	"github.com/mikeszltd/chainsd/outputs"
)

type LinkFactory struct {
	links int
}

func NewLinkFactory() *LinkFactory {
	return &LinkFactory{
		links: 0,
	}
}

func (linkFactory *LinkFactory) Build(config *generic.ConfigurationLink) generic.Link {
	var link generic.Link

	switch (*config)["type"] {
	case "heartbeat":
		link = generators.NewHeartbeat()
	case "fanin":
		link = generators.NewFanin()
	case "httpin":
		link = generators.NewHttpin()
	case "tail":
		link = generators.NewTail()
	case "fanout":
		link = filters.NewFanout()
	case "fileout":
		link = outputs.NewFileout()
	case "mongoout":
		link = outputs.NewMongoout()
	case "hdfsout":
		link = outputs.NewHdfsout()
	case "httpout":
		link = outputs.NewHttpout()
	case "nullout":
		link = outputs.NewNullout()
	case "stdout":
		link = outputs.NewStdout()
	default:
		link = nil
	}

	if link != nil {
		result := link.Config(config)
		if !result {
			return nil
		}

		linkFactory.links++
	}

	return link
}
