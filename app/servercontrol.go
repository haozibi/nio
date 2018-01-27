package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	gg "github.com/haozibi/gglog"
)

func ControlServer() {
	for {
		c, err := TCPl.GetConn()
		if err != nil {
			return
		}
		gg.Debugf("[nio] new conn => %v\n", c.GetRemoteAddr())
		go controlServerApp(c)
	}
}

// 处理每个 Client 的 Conn
func controlServerApp(conn *Conn) {
	// 处理 APP 注册信息
	resquest, err := conn.ReadLine()
	if err != nil {
		gg.Errorf("conn read error,%v\n", err)
		return
	}

	clientResquest := new(ClientControlRequest)
	if err := json.Unmarshal([]byte(resquest), clientResquest); err != nil {
		gg.Errorf("unmarshal err,%v ==> %v\n", err, resquest)
		return
	}

	// 检查 app 信息
	clientResponse := new(ClientControlResponse)
	info, err := checkApp(clientResquest, conn)
	clientResponse.Code = 0
	if err != nil {
		clientResponse.Code = ErrorType
	}

	// 当是 App 注册信息的时候，返回Client一个确认信息
	// 如果是 Client work conn 则直接不回复，直接返回
	// len(info)==0 说明是要回复 app 注册
	if len(info) == 0 {
		defer conn.Close()
		clientResponse.Msg = "hello nio"
		buf, _ := json.Marshal(clientResponse)
		err = conn.Write(string(buf) + "\n")
		if err != nil {
			gg.Errorf("register write response error,%v\n", err)
			// 神奇
			time.Sleep(1 * time.Second)
			return
		}
	} else {
		// import
		// work conn, just return
		return
	}

	s, ok := Servers[clientResquest.AppName]
	if !ok {
		gg.Errorf("app [%v] not exist\n", clientResquest.AppName)
		return
	}

	go readClientHeartBeat(s, conn)

	serverRequest := new(ClientControlRequest)
	serverRequest.Type = WorkConnType
	for {
		// 等待用户连接，一连接则发送工作信息给客户端
		closeFlag := s.WaitUserConn()
		if closeFlag {
			gg.Debugf("app [%v] user conn is closed\n", s.Name)
			break
		}
		buf, _ := json.Marshal(serverRequest)
		err = conn.Write(string(buf) + "\n")
		if err != nil {
			gg.Errorf("app [%v] write to clien error,%v\n", s.Name, err)
			s.Close()
			return
		}
		gg.Debugf("app [%v] write to client to add work conn success", s.Name)
	}
	gg.Infof("app [%v] over", s.Name)
	return
}

// 根据 app 发送的信息检查 app
func checkApp(request *ClientControlRequest, conn *Conn) (info string, err error) {
	s, ok := Servers[request.AppName]
	if !ok {
		info = fmt.Sprintf("app [%v] not exist", request.AppName)
		gg.Errorln(info)
		return info, errors.New(info)
	}
	if request.Passwd != s.Passwd {
		info = fmt.Sprintf("app [%v] passwd not correct", request.AppName)
		gg.Errorln(info)
		return info, errors.New(info)
	}

	if request.Type == ControlConnType {
		// app 第一次的注册连接,并回复client
		if s.Status != IdleType {
			info = fmt.Sprintf("app [%v] already start", request.AppName)
			return info, errors.New(info)
		}

		// 启动 app, 监听端口
		err = s.Start()

		if err != nil {
			info = fmt.Sprintf("app [%v] start error,%v", request.AppName, err)
			gg.Errorln(info)
			return info, err
		}
		info = "hello haozibi"
		gg.Infof("app [%v] start success\n", request.AppName)
		return "", nil
	} else if request.Type == WorkConnType {
		// 正常连接
		if s.Status != WorkingType {
			gg.Errorf("app [%v] not working", request.AppName)
			return
		}
		s.GetNewClientConn(conn)
		return "not need response", nil
	} else {
		info = fmt.Sprintf("app [%v] type [%v] unknow", request.AppName, request.Type)
		gg.Errorln(info)
		return info, errors.New(info)
	}
	return "", nil
}

// 接受 client 的心跳包
func readClientHeartBeat(server *Server, conn *Conn) {
	flag := true
	var err error
	var content string
	var response = &ClientControlResponse{}
	f := func() {
		flag = false
		server.Close()
		gg.Errorf("app [%v], client heartbeat timeout\n", server.Name)
	}
	timer := time.AfterFunc(time.Duration(HeartBeatTimeout)*time.Second, f)
	defer timer.Stop()

	for flag {
		content, err = conn.ReadLine()
		if err != nil {
			if err == io.EOF {
				gg.Errorf("app [%v] client over\n", server.Name)
				server.Close()
				break
			} else if conn == nil || conn.IsClosed() {
				gg.Errorf("app [%v] client connect close\n", server.Name)
				break
			}
			gg.Errorf("app [%v] read error,%v\n", server.Name, err)
			continue
		}
		request := new(ClientControlRequest)
		if err := json.Unmarshal([]byte(content), request); err != nil {
			gg.Errorf("heart beat unmarshal error,%v\n", err)
			continue
		}
		if request.Type == ClientHeartBeat {
			timer.Reset(time.Duration(HeartBeatTimeout) * time.Second)
			response.GeneraResponse.Code = ServerHeartBeat
			responseMsg, err := json.Marshal(response)
			if err != nil {
				gg.Errorf("marshal heart beat response error,%v\n", err)
				continue
			}
			err = conn.Write(string(responseMsg) + "\n")
			if err != nil {
				gg.Errorf("send heart beat response error,%v\n", err)
				continue
			}
		}
	}
}
