//管理session信息
//	pSession = sapi.Plugin("session").(*session.Session)
package session

import (
	"wgf/plugin/cookie"
	"wgf/sapi"

	"encoding/gob"
	"math/rand"
	"net/http"
	"strconv"
	"bytes"
	"time"
)

func uuid() string {
	now := time.Now()
	unixtimestamp := now.Unix()
	rand.Seed(unixtimestamp)

	pre := strconv.FormatInt(unixtimestamp, 36)
	suf := strconv.FormatInt(rand.Int63(), 36)
	return pre + suf
}

type Session struct {
	sapi       *sapi.Sapi
	id         string
	hasStarted bool
	h Handler
}

func (s *Session) Id() string {
	if "" == s.id {
		s.id = uuid()
	}
	return s.id
}

func (s *Session) Start() {
	if s.hasStarted {
		return
	}

	s.hasStarted = true
	ck := s.sapi.Plugin("cookie").(*cookie.Cookie)
	id := ck.Get("SID")

	if id == "" {
		id = uuid()
		newcookie := http.Cookie{Name: "SID", Value: id}
		ck.Set(&newcookie)
	}
	s.id = id
}

func (s *Session) Get(key string, dst interface{}) bool {
	valInStore, ok := s.h.Get(s.id, key)
	if false == ok || nil != gob.NewDecoder(bytes.NewReader(valInStore)).Decode(dst) {
		return false
	}
	return true

}

func (s *Session) Set(key string, value interface{}) bool {
	buf := new(bytes.Buffer)
	if nil != gob.NewEncoder(buf).Encode(value) {
		return false
	}
	return s.h.Set(s.id, key, buf.Bytes())
}

func (s *Session) Del(key string) bool {
	return s.h.Del(s.id, key)
}

func (s *Session) Destory() bool {
	return s.h.Destory(s.id)
}

func sessionCreater() (interface{}, error) {
	ret := &Session{hasStarted: false}
	ret.h = newDefaultHandler(1200)
	return ret, nil
}

func requestInit(s *sapi.Sapi, p interface{}) error {
	session := p.(*Session)
	session.sapi = s
	return nil
}

func init() {
	info := sapi.PluginInfo{}
	info.Creater = sessionCreater
	info.HookPluginRequestInit = requestInit
	info.BasePlugins = []string{"cookie"}
	(&info).Support(sapi.IdHttp, sapi.IdWebsocket)
	sapi.RegisterPlugin("session", info)
}
