package cookie

import (
	"net/http"
	"wgf/sapi"
)

type Cookie struct {
	sapi *sapi.Sapi
}

func (p *Cookie) Get(name string) string {
	c, err := p.sapi.Req.Cookie(name)
	if nil != err {
		return ""
	}
	return c.Value
}

func (p *Cookie) Set(newcookie *http.Cookie) error {
	if nil != p.sapi.Res {
		http.SetCookie(p.sapi.Res, newcookie)
	}
	return nil
}

func requestInit(s *sapi.Sapi, p interface{}) error {
	c := p.(*Cookie)
	c.sapi = s
	return nil
}

func newPlugin() (interface{}, error) {
	return &Cookie{}, nil
}

func init() {
	info := sapi.PluginInfo{}
	info.Creater = newPlugin
	info.HookPluginRequestInit = requestInit
	(&info).Support(sapi.IdHttp, sapi.IdWebsocket)
	sapi.RegisterPlugin("cookie", info)
}
