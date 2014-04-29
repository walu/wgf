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
	if nil == rs {
		rs = new(sql.Rows)
	}

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

func (rs *Rows) FetchInt64(col string) (int64, error) {
	return strconv.ParseInt(rs.currentData[col], 10, 64)
}

func (rs *Rows) FetchUint64(col string) (uint64, error) {
	return strconv.ParseUint(rs.currentData[col], 10, 64)
}

func (rs *Rows) FetchInt32(col string) (ret int32, e error) {
	var r int64
	r, e = strconv.ParseInt(rs.currentData[col], 10, 32)
	if nil != e {
		ret = int32(r)
	}
	return
}

func (rs *Rows) FetchUint32(col string) (ret uint32, e error) {
	var r uint64
	r, e = strconv.ParseUint(rs.currentData[col], 10, 32)
	if nil != e {
		ret = uint32(r)
	}
	return
}

func (rs *Rows) FetchInt16(col string) (ret int16, e error) {
	var r int64
	r, e = strconv.ParseInt(rs.currentData[col], 10, 16)
	if nil != e {
		ret = int16(r)
	}
	return
}


func (rs *Rows) FetchUint16(col string) (ret uint16, e error) {
	var r uint64
	r, e = strconv.ParseUint(rs.currentData[col], 10, 16)
	if nil != e {
		ret = uint16(r)
	}
	return
}


func (rs *Rows) FetchInt8(col string) (ret int8, e error) {
	var r int64
	r, e = strconv.ParseInt(rs.currentData[col], 10, 8)
	if nil != e {
		ret = int8(r)
	}
	return
}


func (rs *Rows) FetchUint8(col string) (ret uint8, e error) {
	var r uint64
	r, e = strconv.ParseUint(rs.currentData[col], 10, 8)
	if nil != e {
		ret = uint8(r)
	}
	return
}

func (rs *Rows) FetchInt(col string) (ret int, e error) {
	var r int64
	r, e = strconv.ParseInt(rs.currentData[col], 10, 64)
	if nil != e {
		ret = int(r)
	}
	return
}

func (rs *Rows) FetchUint(col string) (ret uint, e error) {
	var r uint64
	r, e = strconv.ParseUint(rs.currentData[col], 10, 64)
	if nil != e {
		ret = uint(r)
	}
	return
}

func (rs *Rows) FetchFloat32(col string) (ret float32, e error) {
	var r float64
	r, e = strconv.ParseFloat(rs.currentData[col], 64)
	if nil != e {
		ret = float32(r)
	}
	return
}

func (rs *Rows) FetchFloat64(col string) (float64, error) {
	return strconv.ParseFloat(rs.currentData[col], 64)
}

func (rs *Rows) FetchBool(col string) (bool, error) {
	return strconv.ParseBool(rs.currentData[col])
}

func (rs *Rows) FetchString(col string) (string, error) {
	return rs.currentData[col], nil
}

func (rs *Rows) Fetch() map[string]string {
	return rs.currentData
}

/*
func (rs *Rows) FetchRawBytes(col string) (ret sql.RawBytes, e error) {
}
*/

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
