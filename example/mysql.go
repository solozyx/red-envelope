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
	// "user_root:pass_root@tcp(192.168.174.134:3306)/red_envelope?charset=utf8&parseTime=true&loc=Local"
	db, err := sql.Open("mysql",
		"root:root@tcp(192.168.174.134:3306)/red_envelope?charset=utf8&parseTime=true&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 设置数据库连接池最大空闲连接数
	// 默认值是 0 表示连接池不会保持连接池中连接状态,当连接被返回到连接池时,连接会被关闭
	// 会导致连接在连接池频繁开关,相当于没有使用连接池
	// 空闲连接配置过大 导致大量连接空闲 浪费资源
	db.SetMaxIdleConns(2)
	// 设置数据库最大打开连接数 = 正在使用的连接 + 空闲连接
	// 申请1个连接时,连接池中没用空闲连接或者当前连接达到最大连接数,申请连接的操作会被block直到有可用的连接才会返回
	// 最大连接数配置过大或使用后没及时释放连接,数据库连接过多,导致服务器性能下降 MySQL出现 too many connections 错误
	// 默认值是0  表示无限制
	// 通常是性能工程师评估后配置合理的值 在开发时配置 2 - 3 个连接即可
	db.SetMaxOpenConns(3)
	// 设置闲置连接的最大存活时间
	// MySQL默认的非交互连接的空闲时间是8个小时 连接池中空闲连接超过8小时 会被MySQL自动断开回收失效
	// 所以这里设置连接的最大存活时间要小于8小时 如7个小时
	db.SetConnMaxLifetime(7 * time.Hour)

	fmt.Println(db.Ping())
	fmt.Println(db.Query("select now()"))
}
