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
	ControlConnType = 0
	WorkConnType    = 1
	IdleType        = 0
	WorkingType     = 1
	ErrorType       = 1
)

var (
	userConnTimeOut = CONF.Common.UserConnTimeout
)

func InitLog() {
	if !utils.PathExists(CONF.Log.LogPath) && CONF.Log.LogWay == "file" {
		os.MkdirAll(CONF.Log.LogPath, 0777)
	}
}
