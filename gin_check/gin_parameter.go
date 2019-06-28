package gin_check

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"utils/data_conv/json_lib"
	"utils/http_lib"
)

const (
	ERR_CODE          = "err_code" //
	ERR_MSG           = "err_msg"  //
	HTTP_SUCCESS      = 200        //
	EMPTY             = ""         //
	ERR_LACK_PARAM    = 1001       //缺少参数
	ERR_INVALID_PARAM = 1002       //参数无效
	MSG_IVALID_PARAM  = "无效的参数"    //
)

//POST方式参数检测
func CheckPostParameter(c *gin.Context, params ...string) (strJson string, mapData map[string]interface{}, err error) {
	var buff []byte
	buff, err = c.GetRawData()
	if mapData, err = getBodyData(buff, nil); err == nil {
		strJson = string(buff)
		if err = checkMapData(mapData, params...); err != nil {
			SimpleReturn(ERR_LACK_PARAM, err.Error(), c)
		}
	} else {
		SimpleReturn(ERR_INVALID_PARAM, MSG_IVALID_PARAM, c)
	}
	return
}

//GET方式参数检测
func CheckGetParameter(c *gin.Context, params ...string) (array map[string]string, err error) {
	array, err = http_lib.QueryString(c.Request)
	for _, v := range params {
		if array[v] == EMPTY {
			err = errors.New("lack parameter：" + v)
			SimpleReturn(ERR_LACK_PARAM, err.Error(), c)
			break
		}
	}
	return
}

//参数检测
func checkMapData(v map[string]interface{}, paramName ...string) (err error) {
	for n := range paramName {
		if v[paramName[n]] == nil {
			err = errors.New("lack parameter：" + paramName[n])
			break
		} else {
			switch x := v[paramName[n]].(type) {
			case float64:
				if x <= 0 {
					err = errors.New("parameter value invalid：" + paramName[n])
					break
				}
			case string:
				if x == EMPTY {
					err = errors.New("parameter value invalid：" + paramName[n])
					break
				}
			}
		}
	}
	return
}

//获取body数据(json字符串转成map)
func getBodyData(buffer []byte, e error) (mapData map[string]interface{}, err error) {
	if e == nil {
		mapData = make(map[string]interface{})
		err = json_lib.JsonToObject(string(buffer), &mapData)
	} else {
		err = e
	}
	return
}

func SimpleReturn(code int, msg string, c *gin.Context) {
	if code != 0 {
		c.Status(203)
	}
	obj := gin.H{ERR_CODE: code, ERR_MSG: msg}
	jsonLog, _ := json_lib.ObjectToJson(obj)
	fmt.Println("Return Result->", jsonLog)
	c.JSON(HTTP_SUCCESS, obj)
}
