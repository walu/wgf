package session

import (
	"bytes"
	"testing"
)

func TestCommonAction(t *testing.T) {
	var ok bool
	var val1, val2 []byte
	var sid string

	val1 = []byte("wgf")
	sid = "ABCDEFG"

	h := newDefaultHandler()
	h.Start()

	//set
	if false == h.Set(sid, "name", val1) {
		t.Error("set error")
	}

	//get
	val2, ok = h.Get(sid, "name")
	if false == ok {
		t.Error("value missed")
	}

	if 0 != bytes.Compare(val1, val2) {
		t.Error("value mismatch")
	}

	//del
	h.Set(sid, "name", val1)
	h.Del(sid, "name")
	val2, ok = h.Get(sid, "name")
	if ok {
		t.Error("delete error")
	}

	//destory
	h.Set(sid, "name", val1)
	h.Destory(sid)
	_, ok = h.Get(sid, "name")
	if ok {
		t.Error("Destory error")
	}
}
