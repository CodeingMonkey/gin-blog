package models

/**
struct里面的属性使用了标签，且属性为json，方便在请求返回时json处理方便
*/
type Tag struct {
	Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

/**
return后面为空，是因为函数声明的时候，声明了返回的变量，且在db操作时使用的是tags的指针，
即db操作查询出结果之后，就把结果赋值给了tags，return在把tags返回
*/
func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag) {
	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)

	return
}

func GetTagTotal(maps interface{}) (count int) {
	db.Model(&Tag{}).Where(maps).Count(&count)

	return
}

func ExistTagByName(name string) bool {
	var tag Tag
	//找到符合条件的根据主键排序的第一条记录
	db.Select("id").Where("name = ?", name).First(&tag)
	if tag.ID > 0 {
		return true
	}

	return false
}

func AddTag(name string, state int, createdBy string) bool {
	db.Create(&Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	})

	return true
}

/**
gorm的Callbacks
可以将回调方法定义为模型结构的指针，在创建、更新、查询、删除时将被调用，如果任何回调返回错误，gorm将停止未来操作并回滚所有更改

func和函数名称之间增加了接收者，这种函数叫做方法
指针作为接收者的方法，指针和值都可以调用
值作为接收者的方法，只有值能调用这个方法
*/
//func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
//	scope.SetColumn("CreatedOn", time.Now().Unix())
//
//	return nil
//}
//
//func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
//	scope.SetColumn("ModifiedOn", time.Now().Unix())
//
//	return nil
//}

func ExistTagByID(id int) bool {
	var tag Tag
	db.Select("id").Where("id = ?", id).First(&tag)
	if tag.ID > 0 {
		return true
	}

	return false
}

func DeleteTag(id int) bool {
	db.Where("id = ?", id).Delete(&Tag{})

	return true
}

func EditTag(id int, data interface{}) bool {
	//下面是设置了额外的更新字段modified_on，可以用来验证model.go的代替gorm自带的callback的函数
	//db.Model(&Tag{}).Set("gorm:modified_on", "OPTION (OPTIMIZE FOR UNKNOWN)").Where("id = ?", id).Updates(data)
	db.Model(&Tag{}).Where("id = ?", id).Updates(data)

	return true
}

func CleanAllTag() bool {
	if err := db.Unscoped().Where("deleted_on != ? ", 0).Delete(&Tag{}); err != nil {
		return false
	}

	return true
}

func GetExportTagsData(maps interface{}) (articles []Tag, err error) {

	/**
	Preload就是一个预加载器，通过Preload会查询两条sql，并将查询的tag的结果集嵌入到Article的Tag中
	*/
	err = db.Where(maps).Find(&articles).Error

	return
}

