package routers

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"gogin/controller"
	"log"
	"time"

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
	//httprouter会将所有路由规则构造一颗前缀树 例如有 root and as at cn com

	//路由原理
	//路由器依赖于大量使用公共前缀的树结构，它基本上是一个紧凑的 前缀树（或只是基数树）。具有共同前缀的节点也共享一个共同的父节点。
	//Every*<num>表示处理函数（指针）的内存地址。如果您沿着从树根到叶子的路径穿过树，您将获得完整的路由路径，例如\blog\:post\， where:post只是实际帖子名称的占位符（参数）。与哈希映射不同，树结构还允许我们使用:post参数等动态部分，因为我们实际上匹配路由模式，而不仅仅是比较哈希。正如基准所显示的那样，这非常有效且有效。
	//由于 URL 路径具有分层结构并且仅使用有限的一组字符（字节值），因此很可能存在许多公共前缀。这使我们能够轻松地将路由减少为更小的问题。此外，路由器为每个请求方法管理一个单独的树。一方面，它比在每个节点中保存一个方法->句柄映射更节省空间，它还允许我们在开始查找前缀树之前大大减少路由问题。
	//为了更好的可扩展性，每个树级别上的子节点按优先级排序，其中优先级只是在子节点（子节点、孙子节点等）中注册的句柄数。这有两个方面的帮助：
	//首先评估属于最多路由路径的节点。这有助于使尽可能多的路线尽可能快地到达。
	//这是某种成本补偿。总是可以首先评估最长可达路径（最高成本）。以下方案可视化了树结构。节点从上到下，从左到右进行评估。
	//  Priority   Path             Handle
	//  9          \                *<1>
	//	3          ├s               nil
	//  2          |├earch\         *<2>
	//	1          |└upport\        *<3>
	//	2          ├blog\           *<4>
	//	1          |    └:post      nil
	//  1          |         └\     *<5>
	//	2          ├about-us\       *<6>
	//	1          |        └team\  *<7>
	//	1          └contact\        *<8>
	// 路由器本身实现了http.Handle接口，路由器支持http.Handle和http.HandlerFunc在注册路由时当作httprouter.Handler使用
	// 或者，也可以使用params := r.Context().Value(httprouter.ParamsKey)辅助函数来代替。

	//命名参数
	// router.GET("/hello/:name", Hello) Hello method
	//:name是一个命名参数。这些值可以通过 访问httprouter.Params，它只是httprouter.Params 的一部分。您可以通过切片中的索引或使用ByName(name)方法:name获取参数的值：可以通过ByName("name").
	//当使用http.Handler(using router.Handleror http.HandlerFunc) 而不是使用第三个函数参数的 HttpRouter 的句柄 API 时，命名参数存储在request.Context.
	//命名参数仅匹配单个路径段 由于该路由器只有显式匹配，您不能为同一路径段注册静态路由和参数。例如，您不能同时注册模式/user/new和/user/:user相同的请求方法。不同请求方法的路由是相互独立的。
	//包罗万象参数
	//形式为*name. 顾名思义，它们匹配一切。因此它们必须始终位于模式的末尾

	// 创建路由 默认8080 r := gin.Default()
	router := gin.New()
	//记录器和恢复（无崩溃）中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	gin.SetMode(setting.ServerSetting.RunMode)

	// router.GET(...){...}：创建不同的HTTP方法绑定到Handlers中，也支持POST、PUT、DELETE、PATCH、OPTIONS、HEAD 等常用的Restful方法
	// gin.Context：Context是gin中的上下文，它允许我们在中间件之间传递变量、管理流、验证JSON请求、响应JSON请求等，
	// 在gin中包含大量Context的方法，例如我们常用的DefaultQuery、Query、DefaultPostForm、PostForm等等
	// 绑定路由规则，执行的函数 gin.Context，封装了request和response
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{ //gin.H{...}：就是一个map[string]interface{}
			"message": "test",
		})
	})
	// 对于每个匹配的请求上下文将保存路由定义 localhost:8000/xxx/zhangsan
	router.POST("/user/:name/*action", func(c *gin.Context) {
		// 可以通过Context的Param方法来获取 API参数
		name := c.Param("name")
		action := c.Param("action")
		b := c.FullPath() == "/user/:name/*action" // true
		c.String(http.StatusOK, "路径匹配：%t，获取请求路径参数：%t", b, " name_ "+name+" action_ "+action)
	})

	// 通过 swag init 把 Swagger API 所需要的文件都生成了，那接下来我们怎么访问接口文档呢？其实很简单，我们只需要在 routers 中进行默认初始化和注册对应的路由
	// 从表面上来看，主要做了两件事，分别是初始化 docs 包和注册一个针对 swagger 的路由，而在初始化 docs 包后，
	// 其 swagger.json 将会默认指向当前应用所启动的域名下的 swagger/doc.json 路径，如果有额外需求，可进行手动指定
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 生成token
	router.GET("/auth", controller.GetToken)

	// 设置路由地址
	apiv1 := router.Group("api/v1")
	// 加入jwt鉴权中间件
	apiv1.Use(jwt.JWTAuthMiddleware())
	{
		//注册路由
		api.InitApi(apiv1)
	}

	//重定向
	router.GET("/baidu", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "http://www.baidu.com")
	})

	//goroutine机制可以方便地实现异步处理
	//另外，在启动新的goroutine时，不应该使用原始上下文，必须使用它的只读副本
	// 1.异步
	router.GET("/long_async", func(c *gin.Context) {
		// 需要搞一个副本
		copyContext := c.Copy()
		// 异步处理
		go func() {
			time.Sleep(3 * time.Second)
			log.Println("异步执行：" + copyContext.Request.URL.Path)
		}()
	})
	// 2.同步
	router.GET("/long_sync", func(c *gin.Context) {
		time.Sleep(3 * time.Second)
		log.Println("同步执行：" + c.Request.URL.Path)
	})

	//404 error
	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "404 not found")
	})

	//HTML模板渲染
	//gin支持加载HTML模板, 然后根据模板参数进行配置并返回相应的数据，本质上就是字符串替换
	//LoadHTMLGlob()方法可以加载模板文件
	router.LoadHTMLGlob("tem/*")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"title": "我是测试", "ce": "123456"})
	})
	//如果你需要引入静态文件需要定义一个静态文件目录 r.Static("/assets", "./assets")

	return router
}

/*
    浏览器发出OPTIONS方法，进行跨域CORS
	自动 OPTIONS 响应和 CORS
	可能希望修改对 OPTIONS 请求的自动响应，例如支持CORS 预检请求或设置其他标头。这可以使用Router.GlobalOPTIONS处理程序来实现：
	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", "*")
		}

		// Adjust status code to 204
		w.WriteHeader(http.StatusNoContent)
	})
*/
