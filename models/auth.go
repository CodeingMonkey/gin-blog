package models

type Auth struct {
	ID int `gorm:"primary_key" json:"id"`//gorm:primary_key表明id是主键
	Username string `json:"username"`
	Password string `json:"password"`
}

func CheckAuth(username, password string) bool {
	var auth Auth

	/**
	根据用户名和密码查询对应的用户记录
	 */
	db.Select("id").Where(Auth{Username : username, Password : password}).First(&auth)
	if auth.ID > 0 {
		return true
	}

	return false
}