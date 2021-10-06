package v4

import "github.com/go-examples-with-tests/net/http/v4/cachepb"

type PeerPicker interface {
	PickPeer(key string) (PeerGetter, bool)
}

type PeerGetter interface {
	// Get(group string, key string) ([]byte, error)

	Get(in *cachepb.Request, out *cachepb.Response) error // 替换成使用 protobuf，作为通信信息的格式
}
