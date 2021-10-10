package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/go-examples-with-tests/net/http/v3/gee"
)

func main() {
	engine := gee.Default()

	engine.GET("/json", func(ctx *gee.Context) {
		obj := gee.H{
			"name":     "geektutu",
			"password": "1234",
		}
		array := []int{1, 2, 3}
		array[3] = 4
		ctx.JSON(http.StatusOK, obj)
	})

	// example: curl "http://localhost:9999/postform" -X POST -d 'password=1&name=1'
	engine.POST("/postform", func(ctx *gee.Context) { // 必须是 POST 请求，才能解析出 PostForm 内容
		ctx.JSON(http.StatusOK, gee.H{
			"name":     ctx.PostForm("name"),
			"password": ctx.PostForm("password"),
		})
	})

	// example: curl "http://localhost:9999/query?username=Michoi"
	engine.GET("/query", func(ctx *gee.Context) {
		username := ctx.Query("username")
		ctx.String(http.StatusOK, "Hello, %s!", username)
	})

	engine.GET("/assets/*filepath", func(ctx *gee.Context) {
		ctx.JSON(http.StatusOK, gee.H{"filepath": ctx.Param("filepath")})
	})

	helloGroup := engine.Group("/v1")
	{
		// 获得一个 RouterGroup 后，直接使用该类型注册 pattern
		helloGroup.GET("/:name", func(ctx *gee.Context) {
			ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Param("name"), ctx.Path)
		})
		helloGroup.GET("/geektutu/join", func(ctx *gee.Context) {
			ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
		})
	}

	v2 := engine.Group("/v2")
	{
		v2.GET("/", func(ctx *gee.Context) {
			ctx.String(http.StatusOK, "you're at %s\n", ctx.Path)
		})
		v2.GET("/help", func(ctx *gee.Context) {
			ctx.String(http.StatusOK, "you're at %s\n", ctx.Path)
		})
	}

	v3 := engine.Group("/v3")
	{
		v3.Use(gee.Logger())
		v3.GET("/logger", func(ctx *gee.Context) {
			ctx.String(http.StatusOK, "logger!\n")
		})
	}

	// 或者本地的其他目录
	engine.Static("/assets", "./static")
	engine.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	engine.LoadHTMlGlob("templates/*")

	type Student struct {
		Name string
		Age  int8
	}
	stu1 := &Student{
		Name: "Geektutu",
		Age:  20,
	}
	stu2 := &Student{
		Name: "Jack Ma",
		Age:  22,
	}
	engine.GET("/", func(ctx *gee.Context) {
		ctx.HTML(http.StatusOK, "css.tmpl", nil)
	})

	engine.GET("/students", func(ctx *gee.Context) {
		ctx.HTML(http.StatusOK, "arr.tmpl", gee.H{
			"title":  "gee",
			"stuArr": [2]*Student{stu1, stu2},
		})
	})

	engine.GET("/date", func(ctx *gee.Context) {
		ctx.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
			"title": "geektutu",
			"now":   time.Date(2019, 8, 27, 0, 0, 0, 0, time.UTC),
		})
	})

	engine.Run(":9999")
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}
