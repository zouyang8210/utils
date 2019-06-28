// str_lib
package str_lib

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"net/url"
)

/*
功能:截取字符串
参数:
	str:原字符串
	begin:开始索引
	lenght:长度
返回:截取后的字符串
*/
func SubString(str string, begin, length int) (substr string) {
	lth := len(str)
	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}
	substr = str[begin:end]
	return
}

/*
功能:字符串中插入字符串
参数:
	s:原字符串
	insert_index:插入的位置
	value:插入的字符串
返回:返回插入字符串后的重组字符串
*/
func InsertStr(s string, insertIndex int, value string) (str string) {
	if insertIndex > len(s) {
		str = s + value
	} else {
		str1 := SubString(s, 0, insertIndex) + value
		str2 := SubString(s, insertIndex, len(s)-insertIndex)
		str = str1 + str2
	}
	return
}

//生成GUID
func Guid() (strMd5 string) {
	b := make([]byte, 48)
	_, err := io.ReadFull(rand.Reader, b)
	if err == nil {
		str := base64.URLEncoding.EncodeToString(b)
		h := md5.New()
		h.Write([]byte(str))
		strMd5 = hex.EncodeToString(h.Sum(nil))
	}
	return
}

//urlEncode编码
func UrlToUrlEncode(path string) string {
	return url.QueryEscape(path)
}
