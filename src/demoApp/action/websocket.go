package action

import (
	"demoApp/action/base"
	"wgf/sapi"
)

type WebSocketAction struct {
	base.Action
}

func (p *WebSocketAction) Execute() error {

	for {
		p.Sapi.Stdout.Write([]byte("hello"))

	}
	return nil
}

func init() {
	sapi.RegisterAction("ws", func() sapi.ActionInterface { return &WebSocketAction{} })
}
