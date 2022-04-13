package models

type Tag struct {
	Model            //对于匿名字段，GORM 会将其字段包含在父结构体中
	Name      string `json:"name"`
	State     int    `json:"state"`
	CreatedBy int    `json:"createdBy"`
}
