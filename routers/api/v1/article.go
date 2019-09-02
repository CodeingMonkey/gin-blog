package v1

import (
	"fmt"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-blog/pkg/logging"
	"github.com/gin-blog/pkg/qrcode"
	"github.com/gin-blog/service/article_service"
	"log"
	"net/http"
	"time"

	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-blog/models"
	"github.com/gin-blog/pkg/app"
	"github.com/gin-blog/pkg/e"
	"github.com/gin-blog/pkg/setting"
	"github.com/gin-blog/pkg/util"
	"github.com/gin-gonic/gin"
)

const (
	QRCODE_URL = "https://acg.app"
)

//获取单个文章
//func GetArticle(c *gin.Context) {
//	id := com.StrTo(c.Param("id")).MustInt()
//
//	valid := validation.Validation{}
//	valid.Min(id, 1, "id").Message("ID必须大于0")
//
//	code := e.INVALID_PARAMS
//	var data interface{}
//	if !valid.HasErrors() {
//		if models.ExistArticleByID(id) {
//			data = models.GetArticle(id)
//			code = e.SUCCESS
//		} else {
//			code = e.ERROR_NOT_EXIST_ARTICLE
//		}
//	} else {
//		/**
//		为什会遍历去查询处理错误，valid.Errors是切片类型，整套流程走下来遇到多少验证错误，都会赛到slice中
//		*/
//		for _, err := range valid.Errors {
//			log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
//		}
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"code": code,
//		"msg":  e.GetMsg(code),
//		"data": data,
//	})
//}

//获取单个文章
func GetArticle(c *gin.Context) {
	appG := app.Gin{C: c} //实例化响应对象
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	//验证有没有错误
	if valid.HasErrors() {

		//记录验证错误日志
		app.MarkErrors(valid.Errors)

		//请求返回
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return //？？上一步都返回数据了，为什么还要return
	}

	//验证数据是否存在，查找数据操作移动到service层
	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID() //判断数据库中文章是否存在
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	article, err := articleService.Get() //数据库文章存在的情况，从数据库获取文章
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}

	//响应返回
	appG.Response(http.StatusOK, e.SUCCESS, article)
}

//获取多个文章
func GetArticles(c *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	valid := validation.Validation{}
	time.Sleep(20 * time.Second)

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state

		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	var tagId int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId

		valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	}

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS

		data["lists"] = models.GetArticles(util.GetPage(c), setting.AppSetting.PageSize, maps)
		data["total"] = models.GetArticleTotal(maps)

	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

//新增文章
func AddArticle(c *gin.Context) {
	tagId := com.StrTo(c.PostForm("tag_id")).MustInt()
	title := c.PostForm("title")
	desc := c.PostForm("desc")
	content := c.PostForm("content")
	createdBy := c.PostForm("created_by")
	coverImageUrl := c.PostForm("cover_image_url")
	state := com.StrTo(c.DefaultPostForm("state", "0")).MustInt()

	valid := validation.Validation{}
	//valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(desc, "desc").Message("简述不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	valid.Required(coverImageUrl, "cover_image_url").Message("图片地址不能为空")
	valid.MaxSize(coverImageUrl, 100, "cover_image_url").Message("图片地址不能超过100")

	fmt.Println("图片地址")
	fmt.Println(coverImageUrl)

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistTagByID(tagId) {
			data := make(map[string]interface{})
			data["tag_id"] = tagId
			data["title"] = title
			data["desc"] = desc
			data["content"] = content
			data["created_by"] = createdBy
			data["state"] = state
			data["cover_image_url"] = coverImageUrl

			models.AddArticle(data)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}

//修改文章
func EditArticle(c *gin.Context) {
	valid := validation.Validation{}

	//params函数用来接收url后的/id参数
	id := com.StrTo(c.Param("id")).MustInt()
	tagId := com.StrTo(c.PostForm("tag_id")).MustInt()
	title := c.PostForm("title")
	desc := c.PostForm("desc")
	content := c.PostForm("content")
	modifiedBy := c.PostForm("modified_by")

	var state int = -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Min(id, 1, "id").Message("ID必须大于0")
	valid.MaxSize(title, 100, "title").Message("标题最长为100字符")
	valid.MaxSize(desc, 255, "desc").Message("简述最长为255字符")
	valid.MaxSize(content, 65535, "content").Message("内容最长为65535字符")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if res, _ := models.ExistArticleByID(id); res {
			if models.ExistTagByID(tagId) {
				data := make(map[string]interface{})
				if tagId > 0 {
					data["tag_id"] = tagId
				}
				if title != "" {
					data["title"] = title
				}
				if desc != "" {
					data["desc"] = desc
				}
				if content != "" {
					data["content"] = content
				}

				data["modified_by"] = modifiedBy

				models.EditArticle(id, data)
				code = e.SUCCESS
			} else {
				code = e.ERROR_NOT_EXIST_TAG
			}
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

//删除文章
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if res, _ := models.ExistArticleByID(id); res {
			models.DeleteArticle(id)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

/*func GenerateArticlePoster(c *gin.Context) {
	appG := app.Gin{C: c}
	qrc := qrcode.NewQrCode(QRCODE_URL, 300, 300, qr.M, qr.Auto)
	path := qrcode.GetQrCodeFullPath()//二维码保存的完整路径
	_, _, err := qrc.Encode(path)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
*/

func GenerateArticlePoster(c *gin.Context) {
	appG := app.Gin{C:c}
	article := &article_service.Article{}
	selfQr := qrcode.NewQrCode(QRCODE_URL, 300, 300, qr.M, qr.Auto) // 二维码对象
	posterName := article_service.GetPosterFlag() + "-" + qrcode.GetQrCodeFileName(selfQr.URL) + selfQr.GetQrCodeExt()//获取海报文件名称（文件名+后缀）
	articlePoster := article_service.NewArticlePoster(posterName, article, selfQr)//实例化海报对象（包含海报名称+文章对象+二维码对象）
	articlePosterBgService := article_service.NewArticlePosterBg(//实例化背景对象（背景图片+海报对象）
		"bg.jpeg",
		articlePoster,
		&article_service.Rect{
			X0: 0,
			Y0: 0,
			X1: 550,
			Y1: 700,
		},
		&article_service.Pt{
			X: 125,
			Y: 298,
		},
	)

	_, filePath, err := articlePosterBgService.Generate()//生成海报（二维码图片+背景图片）
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GEN_ARTICLE_POSTER_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"poster_url":      qrcode.GetQrCodeFullUrl(posterName),//用来请求的图片名称
		"poster_save_url": filePath + posterName,//文件保存ID孩子+文件名称
	})
}
