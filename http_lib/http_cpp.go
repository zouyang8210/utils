// http_cpp
package http_lib

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"crypto/tls"
)

//HTTP提交方式
const (
	POST    = "POST"
	GET     = "GET"
	OPTIONS = "OPTIONS"
	PUT     = "PUT"
)

type ResponseInfo struct {
	Head       http.Header
	Cookies    http.Cookie
	Body       string
	BufferBody []byte
	HttpStatus int
}

func HttpsSecureSubmit(method, strUrl, bodyData string, headData http.Header, certFile, keyFile string) (response ResponseInfo, err error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err == nil {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				Certificates:       []tls.Certificate{cert},
			},
		}
		client := &http.Client{
			Transport: tr,
		}
		response, err = httpSubmit(client, method, strUrl, bodyData, headData)
	}
	return
}

/*
功能：https不安全验证提交
参数：
	method:提交方式（POST,GET,OPTIONS...)
	strUlr:地址
	bodyData:body内容
	headData:头信息
*/
func HttpsSubmit(method, strUrl, bodyData string, headData http.Header) (response ResponseInfo, err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
	}
	response, err = httpSubmit(client, method, strUrl, bodyData, headData)
	return
}

/*
功能：http提交
参数：
	method:提交方式（POST,GET,OPTIONS...)
	strUlr:地址
	bodyData:body内容
	headData:头信息
*/
func HttpSubmit(method, strUrl, bodyData string, headData http.Header) (response ResponseInfo, err error) {
	client := &http.Client{}
	response, err = httpSubmit(client, method, strUrl, bodyData, headData)
	return
}

//读取Request body数据
func GetBody(r *http.Request) (buffer []byte, strBody string, err error) {
	buffer, err = ioutil.ReadAll(r.Body)
	if err == nil {
		strBody = string(buffer)
	}

	return
}

func Form(r *http.Request) (mapData map[string]interface{}) {
	r.ParseForm()
	mapData = make(map[string]interface{})
	for k, v := range r.Form {
		mapData[k] = v
	}
	return
}

//从Request中获取get请求参数的值
func QueryString(r *http.Request) (params map[string]string, err error) {
	var queryStr url.Values
	params = make(map[string]string)
	queryStr, err = url.ParseQuery(r.URL.RawQuery)
	for n := range queryStr {
		params[n] = queryStr[n][0]
	}
	return
}

//从url字符串中获取get请求参数的值
func GetUrlParams(strUrl string) (params map[string]string, err error) {
	var queryStr url.Values
	mUrl, err := url.Parse(strUrl)
	params = make(map[string]string)
	if err == nil {
		queryStr, err = url.ParseQuery(mUrl.RawQuery)
		for n := range queryStr {
			params[n] = queryStr[n][0]
		}
	}
	return
}

func httpSubmit(client *http.Client, method, strUrl, bodyData string, headData http.Header) (response ResponseInfo, err error) {
	req, err := http.NewRequest(method, strUrl, strings.NewReader(bodyData))
	if err == nil {
		defer req.Body.Close()
		if headData != nil {
			for n, v := range headData {
				for j := 0; j < len(v); j++ {
					req.Header.Add(n, v[j])
				}
			}
		}
		var resp *http.Response
		resp, err = client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			response.HttpStatus = resp.StatusCode
			response.Head = resp.Header
			response.BufferBody, err = ioutil.ReadAll(resp.Body)
			if err == nil {
				response.Body = string(response.BufferBody)
			}
		}
	}
	return
}
