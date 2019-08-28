package export

import "github.com/gin-blog/pkg/setting"

//获取图片完整地址（域名+文件路径+文件名称）
func GetExcelFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetExcelPath() + name
}

//获取文件路径
func GetExcelPath() string {
	return setting.AppSetting.ExportSavePath
}

//获取完整的文件路径
func GetExcelFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetExcelPath()
}