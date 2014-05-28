
wgf, a wonderful framework in golang.

## 文档(Documentation)

* [中文](<docs/cn.md>)
* [English](<docs/en.md>)

# Wgf

> * 源码: http://github.com/walu/wgf
> * 微博: http://weibo.com/walu

## 是什么？

Hello，[wgf](<http://github.com/walu/wgf>)是基于[Golang](<golang.org>)的的编程框架。目标为提供一个尽可能统一的编程环境，提高工作效率。

目前，wgf已完成对http、cli、socket、websocket的支持，其中socket、websocket的还比较初级。
wgf基于扩展机制，本身也内嵌了httpparam、双向路由、动态模版等扩展，将在文档中一一介绍。

> **强烈建议大家先浏览一下app目录源码（一个index＋login事例），即可对wgf有个大体的了解。**

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

* mvc
	* 可以查看源码中的app目录，查看基本的mvc使用。
* plugin
	* httpparam, 获取http请求中的参数，GET、POST、文件上传等。
	* session (完善中), 处理Session问题。
	* cookie, 获取、设置Cookie。
	* header, 获取、设置Header信息、重定向请求等。
	* router, 根据路由规则分发请求、生成URL等。
	* view, 管理模版文件，无重启更新模版。

