package wechat

import (
	"encoding/xml"
	"fmt"
	"github.com/nanjishidu/gomini/gocrypto"
	"sort"
	"strconv"
	"strings"
	"utils"
	"utils/crypto"
	"utils/data_conv/json_lib"
	"utils/data_conv/number_lib"
	"utils/data_conv/str_lib"
	"utils/http_lib"
)

type RetError struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

/*
交易状态:
SUCCESS—支付成功

REFUND—转入退款

NOTPAY—未支付

CLOSED—已关闭

REVOKED—已撤销（付款码支付）

USERPAYING--用户支付中（付款码支付）

PAYERROR--支付失败(其他原因，如银行返回失败)
*/

//=================================微信支付API返回结构定义=============================

//获取openid时返回的信息
type RetAuthorizeOpenid struct {
	RetError
	AccessToken  string `json:"accessToken"`   //网页授权接口调用凭证,注意：此access_token与基础支持的access_token不同
	RefreshToken string `json:"refresh_token"` //用于刷新AccessToken
	OpenId       string `json:"openid"`        //openid
	Scope        string `json:"scope"`         //获取更多的微信信息的密钥
	ExpiresIn    int    `json:"expiresIn"`     //到期时间
}

//微信平台基本返回内容
type RetBase struct {
	ReturnCode string `xml:"return_code" json:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg" json:"return_msg"`   //返回状态描述
}

//微信平台公共返回
type RetPublic struct {
	AppId      string `xml:"appid" json:"appid"`               //应用APPID
	MchId      string `xml:"mch_id" json:"mch_id"`             //商户号
	NonceStr   string `xml:"nonce_str" json:"nonce_str"`       //随机字符串
	Sign       string `xml:"sign" json:"sign"`                 //签名
	ResultCode string `xml:"result_code" json:"result_code"`   //业务结果码
	ErrCode    string `xml:"err_code" json:"err_code"`         //业务错误代码
	ErrCodeDes string `xml:"err_code_des" json:"err_code_des"` //业务错误代码描述
	DeviceInfo string `xml:"device_info"json:"device_info"`    //自定义参数
}

//下单,微信平台返回
type RetUnifiedOrder struct {
	RetBase
	RetPublic
	TradeType string `xml:"trade_type" json:"trade_type"` //交易类型
	PrepayId  string `xml:"prepay_id" json:"prepay_id"`   //预支付交易会话标识,用于后续接口调用中使用
	CodeUrl   string `xml:"code_url" json:"code_url"`     //二维码连接
	MwebUrl   string `xml:"mweb_url" json:"mweb_url"`     //H5支付跳转码
}

//订单查询返回
type RetQuery struct {
	RetBase
	RetPublic
	Openid         string `xml:"openid" json:"openid"`                     //用户标识
	TradeType      string `xml:"trade_type" json:"trade_type"`             //交易类型
	TradeStatus    string `xml:"trade_state" json:"trade_state"`           //交易状态
	BankType       string `xml:"bank_type" json:"bank_type"`               //付款银行
	TransactionId  string `xml:"transaction_id" json:"transaction_id"`     //微信订单号
	OutTradeNo     string `xml:"out_trade_no" json:"out_trade_no"`         //商户订单号
	TimeEnd        string `xml:"time_end" json:"time_end"`                 //支付完成时间
	TotalFee       int    `xml:"total_fee" json:"total_fee"`               //标价金额
	CashFee        int    `xml:"cash_fee" json:"cash_fee"`                 //现金支付金额
	TradeStateDesc string `xml:"trade_state_desc" json:"trade_state_desc"` //对当前查询订单状态的描述和下一步操作的指引
	IsSubscribe    string `xml:"is_subscribe" json:"is_subscribe"`         //是否关注公众账号
	FeeType        string `xml:"fee_type" json:"fee_type"`                 //货币种类
}

//通过小程序传入的code获取openid
type RetMinProgramOpenId struct {
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
}

//支付码支付返回
type RetMicroPay struct {
	RetBase
	RetPublic
	DeviceInfo         string `xml:"device_info" json:"device_info"`                   //设备号
	OpenId             string `xml:"openid" json:"openid"`                             //用户标识
	IsSubscribe        string `xml:"is_subscribe" json:"is_subscribe"`                 //是否关注公众账号
	TradeType          string `xml:"trade_type" json:"trade_type"`                     //交易类型
	BankType           string `xml:"bank_type" json:"bank_type"`                       //付款银行
	TotalFee           int    `xml:"total_fee" json:"total_fee"`                       //订单金额
	SettlementTotalFee int    `xml:"settlement_total_fee" json:"settlement_total_fee"` //应结订单金额
	FeeType            string `xml:"fee_type" json:"fee_type"`                         //货币种类
	CashFee            int    `xml:"cash_fee" json:"cash_fee"`                         //现金支付金额
	CashFeeType        string `xml:"cash_fee_type" json:"cash_fee_type"`               //现金支付货币类型
	CouponFee          int    `xml:"coupon_fee" json:"coupon_fee"`                     //总代金券金额
	CouponCount        int    `xml:"coupon_count" json:"coupon_count"`                 //代金券使用数量
	TransactionId      string `xml:"transaction_id" json:"transaction_id"`             //微信支付订单号
	OutTradeNo         string `xml:"out_trade_no" json:"out_trade_no"`                 //商户订单号
	Attach             string `xml:"attach" json:"attach"`                             //商家数据包
	TimeEnd            string `xml:"time_end" json:"time_end"`                         //支付完成时间
}

//退款返回
type RetRefund struct {
	RetBase
	RetPublic
	TransactionId string `xml:"transaction_id" json:"transaction_id"` //微信订单号
	OutTradeNo    string `xml:"out_trade_no" json:"out_trade_no"`     //商户订单号
	OutRefundNo   string `xml:"out_refund_no" json:"out_refund_no"`   //商户退款单号
	RefundId      string `xml:"refund_id" json:"refund_id"`           //微信退款单号
	TotalFee      int    `xml:"total_fee" json:"total_fee"`           //标价金额
	RefundFee     int    `xml:"refund_fee" json:"refund_fee"`         //退款金额
	CashFee       int    `xml:"cash_fee" json:"cash_fee"`             //现金支付金额
}

type ReverseResponse struct {
	RetBase
	RetPublic
	ReCall string `xml:"re_call" json:"re_call"` //是否需要继续调用撤销，Y-需要，N-不需要
}

//支付结果异步通知信息
type PaymentNotifyInfo struct {
	RetBase
	RetPublic
	DeviceInfo         string `xml:"device_info" json:"device_info"`                   //设备号
	SignType           string `xml:"sign_type" json:"sign_type"`                       //签名算法
	OpenId             string `xml:"openid" json:"openid"`                             //用户标识
	IsSubscribe        string `xml:"is_subscribe" json:"is_subscribe"`                 //是否关注公众账号
	TradeType          string `xml:"trade_type" json:"trade_type"`                     //交易类型
	BankType           string `xml:"bank_type" json:"bank_type"`                       //付款银行
	TotalFee           int    `xml:"total_fee" json:"total_fee"`                       //订单金额
	SettlementTotalFee int    `xml:"settlement_total_fee" json:"settlement_total_fee"` //应结订单金额
	FeeType            string `xml:"fee_type" json:"fee_type"`                         //货币种类
	CashFee            int    `xml:"cash_fee" json:"cash_fee"`                         //现金支付金额
	CashFeeType        string `xml:"cash_fee_type" json:"cash_fee_type"`               //现金支付货币类型
	CouponFee          int    `xml:"coupon_fee" json:"coupon_fee"`                     //总代金券金额
	CouponCount        int    `xml:"coupon_count" json:"coupon_count"`                 //代金券使用数量
	TransactionId      string `xml:"transaction_id" json:"transaction_id"`             //微信支付订单号
	OutTradeNo         string `xml:"out_trade_no" json:"out_trade_no"`                 //商户订单号
	Attach             string `xml:"attach" json:"attach"`                             //商家数据包
	TimeEnd            string `xml:"time_end" json:"time_end"`                         //支付完成时间
	TradeStateDesc     string `xml:"trade_state_desc" json:"trade_state_desc"`
	TradeState         string `xml:"trade_state" json:"trade_state"`
}

//退款异步通知消息
type RefundNotifyInfo struct {
	ReturnCode string `xml:"return_code" json:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg" json:"return_msg"`   //返回信息
	AppId      string `xml:"appid" json:"appid"`             //公众账号ID
	MchId      string `xml:"mch_id" json:"mch_id"`           //商户号
	NonceStr   string `xml:"nonce_str" json:"nonce_str"`     //随机字符串
	ReqInfo    string `xml:"req_info" json:"req_info"`       //加密信息
	RefundEncryptInfo
}

//退款通知消息加密信息
type RefundEncryptInfo struct {
	TransactionId       string `xml:"transaction_id" json:"transaction_id"`               //微信支付订单号
	OutTradeNo          string `xml:"out_trade_no" json:"out_trade_no"`                   //商户订单号
	RefundId            string `xml:"refund_id" json:"refund_id"`                         //微信退款单号
	OutRefundNo         string `xml:"out_refund_no" json:"out_refund_no"`                 //商户退款订单号
	TotalFee            int    `xml:"total_fee" json:"total_fee"`                         //订单金额
	SettlementTotalFee  int    `xml:"settlement_total_fee" json:"settlement_total_fee"`   //应结订单金额(当该订单有使用非充值券时，返回此字段。应结订单金额=订单金额-非充值代金券金额，应结订单金额<=订单金额)
	RefundFee           int    `xml:"refund_fee" json:"refund_fee"`                       //申请退款金额
	SettlementRefundFee int    `xml:"settlement_refund_fee" json:"settlement_refund_fee"` //退款金额
	RefundStatus        string `xml:"refund_status" json:"refund_status"`                 //退款状态(SUCCESS-退款成功,CHANGE-退款异常,REFUNDCLOSE—退款关闭)
	SuccessTime         string `xml:"success_time" json:"success_time"`                   //退款成功时间
	RefundRecvAccout    string `xml:"refund_recv_accout" json:"refund_recv_accout"`       //退款入账账户
	RefundAccount       string `xml:"refund_account"`                                     //退款资金来源(REFUND_SOURCE_RECHARGE_FUNDS-可用余额退款/基本账户,REFUND_SOURCE_UNSETTLED_FUNDS-未结算资金退款)
	RefundRequestSource string `xml:"refund_request_source" json:"refund_request_source"` //退款发起来源(API-接口,VENDOR_PLATFORM-商户平台)
}

type RetQueryRefund struct {
	RetBase
	RetPublic
	TransactionId string `xml:"transaction_id" json:"transaction_id"` //微信支付订单号
	OutTradeNo    string `xml:"out_trade_no" json:"out_trade_no"`     //商户订单号
	RefundId      string `xml:"refund_id" json:"refund_id"`           //微信退款单号
	OutRefundNo   string `xml:"out_refund_no" json:"out_refund_no"`   //商户退款订单号
	TotalFee      int    `xml:"total_fee" json:"total_fee"`           //订单金额
}

//======================================================================================================================

//==========================================移动端支付所需参数定义========================================================
type RetMeBase struct {
	ErrCode int    `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

/*小程序下单需要参数返回*/
type RetMinProgramPay struct {
	RetMeBase
	AppId      string `json:"appId"`     //小程序ID
	TimeStamp  string `json:"timeStamp"` //时间截
	NonceStr   string `json:"nonceStr"`  //随机字符串
	Package    string `json:"package"`   //数据包,统一下单接口返回的 prepay_id 参数值
	SignType   string `json:"singType"`  //签名方式
	PaySign    string `json:"paySign"`   //签名
	PrepayId   string `json:"prepayId"`
	OutTradeNo string `json:"out_trade_no"`
}

/*APP下单需要参数返回*/
type RetAppPay struct {
	RetMeBase
	AppId      string `json:"appid"`        //小程序ID
	TimeStamp  string `json:"timestamp"`    //时间截
	NonceStr   string `json:"noncestr"`     //随机字符串
	PartnerId  string `json:"partnerid"`    //商户号
	PrepayId   string `json:"prepayid"`     //预支付交易会话ID
	Package    string `json:"package"`      //暂填写固定值Sign=WXPay
	Sign       string `json:"sign"`         //签名
	OutTradeNo string `json:"out_trade_no"` //商户订阅号
}

//===================================================END================================================================

//===================================================小程序=============================================================
//小程序登录返回(获取access token)
type RetMinProgramLogin struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

//=====================================================================================================================

//*****************************************************常量*************************************************************

//字符常量
const (
	SUCCESS      = "SUCCESS"
	CONTENT_TYPE = "application/x-www-form-urlencoded"
	LOCALHOST    = "127.0.0.1"
	EMPTY        = ""
)

//HTTP提交方式
//const (
//	POST = "POST"
//	GET  = "GET"
//)

//数字常量
const (
	FAIL = -1 //失败
	OK   = 0  //完功
	MD5  = "MD5"
)

//微信开放平台API URL接口
const (
	AUTHORIZE_URL         = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code" //获取OPENID
	GET_PAY_CODE_URL      = "https://api.mch.weixin.qq.com/pay/unifiedorder"                                                             //统一下单,返回正确的预支付交易后调起支付
	QUERY_ORDER_URL       = "https://api.mch.weixin.qq.com/pay/orderquery"                                                               //订单状态查询口URL
	MICROPAY_URL          = "https://api.mch.weixin.qq.com/pay/micropay"                                                                 //提交刷卡支付URL
	JSOCDE_TO_SESSION_URL = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"   //小程序程传CODE，获取OPENID
	REFUND_URL            = "https://api.mch.weixin.qq.com/secapi/pay/refund"                                                            //退款接口
	CASH_URL              = "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers"                                        //企业微信转账到个人接口
	GET_ACCESS_TOKEN      = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"                    //获取小程序全局唯一后台接口调用凭据
	SEND_MINI_MESSAGE     = "https://api.weixin.qq.com/cgi-bin/message/wxopen/template/uniform_send?access_token=%s"                     //小程序统一服务消息
	QUERY_REFUND_TRADE    = "https://api.mch.weixin.qq.com/pay/refundquery"                                                              //退款订单查询
	WX_API_REVERSE        = "https://api.mch.weixin.qq.com/secapi/pay/reverse"                                                           //撤销订单
)

//支付交易状态
const (
	PAY_STATUS_SUCCESS    = "SUCCESS"               //支付成功 0
	PAY_STATUS_FEFUND     = "PAY_STATUS_FEFUND"     //转入退示款 1
	PAY_STATUS_NOTPAY     = "PAY_STATUS_NOTPAY"     //未支付 2款
	PAY_STATUS_CLOSED     = "PAY_STATUS_CLOSED"     //已关闭 3
	PAY_STATUS_REVOKED    = "PAY_STATUS_REVOKED"    //已撤销(刷卡支付) 4
	PAY_STATUS_USERPAYING = "PAY_STATUS_USERPAYING" //用户支付中 5
	PAY_STATUS_PAYERROR   = "PAY_STATUS_PAYERROR"   //支付失败 6
)

//微信API参数名
const (
	MCH_ID           = "mch_id"           //商户号
	NONCE_STR        = "nonce_str"        //随机字符串
	NONCESTR         = "noncestr"         //随机字符串
	MIN_NONCESTR     = "nonceStr"         //小程序随机字符串
	SPBILL_CREATE_IP = "spbill_create_ip" //终端IP
	APPID            = "appid"            //微信支付appid
	MIN_APPID        = "appId"            //小程序APPID
	BODY             = "body"             //商品描述
	OUT_TRADE_NO     = "out_trade_no"     //商户订单号
	TOTAL_FEE        = "total_fee"        //标价总金额
	NOTIFY_URL       = "notify_url"       //通知地址
	TRADE_TYPE       = "trade_type"       //交易类型
	OPENID           = "openid"           //用户标识
	TIMESTAMP        = "timestamp"        //时间截
	MIN_TIMESTAMP    = "timeStamp"        //小程序时间截
	PACKAGE          = "package"          //数据包
	SIGNTYPE         = "signType"         //签名方式
	PARTNERID        = "partnerid"        //APP支付的商户号名称
	PREPAYID         = "prepayid"         //预支付交易会话ID
	PREPAY_ID        = "prepay_id="       //
	OUT_REFUND_NO    = "out_refund_no"    //退款订单号
	REFUND_FEE       = "refund_fee"       //退款金额
	AUTH_CODE        = "auth_code"        //授权码
	MICROPAY         = "MICROPAY"
)

//支付方式
const (
	H5PAY    = "MWEB"   //H5支付
	APPPAY   = "APP"    //APP支付
	JSAPIPAY = "JSAPI"  //小程序支付
	PAYCODE  = "NATIVE" //商家支付二维码
)

//**********************************************************************************************************************

//============================================小程序统一服务消息模板====================================================
//小程序统一服务消息
type MiniProgramMessage struct {
	ToUser           string            `json:"touser"`
	WeAppTemplateMsg WeAppMessage      `json:"weapp_template_msg"`
	MpTemplateMsg    MpTemplateMessage `json:"mp_template_msg"`
}

//小程序统一服务消息->小程序消息模板
type WeAppMessage struct {
	TemplateId      string      `json:"template_id"`      //小程序模板ID
	Page            string      `json:"page"`             //小程序页面路径
	FormId          string      `json:"form_id"`          //小程序模板消息formid(订单号或form提交的form_id)
	Data            interface{} `json:"data"`             //小程序模板数据
	EmphasisKeyword string      `json:"emphasis_keyword"` //小程序模板放大关键词
}

//小程序预约模板
type OrderMessage struct {
	Keyword1 MiniDataModule `json:"keyword1"`
	Keyword2 MiniDataModule `json:"keyword2"`
	Keyword3 MiniDataModule `json:"keyword3"`
}

//小程序数据值
type MiniDataModule struct {
	Value string `json:"value"`
}

//小程序统一服务消息->公众号消息模板
type MpTemplateMessage struct {
	AppId       string      `json:"appid"`
	TemplateId  string      `json:"template_id"`
	Url         string      `json:"url"`
	MiniProgram MiniPro     `json:"miniprogram"`
	Data        interface{} `json:"data"`
}

//小程序信息
type MiniPro struct {
	AppId    string `json:"appid"`
	PagePath string `json:"page"`
}

//公众号消息模板数据值
type PublicDataModule struct {
	Value string `json:"value"` //值
	Color string `json:"color"` //字体颜色
}

//=====================================================================================================================

//=====================================================公共接口=========================================================

/*
功能:微信统一下单
参数:
	body:订单描述
	orderId:商户订单号
	notifyUrl:支付状态回调URL
	tradeType:订单类型
	fee:订单价格
	replenish:补充参数
	mchId:商户号
	appId:应用ID
	apiSecret:商户api密钥
返回:统一下单返回,异常
*/
func UnifiedOrder(body, orderId, notifyUrl, tradeType string, fee int, replenish map[string]string,
	mchId, appId, apiSecret string) (result RetUnifiedOrder, err error) {
	params := make(map[string]string)
	PutPublic(params, appId, mchId)
	PutParam(params, BODY, body)
	PutParam(params, OUT_TRADE_NO, orderId)
	PutParam(params, TOTAL_FEE, strconv.Itoa(fee))
	PutParam(params, NOTIFY_URL, notifyUrl)
	PutParam(params, TRADE_TYPE, tradeType)
	if replenish != nil {
		for n, v := range replenish {
			PutParam(params, n, v)
		}
	}
	xmlParam, _ := PutComplete(params, apiSecret) //增加参数完成,返回Post数据和签名
	fmt.Println(xmlParam)
	buffer, _, err := Submit(http_lib.POST, GET_PAY_CODE_URL, xmlParam) //post提交数据
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
func AnalysisWxReturn(info1 RetBase, info2 RetPublic) (int, string) {
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
func AnalysisPayStatus(state string) int {
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
func Submit(submitType, url, data string) ([]byte, int, error) {
	resp, err := http_lib.HttpSubmit(submitType, url, data, nil)
	return resp.BufferBody, resp.HttpStatus, err
}

/*
功能:填充公共参数
参数:
	params:参数集合
	appId:
	mchId:商户号
返回:无
*/
func PutPublic(params map[string]string, appId, mchId string) {
	nonce := str_lib.Guid() //生成随机字符串
	PutParam(params, MCH_ID, mchId)
	PutParam(params, NONCE_STR, nonce)
	//PutParam(params, SPBILL_CREATE_IP, LOCALHOST)
	PutParam(params, APPID, appId)
}

/*
功能:填充单个参数
参数:
	params:参数集合
	param_name:参数名称
	param_value:参数集合
返回:无
*/
func PutParam(params map[string]string, paramName string, paramValue string) {
	params[paramName] = paramValue
}

/*
功能:整理提交的数据并计算签名
参数:
	params:参数集合
	apiSecret:
返回:整理完成的提交数据,签名
*/
func PutComplete(params map[string]string, apiSecret string) (string, string) {
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
	sign += "key=" + apiSecret
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
func ClearParam(params map[string]string) {
	for k := range params {
		delete(params, k)
	}
}

//验签
func VerifySign(v interface{}, sign, apiSecret string) (ret bool) {
	mapData := make(map[string]interface{})
	err := json_lib.ObjectToObject(&mapData, v)
	var stringA = EMPTY
	var signA = EMPTY
	if err == nil {
		var keys []string
		//排序
		for k, v := range mapData {
			if v != EMPTY && k != "sign" {
				keys = append(keys, k)
			}
		}
		sort.Strings(keys)
		for _, n := range keys {
			switch mapData[n].(type) {
			case int, float64, float32:
				x := int(mapData[n].(float64))
				if x > 0 {
					s := number_lib.NumberToStr(x)
					stringA += n + "=" + s + "&"
				}
			case string:
				stringA += n + "=" + mapData[n].(string) + "&"
			default:
				stringA += n + "=" + fmt.Sprintf("%s", mapData[n]) + "&"
			}
		}
		stringA += "key=" + apiSecret
		fmt.Println(stringA)
		signA = crypto.Md5(stringA) //MD5加密
		signA = strings.ToUpper(signA)
		fmt.Println(signA, sign)
		ret = utils.If(signA == sign, true, false).(bool)
	}
	return
}

//解密退款返回的加密数据
func DecodeRefundData(encryptData, apiSecret string) (decodeStr []byte, err error) {
	var b []byte
	b, err = crypto.DecodeBase64(encryptData)
	if err == nil {
		gocrypto.SetAesKey(strings.ToLower(crypto.Md5(apiSecret)))
		decodeStr, err = gocrypto.AesECBDecrypt(b)
	}
	return
}
