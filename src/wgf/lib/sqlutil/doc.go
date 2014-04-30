// Copyright 2014 The Wgf Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
A library helping you deal with databases


sqlutil is a library dealing with database issues, it's a library of wgf, but you can copy and use it independently

it support read and write separation, has a new Rows to drop sql.Rows's Scan method,
and has some productive methods like Insert\Update\Delete\Select, etc.

1. Create *DB

sqlutil wrapped *sql.DB by *DB

if you have only one db sever:

	var dbConn *DB
	dbConn = NewDB(db)

if you have master and slave db servers

	var dbConn *DB
	dbConn = NewRWDB(master, slave1, slave2.....)

2. Query Sqls

2.1 Select

Select provides a simple way to do select works.

	//SELECT * FROM wiki
	dbConn.Select("wiki")

	//SELECT * FROM wiki WHERE id = 1
	dbConn.Select("wiki", "id", 1)

	//SELECT * FROM wiki WHERE id = 1 AND name = wgf
	dbConn.Select("wiki", "id", 1, "name", "wgf")

	//I thinks you might want to use map to pass params
	//I don't recommend you using map, because it is unordered
	//but the order of params is very very very very important.
	params := map[string]interface{}{
		"id": 1,
		"name": "wgf",
	}
	args := WhereMapToSlice(params)
	dbConn.Select("wiki", args...)

2.2 Insert

	data := map[string]interface{}{
		"name": "wgf",
		"age": 1,
	}
	dbConn.Insert("wiki", data)

so easy, doesn't it?

2.3 Update
	
	data := map[string]interface{}{
		"name": "wgf",
		"age": 1,
	}
	
	//UPDATE ...
	dbConn.Update("wiki", data)

	//UPDATE ... WHERE id = 1
	dbConn.Update("wiki", data, "id", "1")

	//UPDATE ... WHERE id = 1 AND name = "wgf"
	dbConn.Update("wiki", data, "id", "1", "name", "wgf")

2.4 Delete 

	//DELETE All
	dbConn.Delete("wiki")

	//Delete ... WHERE id = 1
	dbConn.Delete("wiki", data, "id", "1")

	//Delete ... WHERE id = 1 AND name = "wgf"
	dbConn.Delete("wiki", "id", "1", "name", "wgf")


2.5 Query And Exec

if the methods above can't cover your works(of couse, it's very common),
you shoule use Query And Exec by youself

	dbConn.Query("SELECT * FROM wiki WHERE id = ? LIMIT ?, ?", 1, 100, 50)
	dbConn.Exec("UPDATE wiki WHERE id = ? ", 1)

3 Rows

sqlutil provides a Rows structure to deal with query result.
it's more convenient and easier then sql.Rows, at least I think so.

Rows will read the data into its memory, and close sql.Rows auto, to prevent
you forget.

Select And Query method will return *Rows

	//get a list
	rs, err := dbConn.Select("wiki")
	if nil!=err {
		return
	}

	for rs.Next() {
		tmp = &Wiki{}
		tmp.Id, _ = rs.FetchInt64("id")
		tmp.Type, _ = rs.FetchInt8("type")
		tmp.Body, _ = rs.FetchString("body")
	}

	//some database has its own data type, you can fetch it and parse it by hand
	//sqlutil will provides lot's of helper funcs to deal with this works later
	oriData = tmp.Fetch("pg_json")
*/
package sqlutil
