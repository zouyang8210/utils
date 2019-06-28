// wx_vipcn_cpp
package weixin

import (
	"fmt"
	"time"
	"utils/comm_const"
	//"utils/data_conv"
	"utils/data_conv/json_lib"
	"utils/data_conv/str_lib"
	"utils/data_conv/time_lib"
	"utils/http_lib"
)

/*
功能:登录公众号
参数:
返回:是否成功,异常信息
*/
func (sender *PublicAccount) Login() (result bool, err error) {
	var rInfo RetLogin
	result = false
	tokenUrl := fmt.Sprintf(GET_TOKEN_URL, sender.AppId, sender.AppSecret)

	buff, _, err := sender.submit(GET, tokenUrl, EMPTY)
	if err == nil {
		json_lib.JsonToObject(string(buff), &rInfo)
		if rInfo.ErrCode == 0 {
			result = true
			sender.accessTime = time.Now()
			sender.accessToken = rInfo.AccessToken
			sender.expiresIn = rInfo.ExpiresIn
		}
	}
	return
}

/*
功能:获取观注过公众号的微信的Openid
参数:
	next_id:开始获取的位置,
返回:RetOpenIds,异常信息
*/
func (sender *PublicAccount) GetOpenIdList(nextId string) (result RetOpenIds, err error) {
	sender.checkTimeout()
	openidUrl := fmt.Sprintf(GET_OPENIDS_URL, sender.accessToken, nextId)
	buff, _, err := sender.submit(GET, openidUrl, EMPTY)
	if err == nil {
		json_lib.JsonToObject(string(buff), &result)
	}
	return
}

func (sender *PublicAccount) GetCode(phone string, callbackUrl string) (redirectUrl, submitUrl string) {
	redirectUrl = str_lib.UrlToUrlEncode(callbackUrl)
	submitUrl = fmt.Sprintf(AUTHORIZE_CODE, sender.AppId, redirectUrl, phone)
	return
}

/*
功能:通过code获取到openid
参数:
	code:获取openId的编号
返回:RetAuthorizeOpenid,异常信息
*/
func (sender *PublicAccount) GetAuthorizeOpenid(code string) (result RetAuthorizeOpenid, err error) {
	submitUrl := fmt.Sprintf(AUTHORIZE_OPENID, sender.AppId, sender.AppSecret, code)
	buff, _, err := sender.submit(GET, submitUrl, EMPTY)
	if err == nil {
		json_lib.JsonToObject(string(buff), &result)
	}
	return
}

/*
功能:发送模板消息
参数:
	carport:车位编号
	openid:微信编号
	status:状态
返回:RetError,异常信息
*/
func (sender *PublicAccount) SendModelMessage(carport, openid, status string) (result RetError, err error) {
	var sendInfo ModuleMessage
	sender.checkTimeout()
	sendInfo.ToUser = openid
	sendInfo.TemplateId = "5NGQWfUrq8oWXbCzzt6Hhd8HEpEgu1RFkjIn6UAqewI"
	sendInfo.Data.First.Value = fmt.Sprintf("车位 %s状态改变通知", carport)
	sendInfo.Data.Keyword1.Value = status
	sendInfo.Data.Keyword2.Value = time.Now().Format(comm_const.TIME_yyyyMMddHHmmss)
	sendInfo.Data.Remark.Value = "欢迎使用开能车位锁"
	json, _ := json_lib.ObjectToJson(sendInfo)
	buff, _, err := sender.submit(POST, fmt.Sprintf(SEND_MODEL_MESSAGE_URL, sender.accessToken), json)
	json_lib.JsonToObject(string(buff), &result)
	return
}

func (sender *PublicAccount) checkTimeout() {
	sec := int(time_lib.TimeSub(sender.accessTime, time.Now()))
	if sec > sender.expiresIn || sec-sender.expiresIn < 600 {
		sender.Login()
	}
}

/*
功能:提交数据
参数:
	submitType:提交类型
	url:接口路径
	data:提交的字符串
返回:API返回数据,http状态代码,错误
*/
func (sender *PublicAccount) submit(submitType, url, data string) ([]byte, int, error) {
	resp, err := http_lib.HttpSubmit(submitType, url, data, nil)
	return resp.BufferBody, resp.HttpStatus, err
}
