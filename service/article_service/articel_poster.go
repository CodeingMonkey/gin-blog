package article_service

import (
	"fmt"
	"github.com/gin-blog/pkg/setting"
	"github.com/golang/freetype"
	"image"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
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

	if !a.CheckMergedImage(path) {//检查海报图片是否存在，不存在执行if代码
		fmt.Println("海报图片不存在")
		mergedF, err := a.OpenMergedImage(path)//打开合并后的图片（海报图片），没有文件就创建文件
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

		draw.Draw(jpg, jpg.Bounds(), bgImage, bgImage.Bounds().Min, draw.Over)//把背景图的内容写到RGBA 图像
		draw.Draw(jpg, jpg.Bounds(), qrImage, qrImage.Bounds().Min.Sub(image.Pt(a.Pt.X, a.Pt.Y)), draw.Over)//把二维码内容写到RGBA 图像

		//为图片绘制文字
		err = a.DrawPoster(&DrawText{
			JPG:    jpg,
			Merged: mergedF,

			Title: "极酷社",
			X0:    80,
			Y0:    160,
			Size0: 42,

			SubTitle: "---共生体精神分裂患者",
			X1:       160,
			Y1:       220,
			Size1:    36,
		}, "msyh.ttf")//同时生成最后符合条件的海报图片

		if err != nil {
			return "", "", err
		}

		//jpeg.Encode(mergedF, jpg, nil)//把图片内容写到打开的合并图片里

	}

	return fileName, path, nil
}


type DrawText struct {
	JPG    draw.Image
	Merged *os.File

	Title string
	X0    int
	Y0    int
	Size0 float64

	SubTitle string
	X1       int
	Y1       int
	Size1    float64
}

func (a *ArticlePosterBg) DrawPoster(d *DrawText, fontName string) error {
	fontSource := setting.AppSetting.RuntimeRootPath + setting.AppSetting.FontSavePath + fontName//字体文件路径
	fontSourceBytes, err := ioutil.ReadFile(fontSource)//读取字体文件
	if err != nil {
		return err
	}

	trueTypeFont, err := freetype.ParseFont(fontSourceBytes)//返回给定的字体对象
	if err != nil {
		return err
	}

	fc := freetype.NewContext()//创建一个内容对象
	fc.SetDPI(72)//设置分辨率
	fc.SetFont(trueTypeFont)//设置字体
	fc.SetFontSize(d.Size0)//设置字体大小
	fc.SetClip(d.JPG.Bounds())//设置裁剪矩形
	fc.SetDst(d.JPG)//设置目标图像
	fc.SetSrc(image.Black)//设置绘制图像的源图像

	pt := freetype.Pt(d.X0, d.Y0)
	_, err = fc.DrawString(d.Title, pt)
	if err != nil {
		return err
	}

	fc.SetFontSize(d.Size1)
	_, err = fc.DrawString(d.SubTitle, freetype.Pt(d.X1, d.Y1))
	if err != nil {
		return err
	}

	err = jpeg.Encode(d.Merged, d.JPG, nil)
	if err != nil {
		return err
	}

	return nil
}
