package file

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
)

//获取文件大小
func GetSize(f multipart.File) (int, error) {
	content, err := ioutil.ReadAll(f)

	return len(content), err
}

//获取文件后缀
func GetExt(fileName string) string {
	return path.Ext(fileName)
}

//检测文件是否存在
func CheckExist(src string) bool {
	_, err := os.Stat(src)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

//检测文件是否存在
func CheckNotExist(src string) bool {
	_, err := os.Stat(src)

	return os.IsNotExist(err)
}

//检测权限
func CheckPermission(src string) bool {
	_, err := os.Stat(src)

	return os.IsPermission(err)
}

//判断文件夹是否存在，不存在则创建文件夹
func IsNotExistMkDir(src string) error {
	fmt.Println("目录11", src)
	if exist := CheckExist(src); exist == false {
		if err := MkDir(src); err != nil {
			return err
		}
	}

	return nil
}

//创建文件夹
func MkDir(src string) error {
	fmt.Println("创建目录", src)
	err := os.MkdirAll(src, os.ModePerm)
	fmt.Println("创建目录结果", err)
	if err != nil {
		return err
	}

	return nil
}

//打开文件
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// MustOpen maximize trying to open the file
func MustOpen(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()//返回根路径，当前目录的根路径有多个的话，随机返回一个（为什么会有多个根路径，大概是软连接）
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + "/" + filePath
	perm := CheckPermission(src)//校验权限
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err = IsNotExistMkDir(src)//判断文件夹是否存在，不存在则创建文件夹
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	f, err := Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)//打开文件，权限为0644(中间的参数是根据位运算符|计算得的)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}

	return f, nil
}
