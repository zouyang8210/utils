package tbk_lib

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
	"utils/comm_const"
	"utils/data_conv/json_lib"
	"utils/data_conv/number_lib"
	"utils/data_conv/str_lib"
	"utils/http_lib"
)

func (sender *TaoBaoKey) search(reqStr string) (jsonData string, objData CommodityInfo, err error) {
	resp, err := http_lib.HttpsSubmit(http_lib.GET, reqStr, EMPTY, nil)
	if err == nil {
		//fmt.Println(resp.Body)
		jsonData = sender.conv(resp.BufferBody)
		//json有一个string字段，为关键字。此处替换掉
		jsonData = strings.Replace(jsonData, "string", "small_imgs", 10)
		json_lib.JsonToObject(jsonData, &objData)
		var zk float32 = 0
		for i := 0; i < len(objData.ResultList.MapData); i++ {
			objData.ResultList.MapData[i].CouponPrice = sender.getCouponLaterPrice(objData.ResultList.MapData[i].CouponInfo)
			number_lib.StrToFloat(objData.ResultList.MapData[i].ZkFinalPrice, &zk)
			zk = zk - objData.ResultList.MapData[i].CouponPrice
			number_lib.StrToFloat(fmt.Sprintf("%.2f", zk), &zk)
			objData.ResultList.MapData[i].CouponLater = zk //- objData.ResultList.MapData[i].CouponPrice

			if objData.ResultList.MapData[i].ShortTitle == EMPTY {
				objData.ResultList.MapData[i].ShortTitle = objData.ResultList.MapData[i].Title
			}
		}
		jsonData, _ = json_lib.ObjectToJson(objData)
	}
	return
}
func (sender *TaoBaoKey) getCouponLaterPrice(couponInfo string) float32 {
	var price float32 = 0
	pos1 := strings.Index(couponInfo, "减")
	if pos1 >= 0 {
		pos1 += 3
		pos2 := strings.LastIndex(couponInfo, "元")
		strPrice := str_lib.SubString(couponInfo, pos1, pos2-pos1)
		number_lib.StrToFloat(strPrice, &price)
	}
	return price
}

/*
功能:填充公共参数
参数:
	method:接口名称
	notifyUrl:回调地址
返回:公共参数
*/
func (sender *TaoBaoKey) putPublic(method string) (params map[string]interface{}) {
	params = make(map[string]interface{})
	params["app_key"] = sender.AppKey
	params["method"] = method
	params["sign_method"] = "md5"
	params[SIGN] = EMPTY
	params["format"] = "json"
	params["timestamp"] = time.Now().Format(comm_const.TIME_yyyyMMddHHmmss)
	params["v"] = "2.0"
	//params["simplify"] = false
	return
}

/*
功能:填加参数
参数:
	name:参数名称
	value:参数值
返回:
*/
func (sender *TaoBaoKey) putParam(params map[string]interface{}, name string, value interface{}) {
	params[name] = value
	return
}

/*
功能:整理提交的数据并计算签名
参数:
	params:参数集合
返回:整理完成的提交数据,签名
*/
func (sender *TaoBaoKey) putComplete(params map[string]interface{}) (reqStr string, sign string) {
	//排序
	keys := sender.sort(params)

	//拼连提交数据和待签名字符串
	var waitSign string
	for _, n := range keys {
		if n != SIGN {
			waitSign += fmt.Sprintf("%s%v", n, params[n])
		}
	}
	//fmt.Println("wait sign:", waitSign)
	//签名
	sign = sender.GetMd5String(sender.Secret + waitSign + sender.Secret)
	//签名加入参数集合
	sender.putParam(params, SIGN, sign)

	//加入签名后,再次排序
	keys = make([]string, 0)
	keys = sender.sort(params)

	//参数值转码
	for _, n := range keys {
		reqStr += fmt.Sprintf("%s=%v&", str_lib.UrlToUrlEncode(n), str_lib.UrlToUrlEncode(fmt.Sprintf("%v", params[n])))
	}

	//去除字符串中,最后一个'&'符号
	reqStr = API_GATEWAY + str_lib.SubString(reqStr, 0, len(reqStr)-1)
	//fmt.Println("reqStr:", reqStr)
	return
}

/*
功能:参数名排序
参数:
	params:参数集合
返回:排序后的参数名称数组
*/
func (sender *TaoBaoKey) sort(params map[string]interface{}) (array []string) {
	for k := range params {
		array = append(array, k)
	}
	sort.Strings(array)
	return
}

//生成32位md5字串
func (sender *TaoBaoKey) GetMd5String(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

/*
功能:把接口返回的JSON处理成我们需要样式
参数:
	data:待签名字符串
返回:签名,错误信息
*/
func (sender *TaoBaoKey) conv(buff []byte) (json string) {
	json = string(buff)
	pos1 := strings.Index(json, ":") + 1
	pos2 := strings.LastIndex(json, "}")
	json = str_lib.SubString(json, pos1, pos2-pos1)
	return
}
