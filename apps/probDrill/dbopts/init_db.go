package dbopts

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

//注意方法名大写，就是public
func InitDB() {
	userName := viper.GetString("mysql.userName")
	password := viper.GetString("mysql.password")
	ip := viper.GetString("mysql.ip")
	port := viper.GetString("mysql.port")
	dbName := viper.GetString("mysql.dbName")
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, _ = sql.Open("mysql", path)
	//设置数据库最大连接数
	DB.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(10)
	//验证连接
	if err := DB.Ping(); err != nil {
		fmt.Println("open database fail")
		return
	}
	fmt.Println("connect success")
}
