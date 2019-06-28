// weixin_custom
package weixin

//微信支付对像
type WeiXinPay struct {
	AppId            string `json:"app_id"`     //appid
	MchId            string `json:"mch_id"`     //商户号
	AppSecret        string `json:"app_secret"` //app密钥
	ApiSecret        string `json:"api_secret"` //api密钥
	MinProgramId     string `json:"_"`          //小程序ID
	MinProgramSecret string `json:"_"`          //小程序密钥
}

//=================================微信API返回结构定义=============================
//微信平台基本返回内容
type RetBase struct {
	ReturnCode string `xml:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg"`  //返回状态描述
}

//微信平台公共返回
type RetPublic struct {
	AppId      string `xml:"appid"`        //应用APPID
	MchId      string `xml:"mch_id"`       //商户号
	NonceStr   string `xml:"nonce_str"`    //随机字符串
	Sign       string `xml:"sign"`         //签名
	ResultCode string `xml:"result_code"`  //业务结果码
	ErrCode    string `xml:"err_code"`     //业务错误代码
	ErrCodeDes string `xml:"err_code_des"` //业务错误代码描述
	DeviceInfo string //自定义参数
}

//下单,微信平台返回
type RetUnifiedOrder struct {
	RetBase
	RetPublic
	TradeType string `xml:"trade_type"` //交易类型
	PrepayId  string `xml:"prepay_id"`  //预支付交易会话标识,用于后续接口调用中使用
	CodeUrl   string `xml:"code_url"`   //二维码连接
	MwebUrl   string `xml:"mweb_url"`   //H5支付跳转码
}

//订单查询返回
type RetQuery struct {
	RetBase
	RetPublic
	Openid          string `xml:"openid"`           //用户标识
	TradeType       string `xml:"trade_type"`       //交易类型
	TradeStatus     string `xml:"trade_state"`      //交易状态
	BankType        string `xml:"bank_type"`        //付款银行
	TransactionId   string `xml:"transaction_id"`   //微信订单号
	OutTradeNo      string `xml:"out_trade_no"`     //商户订单号
	TimeEnd         string `xml:"time_end"`         //支付完成时间
	TradeStatusDesc string `xml:"trade_state_desc"` //对当前查询订单状态的描述和下一步操作的指引
	TotalFee        int    `xml:"total_fee"`        //标价金额
	CashFee         int    `xml:"cash_fee"`         //现金支付金额
}

//通过小程序传入的code获取openid
type RetMinProgramOpenId struct {
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
}

//========================================END===================================

//====================================返回给支付端结构定义==========================
//

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

/*APP下单需要参数返回*/
type RetH5Pay struct {
	RetMeBase
	AppId     string `json:"appId"`     //小程序ID
	TimeStamp string `json:"timeStamp"` //时间截
	NonceStr  string `json:"nonceStr"`  //随机字符串
	Package   string `json:"package"`   //数据包,统一下单接口返回的 prepay_id 参数值
	SignType  string `json:"singType"`  //签名方式
	PaySign   string `json:"paySign"`   //签名
}

//订单查询返回
type RetMeQuery struct {
	RetMeBase
	Openid          string `json:"openid"`           //用户标识
	TradeType       string `json:"trade_type"`       //交易类型
	TradeStatus     string `json:"trade_state"`      //交易状态
	BankType        string `json:"bank_type"`        //付款银行
	TransactionId   string `json:"transaction_id"`   //微信订单号
	OutTradeNo      string `json:"out_trade_no"`     //商户订单号
	TimeEnd         string `json:"time_end"`         //支付完成时间
	TradeStatusDesc string `json:"trade_state_desc"` //对当前查询订单状态的描述和下一步操作的指引
	TotalFee        int    `json:"total_fee"`        //标价金额
	CashFee         int    `json:"cash_fee"`         //现金支付金额
}

//========================================END===================================

//**************************************常量************************************

//字符常量
const (
	SUCCESS      = "SUCCESS"
	CONTENT_TYPE = "application/x-www-form-urlencoded"
	LOCALHOST    = "127.0.0.1"
	EMPTY        = ""
)

//HTTP提交方式
const (
	POST = "POST"
	GET  = "GET"
)

//数字常量
const (
	FAIL = -1 //失败
	OK   = 0  //完功
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
)

//支付状态
const (
	PAY_STATUS_SUCCESS    = "SUCCESS"               //支付成功
	PAY_STATUS_FEFUND     = "PAY_STATUS_FEFUND"     //转入退示款
	PAY_STATUS_NOTPAY     = "PAY_STATUS_NOTPAY"     //未支付
	PAY_STATUS_CLOSED     = "PAY_STATUS_CLOSED"     //已关闭
	PAY_STATUS_REVOKED    = "PAY_STATUS_REVOKED"    //已撤销(刷卡支付)
	PAY_STATUS_USERPAYING = "PAY_STATUS_USERPAYING" //用户支付中
	PAY_STATUS_PAYERROR   = "PAY_STATUS_PAYERROR"   //支付失败
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
)

//支付方式
const (
	H5PAY    = "MWEB"  //H5支付
	APPPAY   = "APP"   //APP支付
	JSAPIPAY = "JSAPI" //小程序支付
)

//******************************************************************************
