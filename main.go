package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/chestnutsj/simpleTun/pkg/config"
	"github.com/chestnutsj/simpleTun/pkg/debug"
	"github.com/chestnutsj/simpleTun/pkg/logger"
	"github.com/chestnutsj/simpleTun/pkg/proxy"
)

func main() {
	ctx, cancle := context.WithCancel(context.Background())
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("level " + cfg.Level)
	logger.ChangeLogLevel(cfg.Level)
	if cfg.DebugSer != nil {
		debug.StartApi(*cfg.DebugSer)
	}
	wg := sync.WaitGroup{}
	mSevr := make(map[string]*proxy.Server)
	for k, ss := range cfg.Server {
		svr1 := proxy.NewServer(ctx, ss.Remote, ss.Type)
		err := svr1.SetListen(ss.Local)
		if err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		fmt.Println(k + " start:" + ss.Local)
		go func(pp *proxy.Server, addr string) {
			defer func() {
				wg.Done()
				if cancle != nil {
					cancle()
				}
			}()
			log.Println("start" + addr)
			pp.Start()
		}(svr1, ss.Local)
		mSevr[ss.Local] = svr1
	}
	wg.Wait()
}
