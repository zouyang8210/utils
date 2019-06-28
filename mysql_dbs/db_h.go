// db_h
package mysql_dbs

import (
	"database/sql"
)

const (
	DRIVER_MYSQL = "mysql"
	WHERE        = " WHERE "
	ORDER_BY     = " ORDER BY "
	LIMIT        = " LIMIT "
	EMPTY        = ""
	SPACE        = " "
	ALL_FIELD    = "*"
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

const (
	NORMAL_TIME_FORMAT = "2006-01-02 15:04:05"
	UTC_TIME_FORMAT    = "2006-01-02T15:04:05.000000Z07:00"
	MIL_TIME_FORMAT    = "2006-01-02 15:04:05.000" //带毫秒的时间字符串
)

//数据库对象
type Mysql_Db struct {
	db    *sql.DB
	dns   string
	Debug bool
}

//数据库控制器
type Control struct {
	*Mysql_Db
	tx            *sql.Tx
	lock          string      //数据锁
	TableName     string      //表名
	Fields        []string    //查询时返回的字段
	Condition     Conditions  //条件
	updateData    UpdateValue //更新时数据
	insertData    InsertValue //插入的数据
	Limit         int         //指定返回记录条记
	Limits        Page        //分页时指定返回记录的起始位置和结束位置,例:(0,10)
	OrderByFields []string    //参与排序的字段
	Sort          string      //排序

}

//分页
type Page struct {
	StartRow int //起始行
	EndRow   int //结止行
}

//记录集
type RecordSet struct {
	Count int                      `json:"count"` //记录集行数
	Data  []map[string]interface{} `json:"data"`  //记录集数据
}

//条件
type Conditions struct {
	conditions map[string]interface{}
	operation  []string
}

//更新数据
type UpdateValue struct {
	updateData map[string]interface{}
}

//插入数据
type InsertValue struct {
	insertData map[string]interface{}
}
