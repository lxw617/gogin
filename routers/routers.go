package routers

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"gogin/controller"
	//_ 表示执行init函数时调用改包，需要将这个替换为自己本地的docs目录路径。这个路径是github上别人的docs,此处只是用来测试。
	//此处应该这样写：_ "swagger_demo/docs"
	//上面的swagger_demo为本项目名称，docs就是swag init自动生成的目录，用于存放 docs.go、swagger.json、swagger.yaml 三个文件。
	_ "gogin/docs"
	"gogin/middleware/jwt"
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
	//通过 swag init 把 Swagger API 所需要的文件都生成了，那接下来我们怎么访问接口文档呢？其实很简单，我们只需要在 routers 中进行默认初始化和注册对应的路由
	//从表面上来看，主要做了两件事，分别是初始化 docs 包和注册一个针对 swagger 的路由，而在初始化 docs 包后，
	//其 swagger.json 将会默认指向当前应用所启动的域名下的 swagger/doc.json 路径，如果有额外需求，可进行手动指定
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//生成token
	router.GET("/auth", controller.GetToken)
	//注册路由
	apiv1 := router.Group("api/v1")
	//加入jwt鉴权中间件
	apiv1.Use(jwt.JWTAuthMiddleware())
	{
		api.InitApi(apiv1)
	}

	return router
}
