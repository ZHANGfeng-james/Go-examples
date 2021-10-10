package gee

import (
	"log"
	"time"
)

func Logger() HandleFunc {
	return func(ctx *Context) {
		start := time.Now()

		ctx.Next() // 把控制执行权交给下一个 HandleFunc 实例

		log.Printf("[%v], %v, %v", ctx.StatusCode, ctx.Request.RequestURI, time.Since(start))
	}
}
