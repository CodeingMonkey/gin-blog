package excel_service

import "github.com/gin-blog/pkg/e"

type Excel struct {
	TableName string
	Name      string
	State     int
}

func (excel Excel) Export() (filename string, error e.CustomizeError) {

	switch excel.TableName {
	case "tag": //导出标签
		tag := Tag{Name: excel.Name, State: excel.State}
		filename, error = tag.exportTag()
	case "article": //导出文章
		article := Article{State: excel.State, Name: excel.Name}
		filename, error = article.exportArticle()
	default: //不是指定数据表情况，返回错误信息
		error.ErrorCode = e.ERROR_NOT_TABEL
	}

	return
}
