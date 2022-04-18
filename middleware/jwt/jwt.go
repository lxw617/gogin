package jwt

import (
	"github.com/gin-gonic/gin"
	result "gogin/common"
	"gogin/pkg/util"
	"net/http"
	"strings"
	"time"
)

//基于JWTAuthMiddleware鉴权中间件

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里的具体实现方式要依据你的实际业务情况决定
		//authHeader := c.Request.Header.Get("Authorization")
		//判断是是否放在请求头
		token := c.GetHeader("Authorization")
		if token == "" {
			//判断是否放在请求体
			//判断是否放在url
			token = c.Query("token")
			if token == "" {
				c.JSON(http.StatusUnauthorized, result.ErrSignParam)
				c.Abort()
				return
			}
		}
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 按空格划分为数组，最后提取有效部分
		parts := strings.SplitN(token, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 2004,
				"msg":  "请求头中auth格式有误",
			})
			c.Abort()
			return
		}
		//校验参数
		claims, err := util.ParseToken(parts[1])
		if err != nil {
			//校验失败
			c.JSON(http.StatusUnauthorized, result.ErrAuthCheckTokenFail)
		} else if time.Now().Unix() > claims.ExpiresAt {
			//token参数超时
			c.JSON(http.StatusUnauthorized, result.ErrAuthCheckTokenTimeOut)
		}

		//如果验证失败说明这个token无效,如果验证成功将当前请求的username信息保存到请求的上下文c上
		c.Set("username", claims.Username)
		//调用c.Next()方法才会调用下一层
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}
