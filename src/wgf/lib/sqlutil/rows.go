// Copyright 2014 The Wgf Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sqlutil

import (
	"strconv"
	"database/sql"
	"errors"
)

/*
Rows is instead of sql.Rows

Rows will get all the result into memory of sql.Rows, 
preventing you forget to call the Close() method, the reason let db conn can't be released.
*/
type Rows struct {
	index int
	dataLen int
	currentData map[string]string
	colnames []string

	data []map[string]string
}

/*
Rows will get all the result into memory of sql.Rows, 
preventing you forget to call the Close() method, the reason let db conn can't be released.
*/
func NewRows(rs *sql.Rows) (*Rows, error) {
	if nil == rs {
		rs = new(sql.Rows)
	}
	defer rs.Close()

	var err error
	var tmp map[string]string

	ret := &Rows{}
	ret.currentData = make(map[string]string)
	ret.data = make([]map[string]string, 0)
	ret.colnames, err = rs.Columns()
	if nil != err {
		return nil, err
	}

	for rs.Next() {
		tmp, err = fetchMap(rs)
		if nil != err {
			return nil, err
		}

		ret.data = append(ret.data, tmp)
		ret.dataLen++
	}
	return ret, nil
}

func (rs *Rows) Next() bool {
	if rs.index>=rs.dataLen {
		return false
	}

	rs.currentData = rs.data[rs.index]
	rs.index++
	return true
}

func (rs *Rows) Columns() []string {
	return rs.colnames
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
