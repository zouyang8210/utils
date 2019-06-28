package weixin

import (
	"encoding/xml"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"utils/crypto"
	"utils/data_conv/str_lib"
	"utils/http_lib"
)

func (sender *WeiXinPay) unifiedOrder(body, orderId, notifyUrl, tradeType string, fee int, replenish map[string]string) (result RetUnifiedOrder, err error) {
	params := make(map[string]string)
	sender.putPublic(params)
	sender.putParam(params, BODY, body)
	sender.putParam(params, OUT_TRADE_NO, orderId)
	sender.putParam(params, TOTAL_FEE, strconv.Itoa(fee))
	sender.putParam(params, NOTIFY_URL, notifyUrl)
	sender.putParam(params, TRADE_TYPE, tradeType)
	if replenish != nil {
		for n, v := range replenish {
			sender.putParam(params, n, v)
		}
	}
	xmlParam, _ := sender.putComplete(params)                                  //增加参数完成,返回Post数据和签名
	buffer, _, err := sender.submit(http_lib.POST, GET_PAY_CODE_URL, xmlParam) //post提交数据
	fmt.Println("buffer=", string(buffer))
	if err == nil {
		err = xml.Unmarshal(buffer, &result)
	}
	return
}

/*
功能:解析微返回的状态
参数:
	info:微信平台返回值
返回:状态码,状态描述
*/
func (sender *WeiXinPay) analysisWxReturn(info1 RetBase, info2 RetPublic) (int, string) {
	var errCode = 0
	var errMsg = ""
	if info1.ReturnCode != SUCCESS {
		errCode = FAIL
		errMsg = info1.ReturnMsg
	} else if info2.ResultCode != SUCCESS {
		errCode = FAIL
		errMsg = info2.ErrCodeDes
	}

	return errCode, errMsg
}

/*
功能:支付的字符串状态解析为数字状态
参数:
	state:状态字符串
返回:状态数字码
*/
func (sender *WeiXinPay) analysisPayStatus(state string) int {
	var status = 0
	switch state {
	case PAY_STATUS_SUCCESS:
		status = 0
	case PAY_STATUS_FEFUND:
		status = 1
	case PAY_STATUS_NOTPAY:
		status = 2
	case PAY_STATUS_CLOSED:
		status = 3
	case PAY_STATUS_REVOKED:
		status = 4
	case PAY_STATUS_USERPAYING:
		status = 5
	case PAY_STATUS_PAYERROR:
		status = 6
	default:
		status = -1
	}
	return status
}

/*
功能:提交数据
参数:
	submitType:提交类型
	url:接口路径
	data:提交的字符串
返回:API返回数据,http状态代码,错误
*/
func (sender *WeiXinPay) submit(submitType, url, data string) ([]byte, int, error) {
	resp, err := http_lib.HttpSubmit(submitType, url, data, nil)
	return resp.BufferBody, resp.HttpStatus, err

}

/*
功能:填充公共参数
参数:
	params:参数集合
返回:
*/
func (sender *WeiXinPay) putPublic(params map[string]string) {
	nonce := str_lib.Guid() //生成随机字符串
	sender.putParam(params, MCH_ID, sender.MchId)
	sender.putParam(params, NONCE_STR, nonce)
	sender.putParam(params, SPBILL_CREATE_IP, LOCALHOST)
	sender.putParam(params, APPID, sender.AppId)
}

/*
功能:小程序填充公共参数
参数:
	params:参数集合
返回:
*/
func (sender *WeiXinPay) putMinPublic(params map[string]string) {
	nonce := str_lib.Guid() //生成随机字符串
	sender.putParam(params, APPID, sender.MinProgramId)
	sender.putParam(params, MCH_ID, sender.MchId)
	sender.putParam(params, NONCE_STR, nonce)
	sender.putParam(params, SPBILL_CREATE_IP, "127.0.0.1")

}

/*
功能:填充单个参数
参数:
	params:参数集合
	param_name:参数名称
	param_value:参数集合
返回:
*/
func (sender *WeiXinPay) putParam(params map[string]string, paramName string, paramValue string) {
	params[paramName] = paramValue
}

/*
功能:整理提交的数据并计算签名
参数:
	params:参数集合
返回:整理完成的提交数据,签名
*/
func (sender *WeiXinPay) putComplete(params map[string]string) (string, string) {
	var keys []string
	//排序
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	//拼连提交数据和待签名字符串
	var xmlParam string
	var sign string
	for _, i := range keys {
		xmlParam += "<" + i + ">" + params[i] + "</" + i + ">"
		sign += i + "=" + params[i] + "&"
	}
	sign += "key=" + sender.ApiSecret
	sign = crypto.Md5(sign) //MD5加密
	sign = strings.ToUpper(sign)
	xmlParam = "<xml>" + xmlParam + "<sign>" + sign + "</sign></xml>"
	return xmlParam, sign
}

/*
功能:清空一个 map[string]string
参数:
返回:
*/
func (sender *WeiXinPay) clearParam(params map[string]string) {
	for k := range params {
		delete(params, k)
	}
}
