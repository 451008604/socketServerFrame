package base

import (
	"fmt"
	"net"
	"os"
	"sync"
)

// CustomConnect 自定义连接
type CustomConnect struct {
	net.Conn
	address   string // 服务地址
	port      string // 服务端口
	bufferLen int    // 消息缓冲区长度
}

// 主线程锁
var wg *sync.WaitGroup

var conn *CustomConnect

var messageWhenConnectionIsClosed = make([][]byte, 0)

// NewConnection 新建连接
func NewConnection(address, port string) *CustomConnect {
	// 与服务器请求连接
	serverAddress := address + ":" + port
	dial, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println(fmt.Sprintf("服务器连接失败：%v \n尝试重连中...\n", serverAddress))
		return NewConnection(conn.address, conn.port)
	}

	// 关闭旧的连接
	if conn != nil {
		_ = conn.Close()
	}
	// 创建新的连接
	conn = &CustomConnect{}
	conn.Conn = dial
	conn.address = address
	conn.port = port
	conn.bufferLen = 512

	// 阻塞主进程
	wg = &sync.WaitGroup{}
	wg.Add(1)
	// 监听服务器返回的消息
	go func(conn *CustomConnect) {
		wg.Done()
		for true {
			receiveData := receiveMsg()
			if receiveData == nil {
				return
			}

			// 服务器返回的消息
			fmt.Printf(string(receiveData))
		}
	}(conn)
	wg.Wait()
	wg.Add(1)
	return conn
}

// SetBlocking 阻塞主进程，等待接受消息
func SetBlocking() {
	wg.Wait()
}

// Disconnect 断开连接，结束主进程
func _() {
	wg.Done()
	os.Exit(0)
}

// SendMsg 发送消息到服务器
func SendMsg(msg []byte) {
	if conn == nil {
		return
	}
	_, err := conn.Write(msg)
	if err != nil {
		messageWhenConnectionIsClosed = append(messageWhenConnectionIsClosed, msg)
		fmt.Println("write error err ", err)
	}
}

// receiveMsg 接收服务器消息
func receiveMsg() []byte {
	buf := make([]byte, conn.bufferLen)
	_, err := conn.Read(buf)
	// 现有连接发生错误时尝试重新与服务器建立连接
	if err != nil {
		fmt.Println(err.Error())
		NewConnection(conn.address, conn.port)
		return nil
	}
	return buf
}
