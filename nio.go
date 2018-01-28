//fast reverse proxy

package main

import (
	"flag"
	"fmt"

	gg "github.com/haozibi/gglog"
	"github.com/haozibi/nio/app"
)

var version = "v0.1"

func main() {
	flag.Parse()
	defer gg.Flush()

	app.IsDebug = true

	app.InitLog()

	gg.SetOutLevel(app.CONF.Log.LogLevel)
	if app.IsDebug {
		gg.SetOutType(app.CONF.Log.LogOutType)
	} else {
		gg.SetOutType("SIMPLE")
	}
	gg.SetLogDir(app.CONF.Log.LogPath)
	gg.SetPrefix("[nio] ")

	fmt.Printf("%v\n\n", logo)

	app.StartNio()
}

var logo = `  _   _ ___ ___  
 | \ | |_ _/ _ \ 
 |  \| || | | | |
 | |\  || | |_| |
 |_| \_|___\___/  ` + version
