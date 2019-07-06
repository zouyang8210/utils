package zorm

import "database/sql"

//普通字符串
const (
	EMPTY        = ""
	SPACE        = " "
	ALL_FIELD    = "*"
	DRIVER_MYSQL = "mysql"
)

//关键字
const (
	WHERE    = " WHERE "
	ORDER_BY = " ORDER BY "
	LIMIT    = " LIMIT "
)

//运算符
const (
	AND           = " AND "
	OR            = " OR "
	EQUAL         = " = "
	NOT_EQUAL     = " != "
	LESS          = " < "
	GREATER       = " > "
	LESS_EQUAL    = " <= "
	GREATER_EQUAL = " >= "
	DESC          = "DESC"
	ASC           = "ASC"
)

//时间格式
const (
	NORMAL_TIME_FORMAT = "2006-01-02 15:04:05"
	UTC_TIME_FORMAT    = "2006-01-02T15:04:05.000000+07:00"
	MIL_TIME_FORMAT    = "2006-01-02 15:04:05.000" //带毫秒的时间字符串
)

//数据库对象
type Db struct {
	db    *sql.DB
	dns   string
	Debug bool
}

//事务对象
type Tx struct {
	tx *sql.Tx
}

//记录集
type RecordSet struct {
	Count int                      `json:"count"` //记录集行数
	Data  []map[string]interface{} `json:"data"`  //记录集数据
}

//条件
type Conditions struct {
	conditions []KV
	operation  []string //连接where条件操作符
	limit      string   //查询的记录条数
	store      string   //排序
	pages      string   //分页
}

//where条件值
type KV struct {
	K string      //字段
	S string      //运算符
	V interface{} //值
	B string      //括号
}
