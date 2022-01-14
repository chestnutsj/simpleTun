package proxy

import (
	"context"
	"crypto/tls"
	"net"
	"sync"

	"go.uber.org/zap"
)

type Server struct {
	ctx        context.Context
	cancel     context.CancelFunc
	listener   net.Listener
	remote     string
	remoteType string
	plugin     string
	sync.RWMutex
	tunnelMap map[string]*Tunnel
}

func NewServer(ctxx context.Context, remote string, remoteType string) *Server {
	ctx, cancel := context.WithCancel(ctxx)
	proxy := &Server{
		ctx:        ctx,
		cancel:     cancel,
		remote:     remote,
		remoteType: remoteType,
		plugin:     "tcp",
	}

	return proxy
}

func (p *Server) SetListen(address string) error {
	var err error
	p.listener, err = net.Listen("tcp", address)
	return err
}

func (p *Server) SetListenTls(address, crt, key string) error {
	cer, err := tls.LoadX509KeyPair(crt, key)
	if err != nil {
		return err
	}
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cer}}
	p.listener, err = tls.Listen("tcp", address, tlsConfig)
	p.plugin = "tls"
	return err
}

func (p *Server) Stop() {
	if p.listener != nil {
		p.listener.Close()
		p.listener = nil
	}

	p.cancel()
}

func (p *Server) Start() error {
	defer func() {
		p.Stop()
	}()

	for {
		conn, err := p.listener.Accept()
		if err != nil {
			return err
		}
		go p.onConn(conn)
		select {
		case <-p.ctx.Done():
			return nil
		default:
		}
	}
}

func (p *Server) onConn(conn net.Conn) {
	defer conn.Close()

	key := conn.RemoteAddr().String()
	srcCh := newChannelTcp(conn, p.plugin)
	tunnel, err := NewTunnel(p.ctx, srcCh, p.remoteType, p.remote)
	if err != nil {
		zap.L().Error("err new Tunnel", zap.Error(err))
		return
	}
	p.Lock()
	p.tunnelMap[key] = tunnel
	p.Unlock()

	defer func() {
		p.Lock()
		delete(p.tunnelMap, key)
		p.Unlock()
		tunnel.Stop()
	}()

	tunnel.Start()
}
