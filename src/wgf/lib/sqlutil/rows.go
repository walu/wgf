// Copyright 2014 The Wgf Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sqlutil

import (
	"strconv"
	"database/sql"
	"errors"
)

type Rows struct {
	rows *sql.Rows
	currentData map[string]string
	errorList []error
}

func NewRows(rs *sql.Rows) *Rows {
	ret := &Rows{}
	ret.currentData = make(map[string]string)
	ret.errorList = make([]error, 0)
	ret.rows = rs
	return ret
}

func (rs *Rows) Next() bool {
	var ret bool
	if rs.rows.Next() {
		var err error
		rs.currentData, err = fetchMap(rs.rows)
		if nil == err {
			ret = true
		}
	}
	return ret
}

func (rs *Rows) Close() error {
	return rs.rows.Close()
}

func (rs *Rows) Columns() ([]string, error) {
	return rs.rows.Columns()
}

func (rs *Rows) Err() error {
	return rs.rows.Err()
}

func (rs *Rows) FetchInt(col string) (int64, error) {
	return strconv.ParseInt(rs.currentData[col], 10, 64)
}

func (rs *Rows) FetchUint(col string) (uint64, error) {
	return strconv.ParseUint(rs.currentData[col], 10, 64)
}

func (rs *Rows) FetchString(col string) (string, error) {
	return rs.currentData[col], nil
}

func fetchMap(rs *sql.Rows) (map[string]string, error) {
	var colnames []string
	var err	error

	colnames, err = rs.Columns()
	if nil != err {
		return nil, errors.New("rows to map error: " + err.Error())
	}

	var lenCol int
	lenCol = len(colnames)

	var ret map[string]string
	ret = make(map[string]string)

	var args []sql.RawBytes
	var scanArgs []interface{}

	args = make([]sql.RawBytes, lenCol)
	scanArgs = make([]interface{}, lenCol)
	for i := range args {
		scanArgs[i] = &args[i]
	}

	err = rs.Scan(scanArgs...)
	if nil != err {
		return nil, err
	}

	for index, val := range args {
		ret[colnames[index]] = string(val)
	}

	return ret, nil
}
