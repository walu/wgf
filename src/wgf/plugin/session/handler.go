package session

import (
	"sync"
)

type Handler interface {
	Set(key, value string)
	Get(key string)
}

type defaultHandler struct {
	data map[string][string][]byte
	lock *sync.RWMutex
	ttl int
}

func newDefaultHandler(ttl int) *defaultHandler {
	ret := new(defaultHandler)
	ret.ttl = ttl
	ret.data = make(map[string][string][]byte)
	ret.lock = new(sync.RWMutex)
	return ret
}

func (h *defaultHandler) Set(session, key string, value []byte) bool {
	h.lock.Lock()
	defer h.lock.Unlock()

	ok, _ := h.data[session]
	if !ok {
		h.data = make(map[string][]byte)
	}
	h.data[session][key] = value
	return true
}

func (h *defaultHandler) Get(session, key string) (ok bool, value []byte) {
	h.lock.RLock()
	defer h.lock.RUnlock()

	ok, v := h.data[session]
	if ok {
		value = v[key]
	}
	return
}

func (h *defaultHandler) Flush(session string) bool {
	h.lock.Lock()
	defer h.lock.Unlock()

	delete(h.data, session)
	return true
}
