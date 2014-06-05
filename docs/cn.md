# wgf文档

## 1. wgf开发备注

Hello，wgf是一款基于Golang的编程框架。目标为提供一个尽可能统一的编程环境，提高工作效率。

目前，wgf已完成对http、cli、socket、websocket的支持，其中socket、websocket的还比较初级。

wgf基于扩展机制，本身也内嵌了httpparam、双向路由、动态模版等扩展，将在此文档中一一介绍。

### 1.1 获取wgf

**wgf目前不是library，而是项目工程的一部分，请不要使用go get获取。**

所以，在使用wgf时候，请自行下载到本地，并copy到项目中，或将wgf目录添加到GOPATH中。

```bash

#在project/src下放置自己项目的代码
mkdir -p project/src

#下载wgf，如有以前没有下载的话
git clone git@github.com:walu/wgf
cd wgf
export wgfdir $(pwd)
cd ../

#设置GOPATH
cd project
export GOPATH=$(pwd)
export GOPATH="${GOPATH}:${wgfdir}"
```

> 当然，你的操作完全没有必要向上面这么小白，只要最后一行对了就行了。

> 或者，你可以执行`cp -r $wgfdir/src/wgf $projectdir/src/`


### 1.2 整个wgf框架支持mvc的开发结构

无论是http server，还是cli task，wgf都将他们统一为mvc结构。controller在wgf里叫做action。

请求会有各自server的路由器分配给action，action调用自身及其它逻辑执行任务，执行输出。

建议的文件结构：

* $workspace
	* src/
		* app/
			* action/
			* model/
			* plugin/ (如果有必要的话)
			* lib/ (如果有必要的话)
		* wgf/ (wgf可以直接copy过来，也可以通过GOPATH引入进来)



## 2. 使用wgf进行http编程

http server是wgf第一个支持的server，也将是主要支持的server。


### 2.1 编写第一个action

在wgf http server里，请求都是action来处理的，下面我们编写一个输出hello world的action。

流程：

1. 声明一个action。
2. 在wgf中注册此action。

创建文件：$workspace/app/action/index.go
```go
//app/action/index.go

package action

import (
	"wgf/sapi"
)


//声明IndexAction
type IndexAction struct {
	sapi.Action //内嵌此Action，省得自己实现
}

//可以将逻辑写在Execute方法里。
//也可以写在DoGet\DoPost方法里。
func (p *IndexAction) Execute() {
	p.Sapi.Println("hello world!")
	return nil
}

//注册IndexAction到wgf
func init() {
	sapi.RegisterAction("index", func() sapi.ActionInterface { return &IndexAction{} })
}
````
* 首先，声明了一个IndexAction。
* 然后，实现*IndexAction的一个方法，我们的逻辑写在此方法里。
* 最后，注册此action到wgf中，名为"index"，index则为其路由参数。
	* 路由参数：如果url最后的路由结果与其相等，则调用此Action。

创建文件：$workspace/main.go
```go
//可以仿照wgf源码中src/demo**系列
package main

import (
	_ "app/action" //这不是一个好习惯，暂时为了介绍，姑且这样
	"wgf"
)

func main() {
	wgf.StartHttpServer()
}
````

好了，执行go run src/main.go

打了浏览器访问：127.0.0.1:8080，是不是看到输出了？

## 2.2 使用router

wgf/plugin/router 支持http与websocket server，完整的支持双向路由。

这里我们简单的介绍url中的r参数，router的其它功能将在下文与其文档中介绍。

现在，我们将注册action时的name调整一下，改为"guess"。
```go
func init() {
        sapi.RegisterAction("guess", func() sapi.ActionInterface { return &IndexAction{} })
}
````

运行`go run src/main.go`，再打开浏览器是不是看不到输出了？

现在把url改为：http://127.0.0.1:8080/?r=guess 再试试？

> 默认情况下，wgf/plugin/router 通过url中的r参数来进行路由。

### 2.2.1 url解析与生成

wgf/plugin/router 在wgf中负责为http、websocket server分发请求，生成url。

其支持双向路由，可以按照预订的规则，根据url计算出接收其请求的action，或根据action与参数按照既定的格式生成url。

```go
import "wgf/plugin/router"

func (p *IndexAction) Execute() {

	//生成到IndexAction的Url，没有参数
	url := router.Url("IndexAction", nil)
	p.Sapi.Println(url)

	//生成到IndexAction的url，有参数
	param := map[string]string {
		"uname" : "wgf",
		"repo"  : "github.com/walu/wgf",
	}
	url = router.Url("IndexAction", param)
	p.Sapi.Println(url)

	return nil
}
````

上面，我们生成了简单的url。但他们将输出类似: ?r=IndexAction&uname=wgf。

如果我们对url做了rewrite，wgf是否能直接生成对应的url呢？当然是可以的。

假设我们希望将: /person/wgf 请求转发给IndexAction，并将wgf设置到HttpGet的uname参数中。

很简单，打开$workspace/conf/router.ini，输入：

```ini
/person/#uname# = IndexAction
````

一条配置，不仅将url转发到IndexAction，还会将我们用router.Url生成的url自动格式化成相应的格式。





## 2.3 接收参数

wgf/plugin/httpparam 扩展提供了接收http请求的功能，支持Get\Post\File上传。

下面，我们改写一下上面action的执行逻辑。

```go
//在import中增加以下pkg
import "wgf/plugin/httpparam"

func (p *IndexAction) Execute() {
	pParam := p.Sapi.Plugin("httpparam").(*httpparam.Param)
	
	//获取Get参数
	var uname string
	uname = pParam.Get.Get("uname")

	//获取post参数
	var post string
	post = pParam.Post.Get("post")

	p.Sapi.Println("httpGetParam: uname=", uname)
	p.Sapi.Println("httpGetParam: post=", post)
	return nil
}
````

httpparam.Param还支持文件上传。

通过pParam.File.Get("keyname") 可以获取文件的multipart.File与文件名称。

通过pParam.File.Move可以将上传的文件mv到指定的地方，以供后续操作。


## 2.4 使用模板

好了，上面我们都是通过p.Sapi.Println直接向客户端输出信息。

测试、学习或者特殊用途可以，如果用它来开发web系统，这当然是不适合的。

wgf针对http server提供了wgf/plugin/view，基于html/template，具有完善的语法支持，并且wgf会自动生效改变后的模板（当然，也可以禁掉这个机制）。

下面，我们通过wgf/plugin/view，将用户输入的uname输出到页面中，并用h1标签包括起来。

### 2.4.1 修改action逻辑

我们先修改action逻辑，将数据传给view层，而不是直接输出。

```go
//在import中增加以下pkg
import (
	"wgf/plugin/httpparam"
	"wgf/plugin/view"
)

func (p *IndexAction) Execute() {
	pParam := p.Sapi.Plugin("httpparam").(*httpparam.Param)
	
	//获取Get参数
	var uname string
	uname = pParam.Get.Get("uname")

	pView := p.Sapi.Plugin("view").(*view.View)
	pView.Assign("uname", uname) //Assign支持任意类型的参数

	//系统会去viewDir中寻找uname.tpl并加载。
	//viewDir默认为$workspace/view/
	//也可以修改$workspace/conf/wgf.ini里的wgf.view.dir参数。
	pView.Display("uname.html")
	return nil
}
````

### 2.4.2 编写模板文件

新建文件$workspace/view/uname.html

```html
<html><body>
	<h1>{{.uname}}</h1>
</body></html>
````

重启系统，现在打开浏览器，是否看到一个标准的html页面了？

### 2.4.3 模板函数

wgf/plugin/view支持模板函数，并内嵌了wgfInclude与wgfUrl两个函数用来引入其它模板文件与生成url。

用法分别为：

```
{{wgfInclude "header.tpl" .}}


{{wgfUrl "indexAction"}}
{{wgfUrl "indexAction" "uname" .uname}} //indexAction, get参数为uname={.uname}
````

## 2.5 推荐模式

这里针对action提一句。

一般项目中，action往往会有很多公共的逻辑，所以在app/action/base/base.go里创建一个baseAction(引入sapi.Action)，替代上面sapi.Action的位置，更有利于代码的组织。

另外，业务逻辑的代码，往往存在于app/model中，不要直接写在action中。

这些faq相信熟悉mvc编程的人都了然于心。


# 3. 使用wgf编写cli程序

## 3.1 编写与启动

## 3.2 -a 参数


