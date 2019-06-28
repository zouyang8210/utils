package main

import (
	"fmt"
	"os"
	"os/signal"
	"utils/data_conv/number_lib"
	"utils/net_tcp/tcp_service"
)

var chanService tcp_service.TcpServer

func chanServer() {

	chanService.SetConnectCallback(chanConnEvent) //设置连接事件回调
	chanService.SetDataCallback(chanReceiveData)  //设置数据接收事件回调
	//设置数据接收协议格式
	format := tcp_service.ProtocolFormat{}
	format.Head = []byte{0x02, 0xfd}
	format.End = []byte{0x0a, 0x0d}
	//format.DataSize = 32
	//chanService.SetProtocolFormat(format)
	chanService.SetMaxConnect(30000)
	chanService.SetCheckConn(false)
	//启动服务
	go func() {
		res, err := chanService.Listen(8001)
		if !res {
			fmt.Printf("Listen Fail,err=%v\n", err)
		}
	}()

	//defer release() //释放资源

	//for true {
	//	fmt.Scan(&cmd)
	//	switch cmd {
	//	case "count":
	//		fmt.Printf("count = [%d] \n", chanService.GetConnectCount())
	//	case "send", "close":
	//		tmpCmd = cmd
	//	case "exit":
	//		quit <- os.Kill
	//		break
	//	default:
	//		if tmpCmd == "send" {
	//			chanSend(cmd, "hello")
	//		} else if tmpCmd == "close" {
	//			chanCloseOne(cmd)
	//		}
	//	}
	//}
}

//关闭一个连接
func chanCloseOne(index string) {
	var i int32
	err := number_lib.StrToInt(index, &i)
	if err == nil {
		chanService.CloseOne(int(i))
	}
}

//发送数据
func chanSend(index, content string) {
	var i int32
	err := number_lib.StrToInt(index, &i)
	if err == nil {
		//fmt.Printf("send data [%s]\n", content)
		if count, err := chanService.WriteData([]byte(content), int(i), 0, len(content)); count > 0 {
			//fmt.Println("send success")
		} else {
			fmt.Printf("send fail [%v]\n", err)
		}
	}
}

//数据接收事件
func chanReceiveData(data []byte, dataSize int, clientInfo tcp_service.ClientInfo) {
	chanService.WriteData(data, clientInfo.Index, 0, dataSize)
}

//连接事件
func chanConnEvent(addr string, index int, value bool) {
	msg := ""
	if value {
		msg = "连接成功"
	} else {
		msg = "断开"
	}
	fmt.Printf("连接信息:index = [%d] [%s] [%s] \n", index, addr, msg)
}

//释放
func chanRelease() {
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
	chanService.CloseAll()
	fmt.Println("exit")
	os.Exit(0)
}
