package database_tool

var module = `
import (
	"utils"
	"utils/mysql_dbs"
)


[STRUCT]

//表名
const (
	TB_[CONST_TABLE] = "[tablename]"
)

/*
功能：插入信息
参数：
	obj:数据库结构对像
返回：是否成否，错误信息
*/
func (c *[ClassName]) CreateCommercial(obj [StructName]) (result bool, err error) {
	var count int64
	var ctl mysql_dbs.Control
	ctl, err = c.db.CreateInsertControl([TABLENAME], obj)
	if err == nil {
		count, err = ctl.Insert()
		result = utils.If(count > 0, true, false).(bool)
	}
	return
}

功能：更新数据
参数：
	obj:数据库结构对像
返回：是否成否，错误信息
*/
func (c *[ClassName]) CreateCommercial(obj [StructName],condition mysql_dbs.Conditions) (result bool, err error) {
	var count int64
	var ctl mysql_dbs.Control
	ctl, err = c.db.CreateUpdateControl([TABLENAME], obj,)
	if err == nil {
		count, err = ctl.Insert()
		result = utils.If(count > 0, true, false).(bool)
	}
	return
}

`
