package dbopts

import (
	"database/sql"
	"testing"
)

func Test_testSql(t *testing.T) {
	var user User
	user.Name = "fu"
	user.Phone = "123"
	user.Id = 3
	InitDB()
	//InsertUser(user)
	Query()
	DeleteUser(user)

	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {

		}
	}(DB)
}
