package sapi

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"wgf/sapi/websocket"
)

type WebsocketServerHandler struct {
	disabled bool

	pServer *Server

	//listener
	Ln net.Listener

	maxChildren int64
	currentChildren int64
}

func (p *WebsocketServerHandler) Serve(pServer *Server) {
	p.pServer = pServer
	p.maxChildren = pServer.Conf.Int64("wgf.sapi.maxChildren", 1000)

	var tcpListen string
	var err error

	tcpListen = pServer.Conf.String("wgf.sapi.tcpListen", "")
	p.Ln, err = net.Listen("tcp", tcpListen)
	if nil != err {
		pServer.Logger.Fatalf("cannot listen to %s, error: %s", tcpListen, err.Error())
		return //exit
	}

	//package: net/http
	pWebsocketServer := &websocket.Server{}
	pWebsocketServer.Handler = func(conn *websocket.Conn) {
		p.ServeWebSocket(conn)
	}

	httpServer := &http.Server{}
	httpServer.Handler = pWebsocketServer
	httpServer.Serve(p.Ln)
}

func (p *WebsocketServerHandler) ServeWebSocket(conn *websocket.Conn) {

	if p.disabled {
		return
	}

	if p.currentChildren >= p.maxChildren {
		errorMsg := fmt.Sprintf("currentChildren has reached %d, please raise the wgf.sapi.maxChildren", p.currentChildren)
		p.pServer.Logger.Warning(errorMsg)
		return
	}

	//manage currentChildren
	defer func() { p.currentChildren-- }()
	p.currentChildren++

	sapi := NewWebSocketSapi(p.pServer, conn)
	c := make(chan int)
	go sapi.start(c)
	<-c //blocked here, wait for process finished
}

func (p *WebsocketServerHandler) Shutdown() chan bool {
	c := make(chan bool)
	go p.shutdownWorkder(c)
	return c
}

func (p *WebsocketServerHandler) shutdownWorkder(c chan bool) {
	p.disabled = true
	for p.currentChildren > 0 {
		p.pServer.Logger.Infof(
			"wait for currentChildren stop, remains %d. use [ kill -9 %d ] if you want to kill it at once.",
			p.currentChildren,
			os.Getpid(),
		)
		time.Sleep(1*time.Second)
	}
	p.Ln.Close()
	c<-true
}
