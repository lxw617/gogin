package controller

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	result "gogin/common"
	"gogin/pkg/logging"
	"gogin/pkg/setting"
	"gogin/pkg/util"
	"gogin/service"
	"net/http"
)

/*注解	描述
@Summary	摘要
@Produce	API 可以产生的 MIME 类型的列表，MIME 类型你可以简单的理解为响应类型，例如：json、xml、html 等等
@Param	参数格式，从左到右分别为：参数名、入参类型、数据类型、是否必填、注释
@Success	响应成功，从左到右分别为：状态码、参数类型、数据类型、注释
@Failure	响应失败，从左到右分别为：状态码、参数类型、数据类型、注释
@Router	路由，从左到右分别为：路由地址，HTTP 方法*/

// @Summary Get multiple article tags
// @Produce  json
// @Param name query string false "Name"
// @Param state query int false "State"
// @Success 200 {object} map[string][]string
// @Failure 500 {object} map[string][]string
// @Router /api/v1/tag/tags [get]
func GetTags(c *gin.Context) { //c *gin.Context是Gin很重要的组成部分，可以理解为上下文，它允许我们在中间件之间传递变量、管理流、验证请求的JSON和呈现JSON响应
	//初始化参数
	params := make(map[string]interface{})
	//c.Query使用现有的底层请求对象解析查询字符串参数 请求响应 url 匹配：/tag?name=Jane
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
	tags, count, err := service.GetTags(params, util.GetPage(c), setting.PageSize) //util.GetPage保证了各接口的page处理是一致的
	//在获取标签列表接口中，我们可以根据name、state、page来筛选查询条件，分页的步长可通过app.ini进行配置，以lists、total的组合返回达到分页效果
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.Err)
	}
	c.JSON(http.StatusOK, result.OK.WithData(map[string]interface{}{
		"tags":  tags,
		"count": count}))
}

// @Summary Add article tag
// @Produce  json
// @Param name body string true "Name"
// @Param state body int false "State"
// @Param created_by body int false "CreatedBy"
// @Success 200 {object} map[string][]string
// @Failure 500 {object} map[string][]string
// @Router /api/v1/tag/tag [post]
func AddTag(c *gin.Context) {
	name := c.Query("name")
	//c.DefaultQuery则支持设置一个默认值
	state := com.StrTo(c.DefaultQuery("state", "1")).MustInt()
	createdBy := com.StrTo(c.Query("createdBy")).MustInt()
	//参数校验
	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长不能超过100字符")
	valid.Required(createdBy, "createdBy").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "createdBy").Message("名称最长不能超过100字符")
	valid.Range(state, 0, 1, "status").Message("状态只能为0或者1")
	if !valid.HasErrors() {
		//调用方法
		tag, err := service.AddTag(name, state, createdBy)
		if err != nil {
			c.JSON(http.StatusInternalServerError, result.Err)
		}
		c.JSON(http.StatusOK, result.OK.WithData(tag))
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
		c.JSON(http.StatusBadRequest, result.ErrParam)
	}
}

//删除文章标签

func RemoveTag(c *gin.Context) {
	c.JSON(http.StatusOK, result.OK.WithMsg("删除文章标签成功"))
}

//根据id获取文章标签

func GetTagById(c *gin.Context) {
	c.JSON(http.StatusOK, result.OK.WithMsg("根据id获取文章标签成功"))
}

func EditTag(c *gin.Context) {
	c.JSON(http.StatusOK, result.OK.WithMsg("修改文章标签成功"))
}
