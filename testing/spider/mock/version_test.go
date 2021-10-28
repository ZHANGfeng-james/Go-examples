package mock

import (
	"log"
	"testing"

	"github.com/go-examples-with-tests/testing/spider"
	gomock "github.com/golang/mock/gomock"
)

func TestGetGoVersion(t *testing.T) {
	// 创建一个 spider.Spider 实例
	s := spider.CreateGoVersionSpider()
	if s == nil {
		log.Println("spider is nil")
		return
	}
	v := spider.GetGoVersion(s)
	if v != "go1.8.3" {
		t.Errorf("Get wrong version %s", v)
	}
}

func TestGetGoVersionUseMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSpider := NewMockSpider(ctrl)
	mockSpider.EXPECT().GetBody().Return("go1.8.3")

	// 测试目标是： GetGoVersion 函数，其入参是 spider.Spider 接口
	goVer := spider.GetGoVersion(mockSpider)
	// 为的是验证 GetGoVersion 中的逻辑
	if goVer != "go1.8.3" {
		t.Errorf("Get wrong version %s", goVer)
	}
}
