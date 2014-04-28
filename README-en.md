# Wgf Documentation

> * Source: http://github.com/walu/wgf
> * Weibo: http://weibo.com/walu

## What is it？

[wgf](<http://github.com/walu/wgf>) is a web framework based on [Golang](<golang.org>), with one primary goal:

* To offer a supporting framework for programming Http, Websocket and CLI.

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

## Usage Introduction (In progress)

> For documentation, you may refer to the godoc
> The introduction here is to note some things that do not fit into the godocs.

* Introduction
	* [Introductory Demo](<docs/1-1-the-first-demo.md>)
* Core extensions
	* httpparam
		* Access to Get / Post parameters.
		* Handling of file uploads.
	* session
		* Handles sessions.
		* Replacing session preservation
	* cookie
		* Manages cookies.
	* header
	* router
		* Routing rules
		* URL generation
		* Setting of route properties
	* view
		* Basic usage
		* View configuration parameters

* Advanced Introduction
	* wgf architecture introduction
	* wgf extension development
