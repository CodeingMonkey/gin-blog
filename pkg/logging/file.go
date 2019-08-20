package logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	//日志保存路径
	LogSavePath = "runtime/logs/"

	//日志名称
	LogSaveName = "log"

	//日志后缀
	LogFileExt  = "log"

	//时间格式
	TimeFormat  = "20060102"
)

//获取日志文件路径
func getLogFilePath() string {
	return fmt.Sprintf("%s", LogSavePath)
}

//获取完整的日志文件路径（log开头+时间+文件后缀+文件路径）
func getLogFileFullPath() string {
	prefixPath := getLogFilePath()
	suffixPath := fmt.Sprintf("%s%s.%s", LogSaveName, time.Now().Format(TimeFormat), LogFileExt)

	//根据指定的格式返回一个字符串
	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
}

/**
打开日志文件，返回一个文件类型的指针
 */
func openLogFile(filePath string) *os.File {

	//返回文件信息结构描述文件。如果出现错误，会返回*PathError
	_, err := os.Stat(filePath)
	switch {
	//判断文件或目录不存在，不存在返回false
	case os.IsNotExist(err):
		mkDir()

		//判断权限是否满足，不满足返回false
	case os.IsPermission(err):
		log.Fatalf("Permission :%v", err)
	}

	//调用文件，支持传入文件名称、指定的模式调用文件、文件权限，返回的文件的方法可以用于I/O。如果出现错误，则为*PathError
	//|是位运算符，基于二进制的
	fmt.Println("打开文件标志")
	fmt.Println(os.O_APPEND|os.O_CREATE|os.O_WRONLY)
	handle, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		//打印错误日志，程序会结束
		log.Fatalf("Fail to OpenFile :%v", err)
	}

	return handle
}

func mkDir() {
	//返回与当前目录对应的根路径名
	dir, _ := os.Getwd()

	//创建对应的目录以及所需的子目录，若成功则返回nil，否则返回error，文件已经存在的情况，不作处理，返回nil也不会报错
	err := os.MkdirAll(dir+"/"+getLogFilePath(), os.ModePerm)
	if err != nil {
		//捕获错误
		panic(err)
	}
}
