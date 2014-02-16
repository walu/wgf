# Wgf 文档

## wgf 是什么？

wgf是基于[Golang](<golang.org>)的web框架，目前做wgf的目标只有一点：

* 简单易用、方便灵活。

> **强烈建议大家先浏览一下app目录源码（一个index＋login事例），即可对wgf有个大体的了解。**

灵活的扩展机制使Wgf的一大特点，借鉴了php的设计理念，简化新功能的增添，保持核心结构与整体理念的稳定。

做wgf的一个初衷，是为了将php在web开发中的优势转移到编译型语言中。C之类的
语言是无法想像这件事情的，对java又一点也没感觉，于是go成了最佳的选择。通过
精心的设计，使得框架即有php的灵活之处，又不失golang本身的性能与其它长处。

## wgf支持的功能

* mvc
	* 可以查看源码中的app目录，查看基本的mvc使用。
* plugin
	* httpparam
		* 获取http请求中的参数，GET、POST、文件上传等。
	* session (完善中)
		* 处理Session问题。
	* cookie
		* 获取、设置Cookie。
	* header
		* 获取、设置Header信息、重定向请求等。
	* router
		* 根据路由规则分发请求、生成URL等。
	* view
		* 管理模版文件。

