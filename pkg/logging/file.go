package logging

import (
	"fmt"
	"github.com/gin-blog/pkg/setting"
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
