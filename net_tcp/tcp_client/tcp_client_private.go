package tcp_client

//func (sender *TcpClient) receiveData() {
//	defer sender.Close()
//	buff := make([]byte, 1024)
//	for {
//		_, err := sender.conn.Read(buff)
//		if err != nil {
//			fmt.Println("receive data err:", err)
//			break
//		}
//	}
//}

////功能:按协议格式读取数据
////参数：
////	clientIndex:接收数据客户端索引
////返回：符合协议的完整数据包，错误信息
//func (sender *TcpClient) readDataFormat() (buff []byte, err error) {
//	buff = make([]byte, 1024)
//	offset := 0
//	ok := false
//	if sender.dataFormat.endLen == 0 && sender.dataFormat.headLen == 0 {
//		return nil, errors.New("protocol format invalid,please correct set protoFormat")
//	}
//	//是否有协议头
//	if sender.dataFormat.headLen > 0 {
//		tmpLen, err := sender.ReadData(buff, offset, sender.dataFormat.headLen)
//		if err == nil {
//			offset += tmpLen
//			//判断协议头是否合法
//			if tmpLen != sender.dataFormat.headLen || !sender.headValid(buff) {
//				fmt.Printf("协议头不合法,%x\n", buff[0:tmpLen])
//				return nil, errors.New(fmt.Sprintf("protocol head invalid[%x]", buff[0:tmpLen]))
//			}
//		} else {
//			return nil, err
//		}
//	}
//
//	if sender.dataFormat.endLen > 0 {
//		for true {
//			if ok = sender.endValid(buff, offset); !ok { //没有接收到协议尾,断续接收
//				if offset < 1024 {
//					tmpLen, err := sender.ReadData(buff, offset, 1)
//					if err != nil {
//						return nil, err
//					}
//					offset += tmpLen
//				} else {
//					//数据缓存益出
//					return nil, errors.New("data buffer out")
//				}
//			} else {
//				//数据格式验证合格，作为一个全整的包返回。跳出读取数据
//				break
//			}
//		}
//	}
//	if ok {
//		return buff[0:offset], nil
//	} else {
//		fmt.Printf("错误数据:%x", buff[0:offset])
//		return nil, errors.New(fmt.Sprintf("data error [%x]", buff[0:offset]))
//	}
//}
//
////功能：验证协议起止符是否匹配
////参数：
////	buff:缓存数据
////	headSize:头长度
////返回：是否匹配
//func (sender *TcpClient) headValid(buff []byte) (result bool) {
//	headLen := sender.dataFormat.headLen
//	for i := 0; i < headLen; i++ {
//		if buff[i] != sender.dataFormat.Head[i] {
//			result = false
//			return
//		}
//	}
//	result = true
//	return
//}
//
////功能：验证协议结束符是否匹配
////参数：
////	buff:缓存数据
////	dataSize:buff数据长度
////返回：是否匹配
//func (sender *TcpClient) endValid(buff []byte, dataSize int) (result bool) {
//	headLen := sender.dataFormat.headLen
//	endLen := sender.dataFormat.endLen
//	if (dataSize - headLen) > endLen {
//		for i := 0; i < endLen; i++ {
//			if buff[dataSize-endLen+i] != sender.dataFormat.End[i] {
//				return
//			}
//		}
//	} else {
//		return
//	}
//	result = true
//	return
//}
