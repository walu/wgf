//管理session信息
//	pSession = sapi.Plugin("session").(*session.Session)
package session

import (
	"wgf/plugin/cookie"
	"wgf/sapi"

	"math/rand"
	"net/http"
	"strconv"
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

var sessionMap map[string]map[string]interface{}

type Session struct {
	sapi       *sapi.Sapi
	id         string
	hasStarted bool
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

func (s *Session) Get(key string) interface{} {
	val, ok := sessionMap[s.id]
	if !ok {
		return ""
	}
	return val[key]
}

func (s *Session) Set(key string, value interface{}) {
	_, ok := sessionMap[s.id]
	if !ok {
		sessionMap[s.id] = make(map[string]interface{})
	}
	sessionMap[s.id][key] = value
}

func (s *Session) Del(key string) {
	delete(sessionMap[s.id], key)
}

func (s *Session) Destory() {
	delete(sessionMap, s.id)
}

func sessionCreater() (interface{}, error) {
	return &Session{hasStarted: false}, nil
}

func requestInit(s *sapi.Sapi, p interface{}) error {
	session := p.(*Session)
	session.sapi = s
	return nil
}

func init() {
	sessionMap = make(map[string]map[string]interface{})

	info := sapi.PluginInfo{}
	info.Creater = sessionCreater
	info.HookPluginRequestInit = requestInit
	info.BasePlugins = []string{"cookie"}
	(&info).Support(sapi.IdHttp, sapi.IdWebsocket)
	sapi.RegisterPlugin("session", info)
}
