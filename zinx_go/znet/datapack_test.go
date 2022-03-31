package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

//只是负责测试datapack拆包 封包的单元测试
func TestDataPack(t *testing.T) {
	//模拟的服务器
	//1 创建socketTCP
	listenner, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("serve listen err:", err)
	}

	//创建一个go承载负责从客户端处理业务
	go func() {
		//2.c从客户端读取数据，拆包处理
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("server accept error", err)
			}
			go func(conn net.Conn) {
				//处理客户端请求
				//拆包的过程
				//定义一个拆包的对象dp
				dp := NewDataPack()
				for {
					//1.第一次从conn读 把包的head读出来
					headData := make([]byte, dp.GetHeadLen())
					if _, err := io.ReadFull(conn, headData); err != nil {
						fmt.Println("read head error")
						break
					}
					//2.第二次读从conn读，根据head中dataLen再读取data内容
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err", err)
						break
					}
					if msgHead.GetMsgLen() > 0 {
						//msg是有数据的，需要进行第二次读取
						//2 第二次从conn读，根据head中的dataLen再读取内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err:", err)
							break
						}
						//完整的一个消息已经读取完毕
						fmt.Println("---> Recv MsgID: ", msg.Id, ", dataLen = ", msg.DataLen, "data = ", string(msg.Data))
					}
				}

			}(conn)
		}
	}()

	//模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}
	//创建一个封包对象 dp
	dp := NewDataPack()

	//模拟粘包过程，封装两个msg一同发送
	//封装第一个msg1包
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte("zinx"),
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error", err)
		return
	}
	//封装第二个msg2包
	msg2 := &Message{
		Id:      2,
		DataLen: 8,
		Data:    []byte("helloboy"),
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 error", err)
		return
	}
	msg3 := &Message{
		Id:      3,
		DataLen: 5,
		Data:    []byte("hello"),
	}
	sendData3, err := dp.Pack(msg3)
	if err != nil {
		fmt.Println("client pack msg2 error", err)
		return
	}
	sendData1 = append(sendData1, sendData2...)
	sendData1 = append(sendData1, sendData3...)
	//sendData1 = append(sendData3, sendData4...)
	//将3个包粘在一起
	for {

		//一次性发送给服务端
		conn.Write(sendData1)

		time.Sleep(3 * time.Second)
	}

}
