package sapi

import (
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	//base info
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

func newServer() *Server {
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

func NewServer() *Server {
	p := newServer()
	p.Id = IdHttp
	p.Name = "Http"
	p.FullName = "Wgf Http Server"
	p.Handler = &HttpServerHandler{}
	return p
}

//Get Server For Websocket apps.
func NewWebsocketServer() *Server {
	p := newServer()
	p.Id = IdWebsocket
	p.Name = "Websocket"
	p.FullName = "Wgf Websocket Server"
	p.Handler = &WebsocketServerHandler{}
	return p
}

func NewCliServer() *Server {
	p := newServer()
	p.Id = IdCli
	p.Name = "Cli"
	p.FullName = "Wgf Cli Programe"
	p.Handler = &CliServerHandler{}
	return p
}

func (p *Server) NotifyHandlerFinished() {
	p.handlerFinishedC <- true
}

func (p *Server) Basedir() string {
	return p.basedir
}

func (p *Server) Confdir() string {
	return p.basedir + "/conf/"
}

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

func (p *Server) ServerShutdown() {
	//默认调用Handler的Shutdown后再结束。
	//如果连续按两次SIG_INT信号，则立即结束。
waitLoop:
	for {
		//p.Logger.Sysf("wait for server handler[%s, %p] shutdown, send SIG_INT again to skip this step", p.Name, &p.Handler)
		select {
			case <-p.Handler.Shutdown():
				break waitLoop
			case <-time.After(1 * time.Second):
				p.Logger.Sys("time out")
				if (p.sigIntCount>1) {
					break waitLoop
				}
		}
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
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGUSR1)

	for true {
		s := <-c
		switch s {
		case syscall.SIGINT:
			p.sigIntCount++
			if p.sigIntCount==1 {
				p.sigIntC<-true
			}
		}
	}
}
