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

package parsers

import (
	"github.com/golang/glog"

	"encoding/json"	
	"regexp"
)

type Scraper struct {
	pattern  string
	compiled *regexp.Regexp
}

func NewScraper() *Scraper {
	//pattern := "^(?P<ip>[\\d.]+) (\\S+) (\\S+) \\[(?P<datetime>[\\w:/]+\\s[+\\-]\\d{4})\\] \"(?P<request>.+?)\" (?P<response>\\d{3}) (?P<bytessent>\\d+) \"(?P<referer>[^\"]+)\" \"(?P<browser>[^\"]+)\""

	return &Scraper{
		pattern:  "",
		compiled: nil,
	}
}

func (scraper *Scraper) Pattern(pattern string) bool {
	scraper.pattern = pattern

	var err interface{}

	scraper.compiled, err = regexp.Compile(scraper.pattern)

	if err != nil {
		glog.Errorf("[Parser] Couldn't compile: \"%s\"", scraper.pattern)
		return false
	}

	return true
}

func (scraper *Scraper) Parse(input string) (output string, err error) {
	content := make(map[string]string)

	names := scraper.compiled.SubexpNames()[1:]
	matches := scraper.compiled.FindStringSubmatch(input)[1:]

	for i, m := range matches {
		n := names[i]
		if n != "" {
			content[n] = m
		}
	}

	b, err := json.Marshal(content)	

	output = string(b)

	return
}
