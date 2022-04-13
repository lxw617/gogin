package service

import (
	"github.com/gin-gonic/gin"
	result "gogin/common"
	"gogin/models"
	"net/http"
)

//获取所有文章标签
func GetTags(params map[string]interface{}, pageNum int, pageSize int) ([]models.Tag, int, error) {
	/*scopes := make([]func(*gorm.DB) *gorm.DB, 0)
	name, ok := params["name"].(string)
	if ok {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("name like ?", "%"+name+"%")
		})
	}
	err := models.DB.Model(&models.Tag{}).Scopes(scopes...).Scan(&tags).Error*/
	count := 0
	tags := make([]models.Tag, 0)
	err := models.DB.Model(&models.Tag{}).Where(params).Offset(pageNum).Limit(pageSize).Scan(&tags).Count(&count).Error
	return tags, count, err
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
