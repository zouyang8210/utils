package zorm

import (
	"fmt"
)

//模糊查询
func (c *Conditions) Like(fieldName string, value interface{}) *Conditions {
	c.conditions = append(c.conditions, KV{K: fieldName, V: value, S: " like "})
	return c
}

//
func (c *Conditions) In(fieldName string, value interface{}) *Conditions {
	c.conditions = append(c.conditions, KV{K: fieldName, V: value, S: " in "})
	return c
}

//等于
func (c *Conditions) Equal(fieldName string, value interface{}, bracket ...string) *Conditions {
	b := ""
	if len(bracket) > 0 {
		b = bracket[0]
	}
	c.conditions = append(c.conditions, KV{K: fieldName, V: value, S: " = ", B: b})
	return c
}

//不为空值
func (c *Conditions) NotNil(fieldName string, value interface{}) *Conditions {
	c.conditions = append(c.conditions, KV{K: fieldName, V: value, S: " is not "})
	return c
}

//为空值
func (c *Conditions) EqualNil(fieldName string, value interface{}) *Conditions {
	c.conditions = append(c.conditions, KV{K: fieldName, V: value, S: " is "})
	return c
}

//不等于
func (c *Conditions) NotEqual(fieldName string, value interface{}) *Conditions {
	c.conditions = append(c.conditions, KV{K: fieldName, V: value, S: " != "})
	return c
}

//小于<
func (c *Conditions) Less(fieldName string, value interface{}) *Conditions {
	c.conditions = append(c.conditions, KV{K: fieldName, V: value, S: " < "})
	return c
}

//大于>
func (c *Conditions) Greater(fieldName string, value interface{}) *Conditions {
	c.conditions = append(c.conditions, KV{K: fieldName, V: value, S: " > "})
	return c
}

//大于等于>=
func (c *Conditions) GreaterEqual(fieldName string, value interface{}) *Conditions {
	c.conditions = append(c.conditions, KV{K: fieldName, V: value, S: " >= "})
	return c
}

//小于等于<=
func (c *Conditions) LessEqual(fieldName string, value interface{}, bracket ...string) *Conditions {
	b := ""
	if len(bracket) > 0 {
		b = bracket[0]
	}
	c.conditions = append(c.conditions, KV{K: fieldName, V: value, S: " <= ", B: b})
	return c
}

//设置查记录条数
func (c *Conditions) Limit(value int) {
	c.limit = fmt.Sprintf(" LIMIT %d ", value)
}

//分页查询
func (c *Conditions) Page(pageNo int, pageSize int) {
	startIndex := (pageNo - 1) * pageSize
	c.pages = fmt.Sprintf(" LIMIT %d,%d", startIndex, pageSize)
}

//设置升序排序
func (c *Conditions) StoreAsc(fields ...string) {
	for _, v := range fields {
		c.store += v + ","
	}
	if len(c.store) > 1 {
		c.store = " ORDER BY " + c.store[0:len(c.store)-1] + " ASC"
	}
}

//设置降序排序
func (c *Conditions) StoreDesc(fields ...string) {
	for _, v := range fields {
		c.store += v + ","
	}
	if len(c.store) > 1 {
		c.store = " ORDER BY " + c.store[0:len(c.store)-1] + " DESC"
	}
}

//清除where条件
func (c *Conditions) ClearWhere() {
	c.conditions = []KV{}
	c.operation = []string{}
}

//清除所有条件
func (c *Conditions) ClearAll() {
	c.conditions = []KV{}
	c.operation = []string{}
	c.limit = EMPTY
	c.pages = EMPTY
	c.store = EMPTY
}

//清除where之外的所有条件
func (c *Conditions) ClearOther() {
	c.limit = EMPTY
	c.pages = EMPTY
	c.store = EMPTY
}

//条件连接AND运算符
func (c *Conditions) And() {
	c.operation = append(c.operation, " AND ")
}

//条件连接OR运算符
func (c *Conditions) Or() {
	c.operation = append(c.operation, " OR ")
}
