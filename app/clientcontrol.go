package app

import (
	"encoding/json"
	"errors"
	"io"
	"sync"

	gg "github.com/haozibi/gglog"
)

var connection *Conn = nil

func ControlClient(client *Client, wait *sync.WaitGroup) {
	defer wait.Done()
	// 向 server 上注册 APP
	c, err := registerApp(client)
	if err != nil {
		gg.Errorf("app [%v] register client error,%v\n", client.Name, err)
		return
	}

	connection = c
	defer connection.Close()

	for {
		content, err := connection.ReadLine()
		if err == io.EOF || connection == nil || connection.IsClosed() {
			gg.Debugf("app [%v] server close this control conn", client.Name)
		} else if err != nil {
			gg.Infof("app [%v] read from server error, %v\n", client.Name, err)
		}
		clientCtlResponse := new(ClientControlResponse)
		if err := json.Unmarshal([]byte(content), clientCtlResponse); err != nil {
			gg.Infof("app [%v] parse error,%v\n", err)
			continue
		}

		client.StartTunnel()
	}
}

// 向服务器注册 app
func registerApp(client *Client) (conn *Conn, err error) {
	conn, err = DialServer(ClientServerIP, ClientServerPort)
	if err != nil {
		gg.Errorf("app [%v] register client error,%v\n", client.Name, err)
		return
	}
	req := ClientControlRequest{
		Type:    ControlConn,
		AppName: client.Name,
		Passwd:  client.Passwd,
	}
	buf, _ := json.Marshal(req)
	err = conn.Write(string(buf) + "\n")
	if err != nil {
		gg.Errorf("app [%v] write to server error,%v\n", client.Name, err)
		return
	}
	responseTmp, err := conn.ReadLine()
	if err != nil {
		gg.Errorf("app [%v] read from server error,%v\n", client.Name, err)
		return
	}
	gg.Debugf("app [%v] read from server => %v\n", client.Name, responseTmp)
	reponse := new(ClientControlResponse)
	if err = json.Unmarshal([]byte(responseTmp), reponse); err != nil {
		gg.Errorf("app [%v] unmarshal server response error,%v\n", client.Name, err)
		return
	}
	if reponse.Code != 0 {
		gg.Errorf("app [%v] start app error,%v\n", client.Name, reponse.Msg)
		return conn, errors.New(reponse.Msg)
	}

	// 注册成功，心跳包
	return conn, nil
}
