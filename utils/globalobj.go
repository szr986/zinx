package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"example.com/m/ziface"
)

// 存储Zinx框架全局参数，供其他模块使用
type Globalobj struct {
	// Server
	TcpServer ziface.IServer //当前Zinx全局的Server对象
	Host      string         //当前服务器主机监听的IP
	TcpPort   int            //当前服务器监听的端口号
	Name      string         //服务器名称

	// Zinx
	Version        string //zinx版本号
	Maxconn        int    //允许的最大连接数
	MaxPackageSize uint32 //数据包的最大值
}

// 定义一个全局的对外对象
var GlobalObject *Globalobj

// 从zinx.json加载自定义参数
func (g *Globalobj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		fmt.Println("read zinx.json failed,err:", err)
		return
	}
	// 将json文件数据解析到struct中
	err = json.Unmarshal(data, &g)
	if err != nil {
		return
	}

}

// 提供一个init方法，初始化当前的GlobalObject对象
func init() {
	// 如果配置文件没有加载时的默认值
	GlobalObject = &Globalobj{
		Name:           "ZinxServerApp",
		Version:        "V0.5",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		Maxconn:        1000,
		MaxPackageSize: 4096,
	}

	// 尝试从conf/zinx.json加载自定义参数
	GlobalObject.Reload()
}
