package server

type ServerOption struct {
	Addr string
	Port string
}

func NewServerOption() *ServerOption {
	return &ServerOption{
		Addr: "127.0.0.1",
		Port: ":8081",
	}
}
