/*
作者:邹阳明
描述:所有微信支付功能
*/

package wechat_pay

//微信支付

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"time"
	"utils/data_conv/json_lib"
	"utils/data_conv/str_lib"
	"utils/http_lib"
	. "utils/wechat"
	"utils/wechat/mini_program"
)

//微信支付对像
type WXPay struct {
	AppId            string `json:"app_id"`     //公众号或应用appId
	MchId            string `json:"mch_id"`     //商户号
	AppSecret        string `json:"app_secret"` //公众号或应用密钥
	ApiSecret        string `json:"api_secret"` //商户号密钥
	MinProgramId     string `json:"_"`          //小程序ID(小程序支付时需要)
	MinProgramSecret string `json:"_"`          //小程序密钥(小程序支付时需要)
}

/*
功能:公众号获取openId
参数:
	code:获取openid的验证码
返回
*/
func (sender *WXPay) PublicGetOpenId(code string) (openId string, err error) {
	var info RetAuthorizeOpenid
	var strBody []byte
	var httpStatus int
	strUrl := fmt.Sprintf(AUTHORIZE_URL, sender.AppId, sender.AppSecret, code)
	strBody, httpStatus, err = Submit(http_lib.GET, strUrl, EMPTY)
	if err == nil {
		if httpStatus != 200 {
			err = errors.New(fmt.Sprintf("http status %d", httpStatus))
		} else {
			json_lib.JsonToObject(string(strBody), &info)
			if info.ErrCode != 0 {
				err = errors.New(fmt.Sprintf("errCode=%d,errMsg=%s", info.ErrCode, info.ErrMsg))
			} else {
				openId = info.OpenId
			}
		}
	}
	return
}

/*
功能：公众号下单
参数:
	body:收费显示标题
	orderId:单号
	code:获取openId的代码
	fee:收费金额
返回:公众号支付调起收银台所需要的参数(详细参见 RetMinProgramPay 定义)
*/
func (sender *WXPay) PublicPlaceOrder(body, orderId, notifyUrl, code string, fee int) RetMinProgramPay {
	var params = make(map[string]string)
	var result RetMinProgramPay
	var info RetUnifiedOrder
	var err error
	result.SignType = MD5
	openid, err := sender.PublicGetOpenId(code)
	if err == nil {
		PutParam(params, OPENID, openid)
		info, err = UnifiedOrder(body, orderId, notifyUrl, JSAPIPAY, fee, params, sender.MchId, sender.AppId, sender.ApiSecret)
		if err == nil {
			result.ErrCode, result.ErrMsg = AnalysisWxReturn(info.RetBase, info.RetPublic)
			ClearParam(params)
			nonce := str_lib.Guid()
			timestamp := fmt.Sprint(time.Now().Unix())
			PutParam(params, MIN_APPID, sender.AppId)
			PutParam(params, MIN_TIMESTAMP, timestamp)
			PutParam(params, MIN_NONCESTR, nonce)
			PutParam(params, PACKAGE, PREPAY_ID+info.PrepayId)
			PutParam(params, SIGNTYPE, result.SignType)
			_, sign := PutComplete(params, sender.ApiSecret)

			result.PaySign = sign
			result.AppId = sender.AppId
			result.NonceStr = nonce
			result.TimeStamp = timestamp
			result.Package = PREPAY_ID + info.PrepayId
			result.PrepayId = info.PrepayId
			result.OutTradeNo = orderId
		}
	}
	if err != nil {
		result.ErrCode = FAIL
		result.ErrMsg = err.Error()
	}
	return result
}

//=====================================================================================================================
/*
功能:H5收费下单
参数:
	body:收费显示标题
	order_id:单号
	ip:终端Ip地址
	fee:收费金额
返回:H5支付跳转页面,错误信息
*/
func (sender *WXPay) H5PlaceOrder(body, orderId, notifyUrl, ip string, fee int) (string, error) {
	params := make(map[string]string)
	PutParam(params, SPBILL_CREATE_IP, ip)
	info, err := UnifiedOrder(body, orderId, notifyUrl, H5PAY, fee, params, sender.MchId, sender.AppId, sender.ApiSecret)
	return info.MwebUrl, err
}

/*
功能:获取商家支付二维码
参数:
	body:收费显示标题
	order_id:单号
	ip:终端Ip地址
	fee:收费金额
返回:包含支付二维码的对像,错误信息
*/
func (sender *WXPay) GetPayCode(body, orderId, notifyUrl, ip string, fee int) (ret RetUnifiedOrder, err error) {
	params := make(map[string]string)
	PutParam(params, SPBILL_CREATE_IP, ip)
	ret, err = UnifiedOrder(body, orderId, notifyUrl, PAYCODE, fee, params, sender.MchId, sender.AppId, sender.ApiSecret)
	return
}

/*
功能:app收费下单
参数:
	body:收费显示标题
	order_id:单号
	fee:收费金额
	info:微信返回详细信息
返回:APP调下收银平台所需要的参数(详细参见 RetAppPay 定义)
*/
func (sender *WXPay) AppPlaceOrder(body, orderId, notifyUrl string, fee int) RetAppPay {
	var result RetAppPay
	info, err := UnifiedOrder(body, orderId, notifyUrl, APPPAY, fee, nil, sender.MchId, sender.AppId, sender.ApiSecret)
	if err == nil {
		nonce := str_lib.Guid()
		timestamp := fmt.Sprint(time.Now().Unix())
		params := make(map[string]string)
		PutParam(params, APPID, sender.AppId)
		PutParam(params, PARTNERID, sender.MchId)
		PutParam(params, PREPAYID, info.PrepayId)
		PutParam(params, PACKAGE, result.Package)
		PutParam(params, NONCESTR, nonce)
		PutParam(params, TIMESTAMP, timestamp)
		_, sign := PutComplete(params, sender.ApiSecret)

		result.AppId = sender.AppId
		result.PartnerId = sender.MchId
		result.PrepayId = info.PrepayId
		result.NonceStr = nonce
		result.TimeStamp = timestamp
		result.Sign = sign
		result.OutTradeNo = orderId
	} else {
		result.ErrCode = FAIL
		result.ErrMsg = err.Error()
	}
	return result
}

/*
功能:支付码支付
参数:
	body:收费显示标题
	tradeNo:单号
	notifyUrl:支付状态通知
	authCode:授权码
	fee:收费金额
返回:APP调下收银平台所需要的参数(详细参见 RetMicroPay 定义)
*/
func (sender *WXPay) MicroPay(body, tradeNo, notifyUrl, authCode string, fee int) (ret RetMicroPay, err error) {
	var params = make(map[string]string)
	PutPublic(params, sender.AppId, sender.MchId)
	PutParam(params, OUT_TRADE_NO, tradeNo)
	PutParam(params, BODY, body)
	PutParam(params, TOTAL_FEE, strconv.Itoa(fee))
	PutParam(params, AUTH_CODE, authCode)
	xmlParam, _ := PutComplete(params, sender.ApiSecret) //增加参数完成,返回Post数据和签名
	var buffer []byte
	buffer, _, err = Submit(http_lib.POST, MICROPAY_URL, xmlParam) //post提交数据
	if err == nil {
		err = xml.Unmarshal(buffer, &ret)
	}
	return
}

//==========================================小程序支付=========================================
/*
功能:小程序下单
参数:
	body:收费显示标题
	orderId:单号
	code:获取openId的代码
	fee:收费金额
返回:小程序调起收银台所需要的参数(详细参见 RetMinProgramPay 定义)
*/
func (sender *WXPay) MinProgramPlaceOrder(body, orderId, notifyUrl, code string, fee int) RetMinProgramPay {
	var params = make(map[string]string)
	var result RetMinProgramPay
	result.SignType = MD5
	var mini = mini_program.MiniProgram{MinProgramId: sender.MinProgramId, MinProgramSecret: sender.MinProgramSecret}
	openid, err := mini.GetOpenIdMinProgram(code)
	if err == nil {
		PutParam(params, OPENID, openid)
		fmt.Printf("%#v", params)
		info, err := UnifiedOrder(body, orderId, notifyUrl, JSAPIPAY, fee, params, sender.MchId, sender.MinProgramId, sender.ApiSecret)
		if err == nil {
			result.ErrCode, result.ErrMsg = AnalysisWxReturn(info.RetBase, info.RetPublic)
			ClearParam(params)
			nonce := str_lib.Guid()
			timestamp := fmt.Sprint(time.Now().Unix())
			PutParam(params, MIN_APPID, sender.MinProgramId)
			PutParam(params, MIN_TIMESTAMP, timestamp)
			PutParam(params, MIN_NONCESTR, nonce)

			PutParam(params, PACKAGE, PREPAY_ID+info.PrepayId)
			PutParam(params, SIGNTYPE, result.SignType)
			_, sign := PutComplete(params, sender.ApiSecret)

			result.PaySign = sign
			result.AppId = sender.MinProgramId
			result.NonceStr = nonce
			result.TimeStamp = timestamp
			result.Package = PREPAY_ID + info.PrepayId
			result.PrepayId = info.PrepayId
			result.OutTradeNo = orderId
		}
	}
	if err != nil {
		result.ErrCode = FAIL
		result.ErrMsg = err.Error()
	}
	return result
}

//=================================================================================================================

/*
功能:查询订单是否完成
参数:
	tradeNo:单号
返回:查询结果(详细参见R_Me_Query定义)
*/
func (sender *WXPay) QueryOrder(tradeNo string) (ret RetQuery, err error) {
	var params = make(map[string]string)
	var buffer []byte
	PutPublic(params, sender.AppId, sender.MchId)
	PutParam(params, OUT_TRADE_NO, tradeNo)
	xmlParam, _ := PutComplete(params, sender.ApiSecret) //增加参数完成,返回Post数据和签名
	buffer, _, err = Submit(http_lib.POST, QUERY_ORDER_URL, xmlParam)
	fmt.Printf("qeury:%s", buffer)
	if err == nil {
		xml.Unmarshal(buffer, &ret)
	}
	return

}

/*
功能:退款
参数:
	tradeNo:商户订单单号
	refundTradeNo:退款商户订单号
	totalFee:订单金额
	refundFee:退款金额
返回:查询结果(详细参见R_Me_Query定义)
*/
func (sender *WXPay) Refund(tradeNo, refundTradeNo, notifyUrl string, totalFee, refundFee int,
	certFile, keyFile string) (ret RetRefund, err error) {
	var params = make(map[string]string)
	var resp http_lib.ResponseInfo
	PutPublic(params, sender.AppId, sender.MchId)
	if notifyUrl != EMPTY {
		PutParam(params, NOTIFY_URL, notifyUrl)
	}
	PutParam(params, OUT_TRADE_NO, tradeNo)
	PutParam(params, OUT_REFUND_NO, refundTradeNo)
	PutParam(params, TOTAL_FEE, strconv.Itoa(totalFee))
	PutParam(params, REFUND_FEE, strconv.Itoa(refundFee))

	xmlParam, _ := PutComplete(params, sender.ApiSecret) //增加参数完成,返回Post数据和签名

	resp, err = http_lib.HttpsSecureSubmit(http_lib.POST, REFUND_URL, xmlParam, nil,
		certFile, keyFile)
	if err == nil {
		xml.Unmarshal(resp.BufferBody, &ret)
	}
	return
}

/*
功能:查询退款
参数:
	tradeNo:商户订单单号
	refundTradeNo:退款商户订单号
返回:查询结果(详细参见R_Me_Query定义)
*/
func (sender *WXPay) QueryRefund(outTradeNo, refundTradeNo string) (ret RetQueryRefund, err error) {
	var params = make(map[string]string)
	var buff []byte
	PutPublic(params, sender.AppId, sender.MchId)
	if refundTradeNo == EMPTY {
		PutParam(params, OUT_TRADE_NO, outTradeNo)
	} else {
		PutParam(params, OUT_REFUND_NO, refundTradeNo)
	}
	xmlParam, _ := PutComplete(params, sender.ApiSecret) //增加参数完成,返回Post数据和签名
	buff, _, err = Submit(http_lib.POST, QUERY_REFUND_TRADE, xmlParam)
	if err == nil {
		xml.Unmarshal(buff, &ret)
	}
	return
}

/*
功能:企业微信转账到个人微信
参数:
	amount:退款金额
	tradeNo:订单号
	openId:转账用户的openId
	desc:转账描述
	name:转账者真实姓名(传空时不验证真名姓名)
返回:查询结果(详细参见R_Me_Query定义)
*/
func (sender *WXPay) Cash(amount int, tradeNo, openId, desc, name string) {
	var params = make(map[string]string)
	nonce := str_lib.Guid() //生成随机字符串
	PutParam(params, "mchid", sender.MchId)
	PutParam(params, NONCE_STR, nonce)
	PutParam(params, SPBILL_CREATE_IP, LOCALHOST)
	PutParam(params, "mch_appid", sender.AppId)

	//tradeNo := str_lib.Guid()
	PutParam(params, "partner_trade_no", tradeNo)
	PutParam(params, OPENID, openId)
	if name == EMPTY {
		PutParam(params, "check_name", "NO_CHECK")
	} else {
		PutParam(params, "check_name", "FORCE_CHECK")
		PutParam(params, "re_user_name", name)
	}
	PutParam(params, "amount", strconv.Itoa(amount))
	PutParam(params, "desc", desc)

	xmlParam, _ := PutComplete(params, sender.ApiSecret) //增加参数完成,返回Post数据和签名

	resp, err := http_lib.HttpsSecureSubmit(http_lib.POST, CASH_URL, xmlParam, nil,
		"apiclient_cert.pem", "apiclient_key.pem")
	if err == nil {
		fmt.Println(resp.Body)
	} else {
		fmt.Println("cash error:", err)
	}
}

//撤销订单
func (sender *WXPay) Reverse(outTradeNo string, certFile, keyFile string) (result ReverseResponse, err error) {
	var params = make(map[string]string)
	var resp http_lib.ResponseInfo
	PutPublic(params, sender.AppId, sender.MchId)
	PutParam(params, OUT_TRADE_NO, outTradeNo)
	xmlParam, _ := PutComplete(params, sender.ApiSecret) //增加参数完成,返回Post数据和签名
	if resp, err = http_lib.HttpsSecureSubmit(http_lib.POST, WX_API_REVERSE, xmlParam, nil, certFile, keyFile); err == nil {
		xml.Unmarshal(resp.BufferBody, &result)
	}
	return
}
