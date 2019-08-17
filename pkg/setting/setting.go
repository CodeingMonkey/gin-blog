package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

var (
	Cfg *ini.File

	RunMode string

	HTTPPort int

	/**
	time.Duration底层使用了int64类型
	 */
	ReadTimeout time.Duration
	WriteTimeout time.Duration

	PageSize int
	JwtSecret string
)

func init() {
	var err error

	/**
	加载.ini配置文件，返回两个参数，一个是错误信息，一个是指向文件的指针
	参数可以是字符串类型的文件名或者是字节切片
	 */
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}

	/**
	读取ini默认分区的的key为RUN_MODE的值
	 */
	LoadBase()

	/**
	读取ini文件的server分区下的内容
	[server]：ini文件的分区格式如上，[]中间是分区的名称
	 */
	LoadServer()

	/**
	读取ini文件的app分区下面的内容
	 */
	LoadApp()
}

func LoadBase() {
	/**
	由 Must 开头的方法名允许接收一个相同类型的参数来作为默认值
	当键不存在或者转换失败时，则会直接返回该默认值
	MustString 方法必须传递一个默认值
	*/
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
}

func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}

	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")

	HTTPPort = sec.Key("HTTP_PORT").MustInt(8000)
	/**
	ini文件的读取时间和写入时间都是int类型，读取之后，使用time.Duration处理成time.Duration类型
	 */
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeout =  time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}

func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}

	JwtSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
	PageSize = sec.Key("PAGE_SIZE").MustInt(10)
}