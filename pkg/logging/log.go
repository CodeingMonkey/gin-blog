package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Level int

var (
	F *os.File

	DefaultPrefix      = ""
	DefaultCallerDepth = 1

	logger     *log.Logger
	logPrefix  = ""
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

/**
iota模式，Debug为0，往下+1继续，直到遇到下一个const，iota重新设置为0
 */
const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

/**
main函数之前执行的函数
 */
func Setup() {
	var err error

	//获取日志地址
	filePath := getLogFilePath()

	//获取日志文件名+后缀
	fileName := getLogFileName()

	//打开日志文件
	F, err = openLogFile(fileName, filePath)
	if err != nil {
		log.Fatalln(err)
	}

	logger = log.New(F, DefaultPrefix, log.LstdFlags)
}

/**
可变参数
 */
func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	logger.Println(v)
}

func Info(v ...interface{}) {
	setPrefix(INFO)
	logger.Println(v)
}

func Warn(v ...interface{}) {
	setPrefix(WARNING)
	logger.Println(v)
}

func Error(v ...interface{}) {
	setPrefix(ERROR)
	logger.Println(v)
}

func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	logger.Fatalln(v)
}

/**
设置文件前缀
 */
func setPrefix(level Level) {

	/**
	DefaultCallerDepth：0表示调用runtime.Caller()所在的位置，1表示runtime.Caller()所在函数的调用位置，依此类推
	 */
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	logger.SetPrefix(logPrefix)
}
