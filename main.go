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

	app.InitLog()

	gg.SetOutLevel(app.CONF.Log.LogLevel)
	gg.SetOutSimple(app.CONF.Log.LogSimple)
	gg.SetLogDir(app.CONF.Log.LogPath)

	gg.Debugf("\n%v\n", app.CONF)
	fmt.Printf("\n%v\n", logo)
}

var logo = `  _   _ ___ ___  
 | \ | |_ _/ _ \ 
 |  \| || | | | |
 | |\  || | |_| |
 |_| \_|___\___/  ` + version
