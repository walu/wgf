package sapi

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"


	"wgf/sapi/websocket"
	"wgf/lib/conf"
	"wgf/lib/log"
)

type Server struct {

	disabled bool

	Conf *conf.Conf

	Logger *log.Logger
	LoggerStdout *log.Logger

	basedir         string
	maxChildren     int64
	currentChildren int64

	listener net.Listener

	//for websocket server
	pWebsocketServer *websocket.Server

	//for terminal server
}

func NewServer() *Server {
	p := &Server{}
	p.Logger = log.NewLogger()
	p.LoggerStdout = log.NewLogger()
	p.LoggerStdout.SetLogWriter(os.Stdout)
	return p
}

func (p *Server) Basedir() string {
	return p.basedir
}

func (p *Server) Confdir() string {
	return p.basedir + "/conf/"
}

func (p *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if p.disabled {
		http.Error(res, "the server is shutting down", 503)
		return
	}
	if p.currentChildren >= p.maxChildren {
		errorMsg := fmt.Sprintf("currentChildren has reached %d, please raise the wgf.sapi.maxChildren", p.currentChildren)
		http.Error(res, errorMsg, 503)
		p.Logger.Warning(errorMsg)
		return
	}

	//manage currentChildren
	defer func() { p.currentChildren-- }()
	p.currentChildren++

	sapi := NewSapi(p, res, req)
	c := make(chan int)
	go sapi.start(c)
	select {
		case <-c ://request has been finishied
		case <-res.(http.CloseNotifier).CloseNotify(): //client disconnected
			sapi.Close()
	}
}

func (p *Server) ServeWebSocket(conn *websocket.Conn) {

	if p.disabled {
		return
	}

	if p.currentChildren >= p.maxChildren {
		errorMsg := fmt.Sprintf("currentChildren has reached %d, please raise the wgf.sapi.maxChildren", p.currentChildren)
		p.Logger.Warning(errorMsg)
		return
	}

	//manage currentChildren
	defer func() { p.currentChildren-- }()
	p.currentChildren++

	sapi := NewWebSocketSapi(p, conn)
	c := make(chan int)
	go sapi.start(c)
	<-c //blocked here, wait for process finished
}



func (p *Server) boot(basedir string, pConf *conf.Conf, handler http.Handler) {
	p.basedir = basedir
	p.Conf = pConf
	p.maxChildren = p.Conf.Int64("wgf.sapi.maxChildren", 1000)

	var logWriter io.Writer
	var logFile string
	var err error

	logFile = p.Conf.String("wgf.sapi.logFile", "")
	if "" == logFile {
		logWriter = os.Stdout
		logFile = "stdout"
	} else {
		logWriter, err = os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if nil != err {
			logWriter = os.Stdout
			p.LoggerStdout.Warningf("cannot open log file for write, error: %s, use stdout instead.", err.Error())
			logFile = "stdout"
		}
	}
	p.Logger.SetLogWriter(logWriter)
	log.ConfLogWriter = logWriter

	timezone := p.Conf.String("wgf.sapi.timezone", "Asia/Shanghai")
	p.Logger.SetTimeLocation(timezone)
	p.LoggerStdout.SetTimeLocation(timezone)
	log.ConfTimeLocationName = timezone

	var tcpListen string
	tcpListen = p.Conf.String("wgf.sapi.tcpListen", "")
	p.listener, err = net.Listen("tcp", tcpListen)
	if nil != err {
		p.Logger.Fatalf("cannot listen to %s, error: %s", tcpListen, err.Error())
		return //exit
	}

	//处理退出、info信号
	go p.handleControlSignal()

	//server init
	pluginOrders := GetPluginOrder()
	for _, name := range pluginOrders {
		p.pluginServerInit(name)
	}

	httpServer := &http.Server{}
	httpServer.Handler = handler
	httpServer.Serve(p.listener)

	//server shutdown
	for i := len(pluginOrders) - 1; i >= 0; i-- {
		p.pluginServerShutdown(pluginOrders[i])
	}

	p.Logger.Info("server shutdown\n")
}



func (p *Server) Init(basedir string, pConf *conf.Conf) {
	p.boot(basedir, pConf, p)
}

func (p *Server) InitWebSocket(basedir string, pConf *conf.Conf) {
	pWebsocketServer := &websocket.Server{}
	pWebsocketServer.Handler = func(conn *websocket.Conn) {
		p.ServeWebSocket(conn)
	}
	p.boot(basedir, pConf, pWebsocketServer)
}



func (p *Server) pluginServerInit(name string) {
	info, ok := pluginMap[name]
	if ok {
		if nil != info.HookPluginServerInit {
			info.HookPluginServerInit(p)
		}
	}
}

func (p *Server) pluginServerShutdown(name string) {
	info, ok := pluginMap[name]
	if ok {
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
			p.disabled = true
			for p.currentChildren > 0 {
				p.Logger.Infof(
					"wait for currentChildren stop, remains %d. use [ kill -9 %d ] if you want to kill it at once.", 
					p.currentChildren,
					os.Getpid(),
				)
				time.Sleep(1*time.Second)
			}
			p.listener.Close()
			return
		}
	}
}
