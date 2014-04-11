package socketAction

import (
	"io"

	"wgf/sapi"
)

type IndexAction struct {
	sapi.Action
}

func (p *IndexAction) Execute() error {
	io.Copy(p.Sapi.Stdout, p.Sapi.Stdin)
	p.Sapi.Stdout.Write([]byte("hello \n"))
	return nil
}

func init() {
	sapi.RegisterAction("index", func() sapi.ActionInterface { return &IndexAction{} })
}
