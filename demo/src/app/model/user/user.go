package user

import (
	"app/model/dao"
	"app/model/ldap"
	"fmt"
)

type User struct {
	Uname string
}

//是否存在这个用户
func Exists(uname string) bool {
	u, _ := GetUserByUname(uname)
	return u.Uname != ""
}
func Add(uname string) error {
	return nil
}

func GetUserByUname(uname string) (*User, error) {
	u := &User{}

	db, err := dao.GetDb()
	if nil != err {
		return nil, err
	}
	defer func() { db.Close() }()

	r := db.QueryRow("SELECT uname FROM user WHERE uname = ?", uname)
	err = r.Scan(&u.Uname)
	fmt.Println(err)
	return u, err

}

//通过ldap检查用户名与密码
func LdapCheck(uname string, password string) bool {
	re := ldap.CheckPassword(uname, password)
	return re
}
