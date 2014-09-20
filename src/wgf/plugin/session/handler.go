// Copyright 2014 The Wgf Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package session

import (
	"sync"
	"time"
	"container/list"
)

type Handler interface {

	SetExpireTime(ttl int)

	Start() bool
	Stop() bool

	Set(sid, key string, value []byte) bool
	Get(sid, key string)  ([]byte, bool)
	Del(sid, key string) bool
	Destory(sid string) bool
}

type sessionTime struct {
	sid string
	time time.Time
}

type sessionValue struct {
	value map[string][]byte
	element *list.Element
}

type defaultHandler struct {
	data map[string]sessionValue

	//for expire
	list *list.List
	running bool
	ttl int

	//for data safe
	lock *sync.RWMutex
}

func newDefaultHandler() *defaultHandler {
	ret := new(defaultHandler)
	ret.ttl = 1200 //default 1200
	ret.data = make(map[string]sessionValue)
	ret.lock = new(sync.RWMutex)
	ret.list = list.New()
	return ret
}

func (h *defaultHandler) SetExpireTime(ttl int) {
	h.ttl = ttl
}

func (h *defaultHandler) expire() {
	var now time.Time
	var e *list.Element
	var st sessionTime

	for h.running {
		now = time.Now().Add(time.Duration(-h.ttl)*time.Second)
		for {
			e = h.list.Back()
			if nil == e {
				break
			}

			st=e.Value.(sessionTime)
			if st.time.After(now) {
				break
			}
			h.lock.Lock()
			h.list.Remove(e)
			delete(h.data, st.sid)
			h.lock.Unlock()
		}
		time.Sleep(1*time.Second)
	}
}

func (h *defaultHandler) Start() bool {
	h.running = true
	go h.expire()
	return true
}

func (h *defaultHandler) Stop() bool {
	h.running = false
	return true
}

func (h *defaultHandler) Set(session, key string, value []byte) bool {
	h.lock.Lock()
	defer h.lock.Unlock()

	_, ok := h.data[session]
	if !ok {
		st := sessionTime{session, time.Now()}
		tmp := sessionValue{}
		tmp.value = make(map[string][]byte)
		tmp.element = h.list.PushFront(st)
		h.data[session] = tmp
	}
	h.data[session].value[key] = value
	h.refreshElement(session, true)
	return true
}

func (h *defaultHandler) Del(session, key string) bool {
	h.lock.Lock()
	defer h.lock.Unlock()

	_, ok := h.data[session]
	if ok {
		delete(h.data[session].value, key)
	}
	h.refreshElement(session, ok)
	return true
}

func (h *defaultHandler) Get(session, key string) (value []byte, ok bool) {
	h.lock.RLock()
	defer h.lock.RUnlock()

	v, ok := h.data[session]
	if ok {
		value, ok = v.value[key]
	}
	h.refreshElement(session, ok)
	return
}

func (h *defaultHandler) Destory(session string) bool {
	h.lock.Lock()
	defer h.lock.Unlock()

	_, ok := h.data[session]
	if ok {
		h.list.Remove(h.data[session].element)
		delete(h.data, session)
	}
	return true
}

func (h *defaultHandler) refreshElement(session string, exists bool) {
	if exists {
		h.list.MoveToFront(h.data[session].element)
	}
}
