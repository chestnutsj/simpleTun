package proxy

import (
	"context"
	"log"

	"go.uber.org/zap"
)

const (
	src = iota
	dst
)

type Tunnel struct {
	ctx    context.Context
	cancle context.CancelFunc
	c      [2]*Channel
}

func NewTunnel(ctxx context.Context, srcConn *Channel,
	remoteType string, remote string) (*Tunnel, error) {
	ctx, cancel := context.WithCancel(ctxx)
	t := &Tunnel{
		ctx:    ctx,
		cancle: cancel,
	}
	t.c[src] = srcConn
	var err error
	t.c[dst], err = newChannel(remote, remoteType)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Tunnel) Stop() {
	if t.cancle != nil {
		t.cancle()
	}
	for _, cc := range t.c {
		if cc != nil {
			cc.Close()
		}
	}
}

func (t *Tunnel) Start() {
	defer t.Stop()

	go t.backgroud()
	for {
		data, err := t.c[src].Read()
		if err != nil {
			log.Println("src read failed", err.Error())
			return
		}
		err = t.c[dst].Write(data)
		if err != nil {
			zap.L().Error("dst write failed", zap.Error(err))
			return
		}
		zap.L().Debug("src data" + string(data))
		select {
		case <-t.ctx.Done():
			return
		default:
		}
	}
}
func (t *Tunnel) backgroud() {
	defer t.Stop()
	for {
		data, err := t.c[dst].Read()
		if err != nil {
			zap.L().Error("dst Read failed", zap.Error(err))
			return
		}

		err = t.c[src].Write(data)
		if err != nil {
			zap.L().Error("src write failed", zap.Error(err))
			return
		}
		zap.L().Debug("dst" + string(data))

		select {
		case <-t.ctx.Done():
			return
		default:
		}
	}
}
