# 入门Demo

本文档旨在介绍wgf的最基本的使用方法。

wgf框架源码可以放在任何GOPATH所指的地方，这里为了方便，我们将其与app目录放在一块。

目录格式：

* src
	* app/ <- app目录用来组织应用代码
		* action/
			* index.go <- 这里我们仅用到一个index action，输出一个简单的hello world
	* wgf/
	* main.go

**main.go**

```go
/*
main.go
整个项目的入口文件
*/
package main

import (
	_ "app/action" //其实这个地方应该用app/app.go来组织更好，但这个地方为了方便，简化了。
	"wgf"
)

func main() {
	wgf.StartHttpServer()
}
```

**app/action/index.go**
```go
package action

import (
	"wgf/sapi"
)

type IndexAction struct {
}

//方法主体
func (p *IndexAction) Execute() error {
	p.Sapi.Print("hello world\n")
	p.Sapi.Println(p)
	return nil
}

//将action注册进wgf
func init() {
	sapi.RegisterAction("index", func() sapi.ActionInterface { return &IndexAction{} })
}
```
