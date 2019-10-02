package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	// 引入MySQL驱动程序
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// sql.Open 创建 sql.DB 实例
	db, err := sql.Open("mysql",
		"user_root:pass_root@tcp(192.168.174.134:3306)/red_envelope?charset=utf8&parseTime=true&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 设置数据库连接池最大空闲连接数
	db.SetMaxIdleConns(2)
	// 设置数据库最大打开连接数
	db.SetMaxOpenConns(3)
	// MySQL默认的非交互连接的空闲时间是8个小时 连接池中空闲连接超过8小时 会被MySQL自动断开回收失效
	// 所以这里设置连接的最大存活时间要小于8小时 如7个小时
	db.SetConnMaxLifetime(7 * time.Hour)

	fmt.Println(db.Ping())
	fmt.Println(db.Query("select now()"))
}
