package zorm

import (
	"fmt"
)

//模糊查询
func (c *Conditions) Like(fieldName string, value interface{}) *Conditions {
	c.conditions = append(c.conditions, KV{K: fieldName, V: value, S: LIKE})
	return c
}

//IN
func (c *Conditions) In(fieldName string, value interface{}, bracket ...string) *Conditions {
	b := handleBracket(bracket...)
	c.conditions = append(c.conditions, KV{K: fieldName, V: value, S: IN, B: b})
	return c
}

//等于
func (c *Conditions) Equal(fieldName string, value interface{}, bracket ...string) *Conditions {
	b := handleBracket(bracket...)
	//if len(bracket) > 0 {
	//	b = bracket[0]
	//}
	c.conditions = append(c.conditions, KV{K: fieldName, V: value, S: EQUAL, B: b})
	return c
}

//不为空值
func (c *Conditions) NotNil(fieldName string, value interface{}, bracket ...string) *Conditions {
	b := handleBracket(bracket...)
	c.conditions = append(c.conditions, KV{K: fieldName, V: value, S: IS_NOT, B: b})
	return c
}

//为空值
func (c *Conditions) EqualNil(fieldName string, value interface{}, bracket ...string) *Conditions {
	b := handleBracket(bracket...)
	c.conditions = append(c.conditions, KV{K: fieldName, V: value, S: IS, B: b})
	return c
}

//不等于
func (c *Conditions) NotEqual(fieldName string, value interface{}, bracket ...string) *Conditions {
	b := handleBracket(bracket...)
	c.conditions = append(c.conditions, KV{K: fieldName, V: value, S: NOT_EQUAL, B: b})
	return c
}

//小于<
func (c *Conditions) Less(fieldName string, value interface{}, bracket ...string) *Conditions {
	b := handleBracket(bracket...)
	c.conditions = append(c.conditions, KV{K: fieldName, V: value, S: LESS, B: b})
	return c
}

//大于>
func (c *Conditions) Greater(fieldName string, value interface{}, bracket ...string) *Conditions {
	b := handleBracket(bracket...)
	c.conditions = append(c.conditions, KV{K: fieldName, V: value, S: GREATER, B: b})
	return c
}

//大于等于>=
func (c *Conditions) GreaterEqual(fieldName string, value interface{}, bracket ...string) *Conditions {
	b := handleBracket(bracket...)
	c.conditions = append(c.conditions, KV{K: fieldName, V: value, S: GREATER_EQUAL, B: b})
	return c
}

//小于等于<=
func (c *Conditions) LessEqual(fieldName string, value interface{}, bracket ...string) *Conditions {
	b := handleBracket(bracket...)
	c.conditions = append(c.conditions, KV{K: fieldName, V: value, S: LESS_EQUAL, B: b})
	return c
}

//设置查记录条数
func (c *Conditions) Limit(value int) *Conditions {
	c.limit = fmt.Sprintf(" %s %d ", LIMIT, value)
	return c
}

//分页查询
func (c *Conditions) Page(pageNo int, pageSize int) *Conditions {
	startIndex := (pageNo - 1) * pageSize
	c.pages = fmt.Sprintf(" %s %d,%d", LIMIT, startIndex, pageSize)
	return c
}

//设置升序排序
func (c *Conditions) StoreAsc(fields ...string) *Conditions {
	for _, v := range fields {
		c.store += v + ","
	}
	if len(c.store) > 1 {
		c.store = ORDER_BY + c.store[0:len(c.store)-1] + ASC
	}
	return c
}

//设置降序排序
func (c *Conditions) StoreDesc(fields ...string) *Conditions {
	for _, v := range fields {
		c.store += v + ","
	}
	if len(c.store) > 1 {
		c.store = ORDER_BY + c.store[0:len(c.store)-1] + DESC
	}
	return c
}

//清除where条件
func (c *Conditions) ClearWhere() *Conditions {
	c.conditions = []KV{}
	c.operation = []string{}
	return c
}

//清除所有条件
func (c *Conditions) ClearAll() *Conditions {
	c.conditions = []KV{}
	c.operation = []string{}
	c.limit = EMPTY
	c.pages = EMPTY
	c.store = EMPTY
	return c
}

//清除where之外的所有条件
func (c *Conditions) ClearOther() *Conditions {
	c.limit = EMPTY
	c.pages = EMPTY
	c.store = EMPTY
	return c
}

//条件连接AND运算符
func (c *Conditions) And() *Conditions {
	c.operation = append(c.operation, AND)
	return c
}

//条件连接OR运算符
func (c *Conditions) Or() *Conditions {
	c.operation = append(c.operation, OR)
	return c
}

func (c *Conditions) InnerJoin(tbName, leftField, rightField string) *Conditions {
	c.join = append(c.join, fmt.Sprintf(" INNER JOIN %s ON %s = %s", tbName, leftField, rightField))
	return c
}
