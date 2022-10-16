package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"socketServerFrame/iface"
)

var configPath string // 配置的文件夹路径

type GlobalObj struct {
	ErrInfo     error
	TcpServer   iface.IServer // TCP全局对象
	Host        string        // 当前服务主机IP
	TcpPort     string        // 当前服务端口
	Name        string        // 当前服务名称
	Version     string        // 当前服务版本号
	MaxPackSize uint32        // 传输数据包最大值
	MaxConn     int           // 当前服务允许的最大连接数
}

var globalObject *GlobalObj

func init() {
	globalObject = &GlobalObj{
		TcpServer:   nil,
		Host:        "127.0.0.1",
		TcpPort:     "7777",
		Name:        "socketServerFrame",
		Version:     "v0.1",
		MaxPackSize: 4096,
		MaxConn:     3,
	}

	globalObject.Reload()

	log.Println(fmt.Sprintf("服务配置参数：%v", globalObject))
	if globalObject.ErrInfo != nil {
		panic(fmt.Sprintf("配置文件加载失败：%s", globalObject.ErrInfo.Error()))
	}
}

// GetGlobalObject 获取全局配置对象
func GetGlobalObject() GlobalObj {
	return *globalObject
}

func (o *GlobalObj) Reload() {
	o.ErrInfo = json.Unmarshal(getConfigDataToBytes("config.json"), &globalObject)
}

// 获取配置数据到字节
func getConfigDataToBytes(configName string) []byte {
	if configPath == "" {
		configPath = os.Getenv("GOPATH") + "/src/" + globalObject.Name + "/config/"
	}

	bytes, err := ioutil.ReadFile(configPath + "./" + configName)
	if err != nil {
		panic(fmt.Sprintf("配置文件加载失败：%s", err.Error()))
	}
	return bytes
}
