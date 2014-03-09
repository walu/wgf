package app

import (
	_ "app/action"

	"wgf/plugin/router"
	"wgf/plugin/session"
	"wgf/plugin/view"
	"wgf/sapi"
)

//server-bootstrap操作
//执行一些server级配置与初始化操作
func appServerInit(pServer *sapi.Server) error {

	//view.SetViewDir("/Users/walu/webroot/wiki-public/view/")
	//database init
	//other service init
	return nil
}

//request-bootstrap操作
func appRequestInit(sapi *sapi.Sapi, plugin interface{}) error {
	sess := sapi.Plugin("session").(*session.Session)
	sess.Start()
	uname := sess.Get("uname").(string)

	view := sapi.Plugin("view").(*view.View)
	view.Assign("siteName", "wgf Demo")
	view.Assign("sessionUname", uname)
	view.Assign("urlLogout", sapi.Plugin("router").(*router.Router).Url("logout", nil))
	return nil
}

func appRequestShutdown(sapi *sapi.Sapi, plugin interface{}) error {
	return nil
}

func newAppHook() (interface{}, error) {
	return nil, nil
}

func init() {
	info := sapi.PluginInfo{}
	info.Creater = newAppHook
	info.HookPluginServerInit = appServerInit
	info.HookPluginRequestInit = appRequestInit
	info.BasePlugins = []string{"session", "view", "router"}
	sapi.RegisterPlugin("_app", info)
}
