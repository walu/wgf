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

type defaultHandler struct {
	data map[string]map[string][]byte

	//for expire
	element map[string]*list.Element
	list *list.List
	running bool
	ttl int

	//for data safe
	lock *sync.RWMutex
}

func newDefaultHandler() *defaultHandler {
	ret := new(defaultHandler)
	ret.ttl = 1200 //default 1200
	ret.data = make(map[string]map[string][]byte)
	ret.lock = new(sync.RWMutex)
	ret.list = list.New()
	ret.element = make(map[string]*list.Element)
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
			delete(h.element, st.sid)
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
		h.data[session] = make(map[string][]byte)
	}
	h.data[session][key] = value
	h.refreshElement(session)
	return true
}

func (h *defaultHandler) Del(session, key string) bool {
	h.lock.Lock()
	defer h.lock.Unlock()

	delete(h.data[session], key)
	h.refreshElement(session)
	return true
}

func (h *defaultHandler) Get(session, key string) (value []byte, ok bool) {
	h.lock.RLock()
	defer h.lock.RUnlock()

	v, ok := h.data[session]
	if ok {
		value, ok = v[key]
		h.refreshElement(session)
	}
	return
}

func (h *defaultHandler) Destory(session string) bool {
	h.lock.Lock()
	defer h.lock.Unlock()

	//delete data
	delete(h.data, session)

	//free element in expire list
	v, ok := h.element[session]
	if ok {
		h.list.Remove(v)
		delete(h.element, session)
	}
	return true
}

func (h *defaultHandler) refreshElement(session string) {
	v, ok := h.element[session]
	if ok {
		h.list.MoveToBack(v)
	} else {
		val := sessionTime{session, time.Now()}
		h.element[session] = h.list.PushFront(val)
	}
}
