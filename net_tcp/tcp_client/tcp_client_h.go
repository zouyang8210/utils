package tcp_client

import "net"

type TcpClient struct {
	conn net.Conn //连接对像
	ip   string   //连接的IP地址
	port int      //连接端口
}
