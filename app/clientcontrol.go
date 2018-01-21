package app

import (
	"encoding/json"
	"io"
	"sync"

	gg "github.com/haozibi/gglog"
)

var connection *Conn = nil

func ControlClient(client *Client, wait *sync.WaitGroup) {
	defer wait.Done()
	c, err := DialServer(ClientServerIP, ClientServerPort)
	if err != nil {
		gg.Errorf("[nio] app [%v] dial server failed\n", client.Name)
		return
	}
	connection = c
	defer connection.Close()

	for {
		content, err := connection.ReadLine()
		if err == io.EOF || connection == nil || connection.IsClosed() {
			gg.Debugf("[nio] app [%v] server close this control conn", client.Name)
		} else if err != nil {
			gg.Infof("[nio] app [%v] read from server error, %v\n", client.Name, err)
		}
		clientCtlResponse := new(ClientControlResponse)
		if err := json.Unmarshal([]byte(content), clientCtlResponse); err != nil {
			gg.Infof("[nio] app [%v] parse error,%v\n", err)
			continue
		}

		client.StartTunnel()
	}
}
