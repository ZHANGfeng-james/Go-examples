package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type GobCodec struct {
	conn io.ReadWriteCloser
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
		_ = gob.buf.Flush()
		if err != nil {
			_ = gob.Close()
		}
	}()

	// write header
	if err := gob.enc.Encode(h); err != nil {
		log.Println("rpc codec: gob error encoding header:", err)
		return err
	}
	// write body
	if err := gob.enc.Encode(body); err != nil {
		log.Println("rpc codec: gob error encoding body:", err)
		return err
	}
	return nil
}

func (gob *GobCodec) Close() error {
	return gob.conn.Close()
}
