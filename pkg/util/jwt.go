package util

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/gin-blog/pkg/setting"
)

/**
byte[setting.JwtSecret]将setting.JwtSecret转换成byte字节
将string转换成byte切片
 */
var jwtSecret = []byte(setting.App{}.JwtSecret)

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

/**
生成token
 */
func GenerateToken(username, password string) (string, error) {
	/**
	当前时间
	 */
	nowTime := time.Now()

	/**
	当前时间+3个小时
	 */
	expireTime := nowTime.Add(1 * time.Hour)

	claims := Claims{
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "gin-blog",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

/**
token解析，返回Claims指针+错误值
 */
func ParseToken(token string) (*Claims, error) {
	/**
	初步解析，获取token的实例化
	*/
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	/**
	判断token存在的情况，判断token是否有效
	 */
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
