package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"socketServerFrame/iface"
)

type GlobalObj struct {
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
}

// GetGlobalObject 获取全局配置对象
func GetGlobalObject() GlobalObj {
	return *globalObject
}

func (o *GlobalObj) Reload() {
	bytes, err := ioutil.ReadFile("./config.json")
	if err != nil {
		println(err.Error())
	}
	err = json.Unmarshal(bytes, &globalObject)
	if err != nil {
		println(err.Error())
	}

	if err != nil {
		panic(fmt.Sprintf("配置文件加载失败：%s", err.Error()))
	}
}
