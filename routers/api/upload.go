package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gin-blog/pkg/e"
	"github.com/gin-blog/pkg/logging"
	"github.com/gin-blog/pkg/upload"
)

func UploadImage(c *gin.Context) {
	code := e.SUCCESS
	data := make(map[string]string)

	file, image, err := c.Request.FormFile("image")

	//没有图片，获取图片上传失败
	if err != nil {
		logging.Warn(err)
		code = e.ERROR
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": data,
		})
	}

	//图片不存在
	if image == nil {
		code = e.INVALID_PARAMS
	} else {

		//获取上传图片的图片名称
		imageName := upload.GetImageName(image.Filename)

		fmt.Println("图片名称", imageName)

		//获取图片完整路径
		fullPath := upload.GetImageFullPath()
		fmt.Println("图片完整路径", fullPath)

		//获取图片保存路径
		savePath := upload.GetImagePath()
		fmt.Println("图片保存路径", savePath)

		//获取完整路径+文件名
		src := fullPath + imageName
		fmt.Println("文件名", src)

		//检测文件后缀和文件大小
		if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
			code = e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
		} else {

			//检测图片
			err := upload.CheckImage(fullPath)
			if err != nil {
				logging.Warn(err)
				code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
			} else if err := c.SaveUploadedFile(image, src); err != nil {
				logging.Warn(err)
				code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
			} else {
				data["image_url"] = upload.GetImageFullUrl(imageName)
				data["image_save_url"] = savePath + imageName
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}
