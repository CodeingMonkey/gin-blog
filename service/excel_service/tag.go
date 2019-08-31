package excel_service

import (
	"github.com/gin-blog/pkg/e"
)

type Tag struct {
	State int
	Name  string
}

func (t Tag) exportTag() (string, e.CustomizeError) {
	var filename  string
	var err e.CustomizeError

	return  filename,err
}
