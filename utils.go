package utils

//三目运算
func If(a bool, r1, r2 interface{}) (result interface{}) {
	if a {
		return r1
	} else {
		return r2
	}
}
