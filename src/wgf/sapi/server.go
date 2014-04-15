// Copyright 2014 The Wgf Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sapi

import (
	"io"
	"os"
	"os/signal"
	"syscall"

	"wgf/lib/conf"
	"wgf/lib/log"
)

const (
	IdHttp = 1
	IdWebsocket = 2
	IdSocket = 3
	IdCli = 4
)

type ServerHandler interface {
	Serve(p *Server)
	Shutdown() chan bool
}

type Server struct {
	Id int
	Name string
	FullName string

	Conf *conf.Conf
	Logger *log.Logger
	LoggerStdout *log.Logger

	Handler ServerHandler

	PluginOrder []string

	shutdownC chan bool //限制Shutdown只发生一次
	sigIntC		chan bool //接收SIG_INT信号，用于强制结束程序
	handlerFinishedC chan bool //handler处理完毕，结束程序
	sigIntCount int	//SIG_INT信号次数

	basedir         string

}

//创建一个Server模版。
func NewServer() *Server {
	p := &Server{}
	p.Logger = log.NewLogger()
	p.LoggerStdout = log.NewLogger()
	p.LoggerStdout.SetLogWriter(os.Stdout)

	p.shutdownC = make(chan bool, 1)
	p.shutdownC <- true //执行过程只赋值这一次

	p.sigIntC = make(chan bool)
	p.handlerFinishedC = make(chan bool)
	return p

}

//get server for http apps
func NewHttpServer() *Server {
	p := newServer()
	p.Id = IdHttp
	p.Name = "Http"
	p.FullName = "Wgf Http Server"
	p.Handler = &HttpServerHandler{}
	return p
}

//get server for websocket apps.
func NewWebsocketServer() *Server {
	p := newServer()
	p.Id = IdWebsocket
	p.Name = "Websocket"
	p.FullName = "Wgf Websocket Server"
	p.Handler = &WebsocketServerHandler{}
	return p
}

//get server for cli programe
func NewCliServer() *Server {
	p := newServer()
	p.Id = IdCli
	p.Name = "Cli"
	p.FullName = "Wgf Cli Programe"
	p.Handler = &CliServerHandler{}
	return p
}

//get server for socket apps
func NewSocketServer() *Server {
	p := newServer()
	p.Id = IdSocket
	p.Name = "Socket Server"
	p.FullName = "Wgf Socket Server"
	p.Handler = &SocketServerHandler{}
	return p
}

//告知Server，Handler已经执行完毕。Server将立即启动Shutdown流程
func (p *Server) NotifyHandlerFinished() {
	p.handlerFinishedC <- true
}

func (p *Server) Basedir() string {
	return p.basedir
}

func (p *Server) Confdir() string {
	return p.basedir + "/conf/"
}

/*
启动Server，主进程将阻塞住处理任务，直到接收到SIG_INT信号或者NotifyHandlerFinished被调用。

流程：

1. call ServerInit

2. execute Handler's logic

3. call ServerShutdown
*/
func (p *Server) Boot(basedir string, conf *conf.Conf) {
	p.basedir = basedir
	p.Conf = conf

	p.ServerInit()
	go p.Handler.Serve(p)

	//wait for shutdown
	select {
		case <-p.handlerFinishedC: //Handler执行完毕
			p.Logger.Sys("Shutdown Server: Server handler finished")
		case <-p.sigIntC: //SIG_INT信号
			p.Logger.Sys("Shutdown Server: Signal SIG_INT")
	}
	p.ServerShutdown()
}

/*
执行ServerInit流程

1. 各种初始化。

2. 依次调用各个扩展的HookPluginServerInit方法，如果方法返回非nil值，将产生一个panic终止server运行。
*/
func (p *Server) ServerInit() {
	var logWriter io.Writer
	var logFile string
	var err error

	logFile = p.Conf.String("wgf.sapi.logFile", "")
	if "" == logFile {
		logWriter = os.Stdout
		logFile = "stdout"
	} else {
		logWriter, err = os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
		if nil != err {
			logWriter = os.Stdout
			p.LoggerStdout.Warningf("cannot open log file for write, error: %s, use stdout instead.", err.Error())
			logFile = "stdout"
		}
	}
	p.Logger.SetMinLogLevelName(p.Conf.String("wgf.sapi.minLogLevel", "info"))
	p.Logger.SetLogWriter(logWriter)
	log.ConfLogWriter = logWriter
	log.ConfMinLogLevel = p.Logger.MinLogLevel()

	timezone := p.Conf.String("wgf.sapi.timezone", "Asia/Shanghai")
	p.Logger.SetTimeLocation(timezone)
	p.LoggerStdout.SetTimeLocation(timezone)
	log.ConfTimeLocationName = timezone

	//处理退出、info信号
	go p.handleControlSignal()

	//server init
	p.PluginOrder = GetPluginOrder(p)
	for _, name := range p.PluginOrder {
		p.pluginServerInit(name)
	}
	p.Logger.Sys("ServerInit Done")
}

/*
执行shutdown流程

对于SIG_INT信号：

如果第一次发送，系统启动shutdown流程，调用Handler的Shutdown方法。

如果第二次发送，系统将跳过Handler处理流程，直接执行ServerShutdown流程。
*/
func (p *Server) ServerShutdown() {
	//默认调用Handler的Shutdown后再结束。
	//如果连续按两次SIG_INT信号，则立即结束。
	p.Logger.Sysf("wait for server handler[%s, %p] shutdown, send SIG_INT to skip this step", p.Name, &p.Handler)
	select {
		case <-p.Handler.Shutdown():
		case <-p.sigIntC:
	}

	//server shutdown
	for i := len(p.PluginOrder) - 1; i >= 0; i-- {
		p.pluginServerShutdown(p.PluginOrder[i])
	}
	p.Logger.Sys("Server Shutdown\n")
}

func (p *Server) pluginServerInit(name string) {
	info, ok := pluginMap[name]
	if ok {
		p.Logger.Info("PluginServerInit "+name)
		if nil != info.HookPluginServerInit {
			info.HookPluginServerInit(p)
		}
	}
}

func (p *Server) pluginServerShutdown(name string) {
	info, ok := pluginMap[name]
	if ok {
		p.Logger.Info("PluginServerShutdown "+name)
		if nil != info.HookPluginServerShutdown {
			info.HookPluginServerShutdown(p)
		}
	}
}

//处理退出、info信号
func (p *Server) handleControlSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGUSR1)

	for true {
		s := <-c
		switch s {
		case syscall.SIGINT:
			p.sigIntCount++
			p.sigIntC<-true
		}
	}
}
