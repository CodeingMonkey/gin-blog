package api

import (
	"github.com/Unknwon/com"
	"github.com/gin-blog/pkg/app"
	"github.com/gin-blog/pkg/e"
	"github.com/gin-blog/pkg/export"
	"github.com/gin-blog/service/excel_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ExportCsv(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.PostForm("name")
	tableName := c.PostForm("tableName") //表类型，表明想要导出的是tag Or article
	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	excelService := excel_service.Excel{
		Name:      name,
		State:     state,
		TableName: tableName,
	}

	//service层实现具体的导出操作
	filename, err := excelService.Export()
	if err.ErrorCode != 0 {
		appG.Response(http.StatusOK, err.ErrorCode, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"export_url":      export.GetExcelFullUrl(filename),//导出文件请求路径
		"export_save_url": export.GetExcelPath() + filename,//导出文件保存路径
	})
}