// Copyright 2014 The Wgf Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.
package sapi

/*
测试过静态Action，不过和动态Action相比，性能一样。
后续去掉了这种机制。
*/

import (
	"errors"
)

const (
	USE_EXECUTE = 0
	USE_DOFUNC  = 1
)

type ActionInterface interface {
	SetSapi(p *Sapi)

	UseSpecialMethod() bool
	Execute() error

	DoGet() error
	DoPost() error
}

type StaticActionInterface interface {
	Execute(pSapi *Sapi) error
}

//存储已注册的action
var actionMap map[string]func() ActionInterface

/*
默认action，用于简化app实现逻辑，app在实现自己的action时，可以直接包含此action。
	import "wgf/sapi"
	type IndexAction struct {
		sapi.Action
	}

	func (p *IndexAction) DoGet() error {
	}
*/
type Action struct {
	RunMode int
	Sapi    *Sapi
}

//设置此Action对应的sapi实例
func (p *Action) SetSapi(s *Sapi) {
	p.Sapi = s
}

//是否启用专有方法
func (p *Action) UseSpecialMethod() bool {
	return p.RunMode == USE_DOFUNC
}

//如果未启用专有方法，则将所有类型的请求都转到此方法
func (p *Action) Execute() error {
	p.Sapi.Print("nonsupport.")
	return nil
}

//如果启用了专用方法，且请求类型为GET，则执行此方法
func (p *Action) DoGet() error {
	p.Sapi.Print("nonsupport.")
	return nil
}

//如果启用了专用方法，且请求类型为POST，则执行此方法
func (p *Action) DoPost() error {
	p.Sapi.Print("nonsupport.")
	return nil
}

//注册action
func RegisterAction(name string, creater func() ActionInterface) {
	actionMap[name] = creater
}

//获取已注册的action
func GetAction(name string) (action ActionInterface, err error) {
	creater, ok := actionMap[name]

	if ok {
		action = creater()
	} else {
		err = errors.New("no action named: " + name)
	}

	return action, err
}

func init() {
	actionMap = make(map[string]func() ActionInterface)
}
