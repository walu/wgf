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
	"wgf/lib/conf"
	"wgf/sapi/websocket"
)

type Sapi struct {
	server *Server
	closed bool

	//Log
	Logger *log.Logger

	//conf
	Conf *conf.Conf

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

	//HandlerInfo
	HandlerInfo interface{}

	//Plugins
	plugins map[string]interface{}

	actionName string

	//about runtime
	requestChannel chan int
	actionChannel chan error
}

//Set the actionName
func (p *Sapi) SetActionName(name string) {
	p.actionName = name
}

//中止当前请求，之后的代码不会再执行。但plugin的requestShutdown会执行。
//建议只在action逻辑中执行。
func (p *Sapi) ExitRequest() {
	p.actionChannel <- nil
	runtime.Goexit()
}

func (p *Sapi) Close() {
	p.closed = true
}

func (p *Sapi) IsClosed() bool {
	return p.closed
}

/*
输出内容给客户端。

对于HttpServer，第一次输出之前会先输出header信息
*/
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

func (p *Sapi) start(c chan int) error {

	var err error

	p.actionChannel = make(chan error)

	defer func() {
		r := recover()
		if nil != r {
			p.Logger.Warning(r)
			p.Logger.Print(string(debug.Stack()))
		}
		c<-1
	}()

	pluginOrders := p.server.PluginOrder
	for _, name := range pluginOrders {
		err = p.pluginRequestInit(name)
		if nil!=err {
			p.Logger.Debugf("sapi_request_init_error name[%s] error[%s]", name, err.Error())
			return err
		}
	}

	//execute action, it's will call runtime.Goexit, so fire a new goroutine
	go p.executeAction()
	err = <-p.actionChannel

	//request shutdown
	for i := len(p.server.PluginOrder) - 1; i >= 0; i-- {
		p.pluginRequestShutdown(p.server.PluginOrder[i])
	}
	return err
}

func (p *Sapi) executeAction() {
	staticAction, ok := staticActionMap[p.actionName]
	if ok {
		//static action
		p.actionChannel <- staticAction.Execute(p)
		return
	}

	//dynamic action
	var err error
	action, err := GetAction(p.actionName)
	if nil == err {
		action.SetSapi(p)
		if !action.UseSpecialMethod() {
			err = action.Execute()
		} else {
			switch p.Req.Method {
			case "GET":
				err =action.DoGet()
			case "POST":
				err = action.DoPost()
			}
		}
	}
	p.actionChannel <- err
}

func (p *Sapi) pluginRequestInit(name string) error {
	info, ok := pluginMap[name]
	var err error 
	if ok {
		obj, _ := info.Creater()
		if nil != info.HookPluginRequestInit {
			err = info.HookPluginRequestInit(p, obj)
		}
		p.plugins[name] = obj
	}
	return err
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

func NewHttpSapi(pServer *Server, res http.ResponseWriter, req *http.Request) *Sapi {
	s := &Sapi{}
	s.plugins = make(map[string]interface{})

	s.server = pServer
	s.Logger = pServer.Logger
	s.Conf	 = pServer.Conf
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
	s.Conf	 = pServer.Conf
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
	s.Conf	 = pServer.Conf
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
	s.Conf	 = pServer.Conf

	s.Stdout = conn
	s.Stderr = conn
	s.Stdin = conn

	return s
}

