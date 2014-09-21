// Copyright 2014 The Wgf Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

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

var sessionHandler Handler

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
	valInStore, ok := sessionHandler.Get(s.id, key)
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
	return sessionHandler.Set(s.id, key, buf.Bytes())
}

func (s *Session) Del(key string) bool {
	return sessionHandler.Del(s.id, key)
}

func (s *Session) Destory() bool {
	return sessionHandler.Destory(s.id)
}

func sessionCreater() (interface{}, error) {
	return &Session{hasStarted: false}, nil
}

func serverInit(s *sapi.Server) error {
	if nil != sessionHandler {
		sessionHandler = newDefaultHandler()
	}
	sessionHandler.SetExpireTime(1200)
	sessionHandler.Start()
	return nil
}

func serverShutdown(s *sapi.Server) error {
	sessionHandler.Stop()
	return nil
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
	info.HookPluginServerInit = serverInit
	info.HookPluginServerShutdown = serverShutdown
	info.BasePlugins = []string{"cookie"}
	(&info).Support(sapi.IdHttp, sapi.IdWebsocket)
	sapi.RegisterPlugin("session", info)
}
