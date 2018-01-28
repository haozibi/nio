//nio function

package app

import (
	"os"
	"sync"

	gg "github.com/haozibi/gglog"
	"github.com/haozibi/nio/utils"
)

var (
	ErrorConf = "Config Error"
)

// connection type
const (
	ControlConnType = 0
	WorkConnType    = 1
	IdleType        = 0
	WorkingType     = 1
	ErrorType       = 1
	ClientHeartBeat = 99
	ServerHeartBeat = 100
)

var (
	userConnTimeOut   = CONF.Common.UserConnTimeout
	HeartBeatInterval = 3
	HeartBeatTimeout  = 30
)

var (
	IsServer = false
	IsDebug  = false
)

func InitLog() {
	if !utils.PathExists(CONF.Log.LogPath) && CONF.Log.LogWay == "file" {
		os.MkdirAll(CONF.Log.LogPath, 0777)
	}
	if len(CONF.Log.LogLevel) == 0 || len(CONF.Log.LogOutType) == 0 || len(CONF.Log.LogPath) == 0 || len(CONF.Log.LogWay) == 0 {
		panic("conf log error")
	}

	IsServer = CONF.Common.IsServer

}

func StartNio() {
	if !IsServer {
		// Client
		gg.Infof("start client\n")
		InitClient()
		var wait sync.WaitGroup
		wait.Add(len(CONF.App))
		gg.Infof("add %v app from config\n", len(CONF.App))
		for _, v := range Clients {
			go ControlClient(v, &wait)
		}
		wait.Wait()
		gg.Infof("all app exit !\n")
	}

	if IsServer {
		// Server
		gg.Infof("start server\n")
		InitServer()
		ControlServer()
	}
}
