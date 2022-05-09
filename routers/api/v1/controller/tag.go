package controller

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/unknwon/com"
	result "gogin/common"
	"gogin/models"
	"gogin/pkg/logging"
	gredis "gogin/pkg/redis"
	"gogin/pkg/setting"
	"gogin/pkg/util"
	service2 "gogin/routers/api/v1/service"
	"net/http"
	"reflect"
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
	tags, count, err := service2.GetTags(params, util.GetPage(c), setting.AppSetting.PageSize) //util.GetPage保证了各接口的page处理是一致的
	//将结果使用 hset 格式存入redis缓存中
	//redis.HSet()
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
		tag, err := service2.AddTag(tag)
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
	if _, err := service2.GetTagById(tagId); err != nil {
		c.JSON(http.StatusInternalServerError, result.Err)
	}
	if !valid.HasErrors() {
		//删除文章标签
		if err := service2.DeleteById(tagId); err != nil {
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
		//获取连接池的一个连接
		conn := gredis.RedisPool.Get()
		//获取结果
		res, err := conn.Do("HGET", "tagset", tagId)
		tag := &models.Tag{}
		rebytes, _ := redis.Bytes(conn.Do("HGET", "tagset", tagId))

		//如果存入redis缓存中 string(tagJson)
		//res, err = conn.Do("HSET", "tagset", tagId, string(tagJson))
		//读取时应该使用以下方式
		//reStr, _ := redis.String(conn.Do("HGET", "tagset", tagId))
		//r := []byte(reStr)
		//err = json.Unmarshal(r,tag)

		//格式转换
		json.Unmarshal(rebytes, tag)
		//如果缓存中有
		if res != nil {
			fmt.Println("从缓存中获取")
			if err != nil {
				fmt.Println("redis HGET error:", err)
			} else {
				res_type := reflect.TypeOf(res)
				fmt.Printf("res type : %s \n", res_type)
				fmt.Printf("res  : %s \n", res)
			}
			c.JSON(http.StatusOK, result.OK.WithData(tag))
		} else {
			//没有查数据库
			fmt.Println("从数据库中获取")
			//从数据库根据id获取文章标签
			tag, err := service2.GetTagById(tagId)
			tagJson, err := json.Marshal(tag)
			if err == nil {
				fmt.Println(string(tagJson))
			}
			//存入redis缓存中
			res, err = conn.Do("HSET", "tagset", tagId, tagJson) //string(tagJson)
			if err != nil {
				fmt.Println("redis mset error:", err)
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, result.Err)
			}
			c.JSON(http.StatusOK, result.OK.WithData(tag))
		}

		/*res, err = redis.String(conn.Do("HGET", "tagset", tagId))
		if err != nil {
			fmt.Println("redis HGET error:", err)
		}*/
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
	tag, _ := service2.GetTagById(id)
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
		tag, err := service2.UpdateTag(tag, id)
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
