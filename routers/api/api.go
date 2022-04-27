package api

import (
	"github.com/gin-gonic/gin"
	v1 "gogin/controller"
)

func InitApi(group *gin.RouterGroup) {
	//标签组路由
	tagGroup := group.Group("/tag")
	// {} 是书写规范
	//articleGroup := group.Group("/article")
	//authGroup := group.Group("/auth")
	//获取标签列表
	tagGroup.GET("/getTags", v1.GetTags)
	//获取标签
	tagGroup.GET("/getTag/:id", v1.GetTagById)
	//新建标签
	tagGroup.POST("/addTag", v1.AddTag)
	//更新指定标签
	tagGroup.PUT("/updateTag/:id", v1.EditTag)
	//删除指定标签
	tagGroup.DELETE("/delTag/:id", v1.RemoveTag)

	//文件组路由
	//multipart/form-data格式用于文件上传
	//gin文件上传与原生的net/http方法类似，不同在于gin把原生的request封装到c.Request中
	fileGroup := group.Group("/file")
	fileGroup.POST("/upload")
	fileGroup.GET("/download")
}
