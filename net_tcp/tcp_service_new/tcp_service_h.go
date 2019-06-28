package tcp_service

import "net"

const (
	ERROR = "[ERROR]:"
	IOOUT = "i/o timeout"
)

type TcpServer struct {
	listen       net.Listener   //临听对像
	clientList   []ClientInfo   //客户端列表
	dataCallback ReceiveEvent   //接收数据回调
	connCallback ConnectEvent   //连接信息回调
	maxConn      int            //最大连接数
	connSum      int            //当前连接总数
	protoFormat  ProtocolFormat //协议格式
	chConnChange chan int       //连接数改变通道
	//chConn        chan ClientInfo //客户端通道
	isListen      bool  //监听开关
	isCheck       bool  //连接检测开关
	connTimeout   int64 //连接无数据包超时(秒)
	checkInterval int   //连接检测间隔(秒)

}

//连接服务的客户端信息
type ClientInfo struct {
	Conn            net.Conn //连接句柄
	Addr            string   //客户端地址
	Index           int      //在客户端列表中的索引
	ErrSum          int      //连接错误数据包的次数
	LastPackageTime int64    //最后有效数据的时间截
	Connected       bool     //是否连接
}

//协议格式
type ProtocolFormat struct {
	Head     []byte //头数据
	End      []byte //结尾数据
	headLen  int    //头长度
	endLen   int    //尾长度
	DataSize int    //接收数据量
}

//收到数据事件
type ReceiveEvent func(data []byte, dataSize int, clientInfo ClientInfo)

//连接事件
type ConnectEvent func(addr string, clientIndex int, value bool)
