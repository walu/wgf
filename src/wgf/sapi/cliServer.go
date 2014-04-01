package sapi

type CliServerHandler struct {
	disabled bool

	pServer *Server
}

func (p *CliServerHandler) Shutdown() chan bool {
	p.disabled = true

	c := make(chan bool)
	go func(){c<-true}()
	return c
}

func (p *CliServerHandler) Serve(pServer *Server) {
	p.pServer = pServer

	sapi := NewCliSapi(p.pServer)
	sapi.SetActionName("index")
	c := make(chan int)
	go sapi.start(c)
	<-c
}

