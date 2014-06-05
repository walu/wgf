## 文档(Documentation)

* [中文](<docs/cn.md>)
* [English](<docs/en.md>)

# Wgf Documentation

> * Source: http://github.com/walu/wgf
> * Weibo: http://weibo.com/walu

## What is it？

Hello,[wgf](<http://github.com/walu/wgf>) is a framework based on [Golang](<golang.org>). Objective is to provide a unified programming environment as possible, improve work efficiency. 

Currently, wgf has completed support for http, cli, socket, websocket.

wgf based extension mechanism, itself embedded httpparam, bi-directional routing, dynamic templates and other extensions will be introduced one by one in this document.

> **Everyone is strongly recommended to first look through the app content source code (one index＋login instance), and have a general understanding of wgf.**

## Extension mechanism

wgf offers a flexible extension mechanism. I takes design ideas from PHP, and adds simple new functions, maintaining core structure and an overall stable concept.

## Package dependencies

Package dependencies (From top to bottom):

* app
 * action
 * model
* wgf
 * plugin
 * sapi
 * lib

The dependencies go from top to bottom (the upper levels depend on the lower levels, while the lower levels do not depend on the upper levels).

## Supported features

* MVC
* Can check source app content, and basic mvc use.
* plugins
* httpparam. Get http request parameters, GET, POST, file uploads, etc.
* session (perfected), handle Session problems.
* cookie, getting, installing Cookie.
* header, getting, installing Header information, redirect requests, etc.
* router, distribute requests according to routing rules, generate URL, etc.
* view, manage template files, update templates without reboot.

