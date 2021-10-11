package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type GobCodec struct {
	conn io.ReadWriteCloser // 保存了 net.Conn 实例，保存每一次连接
	buf  *bufio.Writer
	dec  *gob.Decoder
	enc  *gob.Encoder
}

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(conn),
	}
}

func (gob *GobCodec) ReadHeader(h *Header) error {
	return gob.dec.Decode(h)
}

func (gob *GobCodec) ReadBody(body interface{}) error {
	// gob: type mismatch in decoder: want struct type main.Args; got non-struct
	return gob.dec.Decode(body)
}

func (gob *GobCodec) Write(h *Header, body interface{}) (err error) {
	defer func() {
		// buf 此处并没有装载任何数据，调用 Flush() 也无任何实际作用
		_ = gob.buf.Flush()
		if err != nil {
			_ = gob.Close()
		}
	}()

	// write header，仅仅给的是 h *Header 类型实例
	if err := gob.enc.Encode(h); err != nil {
		log.Println("rpc codec: gob error encoding header:", err)
		return err
	}
	// write body，仅仅给的是一个 interface{} 类型的实例
	if err := gob.enc.Encode(body); err != nil {
		log.Println("rpc codec: gob error encoding body:", err)
		return err
	}
	return nil
}

func (gob *GobCodec) Close() error {
	return gob.conn.Close()
}
