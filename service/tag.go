package service

import (
	"github.com/gin-gonic/gin"
	"gogin/models"
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
	if err := models.DB.Model(&models.Tag{}).Where(params).Offset(pageNum).Limit(pageSize).Scan(&tags).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return tags, count, nil
}

//新增文章标签

func AddTag(name string, state int, createdBy int) (*models.Tag, error) {
	tag := &models.Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	}
	if err := models.DB.Model(&models.Tag{}).Create(&tag).Error; err != nil {
		return nil, err
	}
	return tag, nil
}

//删除文章标签

func RemoveTag(c *gin.Context) {

}

//根据id获取文章标签

func GetTagById(c *gin.Context) {

}

//修改文章标签

func EditTag(c *gin.Context) {

}
