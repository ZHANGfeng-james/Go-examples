package gocontext

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	cxt := context.Background()

	delay, cancel := context.WithTimeout(cxt, 3*time.Second)
	defer cancel()

	go func() {
		time.Sleep(10 * time.Second)
	}()

	select {
	case <-delay.Done():
		fmt.Println(delay.Err())
	case <-time.After(6 * time.Second):
		fmt.Println("over")
	}
}
