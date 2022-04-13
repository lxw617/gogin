package routers

import (
	"github.com/gin-gonic/gin"
	"gogin/pkg/setting"
	"gogin/routers/api"
	"net/http"
)

func InitRouter() *gin.Engine {
	router := gin.New()
	//记录器和恢复（无崩溃）中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	gin.SetMode(setting.RunMode)
	//router.GET(...){...}：创建不同的HTTP方法绑定到Handlers中，也支持POST、PUT、DELETE、PATCH、OPTIONS、HEAD 等常用的Restful方法
	//gin.Context：Context是gin中的上下文，它允许我们在中间件之间传递变量、管理流、验证JSON请求、响应JSON请求等，
	//在gin中包含大量Context的方法，例如我们常用的DefaultQuery、Query、DefaultPostForm、PostForm等等
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{ //gin.H{...}：就是一个map[string]interface{}
			"message": "test",
		})
	})
	//对于每个匹配的请求上下文将保存路由定义
	router.POST("/user/:name/*action", func(c *gin.Context) {
		b := c.FullPath() == "/user/:name/*action" // true
		c.String(http.StatusOK, "%t", b)
	})
	//注册路由
	apiv1 := router.Group("api/v1")
	api.InitApi(apiv1)
	return router
}
