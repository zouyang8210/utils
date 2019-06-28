package database_tool

import (
	"fmt"
	"github.com/gohouse/converter"
	"strings"
	"utils/file"
)

type DatabaseTool struct {
	UserName     string
	Password     string
	DatabaseName string
	Address      string
}

const DNS = "%s:%s@tcp(%s:3306)/%s?charset=utf8"

func (c *DatabaseTool) CreateTableStruct(tableName, filePath string) (err error) {
	// 初始化
	t2t := converter.NewTable2Struct()
	// 个性化配置
	t2t.Config(&converter.T2tConfig{
	//RmTagIfUcFirsted: false,// 如果字段首字母本来就是大写, 就不添加tag, 默认false添加, true不添加
	//TagToLower: false, // tag的字段名字是否转换为小写, 如果本身有大写字母的话, 默认false不转
	//UcFirstOnly: false,// 字段首字母大写的同时, 是否要把其他字母转换为小写,默认false不转换
	//SeperatFile: false,// 每个struct放入单独的文件,默认false,放入同一个文件(暂未提供)
	})
	// 开始迁移转换
	if tableName != "" {
		t2t = t2t.Table(tableName) // 指定某个表,如果不指定,则默认全部表都迁移
	}
	//t2t = t2t.Prefix("prefix_")// 表前缀
	t2t = t2t.EnableJsonTag(true) // 是否添加json tag
	//t2t = t2t.PackageName("")     // 生成struct的包名(默认为空的话, 则取名为: package model)
	t2t = t2t.TagKey("orm") // tag字段的key值,默认是orm
	//t2t = t2t.RealNameMethod("TableName")// 是否添加结构体方法获取表名
	t2t = t2t.SavePath(filePath) // 生成的结构体保存路径
	dns := fmt.Sprintf(DNS, c.UserName, c.Password, c.Address, c.DatabaseName)
	fmt.Println(dns)
	t2t = t2t.Dsn(dns) // 数据库dsn,这里可以使用 t2t.DB() 代替,参数为 *sql.DB 对象

	err = t2t.Run() // 执行
	return
}

func (c *DatabaseTool) CreateFunction(tableName, className, savePath string) {
	m := module
	const_tb := strings.ToUpper(tableName)
	m = strings.Replace(m, "[CONST_TABLE]", const_tb, 1)
	m = strings.Replace(m, "[tablename]", tableName, 1)
	m = strings.Replace(m, "[ClassName]", className, 1)
	m = strings.Replace(m, "[StructName]", strings.ToUpper(tableName[0:1])+tableName[1:], 1)
	m = strings.Replace(m, "[TABLENAME]", "TB_"+const_tb, 1)

	if err := c.CreateTableStruct(tableName, "c:/aaa/tmp.json"); err == nil {
		buff, err := file.ReadFile("c:/aaa/tmp.json")
		if err == nil {
			strBuf := string(buff)
			m = strings.Replace(m, "[STRUCT]", strBuf, 1)
		}
	}

	file.WriteFile(savePath, []byte(m))
}
