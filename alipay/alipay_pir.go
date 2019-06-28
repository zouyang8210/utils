// alipay_pir
package alipay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/axgle/mahonia"
	"utils/comm_const"
	"utils/data_conv/json_lib"
	"utils/data_conv/str_lib"
	"utils/http_lib"
)

/*
功能:填充公共参数
参数:
	method:接口名称
	notifyUrl:回调地址
返回:公共参数
*/
func (c *AliPayLib) putPublic(method, notifyUrl string) (params map[string]interface{}) {
	params = make(map[string]interface{})
	params[APP_ID] = c.AppId
	params[METHOD] = method
	params[FORMAT] = "JSON"
	params[CHARSET] = "utf-8"
	params[SIGN_TYPE] = "RSA2"
	params[TIMESTAMP] = time.Now().Format(comm_const.TIME_yyyyMMddHHmmss)
	params[VERSION] = "1.0"
	if notifyUrl != EMPTY {
		params[NOTFIY_URL] = notifyUrl
	}
	return
}

/*
功能:填加参数
参数:
	name:参数名称
	value:参数值
返回:
*/
func (c *AliPayLib) putParam(params map[string]interface{}, name string, value interface{}) {
	params[name] = value
	return
}

/*
功能:整理提交的数据并计算签名
参数:
	params:参数集合
返回:整理完成的提交数据,签名
*/
func (c *AliPayLib) putComplete(params map[string]interface{}) (reqStr string, sign string) {
	//排序
	keys := c.Sort(params)

	//拼连提交数据和待签名字符串
	var waitSign string
	for _, n := range keys {
		waitSign += fmt.Sprintf("%s=%s&", n, params[n])
	}
	//待签名字符串
	waitSign = str_lib.SubString(waitSign, 0, len(waitSign)-1)
	//签名
	sign, _ = c.Rsa2Sign(waitSign)
	//签名加入参数集合
	c.putParam(params, SIGN, sign)

	//加入签名后,再次排序
	keys = make([]string, 0)
	keys = c.Sort(params)

	//参数值转码
	for _, n := range keys {
		reqStr += fmt.Sprintf("%s=%s&", n, str_lib.UrlToUrlEncode(params[n].(string)))
	}
	//去除字符串中,最后一个'&'符号
	reqStr = str_lib.SubString(reqStr, 0, len(reqStr)-1)
	return
}

/*
功能:RSA2签名
参数:
	data:待签名字符串
返回:签名,错误信息
*/
func (c *AliPayLib) Rsa2Sign(data string) (sign string, err error) {
	var privateKey *rsa.PrivateKey
	var bSign []byte
	mPrivate := "-----BEGIN PRIVATE KEY-----\r\n" + c.PrivateKey + "\r\n-----END PRIVATE KEY-----"
	h := sha256.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)
	blockPri, _ := pem.Decode([]byte(mPrivate))
	if blockPri != nil {
		privateKey, err = x509.ParsePKCS1PrivateKey(blockPri.Bytes)
		if err == nil {
			bSign, err = rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
			if err == nil {
				sign = base64.StdEncoding.EncodeToString(bSign)
			}
		}
	} else {
		err = errors.New("密钥不能还原")
	}
	return
}

/*
功能:RSA2验签
参数:
	data:待验签字符串
返回:验签是否成功(验签失败,查看错误信息),错误信息
*/
func (c *AliPayLib) VerifyRas2Sign(data, sign string) (result bool, err error) {
	var pubInterface interface{}
	result = false
	mPublic := "-----BEGIN PUBLIC KEY-----\r\n" + c.PublicKey + "\r\n-----END PUBLIC KEY-----"
	block, _ := pem.Decode([]byte(mPublic))

	h := sha256.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)

	if block != nil {
		pubInterface, err = x509.ParsePKIXPublicKey(block.Bytes)
		if err == nil {
			var bSign []byte
			pub := pubInterface.(*rsa.PublicKey)
			bSign, err = base64.StdEncoding.DecodeString(sign)
			if err == nil {
				err = rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed, bSign)
				if err == nil {
					result = true
				}
			}
		}
	} else {
		err = errors.New("公钥不能还原")
	}
	return
}

/*
功能:参数名排序
参数:
	params:参数集合
返回:排序后的参数名称数组
*/
func (c *AliPayLib) Sort(params map[string]interface{}) (array []string) {
	for k := range params {
		array = append(array, k)
	}
	sort.Strings(array)
	return
}

/*
功能:把接口返回的JSON处理成我们需要样式
参数:
	data:待签名字符串
返回:签名,错误信息
*/
func (c *AliPayLib) conv(buff []byte) (json string, sign string) {
	var mSign RetSign
	json = string(buff)
	pos1 := strings.Index(json, ":") + 1
	pos2 := strings.LastIndex(json, ",")
	json = str_lib.SubString(json, pos1, pos2-pos1)
	json_lib.JsonToObject(string(buff), &mSign)
	sign = mSign.Sign
	return
}

/*
功能:向支付宝平台提交GET请求
参数:
	reqStr:请求字符串
返回:支付宝平台返回信息,错误信息
*/
func (c *AliPayLib) getSubmit(reqStr string) (json string, err error) {
	var sign string
	resp, err := http_lib.HttpsSubmit(GET, fmt.Sprintf("%s?%s", ALIPAY_GATEWAY, reqStr), EMPTY, nil)
	if err == nil && resp.HttpStatus == HTTP_STATUS_OK {
		json, sign = c.conv(resp.BufferBody)
		_, err = c.VerifyRas2Sign(json, sign)
		if err != nil {
			err = errors.New("sign invalid:" + err.Error())
		}
	} else {
		if err == nil {
			err = errors.New(fmt.Sprintf("http status:%d", resp.HttpStatus))
		}
	}
	return
}

/*
功能:请求H5支付
参数:
	reqStr:请求字符串
返回:支付宝平台返回信息,错误信息
*/
func (c *AliPayLib) requestH5(reqStr string) (h5PayPage string, err error) {
	resp, err := http_lib.HttpSubmit(http_lib.GET, fmt.Sprintf("%s?%s", ALIPAY_GATEWAY, reqStr), EMPTY, nil)
	//fmt.Println(fmt.Sprintf("%s?%s", ALIPAY_GATEWAY, reqStr))
	if err == nil && resp.HttpStatus == HTTP_STATUS_OK {
		h5PayPage = resp.Body
	} else {
		if err == nil {
			err = errors.New(fmt.Sprintf("http status:%d", resp.HttpStatus))
		}
	}
	return
}

func (c *AliPayLib) convertToString(src, strCode, tagCode string) string {
	strCoder := mahonia.NewDecoder(strCode)
	strResult := strCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(strResult), true)
	result := string(cdata)
	return result
}
