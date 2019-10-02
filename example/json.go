package main

import (
	"fmt"

	"encoding/json"
)

type User struct {
	// 标识 string 这里是把int类型用string在JSON序列化时编码解码
	Id int `json:"id,string"`
	// 标准库json只能序列化可导出结构体字段 要求字段首字母必须大写
	// 不加json tag 字段Name在JSON序列化的字段是 Name
	// 加 json tag 后在JSON序列化的字段是 username
	Name string `json:"username"`
	// 标识 omitempty 在JSON序列化时忽略空值结构体字段
	// 如 Age=1可以被JSON序列化 Age=0则在JSON序列化时忽略掉Age字段
	Age int `json:"age,omitempty"`
	// 标识 - 中划线 在JSON序列化时忽略该字段,不论空值与否都不会被序列化
	Address string `json:"-"`
}

func main() {
	u1 := User{
		Id:      12,
		Name:    "wendell",
		Age:     1,
		Address: "成都高新区",
	}
	// Marshal 参数是结构体指针
	data, err := json.Marshal(&u1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))

	u2 := &User{}
	err = json.Unmarshal(data, u2)
	if err != nil {
		fmt.Println(err)
	}
	// %+v 把结构体的字段和值都输出
	fmt.Printf("%+v \n", u1)
	fmt.Printf("%+v \n", u2)

	data = []byte(`{"id":"13","username":"user","age":2,"Address":"北京"}`)
	u3 := &User{}
	json.Unmarshal(data, u3)
	fmt.Printf("%+v \n", u3)

	data = []byte(`{"id":1,"username":"user","age":2,"Address":"北京"}`)
	u4 := &User{}
	json.Unmarshal(data, u4)
	fmt.Printf("%+v \n", u4)
}
