package main

import "fmt"

/*
const(
	StatusOk = iota // 0
	StatusFailed    // 1
	StatusTimeout   // 2
)
*/

// 类型别名
type Status int

// 隐式定义枚举 无状态非持久化 可以使用 iota
const (
	// iota 预先声明标识符,Go的常量计数器 只能在常量表达式使用
	StatusOk      Status = iota // iota从0开始
	StatusFailed                // 1
	StatusTimeout               // 2
)

// 显式定义枚举 有状态持久化 不能使用 iota
const (
	// TODO:NOTICE 这里的状态值可以和数据库字段做映射关系,定义好枚举值则不允许改变
	//  如首次定义 SOk=200 在某个时刻被某个程序员改为 SOk=201 导致 数据库200 =/= 201 逻辑判断失败
	SOk      Status = 200
	SFailed  Status = 400
	STimeout Status = 500
)

func main() {
	var s Status
	s = StatusOk
	fmt.Println(s)
	s = StatusFailed
	fmt.Println(s)
	s = StatusTimeout
	fmt.Println(s)

	s = SOk
	fmt.Println(s)
	s = SFailed
	fmt.Println(s)
	s = STimeout
	fmt.Println(s)
}

/*
0
1
2
200
400
500
*/
