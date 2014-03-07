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

	"wgf/conf"
	"wgf/sapi/websocket"
)

type Server struct {

	disabled bool

	Conf *conf.Conf

	//Location
	Location *time.Location

	//log writer
	LogWriter io.Writer

	basedir         string
	maxChildren     int64
	currentChildren int64

	listener net.Listener

	//for websocket server
	pWebsocketServer *websocket.Server

	//for terminal server
}

func (p *Server) Basedir() string {
	return p.basedir
}

func (p *Server) Confdir() string {
	return p.basedir + "/conf/"
}

func (p *Server) Log(log interface{}) {
	logTime := time.Now().In(p.Location).Format(time.RFC3339)
	fmt.Fprintf(p.LogWriter, fmt.Sprintf("%s %s\n", logTime, log))
}

func (p *Server) LogStderr(log interface{}) {
	logTime := time.Now().In(p.Location).Format(time.RFC3339)
	fmt.Fprintf(os.Stderr, fmt.Sprintf("%s %s\n", logTime, log))
}

func (p *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if p.disabled {
		http.Error(res, "the server is shutting down", 503)
		return
	}
	if p.currentChildren >= p.maxChildren {
		errorMsg := fmt.Sprintf("currentChildren has reached %d, please raise the wgf.sapi.maxChildren", p.currentChildren)
		http.Error(res, fmt.Sprintf("currentChildren has reached the max", p.currentChildren), 503)
		p.Log(errorMsg)
		return
	}

	//manage currentChildren
	defer func() { p.currentChildren-- }()
	p.currentChildren++

	sapi := NewSapi(p, res, req)
	c := make(chan int)
	go sapi.start(c)
	<-c //blocked here, wait for process finished
}

func (p *Server) ServeWebSocket(conn *websocket.Conn) {
	if p.currentChildren >= p.maxChildren {
		errorMsg := fmt.Sprintf("currentChildren has reached %d, please raise the wgf.sapi.maxChildren", p.currentChildren)
		//http.Error(conn, fmt.Sprintf("currentChildren has reached the max", p.currentChildren), 503)
		p.Log(errorMsg)
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


func (p *Server) Init(basedir string, pConf *conf.Conf) {
	p.basedir = basedir
	p.Conf = pConf
	p.maxChildren = p.Conf.Int64("wgf.sapi.maxChildren", 1000)

	var logFile string
	var err error

	timezone := p.Conf.String("wgf.sapi.timezone", "Asia/Shanghai")
	p.Location, err = time.LoadLocation(timezone)
	if nil != err {
		p.Log(err)
	}

	logFile = p.Conf.String("wgf.sapi.logFile", "")
	if "" == logFile {
		p.LogWriter = os.Stderr
		logFile = "stderr"
	} else {
		p.LogWriter, err = os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if nil != err {
			p.LogWriter = os.Stderr
			p.Log("can't open log file " + logFile)
			logFile = "stderr"
		}
	}

	p.LogStderr("log will be writed into " + logFile)

	var tcpListen string
	tcpListen = p.Conf.String("wgf.sapi.tcpListen", "")
	p.listener, err = net.Listen("tcp", tcpListen)
	if nil != err {
		p.Log("cannot listen to " + tcpListen)
		return //exit
	}

	//处理退出、info信号
	go p.handleControlSignal()

	//server init
	pluginOrders := GetPluginOrder()
	for _, name := range pluginOrders {
		p.pluginServerInit(name)
	}

	//改到http server试试
	httpServer := &http.Server{}
	httpServer.Handler = p
	httpServer.Serve(p.listener)

	//server shutdown
	for i := len(pluginOrders) - 1; i >= 0; i-- {
		p.pluginServerShutdown(pluginOrders[i])
	}

	p.Log("server shutdown\n")

}

func (p *Server) InitWebSocket(basedir string, pConf *conf.Conf) {
	p.basedir = basedir
	p.Conf = pConf
	p.maxChildren = p.Conf.Int64("wgf.sapi.maxChildren", 1000)

	var logFile string
	var err error

	timezone := p.Conf.String("wgf.sapi.timezone", "Asia/Shanghai")
	p.Location, err = time.LoadLocation(timezone)
	if nil != err {
		p.Log(err)
	}

	logFile = p.Conf.String("wgf.sapi.logFile", "")
	if "" == logFile {
		p.LogWriter = os.Stderr
		logFile = "stderr"
	} else {
		p.LogWriter, err = os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if nil != err {
			p.LogWriter = os.Stderr
			p.Log("can't open log file " + logFile)
			logFile = "stderr"
		}
	}

	p.LogStderr("log will be writed into " + logFile)

	var tcpListen string
	tcpListen = p.Conf.String("wgf.sapi.tcpListen", "")
	p.listener, err = net.Listen("tcp", tcpListen)
	if nil != err {
		p.Log("cannot listen to " + tcpListen)
		return //exit
	}

	//处理退出、info信号
	go p.handleControlSignal()

	//server init
	pluginOrders := GetPluginOrder()
	for _, name := range pluginOrders {
		p.pluginServerInit(name)
	}

	pWebsocketServer := &websocket.Server{}
	pWebsocketServer.Handler = func(conn *websocket.Conn) {
		p.ServeWebSocket(conn)
	}
	httpServer := &http.Server{}
	httpServer.Handler = pWebsocketServer
	httpServer.Serve(p.listener)

	//server shutdown
	//因为httpServer无法接收信号退出，导致这个地方无法执行。。。想办法中。。
	for i := len(pluginOrders) - 1; i >= 0; i-- {
		p.pluginServerShutdown(pluginOrders[i])
	}

	p.Log("server shutdown\n")

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
				p.Log(fmt.Sprintf("wait for currentChildren stop, remains %d", p.currentChildren))
				time.Sleep(1*time.Second)
			}
			p.listener.Close()
			return
		}
	}
}
