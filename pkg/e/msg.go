package e

/**
声明MsgFlags为map（映射），key为int类型，value为string类型
 */
var MsgFlags = map[int]string {
	SUCCESS : "ok",
	ERROR : "fail",
	INVALID_PARAMS : "请求参数错误",
	ERROR_EXIST_TAG : "已存在该标签名称",
	ERROR_NOT_EXIST_TAG : "该标签不存在",
	ERROR_NOT_EXIST_ARTICLE : "该文章不存在",
	ERROR_AUTH_CHECK_TOKEN_FAIL : "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT : "Token已超时",
	ERROR_AUTH_TOKEN : "Token生成失败",
	ERROR_AUTH : "Token错误",
	ERROR_UPLOAD_SAVE_IMAGE_FAIL : "保存图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FAIL : "检测图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT : "校验图片错误，图片格式或大小有问题",
	ERROR_CHECK_EXIST_ARTICLE_FAIL : "文章不存在",
	ERROR_GET_ARTICLE_FAIL : "获取文章失败",
	ERROR_OPERATE_DATABASE : "操作数据库失败",
	ERROR_NOT_TABEL : "数据表不存在",
	ERROR_SAVE_HEADER : "文件头保存失败",
	EROOR_SAVE_FILE : "文件保存失败",
	ERROR_READ_FILE : "读取文件失败",
	ERROR_FORMAT_ERROR : "文件格式错误",
}

/**
传来的code在MsgFlags中不存在的情况，会返回Token错误
 */
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}