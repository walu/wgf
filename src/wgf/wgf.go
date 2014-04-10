// Copyright 2014 The Wgf Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package wgf

import (
	"fmt"
	"flag"
	"os"

	//load sapi module
	"wgf/lib/conf"
	"wgf/sapi"

	//for test only.
	//"runtime/pprof"

	//load all core plugins
	_ "wgf/plugin/cookie"
	_ "wgf/plugin/header"
	_ "wgf/plugin/httpparam"
	_ "wgf/plugin/router"
	_ "wgf/plugin/session"
	_ "wgf/plugin/view"
)

const (
	Version = "0.1"
)

var basedir *string
var conffile *string

//set the main flags for all kinds of servers
func setMainFlags() {
	basedir	 = flag.String("basedir", "", "set the basedir, the `pwd` is default")
	conffile = flag.String("conf", "wgf.ini", "set the default filename located in $basedir/conf/, wgf.ini is the default")
}

func parseArgs() {
	if flag.Parsed() {
		return
	}

	setMainFlags()
	flag.Usage = printHelpInfo
	flag.Parse()

	if "" == *basedir {
		var err error
		*basedir, err = os.Getwd()
		if nil != err {
			fmt.Println("error occurs when getting pwd:", err)
			os.Exit(-1)
		}
	}

	if "" == *conffile {
		*conffile = "wgf.ini"
	}
}

func printHelpInfo() {
	fmt.Fprintf(os.Stderr, "wgf is a framework written in go.\n")
	fmt.Fprintf(os.Stderr, "version: %s\n", Version)
	fmt.Fprintf(os.Stderr, "\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n")
}

func initConf() *conf.Conf {
	var confFile string
	var pConf *conf.Conf
	var err error

	pConf = conf.NewConf()
	confFile = *basedir + "/conf/" + *conffile
	err = pConf.ParseFile(confFile)
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	return pConf
}

//启动Http服务器
func StartHttpServer() {
	parseArgs()
	pConf	:= initConf()
	pServer := sapi.NewServer()
	pServer.Boot(*basedir, pConf)
}

//启动Websocket服务器
func StartWebSocketServer() {
	parseArgs()
	pConf	:= initConf()
	pServer := sapi.NewWebsocketServer()
	pServer.Boot(*basedir, pConf)
}

//启动Cli终端程序
func StartCliServer() {
	parseArgs()
	pConf	:= initConf()
	pServer := sapi.NewCliServer()
	pServer.Boot(*basedir, pConf)
}

