package main

import (
	"fmt"
	"os"
	"os/signal"

	"utils/data_conv/number_lib"
	"utils/net_tcp/tcp_service_new"
)

var threadService tcp_service.TcpServer

func threadSever() {

	threadService.SetConnectCallback(connEvent) //设置连接事件回调
	threadService.SetDataCallback(receiveData)  //设置数据接收事件回调
	//设置数据接收协议格式
	format := tcp_service.ProtocolFormat{}
	format.Head = []byte{0x02, 0xfd}
	format.End = []byte{0x0a, 0x0d}
	//format.DataSize = 32
	//threadService.SetProtocolFormat(format)
	threadService.SetMaxConnect(30000)

	//启动服务
	go func() {
		if !threadService.Listen(8100) {
			fmt.Println("Listen Fail")
		}
	}()

	//defer release() //释放资源

	//for true {
	//	fmt.Scan(&cmd)
	//	switch cmd {
	//	case "count":
	//		fmt.Printf("count = [%d] \n", threadService.GetConnectCount())
	//	case "send", "close":
	//		tmpCmd = cmd
	//	case "exit":
	//		quit <- os.Kill
	//		break
	//	default:
	//		if tmpCmd == "send" {
	//			send(cmd, "hello")
	//		} else if tmpCmd == "close" {
	//			closeOne(cmd)
	//		}
	//	}
	//}
}

//关闭一个连接
func closeOne(index string) {
	var i int32
	err := number_lib.StrToInt(index, &i)
	if err == nil {
		threadService.CloseOne(int(i))
	}
}

//发送数据
func send(index, content string) {
	var i int32
	err := number_lib.StrToInt(index, &i)
	if err == nil {
		//fmt.Printf("send data [%s]\n", content)
		if count, err := threadService.WriteData([]byte("hello"), int(i), 0, len("hello")); count > 0 {
			//fmt.Println("send success")
		} else {
			fmt.Printf("send fail [%v]\n", err)
		}
	}
}

//数据接收事件
func receiveData(data []byte, data_size int, clientInfo tcp_service.ClientInfo) {
	//threadService.WriteData(data, clientInfo.Index, 0, data_size)

	//count := threadService.GetConnectCount()
	//for i := 0; i < count; i++ {
	//if i != clientInfo.Index {
	threadService.WriteData(data, clientInfo.Index, 0, data_size)
	//fmt.Println("rec:", data)
	//}
	//}
}

//连接事件
func connEvent(addr string, index int, value bool) {
	msg := ""
	if value {
		msg = "连接成功"
	} else {
		msg = "断开"
	}
	fmt.Printf("连接信息:index = [%d] [%s] [%s] \n", index, addr, msg)
}

//释放
func release() {
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
	threadService.CloseAll()
	fmt.Println("exit")
	os.Exit(0)
}
