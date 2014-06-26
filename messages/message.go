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

package messages

import (
	"github.com/golang/glog"

	"regexp"
	"strings"
	"time"

	"bytes"
	"encoding/json"
	"text/template"
)

type Message struct {
	Tag       string
	CreatedAt int64
	Content   string

	touchTag     *regexp.Regexp
	touchTagBody *regexp.Regexp
}

func NewMessage(tag string, content string) *Message {
	message := &Message{
		CreatedAt: time.Now().Unix(),
		Tag:       tag,
		Content:   content,
	}

	//TODO: Not so great... optimize it
	message.touchTag = regexp.MustCompile("{{(.+?)}}+")
	message.touchTagBody = regexp.MustCompile("([a-zA-Z0-9.]+)(|[a-zA-Z0-9.]+)?")

	return message
}

func (message *Message) String() string {
	return message.Content
}

func (message *Message) Clone() *Message {
	return &Message{
		CreatedAt: message.CreatedAt,
		Tag:       message.Tag,
		Content:   message.Content,
	}
}

func Optional(args ...interface{}) string {
	if args[0] == nil {
		return args[1].(string)
	}
	return args[0].(string)
}

func (message *Message) prepareContent(content string) string {
	//content has variables in following format:
	// {{key|default}} or just {{key}}, then default is empty string
	output := content

	matches := message.touchTag.FindAllString(content, -1)

	for _, m := range matches {
		tags := message.touchTagBody.FindAllString(m, -1)

		if len(tags) > 1 {
			output = strings.Replace(output, m, "{{ optional ."+tags[0]+" \""+tags[1]+"\" }}", -1)
		} else {
			output = strings.Replace(output, m, "{{ optional ."+tags[0]+" \"\" }}", -1)
		}

	}
	return output
}

//TODO: Optimize it
func (message *Message) Touch(content string) string {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(message.Content), &data)
	if err != nil {
		glog.Warning("[message] Message not in JSON format. Unable to touch.")
		//fmt.Println(message.Content)
		return content
	}
	t := template.New("template_name")

	t = t.Funcs(template.FuncMap{"optional": Optional})
	t, err = t.Parse(message.prepareContent(content))

	if err != nil {
		glog.Error("[message] Error while touching message. Not touched.")
		return content
	}

	var out bytes.Buffer
	t.Execute(&out, data)
	return out.String()
}
