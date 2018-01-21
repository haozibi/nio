package app

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"

	gg "github.com/haozibi/gglog"
)

type Listener struct {
	addr        net.Addr
	tcplistener *net.TCPListener
	conns       chan *Conn
	closeFlag   bool
}

type Conn struct {
	TCPConn   *net.TCPConn
	Reader    *bufio.Reader
	closeFlag bool
}

func Listen(host string, port int64) (tcpl *Listener, err error) {
	tcpAddr, _ := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%v:%v", host, port))
	listenr, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return
	}
	tcpl = &Listener{
		addr:        listenr.Addr(),
		tcplistener: listenr,
		conns:       make(chan *Conn),
		closeFlag:   false,
	}

	go func() {
		for {
			conn, err := tcpl.tcplistener.AcceptTCP()
			if err != nil {
				if tcpl.closeFlag {
					return
				}
				continue
			}
			c := &Conn{
				TCPConn:   conn,
				closeFlag: false,
				Reader:    bufio.NewReader(conn),
			}
			tcpl.conns <- c
		}
	}()

	return tcpl, err
}

func DialServer(host string, port int64) (conn *Conn, err error) {
	conn = new(Conn)
	serverAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		return
	}
	c, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		return
	}
	conn.TCPConn = c
	conn.Reader = bufio.NewReader(c)
	conn.closeFlag = false
	return conn, nil
}

// Listener

// 获得 conn，阻塞等待新的 connection
func (l *Listener) GetConn() (conn *Conn, err error) {
	var ok bool
	conn, ok = <-l.conns
	if !ok {
		return conn, fmt.Errorf("channel close")
	}
	return conn, nil
}

func (l *Listener) Close() {
	if l.tcplistener != nil && l.closeFlag == false {
		l.closeFlag = true
		l.tcplistener.Close()
		close(l.conns)
	}
}

// Conn

func (c *Conn) Close() {
	if c.TCPConn != nil && c.closeFlag == false {
		c.closeFlag = true
		c.TCPConn.Close()
	}
}

func (c *Conn) IsClosed() bool {
	return c.closeFlag
}

func (c *Conn) GetRemoteAddr() (addr string) {
	return c.TCPConn.RemoteAddr().String()
}

func (c *Conn) GetLocalAddr(addr string) {
	return c.TCPConn.LocalAddr().String()
}

func (c *Conn) ReadLine() (buff string, err error) {
	buff, err = c.Reader.ReadString('\n')
	if err == io.EOF {
		c.closeFlag = true
	}
	return
}

func (c *Conn) Write(content string) (err error) {
	_, err = c.TCPConn.Write([]byte(content))
	return
}

func Join(c1 *Conn, c2 *Conn) {
	var wait sync.WaitGroup
	pipe := func(dst *Conn, src *Conn) {
		defer dst.Close()
		defer src.Close()
		defer wait.Done()

		var err error
		// func Copy(dst Writer, src Reader) (written int64, err error)
		_, err = io.Copy(dst.TCPConn, src.TCPConn)
		if err != nil {
			gg.Errorf("join conns error,%v", err)
		}
	}
	wait.Add(2)
	go pipe(c1, c2)
	go pipe(c2, c1)
	wait.Wait()
	return
}
