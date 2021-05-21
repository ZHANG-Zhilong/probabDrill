package dbopts

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

//数据库配置


//Db数据库连接池

type User struct {
	Id    int64
	Name  string
	Age   int8
	Sex   int8
	Phone string
}



//查询操作
func Query() {
	var user User
	rows, e := DB.Query("select * from user where Id in (1,2,3)")
	if e == nil {
		errors.New("query incur error")
	}
	for rows.Next() {
		e := rows.Scan(user.Sex, user.Phone, user.Name, user.Id, user.Age)
		if e != nil {
			fmt.Println(json.Marshal(user))
		}
	}
	rows.Close()
	DB.QueryRow("select * from user where Id=1").Scan(user.Age, user.Id, user.Name, user.Phone, user.Sex)

	stmt, e := DB.Prepare("select * from user where Id=?")
	query, e := stmt.Query(1)
	query.Scan()
}

func DeleteUser(user User) bool {
	//开启事务
	tx, err := DB.Begin()
	if err != nil {
		fmt.Println("tx fail")
	}
	//准备sql语句
	stmt, err := tx.Prepare("DELETE FROM user WHERE Id = ?")
	if err != nil {
		fmt.Println("Prepare fail")
		return false
	}
	//设置参数以及执行sql语句
	res, err := stmt.Exec(user.Id)
	if err != nil {
		fmt.Println("Exec fail")
		return false
	}
	//提交事务
	tx.Commit()
	//获得上一个insert的id
	fmt.Println(res.LastInsertId())
	return true
}

func InsertUser(user User) bool {
	//开启事务
	tx, err := DB.Begin()
	if err != nil {
		fmt.Println("tx fail")
		return false
	}
	//准备sql语句
	stmt, err := tx.Prepare("INSERT INTO user (`Name`, `Phone`) VALUES (?, ?)")
	if err != nil {
		fmt.Println("Prepare fail")
		return false
	}
	//将参数传递到sql语句中并且执行
	res, err := stmt.Exec(user.Name, user.Phone)
	if err != nil {
		fmt.Println("Exec fail")
		return false
	}
	//将事务提交
	tx.Commit()
	//获得上一个插入自增的id
	fmt.Println(res.LastInsertId())
	return true
}
