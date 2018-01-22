package app

import (
	"container/list"
	"os"
	"strconv"
	"sync"
	"time"

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
		s.Status = IdleType
		s.clientConnChan = make(chan *Conn)
		s.controlMsgChan = make(chan int64)
		s.userConnList = list.New()

		Servers[v.Name] = s
	}
}

func (s *Server) Lock() {
	s.mutex.Lock()
}

func (s *Server) UnLock() {
	s.mutex.Unlock()
}

func (s *Server) Start() (err error) {
	s.listenr, err = Listen(s.BindIP, s.ListenPort)
	if err != nil {
		return err
	}
	s.Status = WorkingType

	// 监听用户发送的请求
	go func() {
		for {
			c, err := s.listenr.GetConn()
			if err != nil {
				gg.Errorf("app [%v] listenr is closed\n", s.Name)
				return
			}
			gg.Infof("app [%v] get one new user conn,%v\n", s.Name, c.GetRemoteAddr())
			s.Lock()
			if s.Status != WorkingType {
				gg.Debugf("app [%v] not working,user conn close", s.Name)
				c.Close() // 只是关闭了用户连接
				s.UnLock()
				return
			}
			s.userConnList.PushBack(c)
			s.UnLock()

			s.controlMsgChan <- 1

			// timeout
			time.AfterFunc(time.Duration(userConnTimeOut)*time.Second, func() {
				s.Lock()
				defer s.UnLock()
				ele := s.userConnList.Front()
				if ele == nil {
					return
				}
				userConn := ele.Value.(*Conn)
				if userConn == c {
					gg.Errorf("app [%v] user[ %v] conn time out\n", s.Name, c.GetRemoteAddr())
				}
			})
		}
	}()

	// 用户conn 与客户端的进行交换
	go func() {
		for {
			clientConn, ok := <-s.clientConnChan
			if !ok {
				return
			}
			s.Lock()
			ele := s.userConnList.Front()

			var userConn *Conn
			if ele != nil {
				userConn = ele.Value.(*Conn)
				s.userConnList.Remove(ele)
			} else {
				clientConn.Close()
				s.UnLock()
				// 因为已经是空，所以不需要close
				continue
			}
			s.UnLock()

			// 开始交换
			go Join(clientConn, userConn)
		}
	}()

	return nil
}

func (s *Server) GetNewClientConn(conn *Conn) {
	s.clientConnChan <- conn
}

func (s *Server) WaitUserConn() (closeFlag bool) {
	closeFlag = false

	// start() 中当获得 用户conn 时会向 controlMsgChan 《- 1
	_, ok := <-s.controlMsgChan
	if !ok {
		closeFlag = true
	}
	return
}

func (s *Server) Close() {
	s.Lock()
	s.Status = IdleType
	s.listenr.Close()
	close(s.clientConnChan)
	close(s.controlMsgChan)
	s.userConnList = list.New()
	s.UnLock()
}
