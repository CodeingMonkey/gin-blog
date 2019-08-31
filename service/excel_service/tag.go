package excel_service

import (
	"github.com/gin-blog/pkg/e"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-blog/pkg/export"
	selfFile "github.com/gin-blog/pkg/file"
	"github.com/gin-blog/service/tag_service"
	//"github.com/tealeg/xlsx"
	//"strconv"
	"time"
)

type Tag struct {
	State int
	Name  string
}

func (t Tag) exportTag() (filename string, error e.CustomizeError) {
	tagService := tag_service.Tag{
		Name:  t.Name,
		State: t.State,
	}

	tags, customizeError := tagService.GetExportData()
	if customizeError != nil {
		error.ErrorCode = e.ERROR_OPERATE_DATABASE
		return "", error
	}

	file := excelize.NewFile()
	index := file.NewSheet("标签信息")

	a := 1

	headers := map[string]string{"A1": "ID", "B1": "名称", "C1": "创建人", "D1": "创建时间", "E1": "修改人", "F1": "修改时间"}

	for k, v := range headers {
		file.SetCellValue("标签信息", k, v)
	}

	a++

	file.SetActiveSheet(index)

	for _, value := range tags {
		file.SetCellValue("标签信息", "A"+strconv.Itoa(a), value.ID)
		file.SetCellValue("标签信息", "B"+strconv.Itoa(a), value.Name)
		file.SetCellValue("标签信息", "C"+strconv.Itoa(a), value.CreatedBy)
		file.SetCellValue("标签信息", "D"+strconv.Itoa(a), strconv.Itoa(value.CreatedOn))
		file.SetCellValue("标签信息", "E"+strconv.Itoa(a), value.ModifiedBy)
		file.SetCellValue("标签信息", "F"+strconv.Itoa(a), strconv.Itoa(value.ModifiedOn))
		a++

	}

	time := strconv.Itoa(int(time.Now().Unix()))
	filename = "tags-" + time + ".xlsx"

	fullPath := export.GetExcelFullPath() + filename

	//判断文件夹是否存在，不存在的话，创建文件夹
	selfFile.IsNotExistMkDir(export.GetExcelFullPath())

	err := file.SaveAs(fullPath)
	if err != nil {
		error.ErrorCode = e.EROOR_SAVE_FILE
		return "", error
	}

	return
}
