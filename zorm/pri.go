package zorm

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"utils/data_conv/time_lib"
)

func (c *Db) open() (err error) {
	c.db, err = sql.Open(DRIVER_MYSQL, c.dns)
	return
}

//输出执行sql执行用时
func (c *Db) printDuration(now time.Time) {
	if c.Debug {
		terminal := time.Since(now)
		fmt.Printf("sql exec duration:%v\n", terminal)
	}
}

//输出执行sql语句
func (c *Db) logs(str string) {
	if c.Debug {
		fmt.Printf("%s->%s\n", time.Now().Format(NORMAL_TIME_FORMAT), str)
	}
}

//把sql.Rows转换成RecordSet记录集
func (c *Db) getRecordSet(rows *sql.Rows) (record RecordSet, err error) {
	var columnName []string
	if columnName, err = rows.Columns(); err != nil {
		return
	}
	columnLen := len(columnName)
	values := make([]interface{}, columnLen)
	valuePtr := make([]interface{}, columnLen)
	for rows.Next() {
		for i := 0; i < columnLen; i++ {
			valuePtr[i] = &values[i]
		}
		rows.Scan(valuePtr...)
		entry := make(map[string]interface{})
		for i, col := range columnName {
			if _, ok := values[i].([]byte); ok {
				str := fmt.Sprintf("%s", values[i])
				if tt, err := time.Parse(NORMAL_TIME_FORMAT, str); err == nil {
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
	return
}

//对像转对像
func objectToObject(desc, source interface{}) (err error) {
	var buff []byte
	if buff, err = json.Marshal(source); err != nil {
		return
	}
	err = json.Unmarshal(buff, &desc)
	return
}

//key是否在组组中
func (c *Db) inArray(key string, fields []string) (result bool) {
	for i := 0; i < len(fields); i++ {
		if key == fields[i] {
			result = true
			break
		}
	}
	return
}

//将对像转map数据,以方便拼接sql语句
func (c *Db) createMapData(obj interface{}, flag bool, fields ...string) (data map[string]interface{}, err error) {
	var objMap map[string]interface{}
	if err = objectToObject(&objMap, obj); err != nil {
		return
	}
	data = make(map[string]interface{})
	if len(fields) > 0 {
		if flag {
			for _, n := range fields {
				data[n] = objMap[n]
			}
		} else {
			for n, v := range objMap {
				if !c.inArray(n, fields) {
					data[n] = v
				}
			}
		}
	} else {
		for n, v := range objMap {
			data[n] = v
		}
	}
	return
}

//条接拼接
func (c *Db) createCondition(cond Conditions) (fieldCondition string, err error) {
	//条件组合
	i := 0
	//length := len(condition.operation)
	for _, v := range cond.conditions {
		switch p := v.V.(type) {
		case uint8, uint16, uint32, uint64, int8, int16, int32, int64, int, uint, float32, float64:
			fieldCondition += conditionJoinNumber(v)
		case string:
			fieldCondition += conditionJoinString(v)
		case nil:
			fieldCondition += fmt.Sprintf(" %s %s NULL ", v.K, v.S)
		case time.Time:
			v.V = time_lib.TimeToStr(p, time_lib.TIME_yyyyMMddHHmmss)
			fieldCondition += conditionJoinString(v)
		}
		if i < len(cond.conditions)-1 {
			fieldCondition += SPACE + cond.operation[i] + SPACE
		} else {
			break
		}
		i++
	}

	if cond.store != EMPTY {
		fieldCondition += cond.store
	}
	if cond.limit != EMPTY {
		fieldCondition += cond.limit
	} else if cond.pages != EMPTY {
		fieldCondition += cond.pages
	}
	return
}

//数字条件连接
func conditionJoinNumber(v KV) (str string) {
	if v.B == "(" {
		str += fmt.Sprintf(" %s%s %s %v ", v.B, v.K, v.S, v.V)
	} else if v.B == ")" {
		str += fmt.Sprintf(" %s %s %v%s ", v.K, v.S, v.V, v.B)
	} else {
		str += fmt.Sprintf(" %s %s %v ", v.K, v.S, v.V)
	}
	return
}

//字符条件连接
func conditionJoinString(v KV) (str string) {
	if v.B == "(" {
		str += fmt.Sprintf(" %s%s %s '%v' ", v.B, v.K, v.S, v.V)
	} else if v.B == ")" {
		str += fmt.Sprintf(" %s %s '%v'%s ", v.K, v.S, v.V, v.B)
	} else {
		str += fmt.Sprintf(" %s %s %v ", v.K, v.S, v.V)
	}
	return
}

//拼接更新sql语句
func (c *Db) createUpdateSql(tbName string, updateData map[string]interface{}, cond Conditions) (sql string) {
	//更新字段赋值
	fieldValue := ""
	for n, v := range updateData {
		switch p := v.(type) {
		case uint8, uint16, uint32, uint64, int8, int16, int32, int64, int:
			fieldValue += fmt.Sprintf("%s = %d,", n, p)
		case float32, float64:
			fieldValue += fmt.Sprintf("%s = %f,", n, p)
		case string:
			fieldValue += fmt.Sprintf("%s = '%s',", n, p)
		case nil:
			fieldValue += fmt.Sprintf("%s = NULL,", n)
		case time.Time:
			t := p.Format(MIL_TIME_FORMAT)
			fieldValue += fmt.Sprintf(" %s = '%s',", n, t)
		}
	}
	fieldValue = fieldValue[0 : len(fieldValue)-1]
	//条件组合
	fieldCondition, _ := c.createCondition(cond)
	sql = fmt.Sprintf("UPDATE %s SET %s ", tbName, fieldValue)
	if fieldCondition != EMPTY {
		sql += fmt.Sprintf(" WHERE %s", fieldCondition)
	}
	return
}

//接接插入sql语句
func (c *Db) createInsertSQL(tbName string, insertData map[string]interface{}) (strSql string, err error) {
	var fields, values string
	//组合
	for n, v := range insertData {
		fields += n + ","
		switch p := v.(type) {
		case uint8, uint16, uint32, uint64, uint, int8, int16, int32, int64, int:
			values += fmt.Sprintf("%d,", p)
		case float32, float64:
			values += fmt.Sprintf("%2.f,", p)
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
	fields = fields[0 : len(fields)-1]
	values = values[0 : len(values)-1]
	strSql = fmt.Sprintf("INSERT INTO %s (%s)VALUES(%s)", tbName, fields, values)
	return
}

//创建存储过程值
func (c *Db) createProcedureValue(parameterValue []interface{}) (strValue string) {
	if len(parameterValue) > 0 {
		for i := range parameterValue {
			switch p := parameterValue[i].(type) {
			case uint8, uint16, uint32, uint64, int8, int16, int32, int64, int, uint, float64, float32:
				strValue += fmt.Sprintf("%v,", p)
			case string:
				strValue += fmt.Sprintf("'%s',", p)
			}
		}
		strValue = strValue[0 : len(strValue)-1]
	}
	return
}
