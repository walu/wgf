package sapi

type CliServerHandler struct {
	disabled bool

	pServer *Server
}

func (p *CliServerHandler) shutdown() {
	p.disabled = true
}

func (p *CliServerHandler) Serve(pServer *Server) {
	p.pServer = pServer

	//handle server shutdown
	go func(){
		<-p.pServer.ShutdownNotifyC()
		p.shutdown()
	}()

	sapi := NewCliSapi(p.pServer)
	sapi.SetActionName("index")
	c := make(chan int)
	go sapi.start(c)
	<-c

}

