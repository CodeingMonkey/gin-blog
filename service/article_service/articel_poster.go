package article_service

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"

	"github.com/gin-blog/pkg/file"
	"github.com/gin-blog/pkg/qrcode"
)

type ArticlePoster struct {
	PosterName string//海报
	*Article//文章
	Qr *qrcode.QrCode//二维码
}

//实例化一个新的文章海报
func NewArticlePoster(posterName string, article *Article, qr *qrcode.QrCode) *ArticlePoster {
	return &ArticlePoster{
		PosterName: posterName,
		Article:    article,
		Qr:         qr,
	}
}


func GetPosterFlag() string {
	return "poster"
}

//检查海报图片是否存在
func (a *ArticlePoster) CheckMergedImage(path string) bool {
	if file.CheckNotExist(path+a.PosterName) == true {
		return false
	}

	return true
}

//打开一个海报图片，图片不存在的话，创建一个海报图片
func (a *ArticlePoster) OpenMergedImage(path string) (*os.File, error) {
	f, err := file.MustOpen(a.PosterName, path)
	if err != nil {
		return nil, err
	}

	return f, nil
}

//定义文章的海报背景图
type ArticlePosterBg struct {
	Name string
	*ArticlePoster
	*Rect
	*Pt
}

type Rect struct {
	Name string
	X0   int
	Y0   int
	X1   int
	Y1   int
}

type Pt struct {
	X int
	Y int
}

//创建一个文章的海报背景图的实例化对象
func NewArticlePosterBg(name string, ap *ArticlePoster, rect *Rect, pt *Pt) *ArticlePosterBg {
	return &ArticlePosterBg{
		Name:          name,
		ArticlePoster: ap,
		Rect:          rect,
		Pt:            pt,
	}
}

//创建一个使用了背景图的文章海报二维码
func (a *ArticlePosterBg) Generate() (string, string, error) {
	fullPath := qrcode.GetQrCodeFullPath()//获取完整的二维码存储路径
	fmt.Println("图片保存完整路径", fullPath)

	//可以在此步骤中，查询对应的文章，获取文章的标题或者名称，组合到对应的二维码文件名中，目前代码中所有文件的二维码都是同一个
	fileName, path, err := a.Qr.Encode(fullPath)//生成二维码文件(到指定的文件夹)，子结构嵌入到父结构中，提升了子结构的属性，可以用父结构对象直接去调用
	if err != nil {
		return "", "", err
	}

	if !a.CheckMergedImage(path) {//检查背景图片是否存在，不存在执行if代码
		mergedF, err := a.OpenMergedImage(path)//打开合并后的图片，没有文件就创建文件
		if err != nil {
			return "", "", err
		}
		defer mergedF.Close()//关闭合并后的图片

		bgF, err := file.MustOpen(a.Name, path)//打开背景图片
		if err != nil {
			return "", "", err
		}
		defer bgF.Close()//关闭背景图片

		qrF, err := file.MustOpen(fileName, path)//打开已经生成的二维码图片
		if err != nil {
			return "", "", err
		}
		defer qrF.Close()//关闭生成的二维码图片

		bgImage, err := jpeg.Decode(bgF)//读取背景图片文件的内容
		if err != nil {
			return "", "", err
		}
		qrImage, err := jpeg.Decode(qrF)//读取二维码图片文件的内容
		if err != nil {
			return "", "", err
		}

		jpg := image.NewRGBA(image.Rect(a.Rect.X0, a.Rect.Y0, a.Rect.X1, a.Rect.Y1))//根据指定的边界返回一个新的RGBA 图像

		draw.Draw(jpg, jpg.Bounds(), bgImage, bgImage.Bounds().Min, draw.Over)//把背景图的内容写到新图片里
		draw.Draw(jpg, jpg.Bounds(), qrImage, qrImage.Bounds().Min.Sub(image.Pt(a.Pt.X, a.Pt.Y)), draw.Over)//把二维码内容写到新图片里

		jpeg.Encode(mergedF, jpg, nil)//把图片内容写到打开的合并图片里
	}

	return fileName, path, nil
}