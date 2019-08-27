package logging

import (
	"fmt"
	"github.com/gin-blog/pkg/file"
	"github.com/gin-blog/pkg/setting"
	"os"
	"time"
)

//var (
//	//日志保存路径
//	LogSavePath = "dd"
//
//	//日志名称
//	LogSaveName = setting.AppSetting.LogSaveName
//
//	//日志后缀
//	LogFileExt  = setting.AppSetting.LogFileExt
//
//	//时间格式
//	TimeFormat  = setting.AppSetting.TimeFormat
//)

//获取日志文件路径
func getLogFilePath() string {
	return fmt.Sprintf("%s%s", setting.AppSetting.RuntimeRootPath, setting.AppSetting.LogSavePath)
}

func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s",
		setting.AppSetting.LogSaveName,
		time.Now().Format(setting.AppSetting.TimeFormat),
		setting.AppSetting.LogFileExt,
	)
}

func openLogFile(fileName, filePath string) (*os.File, error) {

	//返回与当前目录对应的根路径名,如果可以通过多个路径访问当前目录（由于符号链接），Getwd可能会返回其中任何一个
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	//根路径+日志目录
	src := dir + "/" + filePath

	//检测文件夹权限
	perm := file.CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	//判断文件夹是否存在，不存在则创建文件夹
	err = file.IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	//打开日志文件
	f, err := file.Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}

	return f, nil
}
