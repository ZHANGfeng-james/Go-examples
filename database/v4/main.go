package main

import (
	"log"
	"net"

	"github.com/go-examples-with-tests/database/v4/pb"
	"github.com/go-examples-with-tests/database/v4/pkg"
	"github.com/go-examples-with-tests/database/v4/server"
	"github.com/go-examples-with-tests/database/v4/store"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func main() {
	log.Println("hello, world!")

	options := pkg.NewMySQLOptions()
	initOptions(options)
	log.Printf("MariaDB options:%s", options)

	// connection to MariaDB
	dbFactory, err := store.GetMySQLFactoryOr(options)
	// close DB connection
	defer dbFactory.Close()
	if err != nil {
		log.Fatal(err)
	}
	store.SetClient(dbFactory)

	serverConf := server.NewServerOption()
	opts := []grpc.ServerOption{grpc.MaxRecvMsgSize(100)}
	grpcServer := grpc.NewServer(opts...)

	cache, err := server.GetCacheInsOr(dbFactory)
	if err != nil {
		log.Println(err.Error())
		return
	}
	pb.RegisterCacheServer(grpcServer, cache)
	reflection.Register(grpcServer)

	Run(serverConf, grpcServer)
}

func Run(opts *server.ServerOption, server *grpc.Server) {
	listener, err := net.Listen("tcp", opts.Addr+opts.Port)
	if err != nil {
		log.Fatalf("failed to listen: %s", listener.Addr())
	}
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to start grpc server: %s", err.Error())
	}
}

func initOptions(options *pkg.MySQLOptions) {
	if options == nil {
		return
	}
	options.Username = "iam"
	options.Password = "iam59!z$"
	options.Database = "iam"
	options.LogLevel = 4
}
