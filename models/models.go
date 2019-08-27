package models

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/gin-blog/pkg/setting"
)

var db *gorm.DB

/**
暂时立即是所有model公用的字段，声明成了一个struct，在后续使用的model声明的时候
只需要把model作为嵌入结构嵌入到新的model的struct里即可
 */
type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
	DeletedOn int `json:"deleted_on"`
}

func Setup() {
	var (
		err                                               error
		dbType, dbName, user, password, host, tablePrefix string
	)

	dbType = setting.DatabaseSetting.Type
	dbName = setting.DatabaseSetting.Name
	user = setting.DatabaseSetting.User
	password = setting.DatabaseSetting.Password
	host = setting.DatabaseSetting.Host
	tablePrefix = setting.DatabaseSetting.TablePrefix

	db, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user,
		password,
		host,
		dbName))

	if err != nil {
		log.Println(err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}

	//全局禁用表名复数，不禁用的type User struct对应users表，开启之后对应的就是user表
	db.SingularTable(true)

	//设置最大空闲连接数
	db.DB().SetMaxIdleConns(10)

	//设置最大的打开连接数
	db.DB().SetMaxOpenConns(100)

	//注册callback代替指定的callback(BeforeCreate/BeforeUpdate)
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)
}

//关闭数据库连接
func CloseDB() {
	defer db.Close()
}

//代替gorm自带的创建时，设置指定字段
func updateTimeStampForCreateCallback(scope *gorm.Scope) {

	//判断当前操作数据库是否有错误
	if !scope.HasError() {

		//返回当前时间的时间戳
		nowTime := time.Now().Unix()

		//获取当前操作的所有字段，判断是否包含验证的字段
		//scope.Fields()其实就是对应的model定义的struct字段
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {

			//字段存在情况下，判断字段的值是否为空
			if createTimeField.IsBlank {

				//值为空的情况，设置字段的值为当前时间
				createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

//代替gorm自带的更新时的callback（更新某字段的时间）
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {

	//这里就是判断，有没有额外的设置modified字段为额外更新字段，没有的话，更新modified_on字段为当前时间
	if value, ok := scope.Get("gorm:modified_on"); !ok {
		fmt.Println(value,1)
		fmt.Println(ok,2)
		scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

//给gorm的删除操作，增加自定义的callback
func deleteCallback(scope *gorm.Scope) {

	//判断数据库操作是否出错
	if !scope.HasError() {
		var extraOption string

		//检测有没有手动指定删除字段delete_option，这里是没有的
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		//获取我们约定的删除字段，若存在则 UPDATE 软删除，若不存在则 DELETE 硬删除
		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")

		//判断是否有软删除的字段
		if !scope.Search.Unscoped && hasDeletedOnField {
			scope.Raw(fmt.Sprintf(//软删除字段存在的情况
				"UPDATE %v SET %v=%v%v%v",

				//返回引用的表名，这个方法 GORM 会根据自身逻辑对表名进行一些处理
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),

				//添加值作为sql的参数，也可以用来防止sql的注入，对应format的第三个%v
				scope.AddToVars(time.Now().Unix()),

				// scope.CombineConditionSql()返回组合好的条件SQL
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {//软删除字段不存在的情况
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

//拼接sql操作
func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
