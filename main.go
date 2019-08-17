package main

import (
	"fmt"
	"github.com/gin-blog/routers"
	"net/http"

	"github.com/gin-blog/pkg/setting"
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

	router := routers.InitRouter()


	//实例化一个服务器（地址/端口号/读取超时/写入超时/header头最大字节数）
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	//启动服务器，监听指定端口号
	s.ListenAndServe()
}