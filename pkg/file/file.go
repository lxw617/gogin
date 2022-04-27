package file

import (
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
)

/*
封装了7个方法
GetSize：获取文件大小
GetExt：获取文件后缀
CheckExist：检查文件是否存在
CheckPermission：检查文件权限
IsNotExistMkDir：如果不存在则新建文件夹
MkDir：新建文件夹
Open：打开文件
在这里我们用到了 mime/multipart 包，它主要实现了 MIME 的 multipart 解析，主要适用于 HTTP 和常见浏览器生成的 multipart 主体
multipart 又是什么，rfc2388 的 multipart/form-data 了解一下 https://datatracker.ietf.org/doc/html/rfc2388
*/

func GetSize(f multipart.File) (int, error) {
	content, err := ioutil.ReadAll(f)

	return len(content), err
}

func GetExt(fileName string) string {
	return path.Ext(fileName)
}

func CheckExist(src string) bool {
	_, err := os.Stat(src)

	return os.IsNotExist(err)
}

func CheckPermission(src string) bool {
	_, err := os.Stat(src)

	return os.IsPermission(err)
}

func IsNotExistMkDir(src string) error {
	if exist := CheckExist(src); exist == false {
		if err := MkDir(src); err != nil {
			return err
		}
	}

	return nil
}

func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}
