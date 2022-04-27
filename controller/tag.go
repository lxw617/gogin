package controller

import (
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	result "gogin/common"
	"gogin/models"
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
// @Router /api/v1/tag/getTags [get]
func GetTags(c *gin.Context) { //c *gin.Context是Gin很重要的组成部分，可以理解为上下文，它允许我们在中间件之间传递变量、管理流、验证请求的JSON和呈现JSON响应

	username, _ := c.Get("username")
	fmt.Sprintln(username)
	//初始化参数
	params := make(map[string]interface{})
	//c.Query使用现有的底层请求对象解析查询字符串参数 请求响应 url 匹配：/tag?name=Jane
	//URL参数 可以通过DefaultQuery()或Query()方法获取，即带问号传参
	//DefaultQuery()若参数不存在，返回默认值，Query()若不存在，返回空串
	name := c.Query("name")
	if name != "" {
		params["name"] = name
	}
	state := com.StrTo(c.DefaultQuery("state", "1")).MustInt()
	/*if arg := c.Query("state"); arg != "" {
		// Convert string to specify type.
		state = com.StrTo(arg).MustInt()
	}*/
	params["state"] = state
	//获取结果
	tags, count, err := service.GetTags(params, util.GetPage(c), setting.AppSetting.PageSize) //util.GetPage保证了各接口的page处理是一致的
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
// @Router /api/v1/tag/addTag [post]
func AddTag(c *gin.Context) {

	/*
		表单传输为post请求，http常见的传输格式为四种：
		application/json
		application/x-www-form-urlencoded
		application/xml
		multipart/form-data
		表单参数可以通过PostForm()方法获取，该方法默认解析的是x-www-form-urlencoded或from-data格式的参数
		types := c.DefaultPostForm("type", "post")
		username := c.PostForm("username")
		password := c.PostForm("password")

		Json 数据解析和绑定 ShouldBindJSON()
		表单数据解析和绑定 Bind()默认解析并绑定form格式 根据请求头中content-type自动推断
		URI数据解析和绑定 ShouldBindUri()
	*/

	//{
	//    "name":"总结",
	//    "createdBy":1
	//}
	//声明接收的变量
	var tag *models.Tag
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&tag); err != nil {
		// 返回错误信息 gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//参数校验
	valid := validation.Validation{}
	valid.Required(tag.Name, "name").Message("名称不能为空")
	valid.MaxSize(tag.Name, 100, "name").Message("名称最长不能超过100字符")
	valid.Required(tag.CreatedBy, "createdBy").Message("创建人不能为空")
	if !valid.HasErrors() {
		//调用方法
		tag, err := service.AddTag(tag)
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
// @Summary Delete article tag
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} map[string][]string
// @Failure 500 {object} map[string][]string
// @Router /api/v1/tag/delTag/{id} [delete]
func RemoveTag(c *gin.Context) {
	tagId := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(tagId, 1, "id").Message("ID必须大于0")
	//判断是否存在
	if _, err := service.GetTagById(tagId); err != nil {
		c.JSON(http.StatusInternalServerError, result.Err)
	}
	if !valid.HasErrors() {
		//删除文章标签
		if err := service.DeleteById(tagId); err != nil {
			c.JSON(http.StatusInternalServerError, result.Err)
		}
		c.JSON(http.StatusOK, result.OK.WithMsg("删除文章标签成功"))
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
		c.JSON(http.StatusBadRequest, result.ErrParam)
	}
}

//根据id获取文章标签
// @Summary Delete article tag
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} map[string][]string
// @Failure 500 {object} map[string][]string
// @Router /api/v1/tag/getTag/{id} [get]
func GetTagById(c *gin.Context) {
	tagId := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(tagId, 1, "id").Message("ID必须大于0")
	if !valid.HasErrors() {
		//删除文章标签
		tag, err := service.GetTagById(tagId)
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

// @Summary Update article tag
// @Produce  json
// @Param id path int true "ID"
// @Param name body string true "Name"
// @Param state body int false "State"
// @Param modified_by body string true "ModifiedBy"
// @Success 200 {object} map[string][]string
// @Failure 500 {object} map[string][]string
// @Router /api/v1/tag/updateTag/{id} [put]
func EditTag(c *gin.Context) {
	name := c.Query("name")
	//c.DefaultQuery则支持设置一个默认值
	state := com.StrTo(c.DefaultQuery("state", "1")).MustInt()
	createdBy := com.StrTo(c.Query("createdBy")).MustInt()
	id := com.StrTo(c.Query("id")).MustInt()
	//获取tag
	tag, _ := service.GetTagById(id)
	//参数校验
	valid := validation.Validation{}
	if name != "" {
		valid.Required(name, "name").Message("名称不能为空")
		valid.MaxSize(name, 100, "name").Message("名称最长不能超过100字符")
		tag.Name = name
	}
	if createdBy > 0 {
		valid.Required(createdBy, "createdBy").Message("创建人不能为空")
		tag.CreatedBy = createdBy
	}
	if id > 0 {
		valid.Required(id, "id").Message("id不能为空")
	}
	if state > 0 && state < 2 {
		valid.Range(state, 0, 1, "status").Message("状态只能为0或者1")
		tag.State = state
	}
	if !valid.HasErrors() {
		//调用方法
		tag, err := service.UpdateTag(tag, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, result.Err)
		}
		c.JSON(http.StatusOK, result.OK.WithData(tag))
		//c.JSON(http.StatusOK, result.OK.WithMsg("修改文章标签成功"))
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
		c.JSON(http.StatusBadRequest, result.ErrParam)
	}
}

/*
返回结果：各种数据格式的响应
json、结构体、XML、YAML类似于java的properties、ProtoBuf
// 多种响应方式
func main() {
    // 1.创建路由
    // 默认使用了2个中间件Logger(), Recovery()
    r := gin.Default()
    // 1.json
    r.GET("/someJSON", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "someJSON", "status": 200})
    })
    // 2. 结构体响应
    r.GET("/someStruct", func(c *gin.Context) {
        var msg struct {
            Name    string
            Message string
            Number  int
        }
        msg.Name = "root"
        msg.Message = "message"
        msg.Number = 123
        c.JSON(200, msg)
    })
    // 3.XML
    r.GET("/someXML", func(c *gin.Context) {
        c.XML(200, gin.H{"message": "abc"})
    })
    // 4.YAML响应
    r.GET("/someYAML", func(c *gin.Context) {
        c.YAML(200, gin.H{"name": "zhangsan"})
    })
    // 5.protobuf格式,谷歌开发的高效存储读取的工具
    // 数组？切片？如果自己构建一个传输格式，应该是什么格式？
    r.GET("/someProtoBuf", func(c *gin.Context) {
        reps := []int64{int64(1), int64(2)}
        // 定义数据
        label := "label"
        // 传protobuf格式数据
        data := &protoexample.Test{
            Label: &label,
            Reps:  reps,
        }
        c.ProtoBuf(200, data)
    })

    r.Run(":8000")
}

*/
/*
Validator 是基于 tag（标记）实现结构体和单个字段的值验证库，它包含以下功能：

使用验证 tag（标记）或自定义验证器进行跨字段和跨结构体验证。
关于 slice、数组和 map，允许验证多维字段的任何或所有级别。
能够深入 map 键和值进行验证。
通过在验证之前确定接口的基础类型来处理类型接口。
处理自定义字段类型（如 sql 驱动程序 Valuer）。
别名验证标记，它允许将多个验证映射到单个标记，以便更轻松地定义结构体上的验证。
提取自定义的字段名称，例如，可以指定在验证时提取 JSON 名称，并在生成的 FieldError 中使用该名称。
可自定义 i18n 错误消息。
Web 框架 gin 的默认验证器。

变量验证
Var 方法使用 tag（标记）验证方式验证单个变量。
func (*validator.Validate).Var(field interface{}, tag string) error
它接收一个 interface{} 空接口类型的 field 和一个 string 类型的 tag，返回传递的非法值得无效验证错误，否则将 nil 或 ValidationErrors 作为错误。如果错误不是 nil，则需要断言错误去访问错误数组，例如：
validationErrors := err.(validator.ValidationErrors)
如果是验证数组、slice 和 map，可能会包含多个错误。

结构体验证
结构体验证结构体公开的字段，并自动验证嵌套结构体，除非另有说明。
func (*validator.Validate).Struct(s interface{}) error
它接收一个 interface{} 空接口类型的 s，返回传递的非法值得无效验证错误，否则将 nil 或 ValidationErrors 作为错误。如果错误不是 nil，则需要断言错误去访问错误数组，例如：
validationErrors := err.(validator.ValidationErrors)
实际上，Struct 方法是调用的 StructCtx 方法


*/
