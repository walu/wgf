// Copyright 2014 The Wgf Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Sapi属于wgf框架的核心逻辑层。

主要分为Server、Sapi、ServerHandler三个部分。

Server作为大的容器，接收任务/请求，调用ServerHandler处理。ServerHandler为每一个任务/请求生成独立的Sapi处理。
*/
package sapi
