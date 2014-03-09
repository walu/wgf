//管理Url解析与生成，支持双向路由.
//
// 从url解析出当前要执行的action。
// 生成对应于action的url。
//
//	//获取方法
//	pRouter = sapi.Plugin("router").(*router.Router)
package router

import (
	"wgf/lib/conf"
	"wgf/plugin/httpparam"
	"wgf/sapi"

	"fmt"
	"net/url"
	"strings"
)

var confEnableRewrite bool
var confDefaultAction string
var confRouterFilePath string
var confRouter *conf.Conf

type Router struct {
	BaseUrl    string
	RequestUrl string
}

func (r *Router) Url(action string, param map[string]string) string {
	re, err := actionToUrl(action, param)
	if nil != err {
		re = "/?r=" + action
		for key, val := range param {
			re = re + "&" + url.QueryEscape(key) + "=" + url.QueryEscape(val)
		}
	}
	return re
}

func serverInit(pServer *sapi.Server) error {
	confDefaultAction = pServer.Conf.String("wgf.router.defaultAction", "index")
	confEnableRewrite = pServer.Conf.Bool("wgf.router.enableRewrite", true)

	confRouterFilePath = pServer.Conf.String("wgf.router.confFile", "router.ini")
	if confRouterFilePath[0] != '/' {
		confRouterFilePath = pServer.Confdir() + confRouterFilePath
	}

	confRouter = conf.NewConf()
	confRouter.ParseFile(confRouterFilePath)

	var err error
	for key, val := range confRouter.Data() {
		err = addRule(key, val)
		if nil != err {
			pServer.Logger.Warningf("router error when addRule, %s=%s, errors: %s", key, val, err.Error())
		}
	}
	return nil
}

func requestInit(pSapi *sapi.Sapi, plugin interface{}) error {
	var action, uri string
	var rewriteParam map[string]string

	param := pSapi.Plugin("httpparam").(*httpparam.Param)
	_, ok := param.Get["r"]
	uri = strings.TrimSpace(pSapi.RequestURI())

	if !ok {
		if uri == "" || uri == "/" {
			action = confDefaultAction
		} else if confEnableRewrite {
			action, rewriteParam = urlToAction(pSapi.RequestURI())
			if len(rewriteParam) > 0 {
				for k, v := range rewriteParam {
					param.Get.Set(k, v)
				}
			}
		}
	} else {
		action = param.Get.Get("r")
		if action == "" {
			action = confDefaultAction
		}
	}

	pSapi.SetActionName(action)
	return nil
}

func newPlugin() (interface{}, error) {
	return &Router{}, nil
}

func init() {
	info := sapi.PluginInfo{}
	info.Creater = newPlugin
	info.HookPluginServerInit = serverInit
	info.HookPluginRequestInit = requestInit
	info.BasePlugins = []string{"httpparam"}
	sapi.RegisterPlugin("router", info)
}
