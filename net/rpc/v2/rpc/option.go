package rpc

import (
	"time"

	"github.com/go-examples-with-tests/net/rpc/v2/codec"
)

const MagicNumber = 0x3bef5c

type Option struct {
	MagicNumber int        // 标记这是 geerpc 的 request
	CodecType   codec.Type // client 还可使用其他的 codec 用于编码 body 部分

	ConnectionTimeout time.Duration
	HandleTimeout     time.Duration
}

var DefaultOption = &Option{
	MagicNumber:       MagicNumber,
	CodecType:         codec.GobType, // 默认情况下，RPC 服务端使用 gob codec 解码
	ConnectionTimeout: 3 * time.Second,
}
