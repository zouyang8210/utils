package zorm

import (
	"database/sql"
	"time"
)

/*
功能:初始化数据库操作对像
参数:
	user:用户名
	pw:密码
	server_addr:数据库服务器地址
	database_name:数据库名称
返回:
*/
func New(user, pw, serverAddr, databaseName string) (db Db, err error) {
	db.dns = user + ":" + pw + "@tcp(" + serverAddr + ":3306)/" + databaseName + "?charset=utf8"
	err = db.open()
	if err != nil {
		panic(err)
	} else {
		db.db.SetMaxOpenConns(50)
		db.db.SetMaxIdleConns(10)
	}
	return
}

//执行查询sql语句
func (c *Db) Query(strSql string) (record RecordSet, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows
	defer c.runTime(time.Now())
	c.logs(strSql)
	stmt, err = c.db.Prepare(strSql)
	if err == nil {
		defer stmt.Close()
		rows, err = stmt.Query()
		if err == nil {
			defer rows.Close()
			record = c.getRecordSet(rows)
		}
	}
	return
}

//执行增,删,改 sql语句
func (c *Db) Execute(strSql string, tx *sql.Tx) (result int, err error) {
	var r sql.Result
	defer c.runTime(time.Now())
	c.logs(strSql)
	if tx == nil {
		r, err = c.db.Exec(strSql)
	} else {
		r, err = tx.Exec(strSql)
	}
	if err == nil {
		ret, _ := r.RowsAffected()
		result = int(ret)
	}
	return
}
