package tag_service

import "github.com/gin-blog/models"

//设置属性
type Tag struct {
	ID            int
	TagID         int
	Title         string
	Name       string
	State         int
	CreatedBy     string
	ModifiedBy    string

	PageNum  int //每页显示条数
	PageSize int //页码
}

//获取导出的数据
func (t *Tag) GetExportData()([]models.Tag, error)  {
	tags, err := models.GetExportTagsData(t.getMaps())

	if err != nil {
		return nil, err
	}
	return tags, nil
}

//封装搜索条件
func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	maps["state"] = 1
	if t.Name != "" {
		maps["Name"] = t.Name
	}

	if t.State != 1 {
		maps["state"] = t.State
	}

	return maps
}

