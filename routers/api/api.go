package api

import (
	"github.com/gin-gonic/gin"
	v1 "gogin/controller"
)

func InitApi(group *gin.RouterGroup) {
	tagGroup := group.Group("/tag")
	//articleGroup := group.Group("/article")
	//authGroup := group.Group("/auth")
	//获取标签列表
	tagGroup.GET("/tags", v1.GetTags)
	//获取标签列表
	tagGroup.GET("/tag", v1.GetTagById)
	//新建标签
	tagGroup.POST("/tag", v1.AddTag)
	//更新指定标签
	tagGroup.PUT("/tag/:id", v1.EditTag)
	//删除指定标签
	tagGroup.DELETE("/tag/:id", v1.RemoveTag)
}
