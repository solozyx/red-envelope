package accounts

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"

	"github.com/solozyx/red-envelope/services"
)

// 持久化对象是ORM映射的基础
// 1. dbx支持自动映射名称，默认是把驼峰命名转换为下划线命名
// 2. 表名默认是结构体名称转换成下划线命名来映射
// 3. 字段名称默认是field name 转换成下划线命名l来映射，字段映射可以用tag描述
// 4. 使用 uni|unique 的tag值来标识字段唯一索引
// 5. 使用 id|pk 的tag值来表示主键
// 6. 使用 omitempty 的tag值来标识字段更新和写入时会被忽略
// 7. 用 — 的tag值来标识字段在更新，写入，查询时会被忽略

// 账户持久化对象
type Account struct {
	// db tag 自定义结构体字段映射数据库表字段名称
	// 账户id omitempty标识在insert update操作时忽略该字段
	Id int64 `db:"id,omitempty"`
	// 账户编号 unique 该字段是唯一索引
	AccountNo string `db:"account_no,unique"`
	// 账户名称
	AccountName string `db:"account_name"`
	// 账户类型
	AccountType int `db:"account_type"`
	// 货币类型编码
	CurrencyCode string `db:"currency_code"`
	// 用户编号
	UserId string `db:"user_id"`
	// 用户编号
	Username sql.NullString `db:"username"`
	// TODO:NOTICE 账户可用余额 避免Go的float32 float64在计算中丢失数据精度
	Balance decimal.Decimal `db:"balance"`
	// 账户状态
	Status int `db:"status"`
	// 账户创建时间
	CreatedAt time.Time `db:"created_at,omitempty"`
	// 账户更新时间
	UpdatedAt time.Time `db:"updated_at,omitempty"`
}

func (po *Account) FromDTO(dto *services.AccountDTO) {
	po.AccountNo = dto.AccountNo
	po.AccountName = dto.AccountName
	po.AccountType = dto.AccountType
	po.CurrencyCode = dto.CurrencyCode
	po.UserId = dto.UserId
	po.Username = sql.NullString{String: dto.Username, Valid: true}
	po.Balance = dto.Balance
	po.Status = dto.Status
	po.CreatedAt = dto.CreatedAt
	po.UpdatedAt = dto.UpdatedAt
}

func (po *Account) ToDTO() *services.AccountDTO {
	dto := &services.AccountDTO{}
	dto.AccountNo = po.AccountNo
	dto.AccountName = po.AccountName
	dto.AccountType = po.AccountType
	dto.CurrencyCode = po.CurrencyCode
	dto.UserId = po.UserId
	dto.Username = po.Username.String
	dto.Balance = po.Balance
	dto.Status = po.Status
	dto.CreatedAt = po.CreatedAt
	dto.UpdatedAt = po.UpdatedAt
	return dto
}
