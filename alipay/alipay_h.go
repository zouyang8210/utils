// alipay_h
package alipay

const (
	ALIPAY_GATEWAY = "https://openapi.alipay.com/gateway.do"
	GET            = "GET"
	SUCCESS        = "Success"
	HTTP_STATUS_OK = 200
	EMPTY          = ""
)

//公共参数
const (
	APP_ID      = "app_id"
	METHOD      = "method"
	FORMAT      = "format"
	CHARSET     = "charset"
	SIGN_TYPE   = "sign_type"
	SIGN        = "sign"
	TIMESTAMP   = "timestamp"
	VERSION     = "version"
	NOTFIY_URL  = "notify_url"
	BIZ_CONTENT = "biz_content"
)

//请求参数
const (
	REQ_BODY           = "body"           //订单商品描述
	REQ_OUT_TRADE_NO   = "out_trade_no"   //商户订单号
	REQ_SUBJECT        = "subject"        //订单标题
	REQ_TOTAL_AMOUNT   = "total_amount"   //订单金额
	REQ_TRADE_NO       = "trade_no"       //支付宝订单号
	REQ_AUTH_CODE      = "auth_code"      //支付授权码
	REQ_SCENE          = "scene"          //支付场景(条码支付，取值：bar_code,声波支付，取值：wave_code)
	REQ_REFUND_AMOUNT  = "refund_amount"  //退款金额
	REQ_OUT_REQUEST_NO = "out_request_no" //退款订单号
)

//支付宝对像
type AliPayLib struct {
	AppId      string
	PrivateKey string
	PublicKey  string
}

//===============================================交易信息==============================================================
//交易基础信息
type DealBaseInfo struct {
	Body      string  `json:"body"`       //对一笔交易的具体描述信息
	Subject   string  `json:"subject"`    //商品标题
	TradeNo   string  `json:"trade_no"`   //商户订单号
	NotifyUrl string  `json:"notify_url"` //交易信息状态回调地址
	TotalFee  float64 `json:"total_fee"`  //交易金额
}

type ScanDealInfo struct {
	DealBaseInfo
	AuthCode string `json:"auth_code"` //授权码
}

//===========================================支付宝平台返回对像定义=========================================

type RetAliPayBase struct {
	Code    string `json:"code"`     //网关返回码
	Msg     string `json:"msg"`      //网关返回码描述
	SubCode string `json:"sub_code"` //业务返回码
	SubMsg  string `json:"sub_msg"`  //业务返回码描述
}

//创建支付二维码返回
type RetCreateCode struct {
	RetAliPayBase
	//Sign         string `json:"sign"`
	OutTradeNo string `json:"out_trade_no"` //商户的订单号
	QrCode     string `json:"qr_code"`      //当前预下单请求生成的二维码码串
}

//订单状态查询返回
type RetQueryTrade struct {
	RetAliPayBase
	TradeNo      string  //支付宝订单号
	OutTradeNo   string  //商户订单号
	BuyerLogonId string  //买家支付定账号
	TradeStatus  string  //订单状态:WAIT_BUYER_PAY（交易创建，等待买家付款）、TRADE_CLOSED（未付款交易超时关闭，或支付完成后全额退款）、TRADE_SUCCESS（交易支付成功）、TRADE_FINISHED（交易结束，不可退款）
	TotalAmount  float32 //订单交易总金额
}

//支付码交易返回
type RetMicroPay struct {
	RetAliPayBase
	TradeNo       string `json:"trade_no"`       //支付宝订单号
	OutTradeNo    string `json:"out_trade_no"`   //商户订单号
	BuyerLogonId  string `json:"buyer_logon_id"` //买家支付定账号
	TotalAmount   string `json:"total_amount"`   //订单交易总金额
	ReceiptAmount string `json:"receipt_amount"` //实收金额
	EndTime       string `json:"gmt_payment"`    //交易支付时间
}

//退款信息
type RetRefundInfo struct {
	RetAliPayBase
	TradeNo    string  `json:"trade_no"`     //支付宝订单号
	OutTradeNo string  `json:"out_trade_no"` //商户订单号
	RefundFee  float64 `json:"refund_fee"`   //退款金额
	EndTime    string  `json:"gmt_payment"`  //退款支付时间
}

//退款查询信息
type RetQueryRefund struct {
	RetAliPayBase
	TradeNo      string  `json:"trade_no"`      //支付宝订单号
	OutTradeNo   string  `json:"out_trade_no"`  //商户订单号
	TotalAmount  float64 `json:"total_amount"`  //该笔退款所对应的交易的订单金额
	RefundAmount float64 `json:"refund_amount"` //本次退款请求，对应的退款金额
}

//返回的签名
type RetSign struct {
	Sign string `json:"sign"` //签名
}

//交易异步通知
type NotifyInfo struct {
	NotifyTime    string  `json:"notify_time"`    //通知时间
	NotifyType    string  `json:"notify_type"`    //通知类型
	NotifyId      string  `json:"notify_id"`      //通知校验ID
	TradeNo       string  `json:"trade_no"`       //支付宝交易号
	AppId         string  `json:"app_id"`         //开发者的app_id
	OutTradeNo    string  `json:"out_trade_no"`   //商户订单号
	BuyerLogonId  string  `json:"buyer_logon_id"` //买家支付宝账号
	TradeStatus   string  `json:"trade_status"`   //交易状态(WAIT_BUYER_PAY-交易创建,TRADE_CLOSED-关闭,TRADE_SUCCESS-完成,TRADE_FINISHED-交易结束,不可退款)
	TotalAmount   float64 `json:"total_amount"`   //订单金额
	ReceiptAmount float64 `json:"receipt_amount"` //实收金额
}
