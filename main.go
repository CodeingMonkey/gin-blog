package main

import (
	"context"
	"fmt"
	"github.com/gin-blog/models"
	"os"
	"os/signal"

	//"github.com/fvbock/endless"
	"github.com/gin-blog/pkg/logging"
	"github.com/gin-blog/pkg/setting"
	"github.com/gin-blog/routers"
	"log"
	"net/http"
	//"os"
	//"os/signal"
	//"syscall"
	"time"
)

func main() {
	//router := gin.Default()
	//
	///**
	//设置接收路由，路由指向一个匿名函数，匿名函数返回指定信息
	//context是gin中的上下文，允许在中间件中传递变量，管理流，传递和管理json请求
	// */
	//router.GET("/test", func(c *gin.Context) {
	//
	//	/**
	//	g.H是一个map映射，key为string类型，value为interface类型，Go中任何类型都实现了空的interface，所以这个例子里，value可以是任何类型，本例中value是string类型
	//	 */
	//	c.JSON(200, gin.H{
	//		"message": "test",
	//	})
	//})

	/**
	endless实现热更新（fork方式，每次fork出一个新的进程，将原来的进程退出，但是已经在处理的进程会在进程处理完成，新的请求连接会进入到新的fork的进程中）'
	但是感觉作用不大，build出的可执行文件可以这样操作，go run 直接进行的话，只能第一次fork成功
	build可执行文件方式，只是重新启动了一个进程，并没有加载新文件的内容
	*/
	//
	//endless.DefaultReadTimeOut = setting.ReadTimeout
	//endless.DefaultWriteTimeOut = setting.WriteTimeout
	//endless.DefaultMaxHeaderBytes = 1 << 20
	//endPoint := fmt.Sprintf(":%d", setting.HTTPPort)
	//
	//server := endless.NewServer(endPoint, routers.InitRouter())
	//server.BeforeBegin = func(add string) {
	//	log.Printf("Actual pid is %d", syscall.Getpid())
	//}
	//
	//err := server.ListenAndServe()
	//if err != nil {
	//	log.Printf("Server err: %v", err)
	//}

	//注册所有路由
	setting.Setup()
	models.Setup()
	logging.Setup()
	router := routers.InitRouter()

	//实例化一个服务器（地址/端口号/读取超时/写入超时/header头最大字节数）
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),
		Handler:        router,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	//启动一个协程，实现端口监听
	go func() {
		fmt.Println("start Listening")

		/**
		listenAndServe端口监听
		*/
		if err := s.ListenAndServe(); err != nil {
			log.Printf("Listen: %s\n", err)
		}
	}()

	//创建一个管道，管道接收系统信号
	quit := make(chan os.Signal)

	//接收到SIGINT信号，管道写入一个数据
	signal.Notify(quit, os.Interrupt)

	//无缓冲通道，在没有SIGINT信号写入管道时，一直都堵在这里
	<-quit

	//接收到SIGINT信号的后续操作
	log.Println("Shutdown Server ...")

	/**
	函数返回的ctx是一个上下文环境，设置了超过5秒后，上下文环境自动销毁
	所以会等sleep结束以后，自动调用cancel函数，结束进程。
	*/
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	//不等待超时，手动关闭
	//cancel()

	//超时情况下，上下文环境会自动销毁
	//time.Sleep(6 * time.Second)
	//fmt.Println(ctx.Err())

	defer cancel()

	//关闭服务，如果上下文环境在关闭服务之前就销毁了，关闭服务会报错，但是服务依然会被关闭掉（没太想明白，为什么上下文环境已经关闭了，但是s.Shutdown没有报错）
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	} else {
		fmt.Println(err)
	}

	log.Println("Server exiting")

	//
	////启动服务器，监听指定端口号
	//s.ListenAndServe()
}
