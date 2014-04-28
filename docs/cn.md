# wgf文档

## 1. wgf开发备注

### 1.1 wgf暂不支持go get，也没有支持的计划。

因为：

1. 主要原因：代码import中的github.com/walu/ 之类的前缀就像裹脚布，又臭又长。
2. 次要原因：没有version控制。

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
