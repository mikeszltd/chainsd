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

package generic

import (
	"fmt"

	"github.com/mikeszltd/chainsd/messages"
)

type LinkType int

const (
	LinkTypeGenerator LinkType = iota
	LinkTypeFilter
	LinkTypeOutput
)

type ConfigurationLink map[string]interface{}

func GetConfigStringValue(config *ConfigurationLink, key string, defaultValue string) (value string, err error) {
	value, exists := (*config)[key].(string)
	if exists {
		return value, nil
	}

	return defaultValue, fmt.Errorf("Link [%s] doesn't exist", key)
}

type Link interface {
	GetType() LinkType
	GetName() string

	Config(config *ConfigurationLink) bool

	Do(input chan messages.Message) chan messages.Message

	Close()
}
