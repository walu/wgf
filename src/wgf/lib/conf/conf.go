package conf

import (
	"path"
	"strconv"
)

type Conf struct {
	data      map[string]string
	orderList []string
}

func (c Conf) String(name string, defaultValue string) string {
	str, ok := c.data[name]
	if !ok {
		str = defaultValue
	}
	return str
}

//1, t, T, TRUE, true, True => true
//all other => false
func (c Conf) Bool(name string, defaultValue bool) bool {
	b, err := strconv.ParseBool(c.data[name])
	if nil != err {
		return defaultValue
	}
	return b
}

//return int64, base 10
//when error, return 0
func (c Conf) Int64(name string, defaultValue int64) int64 {
	i, err := strconv.ParseInt(c.data[name], 10, 64)
	if nil != err {
		i = defaultValue
	}
	return i
}

func (c Conf) Data() map[string]string {
	return c.data
}

func (c Conf) OrderList() []string {
	return c.orderList
}

func (c *Conf) ParseFile(filepath string) error {
	var err error

	//find basedir
	var basedir, filename string
	basedir = path.Dir(filepath)
	filename = path.Base(filepath)
	c.data, c.orderList, err = parseConf(basedir, filename)
	return err
}

func NewConf() *Conf {
	return &Conf{}
}
