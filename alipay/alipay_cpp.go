// alipay_cpp
package alipay

import (
	"errors"
	"fmt"

	"utils/data_conv/json_lib"
	"utils/data_conv/number_lib"
)

/*
功能:生产APP请求字符串(APP支付)
参数:
	body:对一笔交易的具体描述信息
	subject:商品的标题/
	outTradeNo:商户订单号
	notifyUrl:回调地址
	total:交易金额
返回:APP调起收银台所需要的请求字符串
*/
func (c *AliPayLib) AppPay(dealInfo DealBaseInfo) (reqStr string) {
	//公共请求参数
	publicParam := c.putPublic("alipay.trade.app.pay", dealInfo.NotifyUrl)
	//业务参数
	requestParam := make(map[string]interface{})
	c.putParam(requestParam, REQ_TOTAL_AMOUNT, dealInfo.TotalFee)
	c.putParam(requestParam, REQ_OUT_TRADE_NO, dealInfo.TradeNo)
	c.putParam(requestParam, REQ_BODY, dealInfo.Body)
	c.putParam(requestParam, REQ_SUBJECT, dealInfo.Subject)
	bizContent, _ := json_lib.ObjectToJson(requestParam)
	c.putParam(publicParam, BIZ_CONTENT, bizContent)
	reqStr, _ = c.putComplete(publicParam)
	return
}

/*
功能:H5支付(创一个支付页面)
参数:
	body:对一笔交易的具体描述信息
	subject:商品的标题/
	outTradeNo:商户订单号
	notifyUrl:回调地址
	total:交易金额
返回:一个调起收银台页面,错错信息
*/
func (c *AliPayLib) H5Pay(dealInfo DealBaseInfo) (respBody string, err error) {
	//请求参数
	pubParam := c.putPublic("alipay.trade.wap.pay", dealInfo.NotifyUrl)
	//业务参数
	reqParam := make(map[string]interface{})
	c.putParam(reqParam, REQ_BODY, dealInfo.Body)
	c.putParam(reqParam, REQ_SUBJECT, dealInfo.Subject)
	c.putParam(reqParam, REQ_OUT_TRADE_NO, dealInfo.TradeNo)
	c.putParam(reqParam, REQ_TOTAL_AMOUNT, dealInfo.TotalFee)
	c.putParam(reqParam, "product_code", "QUICK_WAP_WAY")
	c.putParam(reqParam, "quit_url", "http://kinot.com")
	//业务参数合并到请求参数中
	bizContent, _ := json_lib.ObjectToJson(reqParam)
	c.putParam(pubParam, BIZ_CONTENT, bizContent)
	//处理请求参,签名
	reqStr, _ := c.putComplete(pubParam)
	//提交请求
	respBody, err = c.requestH5(reqStr)
	//respBody = c.convertToString(respBody, "gbk", "utf-8")
	return
}

/*
功能:生成支付二维码
参数:
	body:对一笔交易的具体描述信息
	subject:商品的标题/
	outTradeNo:商户订单号
	notifyUrl:回调地址
	total:交易金额
返回:RetCreateCode,错误信息
*/
func (c *AliPayLib) CreatePaymentCode(dealInfo DealBaseInfo) (result RetCreateCode, err error) {
	var respStr string
	//请求参数
	pubParam := c.putPublic("alipay.trade.precreate", dealInfo.NotifyUrl)
	//业务参数
	reqParam := make(map[string]interface{})
	c.putParam(reqParam, REQ_BODY, dealInfo.Body)
	c.putParam(reqParam, REQ_SUBJECT, dealInfo.Subject)
	c.putParam(reqParam, REQ_OUT_TRADE_NO, dealInfo.TradeNo)
	c.putParam(reqParam, REQ_TOTAL_AMOUNT, dealInfo.TotalFee)
	//业务参数合并到请求参数中
	bizContent, _ := json_lib.ObjectToJson(reqParam)
	c.putParam(pubParam, BIZ_CONTENT, bizContent)
	//处理请求参,签名
	reqStr, _ := c.putComplete(pubParam)
	//提交请求
	respStr, err = c.getSubmit(reqStr)
	if err == nil {
		json_lib.JsonToObject(respStr, &result)
	}
	return
}

/*
功能:扫码收款
参数:
	body:对一笔交易的具体描述信息
	subject:商品的标题/
	outTradeNo:商户订单号
	authCode:授权码
	notifyUrl:回调地址
	total:交易金额
返回:支付是否成功(不成功时,可查看错误信息,查明原因),错误信息
*/
func (c *AliPayLib) ScanCodePay(dealInfo ScanDealInfo) (result RetMicroPay, err error) {
	var respStr string
	//请求参数
	pubParam := c.putPublic("alipay.trade.pay", dealInfo.NotifyUrl)
	//业务参数
	reqParam := make(map[string]interface{})
	c.putParam(reqParam, REQ_BODY, dealInfo.Body)
	c.putParam(reqParam, REQ_SUBJECT, dealInfo.Subject)
	c.putParam(reqParam, REQ_OUT_TRADE_NO, dealInfo.TradeNo)
	c.putParam(reqParam, REQ_TOTAL_AMOUNT, dealInfo.TotalFee)
	c.putParam(reqParam, REQ_AUTH_CODE, dealInfo.AuthCode)
	c.putParam(reqParam, REQ_SCENE, "bar_code")
	//业务参数合并到请求参数中
	bizContent, _ := json_lib.ObjectToJson(reqParam)
	c.putParam(pubParam, BIZ_CONTENT, bizContent)
	//处理请求参,签名
	reqStr, _ := c.putComplete(pubParam)
	//提交请求
	respStr, err = c.getSubmit(reqStr)
	if err == nil {
		json_lib.JsonToObject(respStr, &result)
		//fmt.Println(respStr)
	}
	return
}

/*
功能:查询订单状态
参数:
	outTradeNo:商户订单号
	tradeNo:支付宝订单号
	notifyUrl:服务器回调地址
注:两个能数只需传一个,另一个传空字符串.两个同时传放时优先out_trade_no
返回:RetQueryTrade,错误信息
*/
func (c *AliPayLib) QueryTrade(outTradeNo, tradeNo string) (result RetQueryTrade, err error) {
	var respStr string
	//请求参数
	pubParam := c.putPublic("alipay.trade.query", EMPTY)
	//业务参数
	reqParam := make(map[string]interface{})
	if outTradeNo == "" {
		c.putParam(reqParam, REQ_TRADE_NO, tradeNo)
	} else {
		c.putParam(reqParam, REQ_OUT_TRADE_NO, outTradeNo)
	}
	//业务参数合并到请求参数中
	bizContent, _ := json_lib.ObjectToJson(reqParam)
	c.putParam(pubParam, BIZ_CONTENT, bizContent)
	//处理请求参,签名
	reqStr, _ := c.putComplete(pubParam)
	//提交请求
	respStr, err = c.getSubmit(reqStr)
	if err == nil {
		json_lib.JsonToObject(respStr, &result)
	}
	return
}

/*
功能:查询订单状态
参数:
	outTradeNo:商户订单号
	tradeNo:支付宝订单号
	notifyUrl:服务器回调地址
注:两个能数只需传一个,另一个传空字符串.两个同时传放时优先out_trade_no
返回:是否成功(不成功时,可查看错误信息,查明原因),错误信息
*/
func (c *AliPayLib) CloseTrade(outTradeNo, tradeNo, notifyUrl string) (result bool, err error) {
	var respStr string
	//请求参数
	pubParams := c.putPublic("alipay.trade.close", notifyUrl)
	//业务参数
	reqParams := make(map[string]interface{})
	if outTradeNo == EMPTY {
		c.putParam(reqParams, REQ_TRADE_NO, tradeNo)
	} else {
		c.putParam(reqParams, REQ_OUT_TRADE_NO, outTradeNo)
	}
	//业务参数合并到请求参数中
	bizContent, _ := json_lib.ObjectToJson(reqParams)
	c.putParam(pubParams, BIZ_CONTENT, bizContent)
	//处理请求参,签名
	reqStr, _ := c.putComplete(pubParams)
	//提交请求
	respStr, err = c.getSubmit(reqStr)

	if err == nil {
		var tmp RetAliPayBase
		json_lib.JsonToObject(respStr, &tmp)
		if tmp.Msg == SUCCESS {
			result = true
		} else {
			result = false
			err = errors.New(fmt.Sprintf("%s , %s", tmp.Msg, tmp.SubMsg))
		}
	}
	return
}

/*
功能:退款
参数:
	outTradeNo:商户订单号
	outRequestNo:退款订单号
	refundAmount:退款金额
返回:退款信息,错误信息
*/
func (c *AliPayLib) Refund(outTradeNo, outRequestNo string, refundAmount float64) (result RetRefundInfo, err error) {
	var respStr string
	//请求参数
	pubParams := c.putPublic("alipay.trade.refund", EMPTY)
	//业务参数
	reqParams := make(map[string]interface{})
	c.putParam(reqParams, REQ_OUT_TRADE_NO, outTradeNo)
	c.putParam(reqParams, REQ_OUT_REQUEST_NO, outRequestNo)
	c.putParam(reqParams, REQ_REFUND_AMOUNT, refundAmount)
	//业务参数合并到请求参数中
	bizContent, _ := json_lib.ObjectToJson(reqParams)
	c.putParam(pubParams, BIZ_CONTENT, bizContent)
	//处理请求参,签名
	reqStr, _ := c.putComplete(pubParams)
	//提交请求
	respStr, err = c.getSubmit(reqStr)
	if err == nil {
		json_lib.JsonToObject(respStr, &result)
	}
	return
}

/*
功能:退款
参数:
	outTradeNo:商户订单号
	outRequestNo:退款订单号
返回:查询的信息,错误信息(注:返回了查询信息,则代表退款成功)
*/
func (c *AliPayLib) QueryRefund(outTradeNo, outRequestNo string) (result RetQueryRefund, err error) {
	var respStr string
	//请求参数
	pubParams := c.putPublic("alipay.trade.fastpay.refund.query", EMPTY)
	//业务参数
	reqParams := make(map[string]interface{})
	c.putParam(reqParams, REQ_OUT_TRADE_NO, outTradeNo)
	c.putParam(reqParams, REQ_OUT_REQUEST_NO, outRequestNo)
	//业务参数合并到请求参数中
	bizContent, _ := json_lib.ObjectToJson(reqParams)
	c.putParam(pubParams, BIZ_CONTENT, bizContent)
	//处理请求参,签名
	reqStr, _ := c.putComplete(pubParams)
	//提交请求
	respStr, err = c.getSubmit(reqStr)
	if err == nil {
		json_lib.JsonToObject(respStr, &result)
	}
	return
}

//解析返回代码
func (c *AliPayLib) AnalysisReturn(base RetAliPayBase) (code int, msg string) {
	code = 0
	msg = "OK"
	if base.Msg != SUCCESS {
		if base.SubMsg == EMPTY {
			number_lib.StrToInt(base.Code, &code)
			msg = base.Msg
		} else {
			number_lib.StrToInt(base.SubCode, &code)
			msg = base.SubMsg
		}
	} else if base.SubMsg != EMPTY {
		number_lib.StrToInt(base.SubCode, &code)
		msg = base.SubMsg

	}
	return
}
