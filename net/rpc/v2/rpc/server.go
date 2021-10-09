package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/go-examples-with-tests/net/rpc/v2/codec"
)

type Server struct {
	serviceMap sync.Map
}

func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

func (server *Server) Register(rcvr interface{}) error {
	s := newService(rcvr)
	if _, dup := server.serviceMap.LoadOrStore(s.name, s); dup {
		return errors.New("rpc: service already defined: " + s.name)
	}
	return nil
}

func Register(rcvr interface{}) error {
	return DefaultServer.Register(rcvr)
}

func (server *Server) findService(serviceMethod string) (svc *service, mtype *methodType, err error) {
	dot := strings.LastIndex(serviceMethod, ".")
	if dot < 0 {
		err = errors.New("rpc service: service/method request ill-formed: " + serviceMethod)
		return
	}
	serviceName, methodName := serviceMethod[:dot], serviceMethod[dot+1:]
	svic, ok := server.serviceMap.Load(serviceName)

	if !ok {
		err = errors.New("rpc server: can not find service " + methodName)
		return
	}
	svc = svic.(*service)
	mtype = svc.method[methodName]
	if mtype == nil {
		err = errors.New("rpc server: can not find method " + methodName)
	}
	return
}

func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept() // net.Listener --> Accept() --> net.Conn
		if err != nil {
			Info("rpc server: accept error, ", err)
			return
		}
		go server.ServeConn(conn)
	}
}

func Accept(lis net.Listener) {
	// net.Listener 从哪里来？ l, err := net.Listen("tcp", ":0")
	DefaultServer.Accept(lis)
}

func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	defer func() {
		_ = conn.Close()
	}()

	var opt Option
	// | Option | Header1 | Body1 | Header2 | Body2 |...
	// option 方面使用 JSON 格式编码，最先解析的是 json 格式的 Option
	//FIXME json.NewDecoder(conn).Decode(&opt) 的工作原理？
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		Info("rpc server: Options error, ", err)
		return
	}
	if opt.MagicNumber != MagicNumber {
		Info("rpc server: invalid magic number %x", opt.MagicNumber)
		return
	}
	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		Info("rpc server: invalid codec type %s", opt.CodecType)
		return
	}
	server.serveCodec(f(conn), &opt)
}

var invalidRequest = struct{}{}

// 也就是一次连接中，紧接在 Option 之后的部分：Header 和 Body
type request struct {
	h            *codec.Header // Header
	argv, replyv reflect.Value // Body

	mtype *methodType
	svc   *service
}

// f(conn) 得到的是一个 codec.Codec 编解码器
func (server *Server) serveCodec(cc codec.Codec, opt *Option) {
	// 注意，此处使用的是 *sync.Mutex 和 *sync.WaitGroup
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for {
		// 读取请求
		req, err := server.readRequest(cc)
		if err != nil {
			if req == nil {
				break
			}

			req.h.Error = err.Error() // string 格式的
			// 回复请求
			server.sendResponse(cc, req.h, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		// 处理请求
		go server.handleRequest(cc, req, sending, wg, opt.HandleTimeout)
	}
	wg.Wait()
	_ = cc.Close()
}

func (server *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			Info("rpc server: read header error: ", err)
		}
		return nil, err
	}
	return &h, nil
}

func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	// | Option | Header1 | Body1 | Header2 | Body2 |...
	h, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}

	req := &request{h: h}
	req.svc, req.mtype, err = server.findService(h.ServiceMethod)
	if err != nil {
		return req, err
	}

	// 依据 methodName 获取对应的入参类型和返回值类型，并构造实例
	req.argv = req.mtype.newArgv()
	req.replyv = req.mtype.newReplyv()

	argvi := req.argv.Interface()
	if req.argv.Type().Kind() != reflect.Ptr {
		// 必须要能够修改这个值
		argvi = req.argv.Addr().Interface()
	}
	// 字节流反序列化成变量值
	if err = cc.ReadBody(argvi); err != nil {
		// gob: type mismatch in decoder: want struct type main.Args; got non-struct
		Info("rpc server: read body err:", err)
		return req, err
	}

	return req, nil
}

func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		Info("rpc server: write response error:", err)
	}
}

func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup, handleTimeout time.Duration) {
	defer wg.Done()

	called := make(chan struct{})
	sent := make(chan struct{})

	// 如果超时，则此处会导致 goroutine 泄漏
	go func() {
		time.Sleep(3 * time.Second)
		err := req.svc.call(req.mtype, req.argv, req.replyv)
		called <- struct{}{}
		if err != nil {
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, sending)
			sent <- struct{}{}
			return
		}
		server.sendResponse(cc, req.h, req.replyv.Interface(), sending)
		sent <- struct{}{}
	}()

	if handleTimeout == 0 {
		<-called
		<-sent
		return
	}

	select {
	case <-time.After(handleTimeout):
		req.h.Error = fmt.Sprintf("rpc server: request handle timeout, expect within:%s", handleTimeout)
		server.sendResponse(cc, req.h, invalidRequest, sending)
	case <-called:
		<-sent
	}
}

const (
	connected        = "200 Connected to Gee RPC"
	defaultRPCPath   = "/_geerpc_"
	defaultDebugPath = "/debug/geerpc"
)

func HandleHTTP() {
	DefaultServer.HandleHTTP()
}

func (server *Server) HandleHTTP() {
	// 启动 HTTP Server 端，同时监听的 path 是：defaultRPCPath 和 defaultDebugPath
	http.Handle(defaultRPCPath, server)
	http.Handle(defaultDebugPath, debugHTTP{server})
}

func (server *Server) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	Info("request.Method:%s", request.Method)

	if request.Method != http.MethodConnect {
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
		rw.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = io.WriteString(rw, "405 must CONNECT method\n")
		return
	}

	conn, _, err := rw.(http.Hijacker).Hijack()
	if err != nil {
		log.Print("rpc hijacker, remote:", request.RemoteAddr, ": ", err.Error())
		return
	}
	_, _ = io.WriteString(conn, "HTTP/1.0 "+connected+"\n\n")
	server.ServeConn(conn)
}
