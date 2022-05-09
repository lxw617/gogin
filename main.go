package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gogin/models"
	"gogin/pkg/logging"
	"gogin/pkg/redis"
	"gogin/pkg/setting"
	"gogin/routers"
	"log"
	"net/http"
)

func init() {
	/*
		配置统管
		在 Go 中，当存在多个 init 函数时，执行顺序为：
		相同包下的 init 函数：按照源文件编译顺序决定执行顺序（默认按文件名排序）
		不同包下的 init 函数：按照包导入的依赖关系决定先后顺序
		所以要避免多 init 的情况，尽量由程序把控初始化的先后顺序
	*/
	setting.Setup()
	//连接数据库
	models.Setup()
	//连接redis
	gredis.Setup()
	//日志
	logging.Setup()
	//util.Setup()
}

func main() {

	gin.SetMode(setting.ServerSetting.RunMode)

	gin.DisableConsoleColor()

	fmt.Println("air running...")

	routersInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info] start http server listening %s", endPoint)

	server.ListenAndServe()
}
