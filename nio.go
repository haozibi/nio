//fast reverse proxy

package main

import (
	"flag"
	"fmt"
	"sync"

	gg "github.com/haozibi/gglog"
	"github.com/haozibi/nio/app"
)

var version = "v0.1"

func main() {
	flag.Parse()
	defer gg.Flush()

	app.InitLog()

	gg.SetOutLevel(app.CONF.Log.LogLevel)
	gg.SetOutType(app.CONF.Log.LogOutType)
	gg.SetLogDir(app.CONF.Log.LogPath)
	gg.SetPrefix("[nio] ")

	fmt.Printf("%v\n\n", logo)

	if !app.CONF.Common.IsServer {
		// Client
		gg.Infof("start client\n")
		app.InitClient()
		var wait sync.WaitGroup
		wait.Add(len(app.CONF.App))
		gg.Infof("add %v app\n", len(app.CONF.App))
		for _, v := range app.Clients {
			go app.ControlClient(v, &wait)
		}
		wait.Wait()
		gg.Infof("all app exit !\n")
	}

	if app.CONF.Common.IsServer {
		// Server
		gg.Infof("start server\n")
		app.InitServer()
		app.ControlServer()
	}
}

var logo = `  _   _ ___ ___  
 | \ | |_ _/ _ \ 
 |  \| || | | | |
 | |\  || | |_| |
 |_| \_|___\___/  ` + version
