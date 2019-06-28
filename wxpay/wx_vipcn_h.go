// wx_vipcn_h
package weixin

import (
	"time"
)

//公众号
type PublicAccount struct {
	AppId       string    //appId
	AppSecret   string    //密钥
	accessToken string    //令牌
	expiresIn   int       //access_token有效时长
	accessTime  time.Time //获取access_token时间
}

type RetError struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

//登录返回
type RetLogin struct {
	RetError
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
}

//获取关注过公众号的微信的openid列表返回
type RetOpenIds struct {
	RetError
	Data       RetOpenidList `json:"data"`
	NextOpenId string        `json:"next_openid"`
	Total      int           `json:"total"`
	Count      int           `json:"count"`
}

//Openid列表
type RetOpenidList struct {
	Openid []string `json:"openid"`
}

//获取openid时返回的信息
type RetAuthorizeOpenid struct {
	RetError
	AccessToken  string `json:"accessToken"`   //网页授权接口调用凭证,注意：此access_token与基础支持的access_token不同
	RefreshToken string `json:"refresh_token"` //用于刷新AccessToken
	OpenId       string `json:"openid"`        //openid
	Scope        string `json:"scope"`         //获取更多的微信信息的密钥
	ExpiresIn    int    `json:"expiresIn"`     //到期时间
}

//消息模块
type ModuleMessage struct {
	ToUser     string     `json:"touser"`      //要发送信息的用户
	TemplateId string     `json:"template_id"` //消息模板ID
	Data       ModuleData `json:"data"`        //数据
}

//模板数据
type ModuleData struct {
	First    ModuleValue `json:"first"`    //头标题
	Keyword1 ModuleValue `json:"keyword1"` //车位名称
	Keyword2 ModuleValue `json:"keyword2"` //设备ID
	Remark   ModuleValue `json:"remark"`   //描述
}

//模板数据值
type ModuleValue struct {
	Value string `json:"value"` //值
	Color string `json:"color"` //字体颜色
}

//接口URL
const (
	GET_TOKEN_URL          = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"                                                    //获取公众号接口访问令牌
	GET_OPENIDS_URL        = "https://api.weixin.qq.com/cgi-bin/user/get?accessToken=%s&next_openid=%s"                                                                   //获取观注公众号的微信OPENID
	AUTHORIZE_CODE         = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_base&state=%s#wechat_redirect" //获取code
	AUTHORIZE_OPENID       = "https://api.weixin.qq.com/sns/oauth2/accessToken?appid=%s&secret=%s&code=%s&grant_type=authorization_code"                                  //通过code 得到access_token和openid
	SEND_MODEL_MESSAGE_URL = "https://api.weixin.qq.com/cgi-bin/message/template/send?accessToken=%s"                                                                     //发送模板消息
)
