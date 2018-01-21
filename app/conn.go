package app

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"

	gg "github.com/haozibi/gglog"
)

type Conn struct {
	TCPConn   *net.TCPConn
	Reader    *bufio.Reader
	closeFlag bool
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

func (c *Conn) Close() {
	if c.TCPConn != nil && c.closeFlag == false {
		c.closeFlag = true
		c.TCPConn.Close()
	}
}

func (c *Conn) IsClosed() bool {
	return c.closeFlag
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
