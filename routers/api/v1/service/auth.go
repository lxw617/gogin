package service

import "gogin/models"

//检查用户是否存在

func CheckAuth(username string, password string) bool {
	var auth models.Auth
	if err := models.DB.Model(&models.Auth{}).Where("username = ?", username).Where("password = ?", password).First(&auth).Error; err != nil {
		return false
	}
	return true
}
