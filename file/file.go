// file
package file

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
	"fmt"
	"utils/comm_const"
)

const middle = "========="

/*
功能：在当程序目录建立Logs目录写日志
	data:日志内容
*/
func WriteLog(data string) {
	mPath := GetCurrentDirectory() + "/Logs/"
	if !DirExists(mPath) {
		Mkdir(mPath)
	}
	mPath += time.Now().Format("20060102") + ".txt"
	Logs(mPath, data)
}

/*
功能：写文件
	file_name:文件全路径
	data:数据流
*/
func WriteFile(fileName string, data []byte) bool {
	err := ioutil.WriteFile(fileName, data, 0666)
	if err == nil {
		return true
	} else {
		fmt.Println(err)
		return false
	}

}

/*
功能：记录日志
参数：
	file_name:文件名全路径
	data:内容
*/
func Logs(file_name, data string) {
	WriteTextFile(file_name, "========================================================================================================================\r\n")
	WriteTextFile(file_name, time.Now().Format(comm_const.TIME_yyyyMMddHHmmss)+":"+data+"\r\n")

}

//追加，写文本文件
func WriteTextFile(path, data string) bool {
	var result = false

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()
	if err == nil {
		_, err = file.WriteString(data)
		if err == nil {
			result = true
		}
	}

	return result
}

//判断文件或目录是否存在
func DirExists(path string) bool {
	var result bool = false
	_, err := os.Stat(path)

	if err == nil {
		result = true
	}
	return result
}

//创建目录
func Mkdir(path string) bool {
	err := os.MkdirAll(path, 07777)
	if err == nil {
		return true
	} else {
		return false
	}
}

//获取当前目录
func GetCurrentDirectory() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

//读取文件
func ReadFile(filepath string) ([]byte, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

//功能:读取配置文件
//参数:
//	node:头
//	key:主键
//	path:配置文件路径
//返回:主键的值
func ReadConfig(node string, key string, path string) string {
	myConfig := new(Config)
	myConfig.initconfig(path)
	return myConfig.read(node, key)
}

type Config struct {
	Mymap  map[string]string
	strcet string
}

//功能:读取配置文件初始化
//参数:
//	path:配置文件路径
//返回:
func (c *Config) initconfig(path string) {
	c.Mymap = make(map[string]string)

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		s := strings.TrimSpace(string(b))
		//fmt.Println(s)
		if strings.Index(s, "#") == 0 {
			continue
		}
		n1 := strings.Index(s, "[")
		n2 := strings.LastIndex(s, "]")
		if n1 > -1 && n2 > -1 && n2 > n1+1 {
			c.strcet = strings.TrimSpace(s[n1+1 : n2])
			continue
		}

		if len(c.strcet) == 0 {
			continue
		}
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}

		frist := strings.TrimSpace(s[:index])
		if len(frist) == 0 {
			continue
		}
		second := strings.TrimSpace(s[index+1:])

		pos := strings.Index(second, "\t#")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " #")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, "\t//")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " //")
		if pos > -1 {
			second = second[0:pos]
		}

		if len(second) == 0 {
			continue
		}
		key := c.strcet + middle + frist
		c.Mymap[key] = strings.TrimSpace(second)
	}
}

/*
功能:读取配置文件
参数:
	node:头
	key:主键
返回:主键的值
*/
func (c *Config) read(node, key string) string {

	key = node + middle + key
	v, found := c.Mymap[key]
	if !found {
		return ""
	}
	return v
}
