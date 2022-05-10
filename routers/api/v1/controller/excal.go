package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	result "gogin/common"
	"gogin/models"
	"gogin/pkg/excel"
	"gogin/routers/api/v1/service"
	"net/http"
	"strconv"
	"time"
)

func ImportExcel(c *gin.Context) {
	//获取文件
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.Err)
	}
	//fmt.Println(file)
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.Err.WithMsg("导入excel文件失败"))
	}

	var tags []*models.Tag
	rows, _ := xlsx.GetRows("标签列表")
	for irow, row := range rows {
		if irow > 0 {
			var data []string
			for _, cell := range row {
				data = append(data, cell)
			}
			createdAt, _ := strconv.Atoi(data[2])
			state, _ := strconv.Atoi(data[3])
			//fmt.Println(data)
			tag := &models.Tag{
				Name:      data[1],
				CreatedBy: createdAt,
				State:     state,
			}
			tags = append(tags, tag)
			//service.AddTag(tag)
		}
	}
	//c.JSON(http.StatusOK, result.OK.WithData(tags))
	c.JSON(http.StatusOK, result.OK.WithMsg("导入excel文件成功"))
	/*// 定义取设备号正则
	reg, _ := regexp.Compile("SD\\d{8}|MN\\w{4}\\d{4}")
	// 读取文件内容
	wb, err := xlsx.OpenFile(header.Filename)
	if err != nil {
		panic(err)
	}
	sh := wb.Sheets[0]
	maxRow := sh.MaxRow
	maxCol := sh.MaxCol
	for row := 0; row < maxRow; row++ {
		for col := 0; col < maxCol; col++ {
			cell := sh.Cell(row, col)
			str, _ := cell.FormattedValue()
			if str == "" {
				continue
			}
			arr := reg.FindAllString(str, -1)
			for _, item := range arr {
				serialMap[item] = true
			}
		}
	}*/

}
func ExportExcel(c *gin.Context) {
	//初始化参数
	params := make(map[string]interface{})
	tags, _, err := service.GetTags(params, -1, -1)
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
		c.JSON(http.StatusInternalServerError, result.Err.WithMsg("导出excel文件失败"))
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
				c.JSON(http.StatusInternalServerError, result.Err.WithMsg("导出excel文件失败"))
				return
			}
		}
	}
	// "E:\\jianyu\\xlsx\\" + fileName
	err = file.Save(fileUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.Err.WithMsg("导出excel文件保存失败"))
		return
	}
	c.JSON(http.StatusOK, result.OK.WithData(fileName))
	return
}
