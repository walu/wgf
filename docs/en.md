>
> This Document was maked with Google Translate
>
> You are welcome to improve it.
>
> **Forgive my poor english, thanks!**
>

# wgf document

Hello, wgf is a programming framework based Golang. Objective is to provide a unified programming environment as possible, improve work efficiency. 

Currently, wgf has completed support for http, cli, socket, websocket which socket, websocket is still relatively junior. 

wgf based extension mechanism, itself embedded httpparam, bi-directional routing, dynamic templates and other extensions will be introduced one by one in this document.

# 1. Wgf development Remarks

### 1.1 get wgf

wgf is not a standalone library, it is a part of your project's source code. When using wgf, please download to your local disk, and copy to the project, or to add wgf path into GOPATH.

```bash
# Place the code in their own projects under project/src
mkdir -p project / src

# Download wgf, if not previously downloaded the words
git clone git@github.com:walu/wgf
cd wgf
export wgfdir $(pwd)
cd ../

# Set GOPATH
cd project
export GOPATH = $(pwd)
export GOPATH = "${GOPATH}:${wgfdir}"
````

> Of course , you are not absolutely necessary to the operation of the above so white, so long as the last line on the line.
> 
> you can also do this: `cp -r $wgfdir/src/wgf $projectdir/src`

### 1.2 supports the development of the structure of the entire wgf mvc framework of

the http server, or cli task, and the rest of patterns, wgf are they unified mvc structure. controller called the action in wgf inside .

Request will be assigned to each server in the router action, action that calls itself and other logic to perform the task execution output.

The proposed structure of the sourcecode:

* $workspace
	* src/
	* app/
		* action/
		* model/
		* plugin/ ( If necessary )
		* lib/ ( If necessary )
	* wgf/ (wgf can copy over, you can come through the introduction GOPATH )

## 2 Writing http servers with wgf

http server is wgf's first supported server, will also be supported by the main server.

### 2.1 Programming a first action

In wgf http server , the request is action to deal with , let's write a hello world output of action.

Process:

1. Declare an action.
2. Sign this action in wgf.

Create a file : $workspace/app/action/index.go

```go
//app/action/index.go

package action

import (
	"wgf/sapi"
)


//Declare IndexAction
type IndexAction struct {
	sapi.Action //embed this Action, save yourself works
}

// You can write the logic in the Execute method .
// and Can also be written in DoGet \ DoPost method in .
func (p * IndexAction) Execute () {
	p.Sapi.Println ("hello world!")
}

// Register IndexAction to wgf
func init () {
	sapi.RegisterAction ("index", func()sapi.ActionInterface{return & IndexAction{}})
}
````

1. First, declare a IndexAction.
2. Then , implement a method * IndexAction we write logic in this approach .
3. Finally , registered in this action to wgf called "index", index by its routing parameters .
	* Routing Parameters: url final routing results if its equal, then call this Action.

Create a file : $workspace/main.go

```go
//Can be modeled wgf source in src/demo** Series
package main

import (
	_ "app/action" // This is not a good habit to introduce temporary , so tentatively
	"wgf"
)

func main () {
	wgf.StartHttpServer ()
}
````

Well , do `go run src/main.go`

Hit the browser to visit : 127.0.0.1:8080, see the output is not it?

### 2.2 Use router

wgf/plugin/router supports http and websocket server, complete support bi-directional routing(parse the url to action with multi rewrite rules and create the url to the defined pattern with actioname and params easily).

Here we introduce a simple url in the r parameter , other features will be presented in the following router and its documentation.

Now, we will adjust the registered name when what action , to "guess".

```go
func init () {
        sapi.RegisterAction ("guess", func () sapi.ActionInterface {return & IndexAction {}})
}
````

Run `go run src/main.go`, and then open the browser can not see the output is not it?

Now the url changed : http://127.0.0.1:8080/?r=guess try again?

By default , wgf/plugin/router via the r parameter to the url routing .


### 2.2.1 url parsing and generation

wgf/plugin/router is responsible for the http in wgf in , websocket server to distribute requests generated url.

It supports two-way route , in accordance with the rules book , calculated in accordance url receives its request action, or generate url based on action in accordance with established formats and parameters .

```go
import "wgf/plugin/router"

func(p *IndexAction) Execute() {

	// Generate the IndexAction the Url, no argument
	url: = router.Url("IndexAction", nil)
	p.Sapi.Println(url)

	// Generate the IndexAction the url, there are arguments
	param: = map[string]string {
		"uname": "wgf",
		"repo": "github.com/walu/wgf",
	}
	url = router.Url("IndexAction", param)
	p.Sapi.Println(url)

	return nil
}
````

Above, we generate a simple url. But they will output similar to: ?r =IndexAction&uname=wgf.

If we do a url rewrite, whether wgf can generate corresponding url direct it? Of course we can .

Suppose we want to: /person/wgf forwards the request to IndexAction, and wgf set to HttpGet the uname parameter.

Very simple, open the $workspace/conf/router.ini, enter :

```ini
/person/#uname# = IndexAction
````

A configuration that will not only be forwarded to the url IndexAction, we will be generated by router.Url url automatically formatted into the appropriate format .


## 2.3 receiver parameters 

wgf/plugin/httpparam receives http requests extension provides functionality to support Get \ Post \ File upload. 

Below, we look at the implementation of rewriting logic of the above action. 

```go 
// Increase in the import of the following pkg 
import "wgf/plugin/httpparam" 

func(p * IndexAction) Execute() {
	pParam  := p.Sapi.Plugin("httpparam").(* httpparam.Param)

	// Get Get parameters 
	var uname string 
	uname = pParam.Get.Get ("uname") 

	//Get post parameters 
	var post string 
	post = pParam.Post.Get ("post") 

	p.Sapi.Println ("httpGetParam: uname =", uname) 
	p.Sapi.Println ("httpGetParam: post =", post) 
	return nil 
} 
```` 

httpparam.Param also supports file upload. 

You can get the file name of the file multipart.File by pParam.File.Get("keyname"). 

Mv files can be uploaded to the designated place by pParam.File.Move for subsequent operations.

## 2.4 Templates 

Well, we are all above output information to the client directly through p.Sapi.Println. 

Test, study or special purposes may, if it is used to develop web systems, which of course is not suitable. 

wgf for http server provides wgf/plugin/view, based on html/template, with perfect grammar support, and wgf template will automatically take effect after the change (of course, this mechanism can also be cut off). 

Below, we wgf/plugin/view, the user input uname output to the page, and include them with the h1 tag.


### 2.4.1 modify action logic 

We modify the action logic, the data to the view layer, rather than a direct output. 

```go 
// Increase in the import of the following pkg 
import (
	"wgf/plugin/httpparam" 
	"wgf/plugin/view" 
) 

func (p *IndexAction) Execute() {
	pParam: = p.Sapi.Plugin.("httpparam").(* httpparam.Param). 

	//Get Get parameters 
	var uname string 
	uname = pParam.Get.Get("uname") 

	pView := p.Sapi.Plugin.("view").(*view.View)
	pView.Assign("uname", uname) //Assign to support any type of parameter 

	// The system will go viewDir find uname.tpl and loaded. 
	// viewDir default is $ workspace / view / 
	// You can also modify $ workspace / conf / wgf.ini Lane wgf.view.dir parameters. 
	pView.Display ("uname.html") 
	return nil 
} 
````

### 2.4.2 template file written 

New File $workspace/view/uname.html 

````html 
<html><body> 
<h1> {{.uname}}</h1> 
</body> </html> 
````

Reboot the system, and now open the browser, whether a standard html page to see the?

### 2.4.3 template function 

wgf/plugin/view support for template functions and embedded wgfInclude and wgfUrl two functions used to introduce other template files and generate url. 

Usage are: 

````
{{wgfInclude "header.tpl".}} 


{{wgfUrl "indexAction"}} 
{{wgfUrl "indexAction" "uname". uname}} // indexAction, get parameters uname = {. uname} 
````


## 2.5 Recommended Mode 

Here for the action to mention one. 

General items, action often have a lot of common logic, so creating a baseAction in app/action/base/base.go Lane (introduced sapi.Action), 
substitute sapi.Action above location more conducive to the organization of the code. 

In addition, the business logic code, often found in app/model, do not write directly in the action. 

Faq believe these people are familiar with the programming mvc clear in mind.


# 3 cli program written using wgf 

## 3.1 Writing and starting 

## 3.2-a parameter
