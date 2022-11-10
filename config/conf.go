package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"socketServerFrame/iface"
	"socketServerFrame/logs"
)

var configPath string // 配置的文件夹路径

type GlobalObj struct {
	Debug            bool          // 是否Debug模式
	TcpServer        iface.IServer // TCP全局对象
	Host             string        // 当前服务主机IP
	TcpPort          string        // 当前服务端口
	Name             string        // 当前服务名称
	Version          string        // 当前服务版本号
	MaxPackSize      uint32        // 传输数据包最大值
	MaxConn          int           // 当前服务允许的最大连接数
	WorkerPoolSize   uint32        // work池大小
	WorkerTaskMaxLen uint32        // work对应的执行队列内任务数量的上限
	MaxMsgChanLen    int           // 读写消息的通道最大缓冲数
}

var globalObject *GlobalObj

func init() {
	globalObject = &GlobalObj{
		Debug:            false,
		TcpServer:        nil,
		Host:             "127.0.0.1",
		TcpPort:          "7777",
		Name:             "socketServerFrame",
		Version:          "v0.1",
		MaxPackSize:      4096,
		MaxConn:          10,
		WorkerPoolSize:   3,
		WorkerTaskMaxLen: 1024,
		MaxMsgChanLen:    100,
	}

	globalObject.Reload()
	logs.SetPrintMode(globalObject.Debug)

	str, _ := json.Marshal(globalObject)
	log.Println(fmt.Sprintf("服务配置参数：%v", string(str)))
}

// GetGlobalObject 获取全局配置对象
func GetGlobalObject() GlobalObj {
	return *globalObject
}

func (o *GlobalObj) Reload() {
	err := json.Unmarshal(getConfigDataToBytes("config.json"), &globalObject)
	logs.PrintLogErrToConsole(err)
}

// 获取配置数据到字节
func getConfigDataToBytes(configName string) []byte {
	if configPath == "" {
		configPath = os.Getenv("GOPATH") + "/src/" + globalObject.Name + "/config/"
	}

	bytes, err := ioutil.ReadFile(configPath + "./" + configName)
	logs.PrintLogPanicToConsole(err)
	return bytes
}
