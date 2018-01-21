package app

import (
	"encoding/json"
	"strconv"

	gg "github.com/haozibi/gglog"
)

type Client struct {
	Name      string
	LocalPort int64
	Passwd    string
}

var (
	ClientServerIP   string
	ClientServerPort int64
	ClientLocalIP    = "127.0.0.1"
)

var Clients map[string]*Client

func InitClient() {
	Clients = make(map[string]*Client)
	if len(CONF.Client.ServerIP) == 0 || len(CONF.Client.ServerPort) == 0 || len(CONF.App) == 0 {
		panic(ErrorConf)
	}
	ClientServerIP = CONF.Client.ServerIP
	ClientServerPort, _ = strconv.ParseInt(CONF.Client.ServerPort, 10, 64)
	for _, v := range CONF.App {
		client := new(Client)

		client.Name = v.Name
		client.Passwd = v.Passwd
		client.LocalPort, _ = strconv.ParseInt(v.LocalPort, 10, 64)

		Clients[v.Name] = client
	}
}

func (c *Client) StartTunnel() (err error) {
	err = nil
	localConn, err := c.GetLocalConn()
	if err != nil {
		return
	}
	remoteConn, err := c.GetRemoteConn()
	if err != nil {
		return
	}
	gg.Debugf("[nio] join two conn\n")
	go Join(localConn, remoteConn)
	return
}

func (c *Client) GetLocalConn() (conn *Conn, err error) {
	conn, err = DialServer(ClientLocalIP, c.LocalPort)
	if err != nil {
		gg.Errorf("[nio] app [%v] connect to local error,%v\n", err)
	}
	return
}

func (c *Client) GetRemoteConn() (conn *Conn, err error) {
	err = nil
	// 如果有错误则关闭连接
	defer func() {
		if err != nil {
			conn.Close()
		}
	}()

	conn, err = DialServer(ClientServerIP, ClientServerPort)
	if err != nil {
		gg.Errorf("[nio] app [%v] connect to remote[%v:%v] error,%v\n", c.Name, ClientServerIP, ClientServerPort, err)
		return
	}

	resquest := &ClientControlRequest{
		Type:    WorkConn,
		AppName: c.Name,
		Passwd:  c.Passwd,
	}

	buf, _ := json.Marshal(resquest)
	err = conn.Write(string(buf) + "\n")
	if err != nil {
		gg.Errorf("[nio] app [%v] write to remote error,%v", err)
		return
	}
	return
}
