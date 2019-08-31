package excel_service

import (
	"github.com/gin-blog/pkg/e"
	"github.com/gin-blog/pkg/export"
	selfFile "github.com/gin-blog/pkg/file"
	"github.com/gin-blog/service/article_service"
	"github.com/tealeg/xlsx"
	"strconv"
	"time"
)

type Article struct {
	State int
	Name  string
}

func (a Article) exportArticle() (filename string, error e.CustomizeError) {
	articleService := article_service.Article{
		Title: a.Name,
		State: a.State,
	}

	articles, customizeError := articleService.GetExportData()
	if customizeError != nil {
		error.ErrorCode = e.ERROR_OPERATE_DATABASE
		return "", error
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("文件信息")
	if err != nil {
		error.ErrorCode = e.ERROR_SAVE_HEADER
		return "", error
	}

	titles := []string{"ID", "标题", "创建人", "创建时间", "修改人", "修改时间", "内容"}
	row := sheet.AddRow()

	var cell *xlsx.Cell
	for _, title := range titles {
		cell = row.AddCell()
		cell.Value = title
	}

	for _, v := range articles {
		values := []string{
			strconv.Itoa(v.ID),
			v.Title,
			v.CreatedBy,
			strconv.Itoa(v.CreatedOn),
			v.ModifiedBy,
			strconv.Itoa(v.ModifiedOn),
			v.Content,
		}

		row = sheet.AddRow()
		for _, value := range values {
			cell = row.AddCell()
			cell.Value = value
		}
	}

	time := strconv.Itoa(int(time.Now().Unix()))
	filename = "articles-" + time + ".xlsx"

	fullPath := export.GetExcelFullPath() + filename

	//判断文件夹是否存在，不存在的话，创建文件夹
	selfFile.IsNotExistMkDir(export.GetExcelFullPath())

	err = file.Save(fullPath)
	if err != nil {
		error.ErrorCode = e.EROOR_SAVE_FILE
		return "", error
	}

	return

}
