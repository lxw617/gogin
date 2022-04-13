package main

import (
	"fmt"
	"gogin/pkg/setting"
	"gogin/routers"
	"net/http"
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

	s.ListenAndServe()
}
