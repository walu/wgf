package base

import (
	"wgf/sapi"

	"wgf/plugin/header"
	"wgf/plugin/router"
	"wgf/plugin/session"
)

type Action struct {
	sapi.Action
}

func (p *Action) GetSessionUser() (uname string) {
	sess := p.Sapi.Plugin("session").(*session.Session)
	return sess.Get("uname").(string)
}

func (p *Action) Logout() {

}

func (p *Action) DenyIfNotLogin() {
	sess := p.Sapi.Plugin("session").(*session.Session)
	sess.Start()
	uname := sess.Get("uname").(string)
	if "" == uname {
		hd := p.Sapi.Plugin("header").(*header.Header)
		pRouter := p.Sapi.Plugin("router").(*router.Router)
		hd.Redirect(pRouter.Url("login", nil))
		p.Sapi.ExitRequest()
	}
}
