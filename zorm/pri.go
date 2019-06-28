package zorm

import (
	"database/sql"
	"fmt"
	"time"
)

func (c *Db) open() (err error) {
	c.db, err = sql.Open(DRIVER_MYSQL, c.dns)
	return
}

//测试sql执行用时
func (c *Db) runTime(now time.Time) {
	if c.Debug {
		terminal := time.Since(now)
		fmt.Printf("Exec Time:%d", terminal)
	}
}

//输出执行sql语句
func (c *Db) logs(str string) {
	if c.Debug {
		fmt.Printf("%s->%s\n", time.Now().Format(NORMAL_TIME_FORMAT), str)
	}
}

//把sql.Rows转换成RecordSet记录集
func (c *Db) getRecordSet(rows *sql.Rows) (record RecordSet) {
	columnName, err := rows.Columns()
	if err == nil {
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
					tt, err := time.Parse(NORMAL_TIME_FORMAT, str)
					if err == nil {
						str = tt.Format(UTC_TIME_FORMAT)
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
