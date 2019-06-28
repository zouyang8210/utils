// number_lib
package number_lib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

//type FLOAT32 float32

const (
	DESC = "DESC" //倒序
	ASC  = "ASC"  //顺序
)

//
//功能:[]byte转对像
//参数:
//	buff:byte 数组
//	orderBy:数据流是正序还是倒序(ASC:正序,DESC:倒序)
//  value:输出参数(不支持传入int,uint和浮点型)
//返回:错误信息
//
func BytesToObject(buff []byte, orderBy string, value interface{}) (err error) {
	buf := bytes.NewBuffer(buff)
	var order binary.ByteOrder = binary.BigEndian
	if orderBy == DESC {
		order = binary.LittleEndian
	}
	err = binary.Read(buf, order, value)
	return

}

/*
功能:[]byte转成float
参数:
	buff:byte数组
	bitSize:数值的位数(取值32,64)
返回:转换的值,错误信息
*/
func BytesToFloat(buff []byte, bitSize int) (value interface{}, err error) {
	buf := bytes.NewBuffer(buff)
	switch bitSize {
	case 32:
		var f32 float32
		err = binary.Read(buf, binary.BigEndian, &f32)
		value = f32
	case 64:
		var f64 uint64
		err = binary.Read(buf, binary.BigEndian, &f64)
		value = f64
	}
	return
}

//功能:数字转[]byte
//参数:
//	number:数值
//	orderBy:数据流排序方式(ASC:正序,DESC:倒序)
//返回:byte数组
func ObjectToBytes(number interface{}, orderBy string) (buff []byte) {
	bytesBuffer := bytes.NewBuffer([]byte{})
	if orderBy == ASC {
		binary.Write(bytesBuffer, binary.BigEndian, number)
	} else if orderBy == DESC {
		binary.Write(bytesBuffer, binary.LittleEndian, number)
	}
	buff = bytesBuffer.Bytes()
	return
}

/*
功能:数值转字符串
参数:
	number:数值
返回:数字的字符串
*/
func NumberToStr(number interface{}) (str string) {
	switch n := number.(type) {
	case uint8, uint16, uint32, uint64, int8, int16, int32, int64, int:
		str = fmt.Sprintf("%d", n)
	case float32, float64:
		str = fmt.Sprintf("%f", n)
	}
	return
}

/*
功能:字符串转int
参数:
	str_number:字符串数值
	value:输出参数,对像地址
返回:错误信息
*/
func StrToInt(strNumber string, value interface{}) (err error) {
	var number interface{}
	number, err = strconv.ParseInt(strNumber, 10, 64)
	switch v := number.(type) {
	case int64:
		switch d := value.(type) {
		case *int64:
			*d = v
		case *int:
			*d = int(v)
		case *int16:
			*d = int16(v)
		case *int32:
			*d = int32(v)
		case *int8:
			*d = int8(v)
		}
	}
	return
}

/*
功能:字符串转uint
参数:
	str_number:字符串数值
	value:输出参数,对像地址
返回:错误信息
*/
func StrToUint(strNumber string, value interface{}) (err error) {
	var number interface{}
	number, err = strconv.ParseUint(strNumber, 10, 64)
	switch v := number.(type) {
	case uint64:
		switch d := value.(type) {
		case *uint64:
			*d = v
		case *uint:
			*d = uint(v)
		case *uint16:
			*d = uint16(v)
		case *uint32:
			*d = uint32(v)
		case *uint8:
			*d = uint8(v)
		}
	}
	return
}

/*
功能:字符串转float
参数:
	str_number:字符串
	[OUT]value:输出转换后的数值
返回:错误信息
*/
func StrToFloat(strNumber string, value interface{}) (err error) {
	var number interface{}
	number, err = strconv.ParseFloat(strNumber, 64)
	switch v := number.(type) {
	case float64:
		switch f := value.(type) {
		case *float64:
			*f = v
		case *float32:
			*f = float32(v)
		}
	}
	return
}

/*
功能:16进制字符串转uint数组
参数:
	str_data:
	start:开始转换的位置
	length:转换的字符串长度
	bitSize:数值的位数(取值0,8,16,32,64)
*/
func HexStrToUints(strData string, start, bitSize int) (num []interface{}, err error) {
	var n uint64
	flag := 8 //初始化标识为32位,8个字符
	if bitSize != 0 {
		flag = bitSize / 4
	}

	strLen := len(strData)            //字符串长度
	arrLen := (strLen - start) / flag //计算出数组最大长度

	num = make([]interface{}, arrLen)

	switch bitSize {
	case 0:
		for i := 0; strLen-start >= flag; i++ {
			n, err = strconv.ParseUint(strData[start:start+flag], 16, bitSize)
			num[i] = uint(n)
			start += flag
		}
	case 8:
		for i := 0; strLen-start >= flag; i++ {
			n, err = strconv.ParseUint(strData[start:start+flag], 16, bitSize)
			num[i] = uint8(n)
			start += flag

		}
	case 16:
		for i := 0; strLen-start >= flag; i++ {
			n, err = strconv.ParseUint(strData[start:start+flag], 16, bitSize)
			num[i] = uint16(n)
			start += flag

		}
	case 32:
		for i := 0; strLen-start >= flag; i++ {
			n, err = strconv.ParseUint(strData[start:start+flag], 16, bitSize)
			num[i] = uint32(n)
			start += flag
		}
	case 64:
		for i := 0; strLen-start >= flag; i++ {
			n, err = strconv.ParseUint(strData[start:start+flag], 16, bitSize)
			num[i] = n
			start += flag
		}
	}
	return
}

/*
功能:16进制字符串转int数组
参数:
	strData:
	start:开始转换的位置
	length:转换的字符串长度
	bitSize:数值的位数(取值0,8,16,32,64)
*/
func HexStrToInts(strData string, start, bitSize int) (num []interface{}, err error) {
	var n int64
	flag := 8 //初始化标识为32位,8个字符
	if bitSize != 0 {
		flag = bitSize / 4
	}

	strLen := len(strData)
	arrLen := (strLen - start) / flag

	num = make([]interface{}, arrLen)

	switch bitSize {
	case 0:
		for i := 0; strLen-start >= flag; i++ {
			n, err = strconv.ParseInt(strData[start:start+flag], 16, bitSize)
			num[i] = int(n)
			start += flag

		}
	case 8:
		for i := 0; strLen-start >= flag; i++ {
			n, err = strconv.ParseInt(strData[start:start+flag], 16, bitSize)
			num[i] = int8(n)
			start += flag

		}
	case 16:
		for i := 0; strLen-start >= flag; i++ {
			n, err = strconv.ParseInt(strData[start:start+flag], 16, bitSize)
			num[i] = int16(n)
			start += flag

		}
	case 32:
		for i := 0; strLen-start >= flag; i++ {
			n, err = strconv.ParseInt(strData[start:start+flag], 16, bitSize)
			num[i] = int32(n)
			start += flag
		}
	case 64:
		for i := 0; strLen-start >= flag; i++ {
			n, err = strconv.ParseInt(strData[start:start+flag], 16, bitSize)
			num[i] = n
			start += flag
		}
	}
	return
}
