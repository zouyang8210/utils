// json
package json_lib

import (
	"database/sql"
	"encoding/json"
)

const (
	TIME_yyyymmddHHmmss = "2006-01-02 15:04:05"
	TIME_UTC_FORMAT     = "2006-01-02T15:04:05.000000Z07:00"
)

/*
功能:对像转json
参数:
	v:对像
返回:json 和错误信息
*/
func ObjectToJson(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err == nil {
		return string(b), err
	} else {
		return string(b), err
	}
}

/*
功能:json 转换成 对像
参数:
	data:json字符串
	v:对像(要传地址)
返回:json 和错误信息
*/
func JsonToObject(data string, v interface{}) (err error) {
	err = json.Unmarshal([]byte(data), &v)
	return
}

/*
功能:对像 转换成 对像
参数:
	desc:目标对像
	source:源对像
返回:json 和错误信息
*/
func ObjectToObject(desc, source interface{}) (err error) {
	json, err := ObjectToJson(source)
	if err == nil {
		err = JsonToObject(json, &desc)
	}
	return
}

/*
功能:sql.Rows转json(数组)
参数:
	rows:记录集
返回:json 字符串
*/
func RowsToJson(rows *sql.Rows) (int, string, error) {
	columns, err := rows.Columns()
	if err != nil {
		return 0, "", err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			switch t := val.(type) {
			case int8, int16, int, int32, int64, uint16, uint, uint32, uint64, float32, float64, bool, uint8:
				v = t
			case []byte:
				v = string(t)
				//如果字符串转时间不出错,则格式化为时间
				//tt, err := time.Parse(TIME_yyyymmddHHmmss, v.(string))
				//if err == nil {
				//	v = tt.Format(TIME_UTC_FORMAT)
				//}
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)

	if err != nil {
		return 0, "", err
	}
	strJson := string(jsonData)
	//strJson = strings.Replace(strJson, "\\", "", len(strJson))
	//strJson = SubString(strJson, 0, len(strJson))
	return len(tableData), strJson, err
}
