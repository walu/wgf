package action

import (
	"app/action/base"

	"wgf/sapi"
	"wgf/plugin/view"
)

type IndexAction struct {
	base.Action
}

func (p *IndexAction) Execute() error {
	p.DenyIfNotLogin()

	view := p.Sapi.Plugin("view").(*view.View)

	view.Assign("title", "首页")

	links := []map[string]string{
		map[string]string{"name":"walu's wiki", "href":"http://www.walu.cc"},
		map[string]string{"name":"golang", "href":"http://golang.org"},
	}
	view.Assign("links", links)
	view.Display("index.tpl")
	return nil
}

func init() {
	sapi.RegisterAction("index", func() sapi.ActionInterface { return &IndexAction{} })
}
