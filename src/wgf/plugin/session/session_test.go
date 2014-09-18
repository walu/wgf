package session

import (
	"testing"
)

//use default session handler
func TestSet(t *testing.T) {
	ret := &Session{hasStarted: false}
	ret.h = newDefaultHandler(1200)

	var ok bool

	//set
	ok = ret.Set("name", "wgf")
	if !ok {
		t.Error("set error")
	}

	//get
	var val string
	ok = ret.Get("name", &val)

	if !ok {
		t.Error("content missed after set")
	}
	if val != "wgf" {
		t.Error("content mismatched after set")
	}
}
