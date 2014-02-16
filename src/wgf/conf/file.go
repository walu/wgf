package conf

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var confRowPattern *regexp.Regexp

func readFileString(filename string) (string, error) {
	body, err := ioutil.ReadFile(filename)
	if nil != err {
		return "", err
	} else {
		return string(body), err
	}
}

func parseConf(basedir, filename string) (map[string]string, []string, error) {
	var body string
	var err error
	var re, reForInclude map[string]string
	var orderList, orderListForInclude []string

	body, err = readFileString(basedir + "/" + filename)
	if nil != err {
		return nil, nil, err
	}

	re = make(map[string]string)
	reForInclude = make(map[string]string)

	var row string
	var result []string

	var name, value, includeFile string
	var includeFileFromGlob []string

	for _, row = range strings.Split(body, "\n") {
		row = strings.TrimSpace(row)

		//ignore comments and invalid rows
		if len(row) < 2 || ';' == row[0] {
			continue
		}

		//deal with include syntax
		if strings.HasPrefix(strings.ToLower(row), "include") {
			//deal with include
			includeFile = strings.TrimSpace(row[7:])
			if "" != includeFile {
				includeFileFromGlob, err = filepath.Glob(basedir + "/" + includeFile)
				if nil != err {
					continue
					//return nil, nil, err
				}

				for _, includeFile = range includeFileFromGlob {
					reForInclude, orderListForInclude, err = parseConf(path.Dir(includeFile), path.Base(includeFile))
					if nil != err {
						return nil, nil, err
					}
					for k, v := range reForInclude {
						re[k] = v
					}
					for _, v := range orderListForInclude {
						orderList = append(orderList, v)
					}
				}
			}
			continue
		}

		//deal with common conf
		result = confRowPattern.FindStringSubmatch(row)
		if 0 == len(result) {
			continue
		}
		name = result[1]
		value = result[2]

		re[name] = value
		orderList = append(orderList, name)
	}
	return re, orderList, nil
}

func init() {
	confRowPattern = regexp.MustCompile("(?i:^([^=\\s]+)\\s*=\\s*(.*)$)")
}
