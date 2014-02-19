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
