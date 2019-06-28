// tcp_demo
package main

import (
	"fmt"
	"os"
	"time"

	"os/signal"

	"utils/data_conv/number_lib"
	"utils/file"
	"utils/net_tcp/tcp_client"
)

var client tcp_client.TcpClient

var quit = make(chan os.Signal)

type DataHead struct {
	Size    int32
	Type    int32
	Status  int32
	Channel int32
	Time    int32
	Date    int32
	Id      int32
	Level   int32
}

const MODEL_TCP_SERVER = 0      //tcp 服务器模式
const MODEL_TCP_CLIENT = 1      //tcp 客户端模式
const MODEL_TCP_SERVER_CHAN = 2 //服务器模式(通道)

var tmpCmd = ""
var cmd string

func main() {
	model := MODEL_TCP_SERVER_CHAN

	//go func() {
	//	signal.Notify(quit, os.Interrupt, os.Kill)
	//	<-quit
	//	if model == MODEL_TCP_SERVER {
	//		threadService.CloseAll()
	//	} else if model == MODEL_TCP_CLIENT {
	//		client.Close()
	//	}
	//	os.Exit(0)
	//}()

	if model == MODEL_TCP_SERVER_CHAN {
		fmt.Println("Tcp Service Run.....")
		go threadSever()

		fmt.Println("Tcp Service Run(chan).....")
		chanServer()
	} else if model == MODEL_TCP_CLIENT {
		fmt.Println("Capture Start Run,Input \"capture\" Capture a Picture ")
		for true {
			fmt.Scan(&cmd)
			switch cmd {
			case "capture":
				fmt.Println("exec capture")
				capture()
			case "exit":
				quit <- os.Kill
			}
		}
	} else if model == MODEL_TCP_SERVER_CHAN {
		fmt.Println("Tcp Service Run(chan).....")
		//chanServer()
	}

	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
	if model == MODEL_TCP_SERVER {
		threadService.CloseAll()
	} else if model == MODEL_TCP_CLIENT {
		client.Close()
	}
	os.Exit(0)
}

//抓拍
func capture() {
	var info DataHead
	info.Size = 32
	info.Type = 0x6550
	buff := number_lib.ObjectToBytes(info, number_lib.DESC)
	defer timeCost(time.Now())

	err := client.Connect("192.168.16.102", 8117)
	if err == nil {
		fmt.Println("Connect Success")
		count, err := client.SendData(buff, 0, len(buff))
		if err == nil {
			fmt.Println("WriteData Capture Data ", count)
			socketRecv()
		} else {
			fmt.Println("WriteData Data Error: ", err)
		}
	} else {
		fmt.Println("Connect Error: ", err)
	}
}

//接收数据
func socketRecv() {
	buff := make([]byte, 102400)
	index := 0 //接收数据buff起始位置
	size := 32 //接收数据大小
	headLen := 32
	receiveCount := 0 //总接收数据量
ReceiveHead:
	fmt.Println("Receive Package Head")
	client.SetReadTimeout(1)
	count, err := client.ReadData(buff, index, size) //接收数据包头
	if err != nil {
		fmt.Println("Receive Package Head Error :", err)
		client.Close()
		return
	} else {
		fmt.Println("Receive Package Head Success")
	}
	receiveCount += count
	var tmp int32
	if count == size {
		number_lib.BytesToObject(buff[index:4], number_lib.DESC, &tmp)
		size = int(tmp)

		if size >= 0 {
			fmt.Println("Should Receive Data ", size)
			for receiveCount < size {
				client.SetReadTimeout(1)
				count, err = client.ReadData(buff, index+receiveCount, size-receiveCount)
				receiveCount += count
				if err == nil {
					fmt.Println("Receive Data ", receiveCount)
				} else {
					fmt.Println("Receive Data Error:", err)
					break
				}
			}
			if receiveCount == size {
				fileName := fmt.Sprintf("capture/%d.jpg", time.Now().Unix())
				if file.WriteFile(fileName, buff[index+headLen:size]) {
					fmt.Printf("Write Image File Success\n")
				}
			} else {
				fmt.Println("No Received Full Data")
				client.Close()
			}
		}
	} else {
		goto ReceiveHead
	}
}

func timeCost(start time.Time) {
	terminal := time.Since(start)
	fmt.Println(terminal)
}
