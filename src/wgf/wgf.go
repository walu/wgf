package wgf

import (
	"fmt"
	"flag"
	"os"

	//load sapi module
	"wgf/lib/conf"
	"wgf/sapi"

	//"runtime/pprof"

	//load all core plugins
	_ "wgf/plugin/cookie"
	_ "wgf/plugin/header"
	_ "wgf/plugin/httpparam"
	_ "wgf/plugin/router"
	_ "wgf/plugin/session"
	_ "wgf/plugin/view"
)

var cliArgs map[string]string

var flagBasedir *string
var flagConf *string

func initCliArgs() {
	flag.Parse()
	if "" == *flagBasedir {
		var err error
		*flagBasedir, err = os.Getwd()
		if nil != err {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	cliArgs["basedir"]	= *flagBasedir
	cliArgs["conf"]		= *flagConf
}

func initConfWithCliArgs() *conf.Conf {
	initCliArgs()
	basedir := cliArgs["basedir"]

	var confFile string
	var pConf *conf.Conf
	var err error

	pConf = conf.NewConf()
	confFile = basedir + "/conf/" + cliArgs["conf"]
	err = pConf.ParseFile(confFile)
	if nil != err {
		fmt.Println(err)
		os.Exit(-1)
	}
	return pConf
}

//启动http(fastcgi)服务器
func StartHttpServer() {
	//load conf file
	pConf	:= initConfWithCliArgs()
	pServer := sapi.NewServer()
	pServer.Boot(cliArgs["basedir"], pConf)
}

func StartWebSocketServer() {
	//load conf file
	pConf	:= initConfWithCliArgs()
	pServer := sapi.NewWebsocketServer()
	pServer.Boot(cliArgs["basedir"], pConf)
}

func StartCliServer() {
	//flagActionName := flag.String("action", "index", "set the action name")

	//load conf file
	pConf	:= initConfWithCliArgs()
	pServer := sapi.NewCliServer()
	pServer.Boot(cliArgs["basedir"], pConf)
}



func init() {
	cliArgs = make(map[string]string)

	flagBasedir	= flag.String("basedir", "", "set the basedir, the `pwd` is default")
	flagConf	= flag.String("conf", "wgf.ini", "set the default filename in conf/, wgf.ini is the default")
}

