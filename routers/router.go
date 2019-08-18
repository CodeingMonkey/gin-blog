package routers

import (
	"fmt"
	"github.com/gin-blog/routers/api/v1"
	"github.com/gin-gonic/gin"

	"github.com/gin-blog/pkg/setting"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "test",
		})
	})

	/**
	创建一个路由组，路由组的路由可以具有相同的路由前缀或者中间件
	 */
	apiV1 := r.Group("api/v1")

	//{}代表作用域，作用内的变量只在作用域内有效，作用域内的a变量是没法在作用域外访问的
	{
		a:=1
		fmt.Println(a)
		//获取标签列表
		apiV1.GET("/tags", v1.GetTags)
		//新建标签
		apiV1.POST("/tags", v1.AddTag)
		//更新指定标签
		apiV1.PUT("/tags/:id", v1.EditTag)
		//删除指定标签
		apiV1.DELETE("/tags/:id", v1.DeleteTag)
	}
	//fmt.Println(a)


	return r
}
