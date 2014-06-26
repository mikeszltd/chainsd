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

type Parser interface {
	Pattern(pattern string) bool
	Parse(input string) (string, error)
}

func ParserFactory(name string) Parser {
	parser := NewScraper()

	switch name {
	case "apache":
		parser.Pattern("^(?P<host>[\\d.]+) (\\S+) (?P<user>\\S+) \\[(?P<date>[\\w:/]+\\s[+\\-]\\d{4})\\] \"(?P<request>.+?)\" (?P<status>\\d{3}) (?P<size>\\d+) \"(?P<referer>[^\"]+)\" \"(?P<agent>[^\"]+)\"")
	case "body":
		parser.Pattern("")
	default:
		parser.Pattern(name)
	}

	return parser
}
