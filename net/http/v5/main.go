package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type Product struct {
	Username    string    `json:"username" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Category    string    `json:"category" binding:"required"`
	Price       int       `json:"price" binding:"required"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type productHandler struct {
	sync.RWMutex
	products map[string]Product
}

func newProductHandler() *productHandler {
	return &productHandler{
		products: make(map[string]Product),
	}
}

func (u *productHandler) Create(c *gin.Context) {
	u.Lock()
	defer u.Unlock()

	// 参数解析
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		// type H map[string]interface{} 初始化？实际上类似于：map[int]string{1:"1"}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		// var value = make(map[int]string)
		value := map[int]string{1: "1", 2: "2"}
		log.Println(value)

		return
	}

	// 参数校验
	if _, ok := u.products[product.Name]; ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("product %s already exist", product.Name)})
		return
	}

	// 逻辑处理和返回结果
	product.CreatedAt = time.Now()

	u.products[product.Name] = product
	log.Printf("register %s product", product.Name)
	c.JSON(http.StatusOK, product)
}

func (u *productHandler) Get(c *gin.Context) {
	u.Lock()
	defer u.Unlock()

	// 参数解析
	product, ok := u.products[c.Param("name")]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("product %s is not exists", c.Param("name"))})
		return
	}

	c.JSON(http.StatusOK, product)
}

func main() {
	log.Println("main.go")

	var eg errgroup.Group
	insecureServer := &http.Server{
		Addr:         ":8080",
		Handler:      router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	secureServer := &http.Server{
		Addr:         ":8443",
		Handler:      router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	eg.Go(func() error {
		err := insecureServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	eg.Go(func() error {
		err := secureServer.ListenAndServeTLS("server.pem", "server.key")
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
}

func router() http.Handler {
	router := gin.Default()
	productHandler := newProductHandler()
	v1 := router.Group("/v1")
	{
		productv1 := v1.Group("/products")
		{
			productv1.POST("", productHandler.Create)
			productv1.GET(":name", productHandler.Get)
		}
	}
	return router
}
