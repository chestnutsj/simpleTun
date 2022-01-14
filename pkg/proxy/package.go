package proxy

import (
	"bufio"
	"crypto/tls"
	"net"
)

const (
	defaultReaderSize = 4096
)
const flag = '\n'

func TcpConnSever(ip string) (conn net.Conn, err error) {

	return net.Dial("tcp", ip)
}

func TcpConnSeverTls(ip string) (*tls.Conn, error) {
	return tls.Dial("tcp", ip, &tls.Config{InsecureSkipVerify: true})
}

type Channel struct {
	c      interface{}
	plugin string
	rb     *bufio.Reader
	wb     *bufio.Writer
}

/*
	tcp or tls
*/
func newChannelTcp(co net.Conn, plugin string) *Channel {
	cc, ok := co.(*net.TCPConn)
	if ok {
		_ = cc.SetNoDelay(false)
	}
	ch := &Channel{
		c:      co,
		plugin: plugin,
		rb:     bufio.NewReaderSize(co, defaultReaderSize),
		wb:     bufio.NewWriterSize(co, defaultReaderSize),
	}
	return ch
}

func newChannel(ip string, plugin string) (*Channel, error) {
	var tcpConn interface{}
	var rb *bufio.Reader
	var wb *bufio.Writer

	switch plugin {
	case "plugin":
		panic("??????? plugin")
	case "tls":
		conn, err := TcpConnSeverTls(ip)
		if err != nil {
			return nil, err
		}
		rb = bufio.NewReaderSize(conn, defaultReaderSize)
		wb = bufio.NewWriterSize(conn, defaultReaderSize)
		tcpConn = conn
	default:
		conn, err := TcpConnSever(ip)
		if err != nil {
			return nil, err
		}
		cc := conn.(*net.TCPConn)
		_ = cc.SetNoDelay(false)
		rb = bufio.NewReaderSize(conn, defaultReaderSize)
		wb = bufio.NewWriterSize(conn, defaultReaderSize)
		tcpConn = conn
	}

	ch := &Channel{
		c:      &tcpConn,
		plugin: plugin,
		rb:     rb,
		wb:     wb,
	}
	return ch, nil
}

func (c *Channel) Read() (data []byte, err error) {
	return c.rb.ReadBytes(flag)
}

func (c *Channel) WriteFlag(data []byte) error {
	data = append(data, flag)
	_, err := c.wb.Write(data)
	if err != nil {
		return err
	}
	return c.wb.Flush()
}

func (c *Channel) Write(data []byte) error {
	_, err := c.wb.Write(data)
	if err != nil {
		return err
	}
	return c.wb.Flush()
}

func (c *Channel) Close() {
	switch c.plugin {
	case "tls":
		cc, ok := c.c.(*tls.Conn)
		if ok {
			cc.Close()
		}
	default:
		cc, ok := c.c.(*net.TCPConn)
		if ok {
			cc.Close()
		}
	}
}
