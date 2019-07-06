package zorm

import (
	"database/sql"
	"fmt"
	"time"
)

//初始化数据库操作对像
func New(userName, password, serverAddr, databaseName string) (db Db, err error) {
	db.dns = userName + ":" + password + "@tcp(" + serverAddr + ":3306)/" + databaseName + "?charset=utf8"
	if err = db.open(); err != nil {
		return
	}
	db.db.SetMaxOpenConns(50)
	db.db.SetMaxIdleConns(10)
	return
}

//执行查询sql语句
func (c *Db) Query(strSql string) (record RecordSet, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows
	defer c.printDuration(time.Now())
	c.logs(strSql)
	if stmt, err = c.db.Prepare(strSql); err != nil {
		return
	}
	defer stmt.Close()
	if rows, err = stmt.Query(); err == nil {
		defer rows.Close()
		record, err = c.getRecordSet(rows)
	}
	return
}

//执行增,删,改 sql语句
func (c *Db) Execute(strSql string, tx *Tx) (result int64, err error) {
	var r sql.Result
	defer c.printDuration(time.Now())
	c.logs(strSql)
	if tx == nil {
		r, err = c.db.Exec(strSql)
	} else {
		r, err = tx.tx.Exec(strSql)
	}
	if err == nil {
		result, err = r.RowsAffected()
	}
	return
}

//创建条件对像
func (c *Db) NewCondition() (cond Conditions) {
	cond = Conditions{}
	return
}

//新增(插入)记录
func (c *Db) Insert(tbName string, obj interface{}, tx *Tx, flag bool, fields ...string) (result int64, err error) {
	data, err := c.createMapData(obj, flag, fields...)
	var strSql string
	if strSql, err = c.createInsertSQL(tbName, data); err == nil {
		result, err = c.Execute(strSql, tx)
	}
	return
}

//更新记录
func (c *Db) Update(tbName string, obj interface{}, cond Conditions, tx *Tx, flag bool, fields ...string) (result int64, err error) {
	data, err := c.createMapData(obj, flag, fields...)
	if strSql := c.createUpdateSql(tbName, data, cond); err == nil {
		result, err = c.Execute(strSql, tx)
	}
	return
}

//删除记录
func (c *Db) Delete(tbName string, cond Conditions, tx *Tx) (result int64, err error) {
	strCond, _ := c.createCondition(cond)
	strSql := fmt.Sprintf("DELETE FROM %s WHERE %s", tbName, strCond)
	result, err = c.Execute(strSql, tx)
	return
}

//查询返回记录集
func (c *Db) QueryRecordSet(tbName string, cond Conditions, fields ...string) (record RecordSet, err error) {
	strCond, _ := c.createCondition(cond)
	strFields := ""
	if len(fields) == 0 {
		strFields = " * "
	} else {
		for _, v := range fields {
			strFields += v + ","
		}
		strFields = strFields[0 : len(strFields)-1]
	}
	strSql := fmt.Sprintf("SELECT %s FROM %s WHERE %s", strFields, tbName, strCond)
	record, err = c.Query(strSql)
	return
}

//查询输出对像
func (c *Db) QueryObject(tbName string, cond Conditions, obj interface{}, fields ...string) (err error) {
	strCond, _ := c.createCondition(cond)
	strFields := ""
	if len(fields) == 0 {
		strFields = " * "
	} else {
		for _, v := range fields {
			strFields += v + ","
		}
		strFields = strFields[0 : len(strFields)-1]
	}
	record := RecordSet{}
	strSql := fmt.Sprintf("SELECT %s FROM %s WHERE %s", strFields, tbName, strCond)
	record, err = c.Query(strSql)
	objectToObject(&obj, record.Data)
	return
}

//存储过程执行查询
func (c *Db) QueryForProcedure(proName string, parameterValue []interface{}, outObject interface{}) (err error) {
	strValue := c.createProcedureValue(parameterValue)
	strSql := ""
	if strValue != EMPTY {
		strSql = fmt.Sprintf("CALL %s(%s)", proName, strValue)
	} else {
		strSql = fmt.Sprintf("CALL %s", proName)
	}
	record := RecordSet{}
	if record, err = c.Query(strSql); err != nil {
		return
	}
	if record.Count > 0 {
		objectToObject(&outObject, record.Data)
	}
	return
}

//存储过程执行修改
func (c *Db) ExecForProcedure(proName string, parameterValue []interface{}, tx *Tx) (result int64, err error) {
	strValue := c.createProcedureValue(parameterValue)
	strSql := ""
	if strValue != EMPTY {
		strSql = fmt.Sprintf("CALL %s(%s)", proName, strValue)
	} else {
		strSql = fmt.Sprintf("CALL %s", proName)
	}
	result, err = c.Execute(strSql, tx)
	return
}

//=================================================事务处理=============================================================

//开始(创建)事务
func (c *Db) Begin() (transaction Tx, err error) {
	transaction.tx, err = c.db.Begin()
	return
}

//事务提交
func (c *Tx) Commit() (err error) {
	err = c.tx.Commit()
	return
}

//事务回滚
func (c *Tx) Rollback() (err error) {
	err = c.tx.Rollback()
	return
}

//=====================================================================================================================
