package router

import (
	"errors"
	"fmt"
	liburl "net/url"
	"regexp"
	"strings"
)

var regexpArgs *regexp.Regexp

type rule struct {
	origin  string
	pattern *regexp.Regexp
	args    []string
	action  string
}

var mapRule map[string]*rule
var listRule []*rule

func actionToUrl(actionName string, params map[string]string) (url string, err error) {
	rule, ok := mapRule[actionName]
	if !ok || len(rule.args) > len(params) {
		return "", errors.New("none or url params is lt than rule.args")
	}

	dupParam := make(map[string]string)
	var moreArgs string
	for k, v := range params {
		dupParam[k] = v
	}

	url = rule.origin
	for _, arg := range rule.args {
		value, ok := params[arg]
		if !ok {
			return "", errors.New("can't find arg")
		}
		tmpRegexp := regexp.MustCompile("#" + arg + "(:[^#]+)?#")
		url = tmpRegexp.ReplaceAllString(url, value)
		delete(dupParam, arg)
	}

	if len(dupParam) > 0 {
		for k, v := range dupParam {
			moreArgs = moreArgs + liburl.QueryEscape(k) + "=" + liburl.QueryEscape(v) + "&"
		}
		url = url + "?" + moreArgs
	}
	return url, nil
}

func addRule(leftPart string, actionName string) error {
	pRule := &rule{}

	re := regexpArgs.FindAllStringSubmatch(leftPart, -1)
	pRule.origin = leftPart
	pRule.action = actionName
	pattern := leftPart

	for _, v := range re {
		pRule.args = append(pRule.args, v[2])
		if v[3] == "" {
			pattern = strings.Replace(pattern, v[1], "(.+)", -1)
		} else {
			pattern = strings.Replace(pattern, v[1], "("+v[3][1:]+")", -1)
		}
	}

	var err error
	pattern = "^" + pattern + "$"
	pRule.pattern, err = regexp.Compile(pattern)
	if nil != err {
		return err
	}

	mapRule[actionName] = pRule
	listRule = append(listRule, pRule)
	return nil
}

func urlToAction(url string) (string, map[string]string) {
	var matched [][]string
	var action string
	var params map[string]string

	params = make(map[string]string)

	for _, r := range listRule {
		matched = r.pattern.FindAllStringSubmatch(url, -1)
		if nil == matched || len(r.args) > len(matched[0]) {
			continue
		}

		for k, v := range r.args {
			params[v] = matched[0][k+1]
		}
		action = r.action
		break
	}
	return action, params
}

func init() {
	regexpArgs = regexp.MustCompile("(#([-\\.a-zA-Z0-9]+)(:([^#]+))?#)")
	mapRule = make(map[string]*rule)
}

func forFmt() {
	fmt.Println("hello world!")
}
