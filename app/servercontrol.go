package app

import (
	gg "github.com/haozibi/gglog"
)

func ControlServer() {
	for {
		c, err := TCPl.GetConn()
		if err != nil {
			return
		}
		gg.Debugf("[nio] new conn => %v\n", c.GetRemoteAddr())
		go controlServerApp(c)
	}
}

// 处理每个 Client 的 Conn
func controlServerApp(conn *Conn) {
	// 第一个 Conn 是 Client 的注册信息
	// 处理 APP 注册信息
	// todo
}
