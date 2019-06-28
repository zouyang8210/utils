package tcp_service

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

//功能:初始化
func (sender *TcpServer) init() {
	sender.isListen = true
	sender.isCheck = true
	//如果未设置最大连接数,则自动设置10000默认值
	if sender.maxConn == 0 {
		sender.maxConn = 10000
	}

	sender.chConnChange = make(chan int)
	//sender.chConn = make(chan ClientInfo)

	sender.connTimeout = 60
	sender.checkInterval = 10

	sender.clientList = make([]ClientInfo, sender.maxConn) //设置连接池大小

	sender.countConnSum() //启动检测连接数量
	//sender.initConnPool() //初化始连接池

	go sender.checkConn()
}

//功能:准备接收数据
func (sender *TcpServer) receiveData(clientIndex int) {
	defer sender.close(clientIndex)
	for {

		buff, err := sender.readDataFormat(clientIndex)
		if err != nil {
			if strings.Contains(err.Error(), ERROR) || strings.Contains(err.Error(), IOOUT) {
				continue
			} else {
				fmt.Printf("read error=%v\n", err)
				break
			}
		} else {
			//fmt.Printf("receive data:%x\n", buff)
			sender.sendReceiveEvent(buff, len(buff), sender.clientList[clientIndex])
		}
	}
	return
}

//功能:读取缓存内的数据
//参数：
//	buff:接收数据缓存
//	clientIndex:接收数据客户端索引
//	offset:buff的偏移量
//	size:接收数据大小
//返回：实际接收数据大小，错误信息
func (sender *TcpServer) read(buff []byte, clientIndex, index, size int) (count int, err error) {
	if len(buff) < index+size {
		return 0, errors.New("buffer cache is too small")
	}
	if sender.clientList[clientIndex].Conn != nil {
		count, err = sender.clientList[clientIndex].Conn.Read(buff[index : index+size])
	} else {
		err = errors.New("connect is nil")
	}
	return
}

//功能:按协议格式读取数据
//参数：
//	clientIndex:接收数据客户端索引
//返回：符合协议的完整数据包，错误信息
func (sender *TcpServer) readDataFormat(clientIndex int) (buff []byte, err error) {
	buff = make([]byte, 1024)
	offset := 0
	//time_lib.RunTime(time.Now(), "read use time:")
	defer func() {
		if sender.clientList[clientIndex].Conn != nil {
			sender.clientList[clientIndex].Conn.SetReadDeadline(time.Time{})
		}
	}()

	if sender.protoFormat.DataSize > 0 { //接收固定长度的数据
		err = sender.readFixedData(buff, clientIndex, &offset, sender.protoFormat.DataSize)
	} else if sender.protoFormat.headLen > 0 && sender.protoFormat.endLen > 0 { //接收有头,有尾的数据
		err = sender.readHeadData(buff, clientIndex, &offset, sender.protoFormat.headLen)
		if err == nil {
			err = sender.readToEnding(buff, clientIndex, &offset, sender.protoFormat.endLen)
		}
	} else if sender.protoFormat.endLen > 0 { //接收只有结尾符的数据
		err = sender.readToEnding(buff, clientIndex, &offset, sender.protoFormat.endLen)
	} else { //接收没有设置接收格式的数据
		offset, err = sender.read(buff, clientIndex, 0, len(buff))
	}
	//offset, err = sender.read(buff, clientIndex, 0, len(buff))
	if err == nil {
		sender.clientList[clientIndex].LastPackageTime = time.Now().Unix()
	}

	return buff[0:offset], err
}

//功能:读取固定长度的信息
//参数：
//	buff:接收数据缓存
//	clientIndex:接收数据客户端索引
//	offset:buff的偏移量
//	count:接收数据长度
//返回：实际接收数据大小，错误信息
func (sender *TcpServer) readFixedData(buff []byte, clientIndex int, offset *int, count int) (err error) {
	sum := 0
	tmpLen := 0
	for true {
		if ok := sum < count; ok { //没有接收到协议尾,断续接收
			if *offset < len(buff) {
				tmpLen, err = sender.read(buff, clientIndex, *offset, count-sum)
				if err != nil {
					err = errors.New("read ending error: " + err.Error())
					break
				}
				*offset += tmpLen
				sum += tmpLen
			} else {
				//数据缓存益出
				err = errors.New("read fixed error: data buffer out")
				break
			}
		} else {
			break
		}

	}
	return
}

//======================================读取数据头====================================================================

func (sender *TcpServer) readHeadData(buff []byte, clientIndex int, offset *int, headLen int) (err error) {
	tmpLen, err := sender.read(buff, clientIndex, *offset, headLen)
	if err == nil {
		*offset += tmpLen
		//判断协议头是否合法
		if tmpLen != headLen || !sender.headValid(buff) {
			//fmt.Printf("协议头不合法,%x\n", buff[0:tmpLen])
			err = errors.New(fmt.Sprintf("%sprotocol head invalid[%x]", ERROR, buff[0:tmpLen]))
		}
	} else {
		err = errors.New("read head error: " + err.Error())
	}
	return
}

//功能：验证协议起止符是否匹配
//参数：
//	buff:缓存数据
//	headSize:头长度
//返回：是否匹配
func (sender *TcpServer) headValid(buff []byte) (result bool) {
	headLen := sender.protoFormat.headLen
	result = true
	for i := 0; i < headLen; i++ {
		if buff[i] != sender.protoFormat.Head[i] {
			result = false
			break
		}
	}
	return
}

//=====================================================读取至数据结止符====================================================

func (sender *TcpServer) readToEnding(buff []byte, clientIndex int, offset *int, headLen int) (err error) {
	var tmpLen int
	for true {
		if ok := sender.endValid(buff, *offset); !ok { //没有接收到协议尾,断续接收
			if *offset < len(buff) {

				sender.clientList[clientIndex].Conn.SetReadDeadline(time.Now().Add(1 * time.Second))

				tmpLen, err = sender.read(buff, clientIndex, *offset, 1)
				if err != nil {
					err = errors.New("read ending error: " + err.Error())
					break
				}
				*offset += tmpLen
			} else {
				//数据缓存益出
				err = errors.New("read ending error: data buffer out")
				break
			}
		} else {
			break
		}

	}
	return
}

//功能：验证协议结束符是否匹配
//参数：
//	buff:缓存数据
//	dataSize:buff数据长度
//返回：是否匹配
func (sender *TcpServer) endValid(buff []byte, dataSize int) (result bool) {
	headLen := sender.protoFormat.headLen
	endLen := sender.protoFormat.endLen
	result = true
	if (dataSize - headLen) > endLen {
		for i := 0; i < endLen; i++ {
			if buff[dataSize-endLen+i] != sender.protoFormat.End[i] {
				result = false
				break
			}
		}
	} else {
		result = false
	}
	return
}

//=====================================================================================================================

//功能:启动连接数量检测
func (sender *TcpServer) countConnSum() {
	go func() {
		for connChange := range sender.chConnChange {
			sender.connSum += connChange
		}
	}()
}

//功能:初始化连接池
//func (sender *TcpServer) initConnPool() {
//	for i := 0; i < sender.maxConn; i++ {
//		go func() {
//			for client := range sender.chConn {
//
//				go sender.sendConnEvent(client.Addr, client.Index, true) //发送连接事件
//				sender.receiveData(client.Index)                         //读取数据
//
//			}
//		}()
//	}
//}

func (sender *TcpServer) startReadData(client ClientInfo) {
	sender.sendConnEvent(client.Addr, client.Index, true) //发送连接事件

	//for true {
	sender.chConnChange <- +1
	sender.receiveData(client.Index)
	sender.chConnChange <- -1
	//}

}

//功能:发送连接状态事件
func (sender *TcpServer) sendConnEvent(addr string, clientIndex int, value bool) {
	if sender.connCallback != nil {
		sender.connCallback(addr, clientIndex, value)
	}
}

//功能:发送接收到有效数据事件
func (sender *TcpServer) sendReceiveEvent(data []byte, size int, client ClientInfo) {
	if sender.dataCallback != nil {
		sender.dataCallback(data, size, client)
	}
}

//功能：接受连接
func (sender *TcpServer) accept(listen net.Listener) {

	for sender.isListen {

		conn, err := sender.listen.Accept()
		if err == nil {
			idx := sender.findFreeIndex()
			if idx >= 0 {
				sender.clientList[idx].Conn = conn
				sender.clientList[idx].Index = idx
				sender.clientList[idx].Addr = conn.RemoteAddr().String()
				sender.clientList[idx].LastPackageTime = time.Now().Unix()
				sender.clientList[idx].Connected = true
				//fmt.Println("连接")
				go sender.startReadData(sender.clientList[idx])
			} else {
				fmt.Println("connect pool filled")
				conn.Close()
			}
		} else {
			fmt.Println("accept connect err:", err)
		}
	}
}

//功能:在客户端连接池中找到没有占用的位置
func (sender *TcpServer) findFreeIndex() (clientIndex int) {
	clientIndex = -1
	for i := range sender.clientList {
		if sender.clientList[i].Conn == nil {
			clientIndex = i
			break
		}
	}
	return
}

//功能:关闭连接
func (sender *TcpServer) close(clientIndex int) {

	if sender.clientList[clientIndex].Conn != nil {

		sender.sendConnEvent(sender.clientList[clientIndex].Addr, clientIndex, false)

		sender.clientList[clientIndex].Conn.Close()
		sender.clientList[clientIndex].Connected = false
		sender.clientList[clientIndex].Conn = nil
		sender.clientList[clientIndex].Addr = ""
		sender.clientList[clientIndex].ErrSum = 0
		sender.clientList[clientIndex].LastPackageTime = 0
	}

}

//功能:检测连接有效性
func (sender *TcpServer) checkConn() {
	count := 0
	for sender.isCheck {
		count = 0
		for i := range sender.clientList {
			if sender.clientList[i].Conn != nil {
				sec := time.Now().Unix() - sender.clientList[i].LastPackageTime
				//fmt.Println(sec)
				count++
				if sec > sender.connTimeout {
					fmt.Println("no heartheat packag disconnect")
					sender.close(sender.clientList[i].Index)
				}
			}
		}
		time.Sleep(time.Duration(sender.checkInterval) * time.Second)
	}
}
