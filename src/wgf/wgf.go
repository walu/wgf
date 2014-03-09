package wgf

import (
	"flag"
	"os"

	//load sapi module
	"wgf/lib/conf"
	"wgf/sapi"

	//load all core plugins
	_ "wgf/plugin/cookie"
	_ "wgf/plugin/header"
	_ "wgf/plugin/httpparam"
	_ "wgf/plugin/router"
	_ "wgf/plugin/session"
	_ "wgf/plugin/view"
)

var cliArgs map[string]string

func showHelpAndExit() {

}

func visitFlags(f *flag.Flag) {
	switch f.Name {
	//common args
	case "basedir":
		cliArgs["basedir"] = f.Value.String()
	default:
		showHelpAndExit()
	}
}

//启动http(fastcgi)服务器
func StartHttpServer() {

	//parse cli params
	flag.Visit(visitFlags)

	var basedir string
	basedir = cliArgs["basedir"]
	if "" == basedir {
		basedir, _ = os.Getwd()
	}

	var confFile string
	var pConf *conf.Conf
	confFile = basedir + "/conf/wgf.ini"
	pConf = conf.NewConf()
	pConf.ParseFile(confFile)

	//load conf file
	server := &sapi.Server{}
	server.Init(basedir, pConf)
}

func StartWebSocketServer() {
	//parse cli params
	flag.Visit(visitFlags)

	var basedir string
	basedir = cliArgs["basedir"]
	if "" == basedir {
		basedir, _ = os.Getwd()
	}

	var confFile string
	var pConf *conf.Conf
	confFile = basedir + "/conf/wgf.ini"
	pConf = conf.NewConf()
	pConf.ParseFile(confFile)

	//load conf file
	server := &sapi.Server{}
	server.InitWebSocket(basedir, pConf)
}

func init() {
	cliArgs = make(map[string]string)
}
