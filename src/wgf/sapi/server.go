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

	shutdownNotifyC chan bool
	sigIntCount int

	basedir         string
}

func newServer() *Server {
	p := &Server{}
	p.Logger = log.NewLogger()
	p.LoggerStdout = log.NewLogger()
	p.LoggerStdout.SetLogWriter(os.Stdout)
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

func (p *Server) Basedir() string {
	return p.basedir
}

func (p *Server) Confdir() string {
	return p.basedir + "/conf/"
}

func (p *Server) ShutdownNotifyC() chan bool {
	if nil == p.shutdownNotifyC {
		p.shutdownNotifyC = make(chan bool)
	}
	return p.shutdownNotifyC
}

func (p *Server) Boot(basedir string, conf *conf.Conf) {
	p.basedir = basedir
	p.Conf = conf

	p.ServerInit()
	p.Handler.Serve(p)
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
	p.Logger.Info("ServerInit Done")
}

func (p *Server) ServerShutdown() {
	if nil != p.shutdownNotifyC && p.sigIntCount<=1 {
		p.Logger.Sysf("wait for server handler[%s, %p] shutdown, send SIG_INT again to skip this step", p.Name, &p.Handler)
		p.shutdownNotifyC <- true
	}

	//server shutdown
	for i := len(p.PluginOrder) - 1; i >= 0; i-- {
		p.pluginServerShutdown(p.PluginOrder[i])
	}
	p.Logger.Sys("server shutdown\n")
	os.Exit(0)
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
//使用SIG_INT来阻止新请求，等待旧请求处理完成后再正式退出。
//使用SIG_KILL来粗暴的直接停掉进程
func (p *Server) handleControlSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGUSR1)

	for true {
		s := <-c
		switch s {
		case syscall.SIGINT:
			p.sigIntCount++
			p.ServerShutdown()
		}
	}
}
