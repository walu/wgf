>
> This Document was maked with Google Translate
>
> You are welcome to improve it.
>
> **Forgive my poor english, thanks!**
>

# wgf document

# 1. Wgf development Remarks

### 1.1 wgf not support go get, there is no support plan yet.

Because :

* The main reason : the prefix of importing path github.com/walu/** is long and boring sort.
* Secondary reason : no version control .

Therefore, when using wgf , please download to a local, and copy to the project, or to add wgf path into GOPATH.

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

Of course , you are not absolutely necessary to the operation of the above so white, so long as the last line on the line.

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
    sapi.Action / / embed this Action, save yourself works
}

// You can write the logic in the Execute method .
// and Can also be written in DoGet \ DoPost method in .
func (p * IndexAction) Execute () {
    p.Sapi.Println ("hello world!")
}

// Register IndexAction to wgf
func init () {
    sapi.RegisterAction ("index", func () sapi.ActionInterface {return & IndexAction {}})
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
    _ "App/action" // This is not a good habit to introduce temporary , so tentatively
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
