//获取http请求的Get、Post参数，以及上传的文件.
//
//	//获取方法
//	param := sapi.Plugin("httpparam").(*httpparam.Param)
package httpparam

import (
	"errors"
	"io"
	"mime/multipart"
	"net/url"
	"os"

	"wgf/sapi"
)

const (
	GET  = "GET"
	POST = "POST"
)

type ParamFiles map[string][]*multipart.FileHeader

//获取请求中名字为key的第一个文件
func (pf ParamFiles) Get(key string) (f multipart.File, name string, err error) {
	fhList, ok := pf[key]
	if !ok || len(fhList) <= 0 {
		err = errors.New("none file named " + key)
		return
	}

	name = fhList[0].Filename
	f, err = fhList[0].Open()
	return
}

//将上传的文件move到指定位置
func (pf ParamFiles) Move(key string, path string) error {
	var f multipart.File
	var err error
	var pFile *os.File

	f, _, err = pf.Get(key)
	if nil != err {
		return err
	}

	pFile, err = os.Create(path)
	if nil != err {
		return err
	}

	_, err = io.Copy(pFile, f)
	if nil != err {
		return err
	}

	return nil

}

type Param struct {
	Get  url.Values
	Post url.Values
	File ParamFiles
}

func requestInit(sapi *sapi.Sapi, plugin interface{}) error {
	p := plugin.(*Param)

	p.Get = url.Values{}
	p.Post = url.Values{}
	p.File = ParamFiles{}

	//http的post编码有两种
	sapi.Req.ParseForm()
	sapi.Req.ParseMultipartForm(1024 * 1024 * 10)
	sapi.Req.Body.Close()

	p.Get, _ = url.ParseQuery(sapi.Req.URL.RawQuery)

	if nil != sapi.Req.MultipartForm {
		for key, val := range sapi.Req.MultipartForm.Value {
			for _, v := range val {
				p.Post.Add(key, v)
			}
		}
		p.File = ParamFiles(sapi.Req.MultipartForm.File)

	} else {
		p.Post = sapi.Req.PostForm
	}

	return nil
}

func newParam() (interface{}, error) {
	return &Param{}, nil
}

func init() {
	info := sapi.PluginInfo{}
	info.Creater = newParam
	info.HookPluginRequestInit = requestInit
	(&info).Support(sapi.IdHttp, sapi.IdWebsocket)
	sapi.RegisterPlugin("httpparam", info)
}
