package registry

import (
	"net/http"
	"time"

	"github.com/go-examples-with-tests/database/v1/log"
)

//FIXME 心跳这种功能，应该是实现在 Server 端，还是在 Registry 端？
func Heartbeat(registry, rpcAddr string, duration time.Duration) {
	if duration == 0 {
		duration = defaultTimeout - time.Duration(1)*time.Minute
	}
	var err error
	err = sendHeartbeat(registry, rpcAddr) // 首次发送一次
	go func() {
		ticker := time.NewTicker(duration)
		for err == nil { // 间隔 duration 持续发送
			<-ticker.C
			err = sendHeartbeat(registry, rpcAddr)
		}
	}()
}

func sendHeartbeat(registry, rpcAddr string) error {
	log.Infof("server:%s send heartbeat signal to registry %s", rpcAddr, registry)
	client := &http.Client{}
	request, _ := http.NewRequest("POST", registry, nil)
	// 为什么在对端接收不到？对端获取 Header 出错
	request.Header.Set("X-Geerpc-Server", rpcAddr)
	if _, err := client.Do(request); err != nil {
		log.Errorf("rpc server: heart beat err:%v", err)
		return err
	}
	return nil
}
