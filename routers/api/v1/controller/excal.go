package controller

import (
	"github.com/gin-gonic/gin"
	result "gogin/common"
	"gogin/pkg/excel"
	service2 "gogin/routers/api/v1/service"
	"net/http"
	"time"
)

func ImportExcel(c *gin.Context) {

}
func ExportExcel(c *gin.Context) {
	//初始化参数
	params := make(map[string]interface{})
	tags, _, err := service2.GetTags(params, -1, -1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.Err)
		return
	}
	// 开始excel文件操作
	tableHead := []interface{}{"序号", "名称", "创建人", "状态"}
	tableName := "标签列表"
	fileName := "标签列表" + time.Now().Format("20060105150405") //+ strconv.FormatInt(time.Now().Local().Unix(), 10)
	sheet, file, fileUrl, fileName, err := excel.GetExcelHeader(fileName, tableHead, tableName)
	if err != nil {
		c.JSON(http.StatusOK, result.Err.WithMsg("导出excel文件失败"))
		return
	}
	if len(tags) != 0 {
		for index, item := range tags {
			row := sheet.AddRow()
			s := []interface{}{
				index + 1,
				item.Name,
				item.CreatedBy,
				item.State,
			}
			if row.WriteSlice(&s, -1) == 0 {
				c.JSON(http.StatusOK, result.Err.WithMsg("导出excel文件失败"))
				return
			}
		}
	}
	// "E:\\jianyu\\xlsx\\" + fileName
	err = file.Save(fileUrl)
	if err != nil {
		c.JSON(http.StatusOK, result.Err.WithMsg("导出excel文件保存失败"))
		return
	}
	c.JSON(http.StatusOK, result.OK.WithData(fileName))
	return
}
