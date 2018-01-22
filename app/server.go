package app

import (
	"container/list"
	"os"
	"strconv"
	"sync"

	gg "github.com/haozibi/gglog"
)

var (
	ServerBindIP   string
	ServerBindPort int64
)

type Server struct {
	Name           string
	Passwd         string
	BindIP         string
	ListenPort     int64
	Status         int64
	listenr        *Listener
	controlMsgChan chan int64
	clientConnChan chan *Conn
	userConnList   *list.List
	mutex          sync.Mutex
}

var Servers map[string]*Server = make(map[string]*Server)
var TCPl *Listener // server 接受用户行为的监听器

func InitServer() {
	if len(CONF.Server.BindIP) == 0 || len(CONF.Server.BindPort) == 0 || len(CONF.App) == 0 {
		panic(ErrorConf)
	}
	ServerBindIP = CONF.Server.BindIP
	ServerBindPort, _ = strconv.ParseInt(CONF.Server.BindPort, 10, 64)
	// 开始监听绑定端口
	// 监听 Client 的请求
	tcpl, err := Listen(ServerBindIP, ServerBindPort)
	if err != nil {
		gg.Errorf("[nio] create server listen error,%v", err)
		os.Exit(-1)
	}

	TCPl = tcpl

	for _, v := range CONF.App {
		s := new(Server)
		s.Name = v.Name
		s.Passwd = v.Passwd
		s.BindIP = v.BindIP
		s.ListenPort, _ = strconv.ParseInt(v.ListenPort, 10, 64)
		s.Status = Idle
		s.clientConnChan = make(chan *Conn)
		s.controlMsgChan = make(chan int64)
		s.userConnList = list.New()

		Servers[v.Name] = s
	}
}
