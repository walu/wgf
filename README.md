
## 废弃


wgf, a wonderful framework in golang.

## 文档(Documentation)

* [中文](<https://github.com/walu/wgf-docs/blob/master/cn.md>)
* [English](<https://github.com/walu/wgf-docs/blob/master/en.md>)

# wgf

> * 源码: http://github.com/walu/wgf
> * 文档: http://github.com/walu/wgf-docs
> * 微博: http://weibo.com/walu

## 是什么？

Hello，[wgf](<http://github.com/walu/wgf>)是基于[Golang](<golang.org>)的的编程框架。目标为提供一个尽可能统一的编程环境，提高工作效率。

目前，wgf已完成对http、cli、socket、websocket的支持。
wgf基于扩展机制，对httpServer内嵌了httpparam、双向路由、动态模版等扩展，将在文档中一一介绍。

## 扩展机制

灵活的扩展机制是Wgf的一大特点，借鉴了php的设计理念，简化新功能的增添，保持核心结构与整体理念的稳定。

## Package的依赖关系

Package dependencies(自上至下):

* app
	* action
	* model
* wgf
	* plugin
	* sapi
	* lib

上层依赖下层，下层不依赖上层。

## 支持的功能

* HttpServer
	* 完善的MVC分层，支持session、cookie、router、view等plugin，方便快速开发。
* SocketServer
	* 通过plugin扩展协议，支持keepalive等属性。
* CliServer
* WebsocketServer[试验]
