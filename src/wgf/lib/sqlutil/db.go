package sqlutil

import (
	"database/sql"
	"strings"
	"errors"
	"fmt"
)

var (
	errArgsMatch = "err_args_not_match"
	errDbQuery = "err_db_query"
)

type DB struct {
	oriMaster *sql.DB
	oriSlaves []*sql.DB

}

//create DB
func NewDB(db *sql.DB) *DB {
	ret := new(DB)

	ret.oriMaster = db
	ret.oriSlaves = []*sql.DB{db}
	return ret
}

//create DB, master and slaves
//INSERT\Delete\Update use master
//Select use slaves randomly
func NewRWDB(master *sql.DB, slaves ...*sql.DB) *DB {
	ret := new(DB)
	ret.oriMaster = master
	ret.oriSlaves = slaves
	return nil
}

/*

help you converting map args to slice args

Why slice was used default?

becaulse golang's map is unordered

like: {id:1, name:wgf}

it will be sql:id=1 and name = wgf or sql:name = wgf and id =1, but the order is very important for database search.

	//If you have to use the map, you can do it like this
	queryParam := map[string]interface{}{
		"id" : 1,
		"name" : "wgf",
	}
	whereArgs := sqlutil.WhereMapToSlice(queryParam)
	p.Select("wiki", whereArgs...)

*/
func WhereMapToSlice(where map[string]interface{}) []interface{} {
	l := len(where)
	if l == 0 {
		return nil
	}

	ret := make([]interface{}, l*2)
	i := 0
	for k, v := range where {
		ret[i] = k
		ret[i+1] = v
		i += 2
	}
	return ret
}

/*

Simple select

	p.Select("wiki")
	select * from wiki

	p.Select("wiki", "id", 1)
	select * from wiki where id = 1

	p.Select("wiki", "id", 1, "name", "wgf") 
	select * from wiki where id = 1 and name = "wgf"

	if args%2==1, return an error

	selece one row only?
	rs, _ := p.Select("wiki", "id", 1)
	rs.Next()
	......

*/
func (p *DB) Select(table string, args ...interface{}) (*Rows, error) {
	var s, sqlWhere string
	var queryArgs []interface{}
	var err error

	sqlWhere, queryArgs, err = parseWhereArgs(args)
	if nil != err {
		return nil, err
	}

	s = "SELECT * FROM " + table + sqlWhere
	return p.Query(s, queryArgs...)
}

func (p *DB) Delete(table string, args ...interface{}) (sql.Result, error) {
	var s, sqlWhere string
	var queryArgs []interface{}
	var err error

	sqlWhere, queryArgs, err = parseWhereArgs(args)
	if nil != err {
		return nil, err
	}

	s = "DELETE FROM " + table + sqlWhere
	return p.Exec(s, queryArgs...)
}

func (p *DB) Insert(table string, row map[string]interface{}) (sql.Result, error) {
	//INSERT INTO TABLE(FILEDS, FIELDS, FIELDS) VALUES()
	var s string
	var queryArgs []interface{}
	var fields []string
	var holders []string

	l := len(row)
	queryArgs = make([]interface{}, l)
	fields = make([]string, l)
	holders = make([]string, l)
	s = "INSERT INTO " + table

	i := 0
	for k, v := range row {
		fields[i] = k
		holders[i] = "?"
		queryArgs[i] = v
		i++
	}

	s = fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", table, strings.Join(fields, ", "), strings.Join(holders, ", "))
	return p.Exec(s, queryArgs...)
}

//update table
//whereArgs only support simple eq relation
func (p *DB) Update(table string, row map[string]interface{}, whereArgs ...interface{}) (sql.Result, error) {
	var s string
	var queryArgs []interface{}
	var fields []string

	l := len(row)
	queryArgs = make([]interface{}, l)
	fields = make([]string, l)

	i := 0
	for k, v := range row {
		fields[i] = k + " = ? "
		queryArgs[i] = v
		i++
	}

	sqlWhere, tmpArgs, e := parseWhereArgs(whereArgs)
	if nil != e {
		return nil, e
	}
	queryArgs = append(queryArgs, tmpArgs...)


	s = fmt.Sprintf("UPDATE  %s SET %s %s", table, strings.Join(fields, ", "), sqlWhere)
	fmt.Println(s, queryArgs)
	return p.Exec(s, queryArgs...)
}


/*
	query one row only?
	rs, _ := p.Query()
	rs.Next()
	......
*/
func (p *DB) Query(s string, args ...interface{}) (*Rows, error) {
	rs, err := p.oriMaster.Query(s, args...)
	if nil != err {
		body := fmt.Sprintf("%s err:%v sql:%s args:%v", errDbQuery, err, s, args)
		log(body)
		return nil, errors.New(body)
	}
	return NewRows(rs), nil
}

func (p *DB) Exec(s string, args ...interface{}) (sql.Result, error) {
	return p.oriMaster.Exec(s, args...)
}



func (p *DB) Close() error {
	return nil
}

func (p *DB) slave() *sql.DB {
	return p.oriSlaves[0]
}

//if parse to where, there will be a space ahead
func parseWhereArgs(where []interface{}) (string, []interface{}, error) {
	l := len(where)
	if l==0 {
		return "", nil, nil
	}

	if l%2!=0 {
		return "", nil, errors.New(fmt.Sprintf("%s %v", errArgsMatch, where))
	}

	var s string
	var queryArgs []interface{}
	queryArgs = make([]interface{}, l/2)

	s = " WHERE "
	for i:=0; i<l; i+=2 {
		s += where[i].(string) + " = ? "
		queryArgs[i/2] = where[i+1]
	}

	return s, queryArgs, nil
}
