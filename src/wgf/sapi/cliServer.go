// Copyright 2014 The Wgf Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sapi

type CliServerHandler struct {
	disabled bool
	pServer *Server

	actionName string
}

func (p *CliServerHandler) SetActionName(actionName string) {
	p.actionName = actionName
}

func (p *CliServerHandler) Shutdown() chan bool {
	p.disabled = true

	c := make(chan bool)
	go func(){c<-true}()
	return c
}

func (p *CliServerHandler) Serve(pServer *Server) {
	defer p.pServer.NotifyHandlerFinished()

	p.pServer = pServer
	sapi := NewCliSapi(p.pServer)

	sapi.SetActionName(p.actionName)
	c := make(chan int)
	go sapi.start(c)
	<-c
}

