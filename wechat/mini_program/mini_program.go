package mini_program

import (
	"errors"
	"fmt"
	"time"
	"utils/data_conv/json_lib"
	"utils/data_conv/time_lib"
	"utils/http_lib"
	. "utils/wechat"
)

type MiniProgram struct {
	PublicAppId      string    `json:"public_app_id"`
	MinProgramId     string    `json:"min_program_id"`     //小程序ID
	MinProgramSecret string    `json:"min_program_secret"` //小程序密钥
	accessToken      string    `json:"access_token"`       //访问令牌
	loginTime        time.Time `json:"login_time"`         //登录时间
	expire           int       `json:"expire"`             //过期时间
}

//小程序登录
func (sender *MiniProgram) MinProgramLogin() (err error) {
	var info RetMinProgramLogin
	var bb []byte
	url := fmt.Sprintf(GET_ACCESS_TOKEN, sender.MinProgramId, sender.MinProgramSecret)
	bb, _, err = Submit(http_lib.GET, url, EMPTY)
	if err == nil {
		err = json_lib.JsonToObject(string(bb), &info)
		if info.ErrCode == OK {
			sender.accessToken = info.AccessToken
			sender.loginTime = time.Now()
			sender.expire = info.ExpiresIn
		} else {
			err = errors.New(fmt.Sprintf("erro info:code=%d,msg=%s\n", info.ErrCode, info.ErrMsg))
		}
	}
	return
}

/*
功能:通过小程序传入的code，调用API获取openid
参数:
	code:小程序生成的临时密钥
返回:openid,错误信息
*/
func (sender *MiniProgram) GetOpenIdMinProgram(code string) (openId string, err error) {
	var info RetMinProgramOpenId
	var bb []byte
	url := fmt.Sprintf(JSOCDE_TO_SESSION_URL, sender.MinProgramId, sender.MinProgramSecret, code)
	bb, _, err = Submit(http_lib.GET, url, EMPTY)
	if err == nil {
		err = json_lib.JsonToObject(string(bb), &info)
		openId = info.OpenId
	}
	return
}

/*
功能:登录凭证校验。通过 wx.login 接口获得临时登录凭证 code 后传到开发者服务器调用此接口完成登录流程
参数:
	code:小程序生成的临时密钥
返回:登录后的信息,错误信息
*/
func (sender *MiniProgram) AuthCodeToSession(code string) (result RetMinProgramOpenId, err error) {
	var bb []byte
	url := fmt.Sprintf(JSOCDE_TO_SESSION_URL, sender.MinProgramId, sender.MinProgramSecret, code)
	bb, _, err = Submit(http_lib.GET, url, EMPTY)
	if err == nil {
		err = json_lib.JsonToObject(string(bb), &result)
	}
	return
}

//发送统一服务消息
func (sender *MiniProgram) SendMiniProgramMsg(moduleData interface{}) (err error) {
	msgJson, _ := json_lib.ObjectToJson(moduleData)
	sender.checkLogin()
	var ret RetError
	bb, _, err := Submit(http_lib.POST, fmt.Sprintf(SEND_MINI_MESSAGE, sender.accessToken), msgJson)
	if err == nil {
		json_lib.JsonToObject(string(bb), &ret)
		if ret.ErrCode != OK {
			err = errors.New(fmt.Sprintf("unify send message fail:code=%d,msg=%s", ret.ErrCode, ret.ErrMsg))
		}
	}
	return
}

//检测令牌是否过期,过期重新登录
func (sender *MiniProgram) checkLogin() {
	now := time.Now()
	sec := time_lib.TimeSub(sender.loginTime, now)
	if int(sec)+120 > sender.expire {
		sender.MinProgramLogin()
	}
}
