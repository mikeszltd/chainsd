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
	"container/list"

	"github.com/mikeszltd/chainsd/generic"
	"github.com/mikeszltd/chainsd/messages"
)

type ChainStore struct {
	links *list.List
}

func NewChainStore() *ChainStore {	
	chain := &ChainStore{
		links: list.New(),
	}
	return chain
}

func (chainStore *ChainStore) Join(link *generic.Link) {
	chainStore.links.PushBack(link)
}

func (chainStore *ChainStore) Validate() bool {
	return true
}

func (chainStore *ChainStore) Close() {
	front := chainStore.links.Front()

	for e := front; e != nil; e = e.Next() {
		link := e.Value.(*generic.Link)
		(*link).Close()		
	}

}

func (chainStore *ChainStore) Run() chan messages.Message {
	front := chainStore.links.Front()

	start := front.Next()

	link := front.Value.(*generic.Link)

	sink := (*link).Do(nil)

	for e := start; e != nil; e = e.Next() {
		link = e.Value.(*generic.Link)
		sink = (*link).Do(sink)
	}
	return sink
}
