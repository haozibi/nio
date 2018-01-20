package app

import (
	"os"
	"strconv"

	"github.com/haozibi/nio/utils"
)

type Client struct {
	Name      string
	LocalPort int64
	Passwd    string
}

var (
	ClientServerIP  string
	ClentServerPort int64
)

var (
	ErrorConf = "Config Error"
)

var Clients map[string]*Client

func InitLog() {
	if !utils.PathExists(CONF.Log.LogPath) && CONF.Log.LogWay == "file" {
		os.MkdirAll(CONF.Log.LogPath, 0777)
	}
}

func initClient() {
	Clients = make(map[string]*Client)
	if len(CONF.Client.ServerIP) == 0 || len(CONF.Client.ServerPort) == 0 || len(CONF.App) == 0 {
		panic(ErrorConf)
	}
	ClientServerIP = CONF.Client.ServerIP
	ClentServerPort, _ = strconv.ParseInt(CONF.Client.ServerPort, 10, 64)
	for _, v := range CONF.App {
		client := new(Client)

		client.Name = v.Name
		client.Passwd = v.Passwd
		client.LocalPort, _ = strconv.ParseInt(v.LocalPort, 10, 64)

		Clients[v.Name] = client
	}
}
