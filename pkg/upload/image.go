package upload

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/gin-blog/pkg/file"
	"github.com/gin-blog/pkg/logging"
	"github.com/gin-blog/pkg/setting"
	"github.com/gin-blog/pkg/util"
)

//获取图片的完整URL
func GetImageFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetImagePath() + name
}

//获取图片名称+后缀
func GetImageName(name string) string {

	//获取文件后缀
	ext := path.Ext(name)

	//处理文件名（如果文件名不包括后缀的话，直接返回文件名），如果文件名包含后缀的话，处理文件名（取出指定的后缀）
	fileName := strings.TrimSuffix(name, ext)

	//文件名md5处理
	fileName = util.EncodeMD5(fileName)

	//返回文件名+后缀
	return fileName + ext
}

//获取图片保存路径
func GetImagePath() string {
	return setting.AppSetting.ImageSavePath
}

//获取图片完整路径
func GetImageFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetImagePath()
}

//检测文件后缀是否符合配置设置的文件后缀
func CheckImageExt(fileName string) bool {
	ext := file.GetExt(fileName)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}

	return false
}

//检测文件大小
func CheckImageSize(f multipart.File) bool {
	size, err := file.GetSize(f)
	fmt.Println("文件大小", size)
	if err != nil {
		log.Println(err)
		logging.Warn(err)
		return false
	}

	return size <= setting.AppSetting.ImageMaxSize
}

//检查图片
func CheckImage(src string) error {

	//获取当前目录所在的根目录
	dir, err := os.Getwd()
	fmt.Println("目录", dir)
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	//判断文件目录是否存在，不存在创建目录
	fmt.Println("真实名录", dir + "/" + src)
	err = file.IsNotExistMkDir(dir + "/" + src)
	fmt.Println("目录是否存在", err)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	//检测权限
	perm := file.CheckPermission(src)
	fmt.Println("权限", perm)
	if perm == true {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	return nil
}
