package v1

import (
	"github.com/EDDYCJY/go-gin-example/service/tag_service"
	"github.com/Unknwon/com"
	"github.com/gin-blog/pkg/app"
	"github.com/gin-blog/pkg/e"
	"github.com/gin-blog/pkg/export"
	"github.com/gin-gonic/gin"
	"net/http"
)

//导出接口，根据各种参数，导出不同的数据
func ExportTag(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.PostForm("name")
	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := tag_service.Tag{
		Name:  name,
		State: state,
	}

	filename, err := tagService.Export()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EXPORT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"export_url":      export.GetExcelFullUrl(filename),
		"export_save_url": export.GetExcelPath() + filename,
	})
}
