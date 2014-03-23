package cliAction

import (
	"wgf/sapi"
)

type IndexAction struct {
	sapi.Action
}

func (p *IndexAction) Execute() error {
	p.Sapi.Println("hello world")
	p.Sapi.Println("Success!!!!!!!!")
	return nil
}

func init() {
	sapi.RegisterAction("index", func() sapi.ActionInterface { return &IndexAction{} })
}
