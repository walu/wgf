package sapi

import (
	"fmt"
)

type PluginInfo struct {
	Creater                   func() (interface{}, error)
	HookPluginServerInit      func(p *Server) error
	HookPluginServerShutdown  func(p *Server) error
	HookPluginRequestInit     func(p *Sapi, plugin interface{}) error
	HookPluginRequestShutdown func(p *Sapi, plugin interface{}) error

	//依赖的其它插件
	//先Init依赖插件
	//后Shutdown依赖插件。
	BasePlugins []string
}

//已注册的plugin
var pluginMap map[string]PluginInfo

//已注册的plugin之间的依赖关系
//与init顺序相同，与shutdown顺序相反
var pluginList []string
var pluginHasOrdered bool

//注册plugin插件。
func RegisterPlugin(name string, hookInfo PluginInfo) {
	pluginMap[name] = hookInfo
}

//获取已注册的plugin列表，已经根据依赖关系排好顺序。排在后面的依赖前面的。
func GetPluginOrder() []string {

	if pluginHasOrdered {
		return pluginList
	}

	pluginHasOrdered = true

	source := map[string][]string{}
	for k, v := range pluginMap {
		source[k] = v.BasePlugins
	}

	var count int
	var canRegister bool
	var lastTry string
	for len(source) > 0 {
		count = 0
		for key, val := range source {
			canRegister = true
			lastTry = fmt.Sprintf("name: %s, relyPlugins: %p", key, val)
			for _, replyPlugin := range val {
				_, ok := source[replyPlugin]
				if true == ok { //doesn't has been registed
					canRegister = false
					break
				}
			}

			if canRegister {
				pluginList = append(pluginList, key)
				delete(source, key)
				count = count + 1
			}
		}

		if count == 0 {
			fmt.Println(fmt.Sprintf("plugin rely relation errors\nplugin registed: %p\nlastTry: %p\n", pluginList, lastTry))
			break
		}
	}

	return pluginList
}

func init() {
	pluginMap = make(map[string]PluginInfo)
}
