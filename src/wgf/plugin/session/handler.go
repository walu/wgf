package session

import (
	"sync"
	"time"
)

type Handler interface {
	Set(sid, key string, value []byte) bool
	Get(sid, key string)  ([]byte, bool)
	Del(sid, key string) bool
	Destory(sid string) bool
}

type defaultHandler struct {
	data map[string]map[string][]byte
	time map[string]time.Time
	lock *sync.RWMutex
	ttl int
}

func newDefaultHandler(ttl int) *defaultHandler {
	ret := new(defaultHandler)
	ret.ttl = ttl
	ret.data = make(map[string]map[string][]byte)
	ret.lock = new(sync.RWMutex)
	return ret
}

func (h *defaultHandler) Set(session, key string, value []byte) bool {
	h.lock.Lock()
	defer h.lock.Unlock()

	_, ok := h.data[session]
	if !ok {
		h.data[session] = make(map[string][]byte)
	}
	h.data[session][key] = value
	return true
}

func (h *defaultHandler) Del(session, key string) bool {
	h.lock.Lock()
	defer h.lock.Unlock()

	delete(h.data[session], key)
	return true
}

func (h *defaultHandler) Get(session, key string) (value []byte, ok bool) {
	h.lock.RLock()
	defer h.lock.RUnlock()

	v, ok := h.data[session]
	if ok {
		value, ok = v[key]
	}
	return
}

func (h *defaultHandler) Destory(session string) bool {
	h.lock.Lock()
	defer h.lock.Unlock()

	delete(h.data, session)
	return true
}
