package action

import (
	"wgf/plugin/view"
	"wgf/plugin/router"
	"wgf/plugin/header"
	"wgf/plugin/httpparam"
	"wgf/plugin/session"
	"wgf/sapi"

	"app/action/base"
)

type LoginAction struct {
	base.Action
}

func (p *LoginAction) UseSpecialMethod() bool {
	return true
}

func (p *LoginAction) DoGet() error {
	pHttpparam := p.Sapi.Plugin("httpparam").(*httpparam.Param)
	uname := pHttpparam.Get.Get("uname")

	pRouter := p.Sapi.Plugin("router").(*router.Router)
	if "" != uname {
		pSession := p.Sapi.Plugin("session").(*session.Session)
		pSession.Set("uname", uname)

		pHeader := p.Sapi.Plugin("header").(*header.Header)
		pHeader.Redirect(pRouter.Url("index", nil))
		return nil
	}

	view := p.Sapi.Plugin("view").(*view.View)
	view.Assign("title", "登录")
	view.Assign("urlLogin", pRouter.Url("login", nil))
	view.Display("login.tpl")
	return nil
}

func (p *LoginAction) DoPost() error {
	pHttpparam := p.Sapi.Plugin("httpparam").(*httpparam.Param)
	uname := pHttpparam.Post.Get("uname")

	pRouter := p.Sapi.Plugin("router").(*router.Router)
	pSession := p.Sapi.Plugin("session").(*session.Session)
	pSession.Set("uname", uname)

	pHeader := p.Sapi.Plugin("header").(*header.Header)
	pHeader.Redirect(pRouter.Url("index", nil))
	return nil
}

func init() {
	sapi.RegisterAction("login", func() sapi.ActionInterface { return &LoginAction{} })
}
