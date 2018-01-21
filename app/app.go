package app

import (
	"os"

	"github.com/haozibi/nio/utils"
)

var (
	ErrorConf = "Config Error"
)

// connection type
const (
	ControlConn = 0
	WorkConn    = 1
)

func InitLog() {
	if !utils.PathExists(CONF.Log.LogPath) && CONF.Log.LogWay == "file" {
		os.MkdirAll(CONF.Log.LogPath, 0777)
	}
}
