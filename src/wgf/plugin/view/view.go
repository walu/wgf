package view

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"wgf/sapi"
)

//view自定义模版函数
//	wgfInclude, 代替{{template }}，实现不重启编译
//	wgfUrl，生成url
var viewFuncMap map[string]interface{}

var viewTemplate *template.Template
var confViewDir string
var confAutoRefresh bool

type viewCache struct {
	tpl     *template.Template
	modTime time.Time
}

var viewCacheMap map[string]*viewCache

type View struct {
	sapi *sapi.Sapi
	data map[string]interface{}
}

func (v *View) Assign(key string, val interface{}) {
	v.data[key] = val
}

func (v *View) Display(viewName string) error {
	tpl, err := getTemplate(viewName)
	if nil != err {
		v.sapi.Logger.Warning(err)
		return err
	}

	v.Assign("_wgf_view_stdoutWriter", v.sapi.Stdout)
	err = tpl.Execute(v.sapi.Stdout, v.data)
	if nil != err {
		v.sapi.Logger.Warning(err)
	}
	return err
}

func requestInit(p *sapi.Sapi, v interface{}) error {
	view := v.(*View)
	view.sapi = p
	return nil
}

//编译模版
//参数为namePath，模版文件的绝对路径。
//没有分离name、path概念，是为了保持统一。
func compileViewFile(path string) (*template.Template, error) {
	var tpl *template.Template
	var err error

	content, err := ioutil.ReadFile(path)
	if nil != err {
		return nil, err
	}

	tpl, err = template.New("view").Funcs(viewFuncMap).Parse(string(content))
	return tpl, err
}

func getTemplate(viewName string) (*template.Template, error) {
	var namePath string
	var tpl *template.Template

	namePath = confViewDir + "/" + viewName

	if confAutoRefresh {
		fi, err := os.Stat(namePath)

		if nil != err {
			return nil, err
		}

		cache, ok := viewCacheMap[viewName]
		if !ok || cache.modTime.Before(fi.ModTime()) {
			tpl, err = compileViewFile(namePath)
			if nil != err {
				return nil, err
			}

			cache = &viewCache{}
			cache.tpl = tpl
			cache.modTime = fi.ModTime()

			viewCacheMap[viewName] = cache
		}
	}
	return viewCacheMap[viewName].tpl, nil
}

func initTemplate(path string) error {
	viewName := path[len(confViewDir):]
	_, err := getTemplate(viewName)
	return err
}

func viewFuncWgfInclude(viewName string, data interface{}) (string, error) {
	viewData := data.(map[string]interface{})
	stdoutWriter, ok := viewData["_wgf_view_stdoutWriter"]
	if !ok {
		return "", errors.New("undefiend _wgf_view_stdoutWriter")
	}

	var tpl *template.Template
	var err error

	tpl, err = getTemplate(viewName)
	if nil != err {
		return "", err
	}

	err = tpl.Execute(stdoutWriter.(io.Writer), data)
	return "", err
}

func SetViewDir(viewDir string) {
	confViewDir = viewDir
}

func serverInit(pServer *sapi.Server) error {
	var err error
	dir := pServer.Conf.String("wgf.view.dir", pServer.Basedir()+"/view/")
	SetViewDir(dir)

	confAutoRefresh = pServer.Conf.Bool("wgf.view.autoRefresh", true)


	_, err = os.Stat(dir)
	if nil != err {
		pServer.Logger.Warning(err)
		return nil
	}

	err = filepath.Walk(
		confViewDir,
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				return initTemplate(path)
			}
			return nil
		},
	)

	if nil != err {
		pServer.Logger.Warning(err)
	}
	return err
}

func pluginCreater() (interface{}, error) {
	v := new(View)
	v.data = make(map[string]interface{})
	return v, nil
}

func init() {
	viewCacheMap = make(map[string]*viewCache)

	viewFuncMap = map[string]interface{}{
		"wgfInclude": viewFuncWgfInclude,
	}

	info := sapi.PluginInfo{}
	info.Creater = pluginCreater
	info.HookPluginServerInit = serverInit
	info.HookPluginRequestInit = requestInit
	sapi.RegisterPlugin("view", info)
}

func noneUseForPackage() {
	fmt.Println("je;;p")
}
