package http_lib

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type HttpClient struct {
}

/*
功能：https带证书提交
参数：
	method:提交方式
	strUlr:url
	bodyData:提交数据
	headData:头部数据
	certFile:证书文件路径
	keyFile:私钥文件路径
返回：
*/
func (sender *HttpClient) HttpsSecureSubmit(method, strUrl, bodyData string, headData http.Header, certFile, keyFile string) (response ResponseInfo, err error) {
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
		response, err = sender.httpSubmit(client, method, strUrl, bodyData, headData)
	}
	return
}

/*
功能：https不验证服务器
参数：
	method:提交方式（POST,GET,OPTIONS...)
	strUlr:地址
	bodyData:body内容
	headData:头信息
*/
func (sender *HttpClient) HttpsSubmit(method, strUrl, bodyData string, headData http.Header) (response ResponseInfo, err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
	}
	response, err = sender.httpSubmit(client, method, strUrl, bodyData, headData)
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
func (sender *HttpClient) HttpSubmit(method, strUrl, bodyData string, headData http.Header) (response ResponseInfo, err error) {
	client := &http.Client{}
	response, err = sender.httpSubmit(client, method, strUrl, bodyData, headData)
	return
}

/*
功能：https form提交
参数：
	strUlr:地址
	formData:表单数据
	headData:头信息
*/
func (sender *HttpClient) HttpsForm(strUrl string, formData map[string]string, headData http.Header) (response ResponseInfo, err error) {
	var data = ""
	if formData != nil {
		for n, v := range formData {
			data += fmt.Sprintf("%s=%s&", n, v)
		}
		if data != "" {
			data = data[0 : len(data)-1]
		}

		headData.Set("Content-CategoryName", "application/x-www-form-urlencoded")
	}
	response, err = HttpsSubmit(POST, strUrl, data, headData)

	return
}

func (sender *HttpClient) httpSubmit(client *http.Client, method, strUrl, bodyData string, headData http.Header) (response ResponseInfo, err error) {
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
