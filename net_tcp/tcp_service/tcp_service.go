package tcp_service

import (
	"errors"
	"fmt"
	"net"
)

//功能:监听
//参数:
//	port:监听端口号
//返回:是否监听成功
func (sender *TcpServer) Listen(port int) (result bool, err error) {
	sender.listen, err = net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err == nil {
		result = true
		sender.init()
		sender.accept(sender.listen)
	}
	return
}

//开启或关闭连接超时检测
func (sender *TcpServer) SetCheckConn(v bool) {
	sender.checkConn = v
}

//开启或关闭是否登录检测
func (sender *TcpServer) SetCheckLogin(v bool) {
	sender.checkLogin = v
}

//设置客户为已登录状态
func (sender *TcpServer) SetLogin(clientIndex int) {
	sender.clientList[clientIndex].IsLogin = true
}

//功能:设置接收数据回调函数
func (sender *TcpServer) SetDataCallback(callback ReceiveEvent) {
	sender.dataCallback = callback
}

//功能:设置连接信息回调函数
func (sender *TcpServer) SetConnectCallback(callback ConnectEvent) {
	sender.connCallback = callback
}

//功能:设置最大连接数量
func (sender *TcpServer) SetMaxConnect(size int) {
	sender.maxConn = size
}

//功能:设置协议格式
func (sender *TcpServer) SetProtocolFormat(protoFormat ProtocolFormat) {
	sender.protoFormat = protoFormat
	sender.protoFormat.headLen = len(sender.protoFormat.Head)
	sender.protoFormat.endLen = len(sender.protoFormat.End)
}

//功能:获取当前连接总数
func (sender *TcpServer) GetConnectCount() (count int) {
	count = sender.connSum
	return
}

//功能：发送数据
//参数：
//	buff:发送数据缓存
//	clientIndex:接收数据客户端索引
//	offset:buff的偏移量
//	size:接收数据大小
//返回：实际接发送数据大小，错误信息
func (sender *TcpServer) WriteData(buff []byte, clientIndex int, offset, size int) (count int, err error) {
	if len(buff) < offset+size {
		return 0, errors.New("buffer cache is too small")
	}
	if sender.clientList[clientIndex].Conn != nil {
		count, err = sender.clientList[clientIndex].Conn.Write(buff[offset:size])
		if err != nil {
			fmt.Println("write data error: ", err)
			sender.close(clientIndex)

		}
	} else {
		err = errors.New("client is nil")
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
func (sender *TcpServer) ReadData(buff []byte, clientIndex, index, size int) (count int, err error) {
	count, err = sender.read(buff, clientIndex, index, size)
	return
}

//功能:关闭一个客户端连接
//参数:
//	index:连接池中的索引号
func (sender *TcpServer) CloseOne(index int) {
	sender.close(index)
}

//功能:关闭所有连接
func (sender *TcpServer) CloseAll() {
	for i := 0; i < sender.maxConn; i++ {
		if sender.clientList[i].Conn != nil {
			sender.close(sender.clientList[i].Index)
		}
	}
}
