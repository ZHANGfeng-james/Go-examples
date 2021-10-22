package server

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/go-examples-with-tests/database/v4/pb"
	"github.com/go-examples-with-tests/database/v4/store"

	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
)

type Cache struct {
	store store.Factory
}

var (
	cacheServer *Cache
	once        sync.Once
)

func GetCacheInsOr(sf store.Factory) (*Cache, error) {
	if sf != nil {
		once.Do(func() {
			cacheServer = &Cache{store: sf}
		})
	}

	if cacheServer == nil {
		return nil, fmt.Errorf("got nil cache server")
	}
	return cacheServer, nil
}

func (cache *Cache) ListUsers(ctx context.Context, request *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	limit, msg := request.Limit, request.Msg
	log.Printf("get Request! limit:%d, msg:%s", *limit, *msg)

	// CRUD，拿着 userStore 就可以和 MariaDB 交互
	userStore := cache.store.Users()
	user, err := userStore.Get(context.TODO(), "admin", metav1.GetOptions{})
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println(user)
	}

	items := make([]*pb.UserInfo, 0)
	items = append(items, &pb.UserInfo{
		Nickname: user.Name,
		Password: user.Password,
		Phone:    user.Phone,
		Email:    user.Email,
	})

	return &pb.ListUsersResponse{
		Count: 10,
		Items: items,
	}, nil
}
