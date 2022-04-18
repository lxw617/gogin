package controller

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	result "gogin/common"
	"gogin/pkg/util"
	"gogin/service"
	"net/http"
)

//生成token
// @Summary Get authToken
// @Produce  json
// @Param username query string false "username"
// @Param password query string false "password"
// @Success 200 {object} map[string][]string
// @Failure 500 {object} map[string][]string
// @Router /auth [get]
func GetToken(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	//参数校验
	valid := validation.Validation{}
	valid.Required(username, "username").Message("名称不能为空")
	valid.Required(password, "password").Message("密码不能为空")

	if !valid.HasErrors() {
		//调用方法
		isExist := service.CheckAuth(username, password)
		if isExist {
			token, err := util.GenerateToken(username, password)
			token = "Bearer " + token
			if err != nil {
				c.JSON(http.StatusOK, result.ErrAuthToken)
				return
			}
			c.JSON(http.StatusOK, result.OK.WithData(token))
		}
	}
}
