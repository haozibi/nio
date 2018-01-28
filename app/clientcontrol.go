package app

import (
	"encoding/json"
	"errors"
	"io"
	"sync"
	"time"

	gg "github.com/haozibi/gglog"
)

var connection *Conn = nil
var heartBeatTimer *time.Timer = nil

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
			for {
				tmpConn, err := registerApp(client)
				if err == nil {
					// 断开重新连接
					connection.Close()
					connection = tmpConn
					break
				}
				time.Sleep(2 * time.Second)
			}
			continue
		} else if err != nil {
			gg.Infof("app [%v] read from server error, %v\n", client.Name, err)
			continue
		}

		clientCtlResponse := new(ClientControlResponse)
		if err := json.Unmarshal([]byte(content), clientCtlResponse); err != nil {
			gg.Infof("app [%v] unmarshal server response error,%v\n", err)
			continue
		}

		// 这是一个心跳包
		if clientCtlResponse.GeneraResponse.Code == ServerHeartBeat {
			// 说明是已经注册成功
			if heartBeatTimer != nil {
				heartBeatTimer.Reset(time.Duration(HeartBeatTimeout) * time.Second)
			} else {
				gg.Errorf("heartBeatTimer is nil\n")
			}
			continue
		}
		// Client 接到 server 的开始工作信息
		// 只有开始工作信息才开始工作
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
		Type:    ControlConnType,
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
	go startHeartbeat(conn, client.Name)

	gg.Debugf("app [%v] connect server success\n", client.Name)
	return conn, nil
}

func startHeartbeat(conn *Conn, name string) {
	f := func() {
		gg.Errorf("[%v] heart beat timeout\n", name)
		if conn != nil {
			conn.Close()
		}
	}
	heartBeatTimer = time.AfterFunc(time.Duration(HeartBeatTimeout)*time.Second, f)
	defer heartBeatTimer.Stop()

	tmpRequest := &ClientControlRequest{
		Type:    ClientHeartBeat,
		AppName: "",
		Passwd:  "",
	}
	request, err := json.Marshal(tmpRequest)
	if err != nil {
		gg.Errorf("marshal error,%v\n", err)
	}
	gg.Debugf("[%v] start heart beat\n", name)
	for {
		time.Sleep(time.Duration(HeartBeatInterval) * time.Second)
		if conn != nil && !conn.IsClosed() {
			err = conn.Write(string(request) + "\n")
			// gg.Infof("send heart beat\n")
			if err != nil {
				gg.Errorf("send heart beat to server error,%v\n", err)
				continue
			}
		} else {
			break
		}
	}
	gg.Debugf("[%v] heart beat exit\n", name)
}
