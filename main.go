package main

import (
	"context"
	"fmt"
	"gogin/pkg/setting"
	"gogin/routers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	//gin.Default()返回Gin的type Engine struct{...}，里面包含RouterGroup，相当于创建一个路由Handlers，可以后期绑定各类的路由规则和函数、中间件等
	//router := gin.Default()
	router := routers.InitRouter()

	//默认情况下，它在 :8080 上提供服务，除非 PORT 环境变量已定义。
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Printf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}
