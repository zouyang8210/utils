// weixin
package weixin

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"time"
	"utils/data_conv/json_lib"
	"utils/data_conv/str_lib"
	"utils/http_lib"
)

//************************************接口函数***********************************

//=============================================公众号支付==============================================================
/*
功能:公众号获取openId
参数:
	code:获取openid的验证码
返回
*/
func (sender *WeiXinPay) PublicGetOpenId(code string) (string, error) {
	var info RetAuthorizeOpenid
	strUrl := fmt.Sprintf(AUTHORIZE_URL, sender.AppId, sender.AppSecret, code)

	strBody, httpStatus, err := sender.submit(http_lib.GET, strUrl, EMPTY)
	if err == nil {
		if httpStatus != 200 {
			err = errors.New(fmt.Sprintf("http status %d", httpStatus))
		}
		json_lib.JsonToObject(string(strBody), &info)
		if info.ErrCode != 0 {
			err = errors.New(fmt.Sprintf("errCode=%d,errMsg=%s", info.ErrCode, info.ErrMsg))
		}
	}
	return info.OpenId, err
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
func (sender *WeiXinPay) PublicPlaceOrder(body, orderId, notifyUrl, code string, fee int) RetMinProgramPay {
	var params = make(map[string]string)
	var result RetMinProgramPay
	var info RetUnifiedOrder
	var err error
	result.SignType = "MD5"
	openid, err := sender.PublicGetOpenId(code)
	if err == nil {
		sender.putParam(params, OPENID, openid)
		info, err = sender.unifiedOrder(body, orderId, notifyUrl, JSAPIPAY, fee, params)
		if err == nil {
			result.ErrCode, result.ErrMsg = sender.analysisWxReturn(info.RetBase, info.RetPublic)
			sender.clearParam(params)
			nonce := str_lib.Guid()
			timestamp := fmt.Sprint(time.Now().Unix())
			sender.putParam(params, MIN_APPID, sender.AppId)
			sender.putParam(params, MIN_TIMESTAMP, timestamp)
			sender.putParam(params, MIN_NONCESTR, nonce)
			sender.putParam(params, PACKAGE, PREPAY_ID+info.PrepayId)
			sender.putParam(params, SIGNTYPE, result.SignType)
			_, sign := sender.putComplete(params)

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
	fee:收费金额
	info:微信返回详细信息
返回:H5支付跳转页面,错误信息
*/
func (sender *WeiXinPay) H5PlaceOrder(body, orderId, notifyUrl, ip string, fee int) (string, error) {
	params := make(map[string]string)
	sender.putParam(params, SPBILL_CREATE_IP, ip)
	info, err := sender.unifiedOrder(body, orderId, notifyUrl, H5PAY, fee, params)
	return info.MwebUrl, err
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
func (sender *WeiXinPay) AppPlaceOrder(body, orderId, notifyUrl string, fee int) RetAppPay {
	var result RetAppPay

	info, err := sender.unifiedOrder(body, orderId, notifyUrl, APPPAY, fee, nil)
	if err == nil {
		nonce := str_lib.Guid()
		timestamp := fmt.Sprint(time.Now().Unix())
		params := make(map[string]string)
		sender.putParam(params, APPID, sender.AppId)
		sender.putParam(params, PARTNERID, sender.MchId)
		sender.putParam(params, PREPAYID, info.PrepayId)
		sender.putParam(params, PACKAGE, result.Package)
		sender.putParam(params, NONCESTR, nonce)
		sender.putParam(params, TIMESTAMP, timestamp)
		_, sign := sender.putComplete(params)

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

//========================================小程序支付=========================================
/*
功能：通过小程序传入的code，调用API获取openid
参数：
	code:小程序生成的临时密钥
返回:openid,错误信息
*/
func (sender *WeiXinPay) GetOpenIdMinProgram(code string) (string, error) {
	var info RetMinProgramOpenId
	var err error
	var bb []byte

	url := fmt.Sprintf(JSOCDE_TO_SESSION_URL, sender.MinProgramId, sender.MinProgramSecret, code)
	//res, err = http.Get(url)
	bb, _, err = sender.submit(http_lib.GET, url, EMPTY)
	if err == nil {
		err = json_lib.JsonToObject(string(bb), &info)
	}
	return info.OpenId, err
}

/*
功能：小程序下单
参数:
	body:收费显示标题
	orderId:单号
	code:获取openId的代码
	fee:收费金额
返回:小程序调起收银台所需要的参数(详细参见 RetMinProgramPay 定义)
*/
func (sender *WeiXinPay) MinProgramPlaceOrder(body, orderId, notifyUrl, code string, fee int) RetMinProgramPay {
	var params = make(map[string]string)
	var result RetMinProgramPay
	result.SignType = "MD5"
	openid, err := sender.GetOpenIdMinProgram(code)
	if err == nil {
		sender.putParam(params, OPENID, openid)
		info, err := sender.unifiedOrder(body, orderId, notifyUrl, JSAPIPAY, fee, params)
		if err == nil {
			result.ErrCode, result.ErrMsg = sender.analysisWxReturn(info.RetBase, info.RetPublic)
			sender.clearParam(params)
			nonce := str_lib.Guid()
			timestamp := fmt.Sprint(time.Now().Unix())
			sender.putParam(params, MIN_APPID, sender.MinProgramId)
			sender.putParam(params, MIN_TIMESTAMP, timestamp)
			sender.putParam(params, MIN_NONCESTR, nonce)
			sender.putParam(params, PACKAGE, PREPAY_ID+info.PrepayId)
			sender.putParam(params, SIGNTYPE, result.SignType)
			_, sign := sender.putComplete(params)

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
	order_id:单号
	info:微信返回详细信息
返回:查询结果(详细参见R_Me_Query定义)
*/
func (sender *WeiXinPay) QueryOrder(orderId string) RetMeQuery {
	var result RetMeQuery
	var params = make(map[string]string)
	var err error
	var buffer []byte
	sender.putPublic(params)
	sender.putParam(params, OUT_TRADE_NO, orderId)

	xmlParam, _ := sender.putComplete(params) //增加参数完成,返回Post数据和签名

	buffer, _, err = sender.submit(http_lib.POST, QUERY_ORDER_URL, xmlParam)
	if err == nil {
		var info RetQuery
		err = xml.Unmarshal(buffer, &info)
		if err == nil {
			result.ErrCode, result.ErrMsg = sender.analysisWxReturn(info.RetBase, info.RetPublic)
			result.TradeStatus = info.TradeStatus
			result.TradeType = info.TradeType
			result.BankType = info.BankType
			result.CashFee = info.CashFee
			result.Openid = info.Openid
			result.TotalFee = info.TotalFee
			result.OutTradeNo = info.OutTradeNo
			result.TransactionId = info.TransactionId
			result.TimeEnd = info.TimeEnd
			result.TradeStatusDesc = info.TradeStatusDesc
		}
	}
	if err != nil {
		result.ErrCode = FAIL
		result.ErrMsg = err.Error()
	}
	return result

}

/*
功能:退款
参数:
	orderId:单号
	refundTradeNo:退款订单号
	totalFee:订单金额
	refundFee:退款金额
返回:查询结果(详细参见R_Me_Query定义)
*/
func (sender *WeiXinPay) Refund(orderId, refundTradeNo string, totalFee, refundFee int, notifyUrl string) {
	var params = make(map[string]string)
	sender.putPublic(params)
	if notifyUrl != EMPTY {
		sender.putParam(params, NOTIFY_URL, notifyUrl)
	}
	sender.putParam(params, OUT_TRADE_NO, orderId)
	sender.putParam(params, OUT_REFUND_NO, refundTradeNo)
	sender.putParam(params, TOTAL_FEE, strconv.Itoa(totalFee))
	sender.putParam(params, REFUND_FEE, strconv.Itoa(refundFee))

	xmlParam, _ := sender.putComplete(params) //增加参数完成,返回Post数据和签名

	resp, err := http_lib.HttpsSecureSubmit(http_lib.POST, REFUND_URL, xmlParam, nil,
		"apiclient_cert.pem", "apiclient_key.pem")
	if err == nil {
		fmt.Println(resp.Body)
	} else {
		fmt.Println("refund error:", err)
	}

}

func (sender *WeiXinPay) Cash(amount int, tradeNo, openId, desc, name string) {
	var params = make(map[string]string)

	nonce := str_lib.Guid() //生成随机字符串
	sender.putParam(params, "mchid", sender.MchId)
	sender.putParam(params, NONCE_STR, nonce)
	sender.putParam(params, SPBILL_CREATE_IP, LOCALHOST)
	sender.putParam(params, "mch_appid", sender.AppId)

	//tradeNo := str_lib.Guid()
	sender.putParam(params, "partner_trade_no", tradeNo)
	sender.putParam(params, OPENID, openId)
	if name == EMPTY {
		sender.putParam(params, "check_name", "NO_CHECK")
	} else {
		sender.putParam(params, "check_name", "FORCE_CHECK")
		sender.putParam(params, "re_user_name", name)
	}
	sender.putParam(params, "amount", strconv.Itoa(amount))
	sender.putParam(params, "desc", desc)

	xmlParam, _ := sender.putComplete(params) //增加参数完成,返回Post数据和签名

	resp, err := http_lib.HttpsSecureSubmit(http_lib.POST, CASH_URL, xmlParam, nil,
		"apiclient_cert.pem", "apiclient_key.pem")
	if err == nil {
		fmt.Println(resp.Body)
	} else {
		fmt.Println("cash error:", err)
	}
}

//***************************************************************************************************************
