package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-blog/models"
	"github.com/gin-blog/pkg/e"
	"github.com/gin-blog/pkg/util"
	"github.com/gin-gonic/gin"
)

/**
生米高结构的时候，直接打上标记，valid，必传参数，最大长度为50
 */
type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetAuth(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	valid := validation.Validation{}
	a := auth{Username: username, Password: password}

	/**
	这种方式就直接验证了username和password必传，且最大长度为50
	 */
	ok, _ := valid.Valid(&a)

	data := make(map[string]interface{})
	code := e.INVALID_PARAMS
	if ok {
		isExist := models.CheckAuth(username, password)
		if isExist {

			/**
			生成token，生成的时候，同时设置了token有效期，
			 */
			token, err := util.GenerateToken(username, password)
			fmt.Println(err)
			if err != nil {
				code = e.ERROR_AUTH_TOKEN
			} else {
				data["token"] = token

				code = e.SUCCESS
			}

		} else {
			code = e.ERROR_AUTH
		}
	} else {

		/**
		验证出现错误的情况，遍历打印错误信息
		 */
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}
