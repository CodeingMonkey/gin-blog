package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/gin-blog/pkg/e"
	"github.com/gin-blog/pkg/util"
)

/**
返回一个HandlerFunc类型值
 */
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		/**
		初始化返回码200
		 */
		code = e.SUCCESS
		token := c.Query("token")

		/**
		判断token是否为空
		 */
		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			/**
			解析token，验证token是否合法，获取错误信息+Claims实例化（包含用户名/密码/过期时间）
			 */
			_, err := util.ParseToken(token)

			/**
			token过期情况，err依然不是nil
			 */
			if err != nil {

				/**
				判断错误类型
				 */
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired://token超时
					code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
				default://token超时以外的其他错误
					code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
				}
			}
		}

		/**
		验证不通过的情况，直接返回错误信息
		 */
		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
