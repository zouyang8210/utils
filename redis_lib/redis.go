package redis_lib

import (
	"database/sql"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"reflect"
	"utils/data_conv/number_lib"
	"utils/data_conv/time_lib"
)

//redis动作
const (
	GET     = "GET"
	SET     = "SET"
	DEL     = "DEL"
	EXISTS  = "EXISTS"
	EXPIRE  = "EXPIRE"
	HMSET   = "HMSET"
	HGET    = "HGET"
	HGETALL = "HGETALL"
)

//格式化占位符
const (
	FS = "%s" //格式化为字符串
)

const OK = "OK"

type Redis struct {
	network string
	address string
	conn    redis.Conn
}

/*
功能:连接redis服务器
参数:
	network:
	address:
返回:错误信息
*/
func (sender *Redis) Connect(network, address string) (err error) {
	c, err := redis.Dial(network, address)
	if err == nil {
		sender.conn = c
		sender.network = network
		sender.address = address
	}
	return
}

/*
功能:关闭redis服务器的连接
参数:
返回:错误信息
*/
func (sender *Redis) Close() (err error) {
	err = sender.conn.Close()
	return
}

//================================================基本操作===================================================
/*
功能:判断key是否存在
参数:
	key:
返回:是否存在,错误信息
*/
func (sender *Redis) ExistsKey(key string) (result bool, err error) {
	exists, err := sender.conn.Do(EXISTS, key)
	result = exists.(int64) == 1
	return
}

/*
功能:删除一个key
参数:
	key:
返回:是否成功,错误信息
*/
func (sender *Redis) DeleteKey(key string) (result bool, err error) {
	bit, err := sender.conn.Do(DEL, key)
	result = bit.(int64) == 1
	return
}

/*
功能:设置key可保存时间(秒)
参数:
	key:
	seconds:秒数
返回:是否成功,错误信息
*/
func (sender *Redis) ExpireKey(key string, seconds int) (result bool, err error) {
	bit, err := sender.conn.Do(EXPIRE, key, seconds)
	result = bit.(int64) == 1
	return
}

//===================================================字符串操作====================================================
/*
功能:获取key的值
参数:
	key:
返回:值,错误信息
*/
func (sender *Redis) Get(key string) (value string, err error) {
	bit, err := sender.conn.Do(GET, key)
	if err == nil {
		if bit != nil {
			value = bit.(string)
		}
	}
	return
}

/*
功能:设置key的值
参数:
	key:
	value:
返回:是否成功,错误信息
*/
func (sender *Redis) Set(key, value string) (result bool, err error) {
	bit, err := sender.conn.Do(SET, key, value)
	result = bit.(string) == OK
	return
}

//==========================================哈希操作==================================================

func (sender *Redis) HmSetRows(key string, rows sql.Rows) (result bool, err error) {

	return
}

/*
功能:把一个对存储为哈希值
参数:
	key:
	obj:
返回:是否成功,错误信息
*/
func (sender *Redis) HmSetObj(key string, obj interface{}) (result bool, err error) {
	mValue := reflect.ValueOf(obj)
	numField := mValue.NumField()
	l := numField + 1
	if l%2 == 0 {
		l++
	}
	inputValue := make([]interface{}, l)
	inputValue[0] = key
	for i := 1; i <= numField; i++ {
		inputValue[i] = mValue.Field(i - 1)
	}
	bit, err := sender.conn.Do(HMSET, inputValue...)
	result = bit.(string) == OK
	return
}

/*
功能:获取哈希值,转为一个具体对像输出
参数:
	key:
	v:对像,需传入指针[输出参数]
返回:是否成功,错误信息
*/
func (sender *Redis) HmGetAll(key string, v interface{}) (result bool, err error) {
	result = true
	bit, err := sender.conn.Do(HGETALL, key)
	mValue := reflect.ValueOf(v).Elem()
	numField := mValue.NumField()
	value := bit.([]interface{})
	l := 0
	if len(value) > numField {
		l = numField
	} else {
		l = len(value)
	}
	for i := 0; i < l; i++ {
		switch mValue.Field(i).Type().String() {
		case "string":
			if mValue.Field(i).CanSet() {
				mValue.Field(i).SetString(fmt.Sprintf(FS, value[i]))
			}
		case "int", "int32", "int16", "int64":
			if mValue.Field(i).CanSet() {
				var iTmp int64
				number_lib.StrToInt(fmt.Sprintf(FS, value[i]), &iTmp)
				mValue.Field(i).SetInt(iTmp)
			}
		case "uint", "uint32", "uint16", "uint64":
			if mValue.Field(i).CanSet() {
				var iTmp uint64
				number_lib.StrToUint(fmt.Sprintf(FS, value[i]), &iTmp)
				mValue.Field(i).SetUint(iTmp)
			}
		case "time.Time":
			if mValue.Field(i).CanSet() {
				tTmp := time_lib.StrToTime(fmt.Sprintf(FS, value[i])[0:26])
				mValue.Field(i).Set(reflect.ValueOf(tTmp))
			}
		case "uint8", "byte":
			if mValue.Field(i).CanSet() {
				mValue.Field(i).Set(reflect.ValueOf(value[i].(uint8)))
			}
		case "int8":
			if mValue.Field(i).CanSet() {
				mValue.Field(i).Set(reflect.ValueOf(value[i].(int8)))
			}
		}
	}
	return
}
