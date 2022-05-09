package controller

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	result "gogin/common"
	img "gogin/pkg/image"
	"net/http"
	"os"
)

func TestFile(c *gin.Context) {
	file, image, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	//16cfd547d837272e28cad04b5c42fa6e.png
	imageName := img.GetImageName(image.Filename)
	fullPath := img.GetImageFullPath()
	savePath := img.GetImagePath()
	fmt.Println(file)
	fmt.Println(image)
	fmt.Println(imageName)
	fmt.Println(fullPath)
	fmt.Println(savePath)
}

//文件上传

func UploadOSS(c *gin.Context) {
	_, image, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	//更改后的图片名称
	imageName := img.GetImageName(image.Filename)
	//原图片名称
	originalImageName := image.Filename

	// 创建OSSClient实例。
	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	client, err := oss.New("<yourEndpoint>", "<yourAccessKeyId>", "<yourAccessKeySecret>")

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 填写存储空间名称，例如examplebucket。
	bucket, err := client.Bucket("srb-1130")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	url := "E:\\jianyu\\image\\" + originalImageName
	// 依次填写Object的完整路径（例如exampledir/exampleobject.txt）和本地文件的完整路径（例如D:\\localpath\\examplefile.txt）。
	err = bucket.PutObjectFromFile("gogin/"+imageName, url)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	//https://srb-1130.oss-cn-shenzhen.aliyuncs.com/gogin/2420AAD5-903F-4374-953E-82126615AEEA.png
	urlResult := "https://srb-1130.oss-cn-shenzhen.aliyuncs.com/gogin/" + imageName
	c.JSON(http.StatusOK, result.OK.WithData(urlResult))
}

//文件下载

func DownloadOSS(c *gin.Context) {
	image := c.Param("image")
	valid := validation.Validation{}
	valid.Required(image, "image").Message("图片地址不能为空")
	if valid.HasErrors() {
		c.JSON(http.StatusOK, result.OK.WithMsg("图片地址参数错误"))
	}
	// 创建OSSClient实例。
	// yourEndpoint填写Bucket所在地域对应的Endpoint。以华东1（杭州）为例，Endpoint填写为https://oss-cn-hangzhou.aliyuncs.com。
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	client, err := oss.New("<yourEndpoint>", "<yourAccessKeyId>", "<yourAccessKeySecret>")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 填写Bucket名称，例如examplebucket。
	bucket, err := client.Bucket("srb-1130")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 下载文件到本地文件，并保存到指定的本地路径中。如果指定的本地文件存在会覆盖，不存在则新建。
	// 如果未指定本地路径，则下载后的文件默认保存到示例程序所属项目对应本地路径中。
	// 依次填写Object完整路径（例如exampledir/exampleobject.txt）和本地文件的完整路径(例如D:\\localpath\\examplefile.txt)。Object完整路径中不能包含Bucket名称。

	err = bucket.GetObjectToFile("gogin/"+image, "E:\\jianyu\\image")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
}

//文件删除

func DeleteOSS(c *gin.Context) {
	imageName := c.Param("imageNamw")
	valid := validation.Validation{}
	valid.Required(imageName, "image").Message("图片名称不能为空")
	if valid.HasErrors() {
		c.JSON(http.StatusOK, result.OK.WithMsg("图片名称参数错误"))
	}
	// 创建OSSClient实例。
	// yourEndpoint填写Bucket所在地域对应的Endpoint。以华东1（杭州）为例，Endpoint填写为https://oss-cn-hangzhou.aliyuncs.com。
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	client, err := oss.New("<yourEndpoint>", "<yourAccessKeyId>", "<yourAccessKeySecret>")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 填写Bucket名称，例如examplebucket。
	bucket, err := client.Bucket("srb-1130")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 下载文件到本地文件，并保存到指定的本地路径中。如果指定的本地文件存在会覆盖，不存在则新建。
	// 如果未指定本地路径，则下载后的文件默认保存到示例程序所属项目对应本地路径中。
	// 依次填写Object完整路径（例如exampledir/exampleobject.txt）和本地文件的完整路径(例如D:\\localpath\\examplefile.txt)。Object完整路径中不能包含Bucket名称。

	//"https://srb-1130.oss-cn-shenzhen.aliyuncs.com/gogin/16cfd547d837272e28cad04b5c42fa6e.png"
	//err = bucket.DeleteObject(key, oss.VersionId("yourObjectVersionId"))
	err = bucket.DeleteObject("gogin/" + imageName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	c.JSON(http.StatusOK, result.OK.WithMsg("删除成功"))
}
