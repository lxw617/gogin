package excel

import (
	"github.com/tealeg/xlsx"
	"os"
	"path/filepath"
	"strings"
)

// 生成excel表头,并返回文件路径以及名字
func GetExcelHeader(fileName string, values []interface{}, tableName string) (sheet *xlsx.Sheet, file *xlsx.File, url string, name string, err error) {
	sheet, file, url, name, err = GetExcel(fileName, tableName)
	if err != nil {
		return nil, nil, "", "", err
	}
	sheet.AddRow().WriteSlice(&values, -1)
	return
}

func GetExcel(fileName string, tableName string) (sheet *xlsx.Sheet, file *xlsx.File, url string, name string, err error) {
	//获取当前项目路径
	root, _ := os.Getwd()
	//获取项目上一层目录
	root = filepath.Dir(root)
	//对路径进行处理 "E:\\jianyu\\xlsx\\" + fileName
	root = strings.Replace(root, "\\", "\\\\", -1)
	path := root + "\\\\xlsx"
	//os.Stat：返回文件信息结构描述文件。如果出现错误，会返回*PathError
	_, err = os.Stat(path)
	if exist := os.IsNotExist(err); exist == true {
		//os.MkdirAll：创建对应的目录以及所需的子目录，若成功则返回nil，否则返回error
		err = os.MkdirAll(path, os.ModePerm)
	}
	name = fileName + ".xlsx"
	url = path + "\\\\" + name
	file = xlsx.NewFile()
	sheet, err = file.AddSheet(tableName)
	return
}
