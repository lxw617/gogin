package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	result "gogin/common"
	"gogin/pkg/setting"
	"gogin/pkg/util"
	"gogin/service"
	"net/http"
)

//获取所有文章标签
func GetTags(c *gin.Context) {
	//初始化参数
	params := make(map[string]interface{})
	//使用现有的底层请求对象解析查询字符串参数 请求响应 url 匹配：/tag?name=Jane
	name := c.Query("name")
	if name != "" {
		params["name"] = name
	}
	state := -1
	if arg := c.Query("state"); arg != "" {
		// Convert string to specify type.
		state = com.StrTo(arg).MustInt()
		params["state"] = state
	}
	//获取结果
	tags, count, err := service.GetTags(params, util.GetPage(c), setting.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.Err)
	}
	c.JSON(http.StatusOK, result.OK.WithData(map[string]interface{}{
		"tags":  tags,
		"count": count}))
}

//新增文章标签
func AddTag(c *gin.Context) {
	c.JSON(http.StatusOK, result.OK.WithMsg("新增文章标签成功"))
}

//删除文章标签
func RemoveTag(c *gin.Context) {
	c.JSON(http.StatusOK, result.OK.WithMsg("删除文章标签成功"))
}

//根据id获取文章标签
func GetTagById(c *gin.Context) {
	c.JSON(http.StatusOK, result.OK.WithMsg("根据id获取文章标签成功"))
}

//修改文章标签
func EditTag(c *gin.Context) {
	c.JSON(http.StatusOK, result.OK.WithMsg("修改文章标签成功"))
}
