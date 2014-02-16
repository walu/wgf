//用于读写http请求中的header信息
package header

import (
	"net/http"
	"wgf/sapi"
)

type Header struct {
	sapi *sapi.Sapi
}

func (h *Header) Redirect(loc string) {
	http.Redirect(h.sapi.Res, h.sapi.Req, loc, 302)
}

func (h *Header) Set(key, val string) {
	h.sapi.Res.Header().Set(key, val)
}

func requestInit(sapi *sapi.Sapi, plugin interface{}) error {
	p := plugin.(*Header)
	p.sapi = sapi
	return nil
}

func newPlugin() (interface{}, error) {
	return &Header{}, nil
}

func init() {
	info := sapi.PluginInfo{}
	info.Creater = newPlugin
	info.HookPluginRequestInit = requestInit
	sapi.RegisterPlugin("header", info)
}
