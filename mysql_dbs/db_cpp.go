// db_cpp
package mysql_dbs

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"utils/data_conv/json_lib"
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
func (c *Mysql_Db) InitDB(user, pw, serverAddr, databaseName string) {
	c.dns = user + ":" + pw + "@tcp(" + serverAddr + ":3306)/" + databaseName + "?charset=utf8"
	//fmt.Println(c.dns)
	err := c.open()
	if err != nil {
		panic(err)
	} else {
		c.db.SetMaxOpenConns(50)
		c.db.SetMaxIdleConns(10)
	}
}

func (c *Mysql_Db) GetOperationObject() (clt Control) {
	clt.Mysql_Db = c
	return
}

/*
功能:创建插入控制器
参数:
	tbName:表名
	obj:表结构对像
	flag:(为真时fields内的字段为插入字段,为假时fields内的字段为过虑字段
	fields:字段名组数
*/
func (c *Mysql_Db) CreateInsertControl(tbName string, obj interface{}, flag bool, insertFields ...string) (ctl Control, err error) {
	var objMap map[string]interface{}
	err = json_lib.ObjectToObject(&objMap, obj)
	ctl = c.GetOperationObject()
	ctl.TableName = tbName
	if err == nil {
		if len(insertFields) > 0 {
			if flag {
				for i := range insertFields {
					ctl.SetInsertData(insertFields[i], objMap[insertFields[i]])
				}
			} else {
				for n, v := range objMap {
					if !c.inArray(n, insertFields) {
						ctl.SetInsertData(n, v)
					}
				}
			}
		} else {
			for n, v := range objMap {
				ctl.SetInsertData(n, v)
			}
		}
	}
	return
}

/*
功能:可控创建插入控制器
参数:
	tbName:表名
	obj:表结构对像
	flag:(为真时fields内的字段为插入字段,为假时fields内的字段为过虑字段
	fields:字段名组数
*/
func (c *Mysql_Db) TransactionInsertControl(ctl *Control, tbName string, obj interface{}, flag bool, insertFields ...string) (err error) {
	var objMap map[string]interface{}
	err = json_lib.ObjectToObject(&objMap, obj)
	ctl.TableName = tbName
	if err == nil {
		if len(insertFields) > 0 {
			if flag {
				for i := range insertFields {
					ctl.SetInsertData(insertFields[i], objMap[insertFields[i]])
				}
			} else {
				for n, v := range objMap {
					if !c.inArray(n, insertFields) {
						ctl.SetInsertData(n, v)
					}
				}
			}
		} else {
			for n, v := range objMap {
				ctl.SetInsertData(n, v)
			}
		}
	}
	return
}

/*
功能:创建更新控制器
参数:
	tbName:表名
	obj:表结构对像
	flag:(为真时fields内的字段为插入字段,为假时fields内的字段为过虑字段
	fields:字段名组数
*/
func (c *Mysql_Db) CreateUpdateControl(tbName string, obj interface{}, condition Conditions, flag bool, updateFields ...string) (ctl Control, err error) {
	var objMap map[string]interface{}
	err = json_lib.ObjectToObject(&objMap, obj)
	ctl = c.GetOperationObject()
	if err == nil {
		ctl.TableName = tbName
		if len(updateFields) > 0 {
			if flag {
				for i := range updateFields {
					ctl.SetUpdateData(updateFields[i], objMap[updateFields[i]])
				}
			} else {
				for n, v := range objMap {
					if !c.inArray(n, updateFields) {
						ctl.SetUpdateData(n, v)
					}
				}
			}
		} else {
			for n, v := range objMap {
				ctl.SetUpdateData(n, v)
			}
		}
		ctl.Condition = condition
	}
	return
}

/*
功能:可控创建更新控制器
参数:
	tbName:表名
	obj:表结构对像
	flag:(为真时fields内的字段为插入字段,为假时fields内的字段为过虑字段
	fields:字段名组数
*/
func (c *Mysql_Db) TransactionUpdateControl(ctl *Control, tbName string, obj interface{}, flag bool, updateFields ...string) (err error) {
	var objMap map[string]interface{}
	err = json_lib.ObjectToObject(&objMap, obj)
	if err == nil {
		ctl.TableName = tbName
		if flag {
			for i := range updateFields {
				ctl.SetUpdateData(updateFields[i], objMap[updateFields[i]])
			}
		} else {
			for n, v := range objMap {
				if !c.inArray(n, updateFields) {
					ctl.SetUpdateData(n, v)
				}
			}
		}
	}
	return
}

func (c *Mysql_Db) TransactionQuery(ctl *Control, tbName string, obj interface{}, fields ...string) (err error) {
	ctl.TableName = tbName
	ctl.Fields = fields
	err = ctl.QueryGetObject(&obj)
	return
}

func (c *Mysql_Db) inArray(key string, fields []string) (result bool) {
	for i := 0; i < len(fields); i++ {
		if key == fields[i] {
			result = true
			break
		}
	}
	return
}

func (c *Mysql_Db) defaultValue(v interface{}) (b bool) {
	switch n := v.(type) {
	case string:
		b = n == "NULL"
	case int, int32:
		b = n == math.MaxInt32
	case int16:
		b = n == math.MaxInt16
	case int64:
		b = n == math.MaxInt64
	case uint8:
		b = n == math.MaxUint8
	case uint, uint32:
		b = n == math.MaxUint32
	case uint16:
		b = n == math.MaxUint16
	case uint64:
		b = n == math.MaxUint64
	}
	return
}

//=====================================================================================================================

/*
功能:开始事条
参数:
返回:
*/
func (c *Control) Begin() (tx *sql.Tx, err error) {
	c.tx, err = c.db.Begin()
	return
}

/*
功能:提交事务
参数:
返回:信息信息
*/
func (c *Control) Commit() (err error) {
	err = c.tx.Commit()
	return
}

/*
功能:回滚事务
参数:
返回:信息信息
*/
func (c *Control) Rollback() (err error) {
	err = c.tx.Rollback()
	return
}

/*
功能:执行原生的SQLf增删改语句
参数:
	strSql：原生的SQL语句
返回:信息信息
*/
func (c *Control) ExecIUDSql(strSql string) (result int64, err error) {
	result, err = c.execSql(strSql, c.tx)
	return
}

/*
功能:执行原生的SQL查询语句
参数:
	strSql：原生的SQL语句
返回:信息信息
*/
func (c *Control) ExecQuery(strSql string) (result RecordSet, err error) {
	result, err = c.execQuerySql(strSql)
	return
}

//=====================================执行存储过程===============================

/*
功能:执行查询型,存储过程
参数:
	pro_name:存储过程名
	param_value:参数值
返回:记录集,错误信息
*/
func (c *Control) ExecQueryPro(proName string, paramValues []interface{}) (record RecordSet, err error) {
	//var rows *sql.Rows
	strValue := c.proParameterJoint(paramValues)
	strSql := fmt.Sprintf("CALL %s(%s)", proName, strValue)
	record, err = c.execQuerySql(strSql)
	return
}

func (c *Control) ExecQueryPro2(proName string, paramValues []interface{}) (record RecordSet, err error) {
	//var rows *sql.Rows
	strValue := c.proParameterJoint(paramValues)
	strSql := fmt.Sprintf("CALL %s(%s)", proName, strValue)
	record, err = c.execQuerySql2(strSql)
	return
}

/*
功能:执行查插入,更新,删除型 ,无返回值的存储过程
参数:
	pro_name:表名
	param_value:参数值
返回:影响的行数,错误信息
*/
func (c *Control) ExecPro(proName string, paramValues []interface{}) (result int64, err error) {
	strValue := c.proParameterJoint(paramValues)
	strSql := fmt.Sprintf("CALL %s(%s)", proName, strValue)
	result, err = c.execSql(strSql, c.tx)
	return
}

/*
功能:执行查插入,更新,删除型 ,有一个整形返回值的存储过程
参数:
	pro_name:表名
	param_value:参数值
返回:返回值,错误信息
*/
func (c *Control) ExecProOut(proName string, paramValues []interface{}) (returnVal int64, err error) {
	var record RecordSet
	outParam := "ret"
	strValue := c.proParameterJoint(paramValues)
	strSql := fmt.Sprintf("CALL %s(%s ,@%s)", proName, strValue, outParam)

	if err == nil {
		_, err = c.execSql(strSql, c.tx)
		if err == nil {
			strSql = fmt.Sprintf("SELECT @%s as %s", outParam, outParam)
			record, err = c.execQuerySql(strSql)
			if record.Count > 0 {
				returnVal = record.Data[0][outParam].(int64)
			} else {
				err = errors.New("the return value is not get")
			}
		}
	}
	return
}

//========================================RecordSet==============================================
func (c *RecordSet) GetObject(s interface{}) (err error) {
	var b []byte
	b, err = json.Marshal(c.Data)
	if err == nil {
		err = json.Unmarshal(b, &s)
	}
	return
}

//========================================Conditions==============================================

//等于
func (c *Conditions) Like(fieldName string, value interface{}) {
	if c.conditions == nil {
		c.conditions = make(map[string]interface{})
	}
	c.conditions[fieldName+" like "] = value
}

//等于
func (c *Conditions) Equal(fieldName string, value interface{}) {
	if c.conditions == nil {
		c.conditions = make(map[string]interface{})
	}
	c.conditions[fieldName+EQUAL] = value
}

func (c *Conditions) NotNil(fieldName string, value interface{}) {
	if c.conditions == nil {
		c.conditions = make(map[string]interface{})
	}
	c.conditions[fieldName+" is not "] = value
}

func (c *Conditions) EqualNil(fieldName string, value interface{}) {
	if c.conditions == nil {
		c.conditions = make(map[string]interface{})
	}
	c.conditions[fieldName+" is "] = value
}

//不等于
func (c *Conditions) NotEqual(fieldName string, value interface{}) {
	if c.conditions == nil {
		c.conditions = make(map[string]interface{})
	}
	c.conditions[fieldName+NOT_EQUAL] = value
}

//小于<
func (c *Conditions) Less(fieldName string, value interface{}) {
	if c.conditions == nil {
		c.conditions = make(map[string]interface{})
	}
	c.conditions[fieldName+LESS] = value
}

//大于>
func (c *Conditions) Greater(fieldName string, value interface{}) {
	if c.conditions == nil {
		c.conditions = make(map[string]interface{})
	}
	c.conditions[fieldName+GREATER] = value
}

//大于等于>=
func (c *Conditions) GreaterEqual(fieldName string, value interface{}) {
	if c.conditions == nil {
		c.conditions = make(map[string]interface{})
	}
	c.conditions[fieldName+GREATER_EQUAL] = value
}

//小于等于<=
func (c *Conditions) LessEqual(fieldName string, value interface{}) {
	if c.conditions == nil {
		c.conditions = make(map[string]interface{})
	}
	c.conditions[fieldName+LESS_EQUAL] = value
}

//条件连接运算符
func (c *Conditions) JointCondition(operation string) {
	if c.operation == nil {
		c.operation = make([]string, 0)
	}
	c.operation = append(c.operation, operation)
}

//=========================================UpdateValue==================================

func (c *UpdateValue) SetData(fieldName string, value interface{}) {
	if c.updateData == nil {
		c.updateData = make(map[string]interface{})
	}
	c.updateData[fieldName] = value
}

//=========================================InsertValue==================================

func (c *InsertValue) SetData(fieldName string, value interface{}) {
	if c.insertData == nil {
		c.insertData = make(map[string]interface{})
	}
	c.insertData[fieldName] = value
}

//=========================================Control==================================

/*
功能:设置插入数据
*/
func (c *Control) SetInsertData(fieldName string, value interface{}) {
	c.insertData.SetData(fieldName, value)
}

/*
功能:设置更新数据
*/
func (c *Control) SetUpdateData(fieldName string, value interface{}) {
	c.updateData.SetData(fieldName, value)
}

/*
功能:分页时指定取数据开始与结束位置
参数:
返回:
*/
func (c *Control) Page(startPos, endPos int) {
	c.Limits.StartRow = startPos
	c.Limits.EndRow = endPos
}

/*
功能:执行语句上锁
参数:
	t:锁类型(1=共享锁 2=排它锁)
返回:错误信息
*/
func (c *Control) OnLock(t byte) {
	if t == 1 {
		c.lock = " LOCK IN SHARE MODE" //共享锁
	} else {
		c.lock = " FOR UPDATE" //排它锁
	}
}

/*
功能:执行查询
参数:
	v:[input/out]具体数据类型切片
返回:错误信息
*/
func (c *Control) QueryGetObject(v interface{}) (err error) {
	strSql, err := c.createQuerySQL()
	if err == nil {
		err = c.execQueryObj(strSql, &v)
	}
	//c.ClearConditionData()
	return
}

/*
功能:执行查询
参数:
返回:RecordSet,错误信息
*/
func (c *Control) Query() (record RecordSet, err error) {
	strSql, err := c.createQuerySQL()
	if err == nil {
		record, err = c.execQuerySql(strSql)
	}
	//c.ClearConditionData()
	return
}

/*
功能:执行分页查询
参数:
返回:RecordSet,错误信息
*/
func (c *Control) QueryPage() (record RecordSet, err error) {
	strSql, err := c.createQuerySQL()
	if err == nil {
		record, err = c.execQuerySql(strSql)
	}
	//c.ClearConditionData()
	return
}

/*
功能:执行更新
参数:
返回:影响的记录条数,错误信息
*/
func (c *Control) Update() (result int64, err error) {
	strSql, err := c.createUpdateSQL()
	if err == nil {
		result, err = c.execSql(strSql, c.tx)
	}
	//c.ClearConditionData()
	//c.ClearUpdateData()
	return
}

/*
功能:执行插入
参数:
返回:影响的行数,错误信息
*/
func (c *Control) Insert() (result int64, err error) {
	strSql, err := c.createInsertSQL()
	if err == nil {
		result, err = c.execSql(strSql, c.tx)
	}
	return
}

/*
功能:执行删除
参数:
返回:影响的行数,错误信息
*/
func (c *Control) Delete() (result int64, err error) {
	cond := c.conditionJoint(c.Condition)
	strSql := fmt.Sprintf("DELETE FROM %s WHERE %s", c.TableName, cond)
	result, err = c.execSql(strSql, c.tx)
	//c.ClearConditionData()
	return
}

/*
功能:清除插入数据
*/
func (c *Control) ClearInsertData() {
	for n := range c.insertData.insertData {
		delete(c.insertData.insertData, n)
	}
	c.insertData.insertData = nil
}

/*
功能:清除更新数据
*/
func (c *Control) ClearUpdateData() {
	for n := range c.updateData.updateData {
		delete(c.updateData.updateData, n)
	}
	c.updateData.updateData = nil
}

/*
功能:清除条件数据
*/
func (c *Control) ClearConditionData() {
	for n := range c.Condition.conditions {
		delete(c.Condition.conditions, n)
	}
	c.Condition.conditions = nil
	c.Condition.operation = make([]string, 0)
}
