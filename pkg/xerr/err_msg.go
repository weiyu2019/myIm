package xerr

var codeToText map[int]string = map[int]string{
	SERVER_COMMON_ERROR: "服务器异常",
	REQUEST_PARAM_ERROR: "请求参数有误",
	DB_ERROR:            "数据库繁忙",
}

func ErrMsg(errCode int) string {
	if msg, ok := codeToText[errCode]; ok {
		return msg
	}
	return ""
}
