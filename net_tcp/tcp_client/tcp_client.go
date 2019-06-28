package tcp_client

import (
	"errors"
	"fmt"
	"net"
	"time"
)

//功能:连接TCP服务器
//参数:
//	ip:服务器IP地址
//	port:服务器监听端口
//返回:错误信息
func (sender *TcpClient) Connect(ip string, port int) (err error) {
	sender.conn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), time.Duration(1)*time.Second)
	return
}

//功能:字符地址连接TCP服务器
//参数:
//	address:服务器IP地址+端口号
//返回:错误信息
func (sender *TcpClient) ConnectStr(address string) (err error) {
	sender.conn, err = net.DialTimeout("tcp", address, time.Duration(1)*time.Second)
	return
}

func (sender *TcpClient) SetReadTimeout(second int) (err error) {
	if sender.conn != nil {
		if second > 0 {
			err = sender.conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(second)))
		} else {
			err = sender.conn.SetReadDeadline(time.Time{})
		}
	} else {
		err = errors.New("socket is closed")
	}
	return
}

func (sender *TcpClient) SetWriteTimeout(second int) (err error) {
	if sender.conn != nil {
		if second > 0 {
			err = sender.conn.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(second)))
		} else {
			err = sender.conn.SetWriteDeadline(time.Time{})
		}
	} else {
		err = errors.New("socket is closed")
	}
	return
}

//功能:发送数据给服务器
//参数:
//	buff:数据缓存
//	index:开始发送的位置
//	size:发送数据大小
//返回:发送数据数量,错误信息
func (sender *TcpClient) SendData(buff []byte, index, size int) (count int, err error) {
	if sender.conn != nil {
		count, err = sender.conn.Write(buff[index : index+size])
	} else {
		err = errors.New("invalid connect")
	}
	return
}

//功能:写取服务器发送来的数据
//参数:
//	buff:数据缓存
//	index:读取数据到到buff的开始位置
//	size:读取数据大小
//返回:读取数据数量,错误信息
func (sender *TcpClient) ReadData(buff []byte, index, size int) (count int, err error) {
	if sender.conn != nil {
		count, err = sender.conn.Read(buff[index : index+size])
	} else {
		err = errors.New("invalid connect")
	}
	return
}

//功能:关闭TCP连接
//返回:是否成功
func (sender *TcpClient) Close() (result bool) {
	if sender.conn != nil {
		err := sender.conn.Close()
		if err == nil {
			result = true
		}
	} else {
		result = true
	}
	return
}
