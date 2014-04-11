// Copyright 2014 The Wgf Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sapi

import (
	"fmt"
	"io"
	"os"
	"net"
	"net/http"
	"runtime"
	"runtime/debug"

	"wgf/lib/log"
	"wgf/sapi/websocket"
)

type Sapi struct {
	server *Server
	closed bool

	//Log
	Logger *log.Logger

	//by golang http package
	//不要用这两个属性，不保证兼容性，一旦解决，立马变为不再导出。
	Res http.ResponseWriter
	Req *http.Request

	//IO
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	//Response Header Status
	Status int

	//Plugins
	plugins map[string]interface{}

	actionName string

	//about runtime
	requestChannel chan int
}

func (p *Sapi) SetActionName(name string) {
	p.actionName = name
}

//中止当前请求，之后的代码不会再执行。但plugin的requestShutdown会执行。
//建议只在action逻辑中执行。
func (p *Sapi) ExitRequest() {
	p.requestChannel <- 1
	runtime.Goexit()
}

func (p *Sapi) Close() {
	p.closed = true
}

func (p *Sapi) IsClosed() bool {
	return p.closed
}

//输出内容给客户端，第一次输出之前会先输出header信息
func (p *Sapi) Print(val interface{}) (int, error) {
	return fmt.Fprint(p.Stdout, val)
}

func (p *Sapi) Println(val interface{}) (int, error) {
	return fmt.Fprintln(p.Stdout, val)
}

//获取与当前请求相关的plugin，通常是指针。
func (p *Sapi) Plugin(name string) interface{} {
	return p.plugins[name]
}

func (p *Sapi) RequestURI() string {
	if nil == p.Req {
		return ""
	}
	return p.Req.URL.Path
}

func (p *Sapi) start(c chan int) {
	p.requestChannel = c
	defer p.ExitRequest()

	defer func() {
		r := recover()
		if nil != r {
			p.Logger.Warning(r)
			p.Logger.Print(string(debug.Stack()))
		}
	}()

	pluginOrders := p.server.PluginOrder
	for _, name := range pluginOrders {
		p.pluginRequestInit(name)
	}

	//execute action
	action, actionErr := GetAction(p.actionName)
	if nil != actionErr {
		p.Logger.Debug("ROUTER[" + p.RequestURI() + "] " + actionErr.Error())
		return
	}

	action.SetSapi(p)
	if !action.UseSpecialMethod() {
		action.Execute()
	} else {
		switch p.Req.Method {
		case "GET":
			action.DoGet()
		case "POST":
			action.DoPost()
		}
	}

	//request shutdown
	for i := len(p.server.PluginOrder) - 1; i >= 0; i-- {
		p.pluginRequestShutdown(p.server.PluginOrder[i])
	}
}

func (p *Sapi) pluginRequestInit(name string) {
	info, ok := pluginMap[name]
	if ok {
		obj, _ := info.Creater()
		if nil != info.HookPluginRequestInit {
			info.HookPluginRequestInit(p, obj)
		}
		p.plugins[name] = obj
	}
}

func (p *Sapi) pluginRequestShutdown(name string) {
	info, ok := pluginMap[name]
	if ok {
		obj, _ := p.plugins[name]
		if nil != info.HookPluginRequestShutdown {
			info.HookPluginRequestShutdown(p, obj)
		}
		delete(p.plugins, name)
	}
}

func NewSapi(pServer *Server, res http.ResponseWriter, req *http.Request) *Sapi {
	s := &Sapi{}
	s.plugins = make(map[string]interface{})

	s.server = pServer
	s.Logger = pServer.Logger
	s.Res = res
	s.Req = req

	s.Stdout = res
	s.Stderr = res
	s.Stdin = req.Body

	return s
}

func NewWebSocketSapi(pServer *Server, conn *websocket.Conn) *Sapi {
	s := &Sapi{}
	s.plugins = make(map[string]interface{})

	s.server = pServer
	s.Logger = pServer.Logger
	s.Res = nil
	s.Req = conn.Request()

	s.Stdout = conn
	s.Stderr = conn
	s.Stdin = conn

	return s
}

func NewCliSapi(pServer *Server) *Sapi {
	s := &Sapi{}
	s.plugins = make(map[string]interface{})

	s.server = pServer
	s.Logger = pServer.Logger

	s.Stdout = os.Stdout
	s.Stderr = os.Stderr
	s.Stdin = os.Stdin

	return s
}

func NewSocketSapi(pServer *Server, conn net.Conn) *Sapi {
	s := &Sapi{}
	s.plugins = make(map[string]interface{})

	s.server = pServer
	s.Logger = pServer.Logger

	s.Stdout = conn
	s.Stderr = conn
	s.Stdin = conn

	return s
}

