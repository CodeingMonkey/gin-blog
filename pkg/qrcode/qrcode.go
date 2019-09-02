package qrcode

import (
	"image/jpeg"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"

	"github.com/gin-blog/pkg/file"
	"github.com/gin-blog/pkg/setting"
	"github.com/gin-blog/pkg/util"
)

type QrCode struct {
	URL    string//二维码图片地址
	Width  int//宽度
	Height int//高度
	Ext    string//后缀
	Level  qr.ErrorCorrectionLevel//二维码存储备份数据的数量
	Mode   qr.Encoding//二维码编码格式
}

//二维码图片后缀
const (
	EXT_JPG = ".jpg"
)

//创建一张二维码
func NewQrCode(url string, width, height int, level qr.ErrorCorrectionLevel, mode qr.Encoding) *QrCode {
	return &QrCode{
		URL:    url,
		Width:  width,
		Height: height,
		Level:  level,
		Mode:   mode,
		Ext:    EXT_JPG,
	}
}

//获取二维码保存路径
func GetQrCodePath() string {
	return setting.AppSetting.QrCodeSavePath
}

//获取二维码完整的保存路径
func GetQrCodeFullPath() string {
	return setting.AppSetting.RuntimeRootPath + setting.AppSetting.QrCodeSavePath
}

//获取图片请求路径（加上域名或者IP）
func GetQrCodeFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetQrCodePath() + name
}

//获取二维码文件名称（MD5处理过）
func GetQrCodeFileName(value string) string {
	return util.EncodeMD5(value)
}

//获取文件后缀（文件后缀这个应该是可以写到app.ini配置文件的）
func (q *QrCode) GetQrCodeExt() string {
	return q.Ext
}

//检测二维码文件是否存在
func (q *QrCode) CheckEncode(path string) bool {
	src := path + GetQrCodeFileName(q.URL) + q.GetQrCodeExt()
	if file.CheckNotExist(src) == true {
		return false
	}

	return true
}

//生成二维码
func (q *QrCode) Encode(path string) (string, string, error) {
	name := GetQrCodeFileName(q.URL) + q.GetQrCodeExt()//二维码文件名
	src := path + name
	if file.CheckNotExist(src) == true {//校验文件是否不存在，不存在的话，创建文件
		code, err := qr.Encode(q.URL, q.Level, q.Mode)//返回使用给定内容，给定等级，给定编码方式处理过的条形码
		if err != nil {
			return "", "", err
		}

		code, err = barcode.Scale(code, q.Width, q.Height)//返回根据宽度和高度处理过的条形码
		if err != nil {
			return "", "", err
		}

		f, err := file.MustOpen(name, path)//打开文件，没有文件则创建文件
		if err != nil {
			return "", "", err
		}
		defer f.Close()//函数结束之后，关闭打开的文件

		err = jpeg.Encode(f, code, nil)//条形码写入文件
		if err != nil {
			return "", "", err
		}
	}

	return name, path, nil
}
