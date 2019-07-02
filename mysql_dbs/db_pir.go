// db_pir
package mysql_dbs

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
	"utils/data_conv/json_lib"

	_ "github.com/go-sql-driver/mysql"
	"utils/data_conv/str_lib"
	"utils/data_conv/time_lib"
)

/*
功能:开始数据库
参数:
返回:错误信息,nil为正常关闭
*/
func (c *Mysql_Db) open() (err error) {
	c.db, err = sql.Open(DRIVER_MYSQL, c.dns)
	return
}

/*
功能:关闭数据连接
参数:
返回:错误信息,nil为正常关闭
*/
func (c *Mysql_Db) close() (err error) {
	if c.db != nil {
		err = c.db.Close()
		c.db = nil
	}
	return
}

/*
功能:执行查询sql语句
参数:
	str_sql:语句
返回:记录集,信息信息
*/
func (c *Mysql_Db) execQuerySql(strSql string) (record RecordSet, err error) {
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

/*
功能:执行查询sql语句
参数:
	str_sql:语句
返回:记录集,信息信息
*/
func (c *Mysql_Db) execQuerySql2(strSql string) (record RecordSet, err error) {
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
	if rows.NextResultSet() {
		for rows.Next() {
			rows.Scan(&record.Count)
		}
	}
	return
}

/*
功能:执行查询sql语句
参数:
	str_sql:语句
	v:输出参数,一个数据结构的切片
返回:错误信息
*/
func (c *Mysql_Db) execQueryObj(strSql string, v interface{}) (err error) {
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
			err = c.getRowToObj(rows, &v)
		}
	}
	return
}

/*
功能:sql.Rows转struct(获取具体数据类型的切片)
参数:
	rows:数据列表
	v:输出参数,一个数据结构的切片
返回:错误信息
*/
func (c *Mysql_Db) getRowToObj(rows *sql.Rows, v interface{}) (err error) {
	count, json, err := json_lib.RowsToJson(rows)
	if err == nil {
		if count > 0 {
			err = json_lib.JsonToObject(json, &v)
		}
	}
	return
}

/*
功能:执行增删改sql语句
参数:
	strSql:语句
返回:影响的记集数,信息信息
*/
func (c *Mysql_Db) execSql(strSql string, tx *sql.Tx) (ret int64, err error) {
	var r sql.Result
	defer c.runTime(time.Now())
	c.logs(strSql)
	if tx == nil {
		r, err = c.db.Exec(strSql)
	} else {
		r, err = tx.Exec(strSql)
	}
	if err == nil {
		ret, _ = r.RowsAffected()
	}
	return
}

/*
功能:组合查询条件
参数:
	conditions:查询条件
	operation:合条件的关系运算符 例 "and,or等"
返回:条件字符串
*/
func (c *Mysql_Db) conditionJoint(condition Conditions) (str string) {

	i := 0
	//length := len(condition.operation)
	for n, v := range condition.conditions {
		switch p := v.(type) {
		case uint8, uint16, uint32, uint64, int8, int16, int32, int64, int:
			str += fmt.Sprintf(" %s%d ", n, p)
		case float32, float64:
			str += fmt.Sprintf(" %s%f ", n, p)
		case string:
			str += fmt.Sprintf(" %s'%s' ", n, p)
		case nil:
			str += fmt.Sprintf(" %s  null ", n)
		case time.Time:
			t := time_lib.TimeToStr(p, time_lib.TIME_yyyyMMddHHmmss)
			str += fmt.Sprintf(" %s'%s' ", n, t)
		}
		if i < len(condition.conditions)-1 {
			str += SPACE + condition.operation[i] + SPACE
		} else {
			break
		}
		i++
	}
	return
}

/*
功能:存储过程,参数值组合
参数:
	param_value:参数值
返回:参数值字符串
*/
func (c *Mysql_Db) proParameterJoint(paramValues []interface{}) (str string) {
	var i = 0
	for i = range paramValues {
		switch p := paramValues[i].(type) {
		case uint8, uint16, uint32, uint64, int8, int16, int32, int64, int:
			str += fmt.Sprintf("%d,", p)
		case float32, float64:
			str += fmt.Sprintf("%f,", p)
		case string:
			str += fmt.Sprintf("'%s',", p)
		}
	}
	str = str_lib.SubString(str, 0, len(str)-1) //减掉最后一个逗号
	return
}

/*
功能:更新字段组成
参数:
	update_data:更新的字段
返回:更新字段字符串 例:a=1,b=2,c=3
*/
func (c *Mysql_Db) updateJoint(updateData UpdateValue) (str string) {
	for n, v := range updateData.updateData {
		switch p := v.(type) {
		case uint8, uint16, uint32, uint64, int8, int16, int32, int64, int:
			str += fmt.Sprintf("%s = %d,", n, p)
		case float32, float64:
			str += fmt.Sprintf("%s = %f,", n, p)
		case string:
			str += fmt.Sprintf("%s = '%s',", n, p)
		case nil:
			str += fmt.Sprintf("%s = NULL,", n)
		case time.Time:
			t := p.Format(MIL_TIME_FORMAT) //time_lib.Time_format_str24(p, 1, 1)
			str += fmt.Sprintf(" %s = '%s',", n, t)
		}
	}
	str = str_lib.SubString(str, 0, len(str)-1) //减掉最后一个逗号
	return
}

/*
功能:把sql.Rows转换成RecordSet记录集
参数:
	rows:sql.Rows记录集
返回:RecordSet
*/
func (c *Mysql_Db) getRecordSet(rows *sql.Rows) (record RecordSet) {
	columnName, err := rows.Columns()
	if err == nil {
		columnLen := len(columnName)
		values := make([]interface{}, columnLen)
		valuePtrs := make([]interface{}, columnLen)
		for rows.Next() {
			for i := 0; i < columnLen; i++ {
				valuePtrs[i] = &values[i]
			}
			rows.Scan(valuePtrs...)
			entry := make(map[string]interface{})
			for i, col := range columnName {
				if _, ok := values[i].([]byte); ok {
					str := fmt.Sprintf("%s", values[i])
					tt, err := time.Parse(NORMAL_TIME_FORMAT, str)
					if err == nil {
						//str = tt.Format(UTC_TIME_FORMAT)
						str = tt.Format(NORMAL_TIME_FORMAT)
					}
					entry[col] = str

				} else {
					entry[col] = values[i]
				}
			}
			record.Count++
			record.Data = append(record.Data, entry)
		}
	}
	return
}

/*
功能:debug模式下输出信息
参数:
	str:输出的信息
返回:
*/
func (c *Mysql_Db) logs(str string) {
	if c.Debug {
		fmt.Printf("%s->%s\n", time.Now().Format(NORMAL_TIME_FORMAT), str)
	}
}

//测试代码执行用时
func (c *Mysql_Db) runTime(now time.Time) {
	if c.Debug {
		//terminal := time.Since(now)
		//fmt.Println("Exec Time:", terminal)
	}
}

func (c *Control) createInsertSQL() (strSql string, err error) {
	var fields, values string
	//组合
	for n, v := range c.insertData.insertData {
		fields += n + ","
		switch p := v.(type) {
		case uint8, uint16, uint32, uint64, int8, int16, int32, int64, int:
			values += fmt.Sprintf("%d,", p)
		case float32, float64:
			values += fmt.Sprintf("%f,", p)
		case string:
			values += fmt.Sprintf("'%s',", p)
		case time.Time:
			t := p.Format(MIL_TIME_FORMAT)
			values += fmt.Sprintf("'%s',", t)
		case bool:
			if p {
				values += fmt.Sprintf("%d,", 1)
			} else {
				values += fmt.Sprintf("%d,", 0)
			}

		}
	}
	fields = subLastChar(fields)
	values = subLastChar(values)
	strSql = fmt.Sprintf("INSERT INTO %s (%s)VALUES(%s)", c.TableName, fields, values)
	return
}

func (c *Control) createUpdateSQL() (strSql string, err error) {
	uData := c.updateJoint(c.updateData)
	cond := c.conditionJoint(c.Condition)
	strSql = fmt.Sprintf("UPDATE %s SET %s ", c.TableName, uData)
	if cond != EMPTY {
		strSql += fmt.Sprintf(" WHERE %s", cond)
	}
	if c.lock != EMPTY {
		strSql += c.lock
	}
	return
}

func (c *Control) createQuerySQL() (strSql string, err error) {
	var fields string //查询字段字符串
	var cond string   //条件字符串
	if c.Fields != nil {
		//组合查询字段
		for i := 0; i < len(c.Fields); i++ {
			fields += c.Fields[i] + ","
		}
		fields = subLastChar(fields)
	} else {
		fields = ALL_FIELD
	}
	//条件组合
	cond = c.conditionJoint(c.Condition)
	//组合完整SQL简单语句
	strSql = fmt.Sprintf("SELECT %s FROM %s", fields, c.TableName)
	//拼接查询条件
	if cond != EMPTY {
		strSql += WHERE + cond
	}
	//拼接排序
	if c.OrderByFields != nil {
		for i := 0; i < len(c.OrderByFields); i++ {
			strSql += ORDER_BY + c.OrderByFields[i] + ","
		}
		strSql = subLastChar(strSql) + SPACE + c.Sort

	}
	//拼接返回记录条数
	if c.Limit > 0 {
		strSql += LIMIT + strconv.Itoa(c.Limit)
		//sql语句完成,归零
		c.Limit = 0
	} else if c.Limits.EndRow > 0 {
		strSql += LIMIT + fmt.Sprintf(" %d,%d ", c.Limits.StartRow, c.Limits.EndRow)
		//sql语句完成,归零
		c.Limits.StartRow = 0
		c.Limits.EndRow = 0
	}
	//上锁
	if c.lock != EMPTY {
		strSql += c.lock
	}
	return
}

func subLastChar(str string) (result string) {
	str = str_lib.SubString(str, 0, len(str)-1)
	result = str
	return
}
