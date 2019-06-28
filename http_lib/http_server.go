package http_lib

import (
	"net/http"
	"time"
	"utils/data_conv/number_lib"
)

type HttpServer struct {
	mux    *http.ServeMux
	server *http.Server
}

func (sender *HttpServer) Router(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	if sender.mux == nil {
		sender.mux = http.NewServeMux()
	}
	sender.mux.HandleFunc(pattern, handler)
}

/*
功能:开启http服务
参数:
	port:服务端口
返回:错误信息
*/
func (sender *HttpServer) HttpListen(port int, timeOut int) (err error) {
	strPort := ":" + number_lib.NumberToStr(port)
	sender.server = &http.Server{
		Addr:         strPort,
		WriteTimeout: time.Second * time.Duration(timeOut),
		ReadTimeout:  time.Second * time.Duration(timeOut),
		Handler:      sender.mux,
	}
	err = sender.server.ListenAndServe()
	return
}

/*
功能:开启https服务
参数:
	port:服务端口
	certFile:证书文件路径
	keyFile:私钥文件路径
返回:错误信息
*/
func (sender *HttpServer) HttpsListen(port, timeOut int, certFile, keyFile string) (err error) {
	strPort := ":" + number_lib.NumberToStr(port)
	sender.server = &http.Server{
		Addr:         strPort,
		WriteTimeout: time.Second * time.Duration(timeOut),
		ReadTimeout:  time.Second * time.Duration(timeOut),
		Handler:      sender.mux,
	}
	err = sender.server.ListenAndServeTLS(certFile, keyFile)
	return
}
