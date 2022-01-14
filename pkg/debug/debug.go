package debug

import (
	"net/http"
	_ "net/http/pprof"

	"go.uber.org/zap"
)

func StartApi(addr string) {
	//http://127.0.0.1:8899/debug/pprof/

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		zap.L().Error(" api server start failed ")
	}
}
