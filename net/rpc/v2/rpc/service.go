package rpc

import (
	"go/ast"
	"log"
	"reflect"
	"sync/atomic"
)

// 一个 method 的所有完整信息
type methodType struct {
	// func (t *T)MethodName(argType T1, replyType *T2) error
	// 一个 method 的所有信息包括：方法名（统一到 Func 这种类型值上），入参，返回值
	method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
	numCalls  uint64
}

func (m *methodType) NumCalls() uint64 {
	return atomic.LoadUint64(&m.numCalls)
}

func (m *methodType) newArgv() reflect.Value {
	var argv reflect.Value

	// 准备接收 reqeust.argv
	if m.ArgType.Kind() == reflect.Ptr {
		argv = reflect.New(m.ArgType.Elem()) // reflect.Type.Elem()
	} else {
		argv = reflect.New(m.ArgType).Elem() // reflect.Value.Elem()
	}

	return argv
}

func (m *methodType) newReplyv() reflect.Value {
	// reply must be a pointer type，这是 RPC 协议规定的
	replyv := reflect.New(m.ReplyType.Elem())

	// 依据 Kind 类型的不同，对 request.replyv 初始化
	switch m.ReplyType.Elem().Kind() {
	case reflect.Map:
		replyv.Elem().Set(reflect.MakeMap(m.ReplyType.Elem())) // reflect.Value
	case reflect.Slice:
		replyv.Elem().Set(reflect.MakeSlice(m.ReplyType.Elem(), 0, 0))
	}
	return replyv
}

type service struct {
	name   string
	typ    reflect.Type           // 结构体类型
	rcvr   reflect.Value          // 这个机构体的实例，后面调用结构体方法时，作为第一个参数
	method map[string]*methodType // 这个结构体的所有方法列表
}

func newService(receiver interface{}) *service {
	s := new(service)

	s.rcvr = reflect.ValueOf(receiver)

	s.typ = reflect.TypeOf(receiver)
	s.name = reflect.Indirect(s.rcvr).Type().Name() // *Foo --> Foo

	// 判断 struct name 是否是可导出的
	if !ast.IsExported(s.name) {
		log.Fatalf("rpc server: %s is not a valid service name", s.name)
	}
	log.Printf("new Service name:%s", s.name)
	s.registerMethods()

	return s
}

func (s *service) registerMethods() {
	s.method = make(map[string]*methodType)
	// 使用 s.typ 获取这个 reflect.Type 下的方法信息
	for i := 0; i < s.typ.NumMethod(); i++ {
		method := s.typ.Method(i) // reflect.Method
		mType := method.Type      // reflect.Type --> Func

		// 方法的第一个入参是接收者本身 mType.NumIn 和 mType.NumOut 在调用时 Kind 必须是 Func
		if mType.NumIn() != 3 || mType.NumOut() != 1 {
			continue
		}
		if mType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			continue
		}

		// 方法的参数，一定是从 index = 1 开始，index = 0 位置的参数是方法接收者
		argType, replyType := mType.In(1), mType.In(2)
		if !isExportedOrBuiltinType(argType) || !isExportedOrBuiltinType(replyType) {
			continue
		}

		// Sum
		s.method[method.Name] = &methodType{
			method:    method,
			ArgType:   argType,
			ReplyType: replyType,
		}
		log.Printf("rpc server: register %s.%s\n", s.name, method.Name)
	}
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	return ast.IsExported(t.Name()) || t.PkgPath() == ""
}

func (s *service) call(m *methodType, argv, replyv reflect.Value) error {
	atomic.AddUint64(&m.numCalls, 1)

	// m.method 是 reflect.Method 类型 --> reflect.Value 类型，且其 Kind 是 Func
	f := m.method.Func

	returnValues := f.Call([]reflect.Value{s.rcvr, argv, replyv})
	if errInter := returnValues[0].Interface(); errInter != nil {
		return errInter.(error)
	}
	return nil
}
