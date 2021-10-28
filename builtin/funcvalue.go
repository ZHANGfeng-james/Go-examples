package builtin

import "log"

type SMStartFunc func() error

func (f SMStartFunc) GetName() string {
	return "A Function Value"
}

func newFuncValue() {
	// 创建了一个方法类型的实例，该类型是 SMStartFunc 类型
	funcValue := SMStartFunc(func() error {
		log.Println("create a instance of SMStartFunc type")
		return nil
	})
	log.Println(funcValue.GetName())
}
