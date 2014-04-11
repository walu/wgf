package main

import (
	"wgf"
	_ "demoApp/socketAction"
)

func main() {
	wgf.StartSocketServer()
}
