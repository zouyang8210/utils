package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"utils/data_conv/json_lib"
	"utils/data_conv/str_lib"
	"utils/http_lib"
)

type SmartAuth struct {
	authorizationPage string
	waitAuthList      map[string]OAuthInfo
	checkClientCall   CheckClientCallback
	loginCall         LoginCallback
	getTokenCall      AccessTokenCallback
}

//请求鉴权信息
type OAuthInfo struct {
	RedirectUrl  string `json:"redirect_url"`  //跳转页
	ClientId     string `json:"client_id"`     //客户ID
	ResponseType string `json:"response_type"` //请求鉴权类型
	State        string `json:"state"`         //用户自义信息
	Code         string `json:"code"`          //换取access token的代码
}

type RetOAuth struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int    `json:"expires_in"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

//======================================================================================================================
type TmBaseInfo struct {
	Header  TmHeader      `json:"header"`
	Payload tmBasePayLoad `json:"payload"`
}

type TmInfo struct {
	Header  TmHeader  `json:"header"`
	Payload TmPayLoad `json:"payload"`
}

type TmHeader struct {
	Namespace      string `json:"namespace"`      //消息命名空间
	Name           string `json:"name"`           //
	MessageId      string `json:"messageId"`      //用于跟踪请求
	PayLoadVersion int    `json:"payLoadVersion"` //payload的版本,目前版本为1
}

type tmBasePayLoad struct {
	AccessToken string    `json:"accessToken"`
	DeviceId    string    `json:"deviceId"`
	DeviceType  string    `json:"deviceType"`
	Attribute   string    `json:"attribute"`
	Value       string    `json:"value"`
	Extensions  Extension `json:"extensions"`
	ErrorCode   string    `json:"errorCode"`
	Message     string    `json:"message"`
}

type TmPayLoad struct {
	AccessToken string       `json:"accessToken"`
	Devices     []DeviceInfo `json:"device_platform"`
}

type Extension struct {
	Extension1 string `json:"extension1"`
	Extension2 string `json:"extension2"`
}

type DeviceProperties struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type DeviceInfo struct {
	DeviceId   string             `json:"deviceId"`
	DeviceName string             `json:"deviceName"`
	DeviceType string             `json:"deviceType"`
	Zone       string             `json:"zone"`
	Brand      string             `json:"brand"`
	Model      string             `json:"model"`
	Icon       string             `json:"icon"`
	Properties []DeviceProperties `json:"properties"`
	Actions    []string           `json:"actions"`
	Extensions Extension          `json:"extensions"`
}

//======================================================================================================================

//回调定义
type (
	CheckClientCallback func(clientId string) bool         //client id 是否存
	LoginCallback       func(uid string, pw string) bool   //用户登录
	AccessTokenCallback func(clientId, secret string) bool //换取token回调
)

const (
	REDIRECT_URI  = "redirect_uri"
	CLIENT_ID     = "client_id"
	RESPONSE_TYPE = "response_type"
	STATE         = "state"
	PASSWORD      = "pw"
	USERID        = "uid"
	EMPTY         = ""
	CODE          = "code"
	SECRET        = "secret"
)

const (
	//CERT_FILE_PATH = "cert/fullchain.pem" //中间证书文件
	//KEY_FILE_PATH  = "cert/privkey.pem"   //密钥文件

	CERT_FILE_PATH = "cert/fullchain.pem"
	KEY_FILE_PATH  = "cert/privkey.pem"
)

func (sender *SmartAuth) Initialize() {

	sender.authorizationPage = `<html><body>
		<a href="/TM/login?client_id=7e9ef2fcc2a74c6fb0437a83e99972cd&uid=1001&pw=123456">Log in </a>
		</body></html>`
	sender.waitAuthList = make(map[string]OAuthInfo)
	http.HandleFunc("/TM/auth2/authorization", sender.auth2Entry)
	http.HandleFunc("/TM/login", sender.loginApplication)
	http.HandleFunc("/TM/getAccessToken", sender.getAccessToken)
	http.HandleFunc("/TM/gateway", sender.gateway)

	err := http.ListenAndServeTLS(":444", CERT_FILE_PATH, KEY_FILE_PATH, nil)

	if err != nil {
		fmt.Println(err)
	}
}

//设置应用ID测查回调
func (sender *SmartAuth) SetCheckClient(callback CheckClientCallback) {
	sender.checkClientCall = callback
}

//设置登录验证回调
func (sender *SmartAuth) SetLogin(callback LoginCallback) {
	sender.loginCall = callback
}

//设置获取token验证回调
func (sender *SmartAuth) SetAccessToken(callback AccessTokenCallback) {
	sender.getTokenCall = callback
}

func (sender *SmartAuth) gateway(w http.ResponseWriter, r *http.Request) {
	_, strData, err := http_lib.GetBody(r)
	if err == nil {
		fmt.Println(strData)
		var requestInfo TmBaseInfo
		json_lib.JsonToObject(strData, &requestInfo)
		switch requestInfo.Header.Namespace {
		case "AliGenie.Iot.Device.Discovery":
			switch requestInfo.Header.Name {
			case "DiscoveryDevices":
				var deviceInfo TmInfo
				requestInfo.Header.Name = "DiscoveryDevicesResponse"
				deviceInfo.Header = requestInfo.Header

				deviceInfo.Payload.Devices = make([]DeviceInfo, 1)

				deviceInfo.Payload.Devices[0].DeviceId = "1234567890"
				deviceInfo.Payload.Devices[0].DeviceName = "灯"
				deviceInfo.Payload.Devices[0].DeviceType = "light"
				deviceInfo.Payload.Devices[0].Icon = "http://img.zcool.cn/community/01120455447f0f0000019ae94f9713.jpg@1280w_1l_2o_100sh.jpg"
				deviceInfo.Payload.Devices[0].Properties = make([]DeviceProperties, 1)
				deviceInfo.Payload.Devices[0].Properties[0].Name = "powerstate"
				deviceInfo.Payload.Devices[0].Properties[0].Value = "off"
				deviceInfo.Payload.Devices[0].Actions = make([]string, 2)
				deviceInfo.Payload.Devices[0].Actions[0] = "TurnOn"
				deviceInfo.Payload.Devices[0].Actions[1] = "TurnOff"
				deviceInfo.Payload.Devices[0].Model = "bbb"
				deviceInfo.Payload.Devices[0].Brand = "Joyoung 九阳"
				json, _ := json_lib.ObjectToJson(deviceInfo)
				fmt.Println(json)
				w.Write([]byte(json))
			}
		case "AliGenie.Iot.Device.Control":
			switch requestInfo.Header.Name {
			case "TurnOn", "TurnOff":
				var reply TmBaseInfo
				requestInfo.Header.Name = "TurnOnResponse"
				reply.Header = requestInfo.Header
				reply.Payload.DeviceId = requestInfo.Payload.DeviceId
				json, _ := json_lib.ObjectToJson(reply)
				fmt.Println(json)
				w.Write([]byte(json))
			}
		}
	} else {
		fmt.Println(err)
	}

}

//oauth2签权入口
func (sender *SmartAuth) auth2Entry(w http.ResponseWriter, r *http.Request) {
	mapParam, err := http_lib.QueryString(r)
	if err == nil {
		err = sender.checkStrMap(mapParam, REDIRECT_URI, CLIENT_ID, RESPONSE_TYPE, STATE)
		if err == nil {
			if sender.checkClientCall != nil && sender.checkClientCall(mapParam[CLIENT_ID]) {
				var info OAuthInfo
				info.RedirectUrl = mapParam[REDIRECT_URI]
				info.ClientId = mapParam[CLIENT_ID]
				info.ResponseType = mapParam[RESPONSE_TYPE]
				info.State = mapParam[STATE]
				//fmt.Println(info)
				sender.waitAuthList[info.ClientId] = info
				//fmt.Println(sender.waitAuthList[info.ClientId])
				if sender.authorizationPage == EMPTY {
					fmt.Println("authorization page is not exist")
				} else {
					fmt.Fprintf(w, sender.authorizationPage)
				}
			} else {
				fmt.Println("client id is not exist or not set checkClientCall")
			}
		} else {
			fmt.Println("gin_check parameter error:", err)
		}
	} else {
		fmt.Println("get parameter error:", err)
	}

}

//鉴权登录
func (sender *SmartAuth) loginApplication(w http.ResponseWriter, r *http.Request) {
	mapParam, err := http_lib.QueryString(r)
	if err == nil {
		err = sender.checkStrMap(mapParam, CLIENT_ID, USERID, PASSWORD)
		if err == nil {
			clientId := mapParam[CLIENT_ID]
			if sender.loginCall != nil && sender.loginCall(mapParam[USERID], mapParam[PASSWORD]) && sender.waitAuthList[clientId].ClientId != EMPTY {
				code := fmt.Sprintf("%d", sender.random())
				redirectUlr := fmt.Sprintf("%s?code=%s&state=%s", sender.waitAuthList[clientId].RedirectUrl, code, sender.waitAuthList[clientId].State)
				tmp := sender.waitAuthList[clientId]
				tmp.Code = code
				sender.waitAuthList[clientId] = tmp
				fmt.Println("redirect url:", redirectUlr)
				http.Redirect(w, r, redirectUlr, http.StatusTemporaryRedirect)
			} else {
				fmt.Println("incorrect user name or password")
			}
		} else {
			fmt.Println("gin_check parameter error:", err)
		}
	} else {
		fmt.Println("get parameter error:", err)

	}
}

//code换取access token
func (sender *SmartAuth) getAccessToken(w http.ResponseWriter, r *http.Request) {
	var info RetOAuth
	_, body, err := http_lib.GetBody(r)
	if err == nil {
		fmt.Println("get access token context:", body)
		mapParam := make(map[string]string)
		mapParam, err = http_lib.GetUrlParams("http://127.0.0.1?" + body)
		clientId := mapParam[CLIENT_ID]
		if mapParam[CODE] == sender.waitAuthList[clientId].Code {
			if sender.getTokenCall != nil && sender.getTokenCall(clientId, mapParam[SECRET]) {
				info.AccessToken = str_lib.Guid()
				info.RefreshToken = str_lib.Guid()
				info.ExpiresIn = 7200
				delete(sender.waitAuthList, clientId)
				json, _ := json_lib.ObjectToJson(info)
				w.Write([]byte(json))
			} else {
				fmt.Println("incorrect secret")
			}
		} else {
			fmt.Println("code is invalid")
		}
	} else {
		fmt.Println("access token error:", err)
	}

}

//参数检测
func (sender *SmartAuth) checkStrMap(v map[string]string, paramName ...string) (err error) {
	for n := range paramName {
		if v[paramName[n]] == EMPTY {
			err = errors.New("invalid parameter:" + paramName[n])
			break
		}
	}
	return
}

//随机数
func (sender *SmartAuth) random() (n int) {
	rand.Seed(time.Now().UnixNano())
	n = rand.Intn(999999)
	return
}

func main() {
	var oauth SmartAuth
	oauth.SetCheckClient(checkClientId)
	oauth.SetLogin(login)
	oauth.SetAccessToken(getToken)
	oauth.Initialize()
}

func checkClientId(clientId string) bool {
	return true
}

func login(uid, pw string) bool {
	return true
}

func getToken(clientId, secret string) bool {
	return true
}
