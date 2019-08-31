package article_service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-blog/models"
	"github.com/gin-blog/pkg/gredis"
	"github.com/gin-blog/pkg/logging"
	"github.com/gin-blog/service/cache_service"
)

//设置属性
type Article struct {
	ID            int
	TagID         int
	Title         string
	Desc          string
	Content       string
	CoverImageUrl string
	State         int
	CreatedBy     string
	ModifiedBy    string

	PageNum  int //每页显示条数
	PageSize int //页码
}

//获取文章（redis不存在读取数据库）
func (a *Article) Get() (*models.Article, error) {
	var cacheArticle *models.Article

	cache := cache_service.Article{ID: a.ID}
	key := cache.GetArticleKey() //获取redisKey
	if gredis.Exists(key) {      //判断key是否存在

		fmt.Println("key is exist")
		data, err := gredis.Get(key) //key存在的情况，从redis获取数据
		if err != nil {              //redis数据不存在的情况，记录日志
			logging.Info(err)
		} else { //redis存在，直接从redis中获取数据，解析json，并且返回
			json.Unmarshal(data, &cacheArticle)
			return cacheArticle, nil
		}
	}

	//redis不存在的情况，从数据库查询文章
	article, err := models.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}

	//设置redis数据，有效期为1个小时
	gredis.Set(key, article, 3600)
	return article, nil
}

func (a *Article) ExistByID() (bool, error) {
	return models.ExistArticleByID(a.ID)
}

//获取导出的数据
func (a *Article) GetExportData()([]models.Article, error)  {
	article, err := models.GetExportData(a.getMaps())

	if err != nil {
		return nil, err
	}
	return article, nil
}

//封装搜索条件
func (a *Article) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	maps["state"] = 1
	if a.Title != "" {
		maps["title"] = a.Title
	}

	if a.State != 1 {
		maps["state"] = a.State
	}

	return maps
}
