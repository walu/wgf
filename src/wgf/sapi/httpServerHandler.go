package sapi

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

type HttpServerHandler struct {
	disabled bool

	pServer *Server

	//listener
	Ln net.Listener

	maxChildren int64
	currentChildren int64
}

func (p *HttpServerHandler) Shutdown() chan bool{
	c := make(chan bool)
	go p.shutdownWorker(c)
	return c
}

func (p *HttpServerHandler) shutdownWorker(c chan bool){
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


func (p *HttpServerHandler) Serve(pServer *Server) {
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
	httpServer := &http.Server{}
	httpServer.Handler = p
	httpServer.Serve(p.Ln)

	//notifyHandlerFinished
	pServer.NotifyHandlerFinished()
}

func (p *HttpServerHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	if p.disabled {
		http.Error(res, "the server is shutting down", 503)
		return
	}

	if p.currentChildren >= p.maxChildren {
		errorMsg := fmt.Sprintf("currentChildren has reached %d, please raise the wgf.sapi.maxChildren", p.currentChildren)
		http.Error(res, errorMsg, 503)
		p.pServer.Logger.Warning(errorMsg)
		return
	}

	//manage currentChildren
	defer func() { p.currentChildren-- }()
	p.currentChildren++

	sapi := NewSapi(p.pServer, res, req)
	c := make(chan int)
	go sapi.start(c)
	select {
		case <-c ://request has been finishied
		case <-res.(http.CloseNotifier).CloseNotify(): //client disconnected
			sapi.Close()
	}
}

